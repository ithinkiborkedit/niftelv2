package codegen

import (
	"fmt"
	"strings"

	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	tokens "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
	"github.com/ithinkiborkedit/niftelv2.git/internal/tokentoval"
)

type Codegen struct {
	entryBuilder strings.Builder
	builder      strings.Builder
	strings      map[string]string
	structs      structTypes
	nextStrIndex int
	symbols      map[string]VariableInfo
	// symbolTypes  map[string]string
	nextReg int
}

func NewCodeGen() *Codegen {
	return &Codegen{
		strings: make(map[string]string),
		symbols: make(map[string]VariableInfo),
		structs: NewStructTypes(),
		nextReg: 0,
	}
}

func (c *Codegen) RegisterStruct(name string, fieldNames []string, fieldTypes []string) {
	fieldIndices := make(map[string]int)
	for i, fname := range fieldNames {
		fieldIndices[fname] = i
	}
	info := &StructTypeInfo{
		Name:         name,
		FieldNames:   fieldNames,
		FieldTypes:   fieldTypes,
		FieldIndices: fieldIndices,
		Emitted:      false,
	}
	c.structs.Register(info)
}

func (c *Codegen) EmitStructDefinitions() {
	for _, info := range c.structs {
		if info.Emitted {
			continue
		}
		c.emitStructDefinition(info)
		c.structs.MarkEmitted(info.Name)
	}
}

func (c *Codegen) emitStructDefinition(info *StructTypeInfo) {
	fieldsLLVMTypes := make([]string, len(info.FieldTypes))
	copy(fieldsLLVMTypes, info.FieldTypes)
	c.builder.WriteString(fmt.Sprintf("%%%s = type { %s }\n", info.Name, strings.Join(fieldsLLVMTypes, ", ")))
}

func (c *Codegen) GenerateLLVM(stmts []ast.Stmt) (string, error) {
	c.emitPreamble()
	for _, stmt := range stmts {
		if s, ok := stmt.(*ast.StructStmt); ok {
			c.emitStructStmt(s)
		} else {
			c.collectStringsStmt(stmt)
		}
	}

	c.emitStringConstants()
	c.emitMainFunction(stmts)
	return c.builder.String(), nil
}

func (c *Codegen) registerStringLiteral(s string) string {
	if name, ok := c.strings[s]; ok {
		return name
	}
	name := fmt.Sprintf("@.str%d", c.nextStrIndex)
	c.strings[s] = name
	c.nextStrIndex++
	return name
}

func (c *Codegen) emitStringConstants() {
	for s, name := range c.strings {
		length := len(s) + 1
		escaped := escapeStringForLLVM(s)
		c.builder.WriteString(fmt.Sprintf("%s = private constant [%d x i8] c\"%s\\00\"\n", name, length, escaped))
	}
}

func (c *Codegen) freshReg() string {
	reg := fmt.Sprintf("%%t%d", c.nextReg)
	c.nextReg++

	return reg
}

func (c *Codegen) emitExpr(e ast.Expr) (string, string) {
	switch expr := e.(type) {
	case *ast.LiteralExpr:
		val, err := tokentoval.Convert(expr.Value)
		if err != nil {
			panic(err)
		}
		llvmLiteral := c.emitValueLiteral(val)
		parts := strings.SplitN(llvmLiteral, " ", 2)
		if len(parts) != 2 {
			panic("invalid llvm literal")
		}

		return parts[1], parts[0]
	case *ast.VariableExpr:
		reg, typ := c.loadVariable(expr.Name.Lexeme)
		return reg, typ.LLVMType
	case *ast.StructLiteralExpr:
		return c.emitStructLiteralExpr(expr)
	default:
		panic(fmt.Sprintf("unsupported expr %T", expr))
	}
}

func (c *Codegen) loadVariable(name string) (string, VariableInfo) {
	allocaReg, ok := c.symbols[name]
	if !ok {
		panic("undefined varaible: " + name)
	}
	llvmType := c.symbols[name]
	loadReg := c.freshReg()

	c.builder.WriteString(fmt.Sprintf(" %s = load %s, %s* %s\n", loadReg, llvmType, llvmType, allocaReg))
	return loadReg, llvmType
}

func (c *Codegen) emitVarStmt(s *ast.VarStmt) {
	name := s.Names[0].Lexeme
	llvmType := c.llvmTypeForTypeExpr(s.Type)

	allocaReg := c.freshReg()
	c.builder.WriteString(fmt.Sprintf(" %s = alloca %s\n", allocaReg, llvmType))

	initValReg, initValType := c.emitExpr(s.Init)

	c.builder.WriteString(fmt.Sprintf(" store %s %s, %s* %s\n", initValType, initValReg, initValType, allocaReg))

	c.symbols[name] = VariableInfo{
		LLVMName: allocaReg,
		LLVMType: llvmType,
	}
}

func (c *Codegen) collectStringsStmt(s ast.Stmt) {
	switch stmt := s.(type) {
	case *ast.PrintStmt:
		c.collectStringsExpr(stmt.Expr)
	case *ast.VarStmt:
		c.collectStringsExpr(stmt.Init)
	}
}

func (c *Codegen) collectStringsExpr(e ast.Expr) {
	switch expr := e.(type) {
	case *ast.LiteralExpr:
		if expr.Value.Type == tokens.TokenString {
			c.registerStringLiteral(expr.Value.Lexeme)
		}
	case *ast.StructLiteralExpr:
		for _, fieldExpr := range expr.Fields {
			c.collectStringsExpr(fieldExpr)
		}
	}
}

func escapeStringForLLVM(s string) string {
	var sb strings.Builder

	for _, r := range s {
		if r == '\\' || r == '"' {
			sb.WriteRune('\\')
			sb.WriteRune(r)
		} else if r < 32 || r > 126 {
			sb.WriteString(fmt.Sprintf("\\%02X", r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func (c *Codegen) emitPreamble() {
	c.builder.WriteString(`declare i32 @printf(i8*,...)
	@print_str_format = constant [4 x i8] c"%s\0A\00"
	@print_int_format = constant [4 x i8] c"%d\0A\00"
	@print_float_format = constant [4 x i8] c"%f\0A\00"
	@print_str_open_brace = constant [2 x i8] c"}\00"
	@print_str_close_brace = constant [2 x i8] c"{\00"
	@print_str_comma = constant [3 x i8] c", \00"
	`)
}

func (c *Codegen) emitMainFunction(stmts []ast.Stmt) {
	c.builder.WriteString("define i32 @main(){\n")
	c.builder.WriteString("entry:\n")

	for _, stmt := range stmts {
		c.emitStmt(stmt)
	}

	c.builder.WriteString(" ret i32 0\n")
	c.builder.WriteString("}\n")
}

func (c *Codegen) llvmTypeForTypeExpr(typ *ast.TypeExpr) string {
	switch typ.Name.Lexeme {
	case "int":
		return "i64"
	case "float":
		return "double"
	case "bool":
		return "i1"
	case "string":
		return "i8*"
	default:
		return "%" + typ.Name.Lexeme
	}
}

func (c *Codegen) emitStructLiteralExpr(expr *ast.StructLiteralExpr) (string, string) {
	structName := expr.TypeName.Name.Lexeme
	structType := "%" + structName

	ptrReg := c.freshReg()
	c.builder.WriteString(fmt.Sprintf(" %s = alloca %s\n", ptrReg, structType))

	info, ok := c.structs[structName]
	if !ok {
		panic(fmt.Sprintf("unknown struct type %s", structName))
	}

	for fieldName, fieldExpr := range expr.Fields {
		fieldIndex, exists := info.FieldIndices[fieldName]
		if !exists {
			panic(fmt.Sprintf("unknown field %s in struct %s", fieldName, structName))
		}

		valReg, valType := c.emitExpr(fieldExpr)

		gepReg := c.freshReg()

		c.builder.WriteString(fmt.Sprintf(
			"%s = getelementptr %s, %s* %s, i32 0, i32 %d\n",
			gepReg, structType, structType, ptrReg, fieldIndex,
		))

		c.builder.WriteString(fmt.Sprintf(
			"store %s %s, %s* %s\n",
			valType, valReg, valType, gepReg,
		))

	}

	return ptrReg, structType + "*"
}

func (c *Codegen) emitStructStmt(s *ast.StructStmt) {
	if _, exists := c.structs[s.Name.Lexeme]; exists {
		return
	}

	fieldNames := []string{}
	fieldTypes := []string{}

	fieldIndices := make(map[string]int)

	for i, f := range s.Fields {
		fieldName := f.Names[0].Lexeme
		fieldType := c.llvmTypeForTypeExpr(f.Type)

		fieldNames = append(fieldNames, fieldName)
		fieldTypes = append(fieldTypes, fieldType)
		fieldIndices[fieldName] = i
	}

	st := &StructTypeInfo{
		Name:         s.Name.Lexeme,
		FieldNames:   fieldNames,
		FieldTypes:   fieldTypes,
		FieldIndices: fieldIndices,
		Emitted:      false,
	}
	c.structs.Register(st)
	c.emitStructDefinition(st)
	st.Emitted = true
}

func (c *Codegen) RegisterStructFromAST(s *ast.StructStmt) {
	if _, exists := c.structs[s.Name.Lexeme]; exists {
		return
	}
	var fieldNames []string
	var fieldTypes []string

	for _, f := range s.Fields {
		fieldNames = append(fieldNames, f.Names[0].Lexeme)
		fieldTypes = append(fieldTypes, c.llvmTypeForTypeExpr(f.Type))
	}
	c.RegisterStruct(s.Name.Lexeme, fieldNames, fieldTypes)
}

func (c *Codegen) emitStmt(s ast.Stmt) {
	switch stmt := s.(type) {
	case *ast.PrintStmt:
		c.emitPrint(stmt)
		c.builder.WriteString("\n")
	case *ast.StructStmt:
		c.emitStructStmt(stmt)
	case *ast.VarStmt:
		c.emitVarStmt(stmt)
	default:
		fmt.Printf("warning: unsupported statement type: %T\n", stmt)
	}
}

func (c *Codegen) lookupVariable(name string) (VariableInfo, bool) {
	info, ok := c.symbols[name]
	return info, ok
}

func (c *Codegen) emitPrint(s *ast.PrintStmt) {
	switch expr := s.Expr.(type) {
	case *ast.LiteralExpr:
		pl := &printableLiteralExp{lit: expr}
		if err := pl.EmitPrint(c); err != nil {
			fmt.Printf("error emitting print: %v", err)
		}
	case *ast.VariableExpr:
		pv := &printableVariableExpr{varExpr: expr}
		if err := pv.EmitPrint(c); err != nil {
			fmt.Printf("error emitting print: %v\n", err)
		}
	default:
		fmt.Printf("unsupported print expression type: %T\n", expr)
	}
}
