declare i32 @printf(i8*,...)
	@print_str_format = constant [4 x i8] c"%s\0A\00"
	@print_int_format = constant [4 x i8] c"%d\0A\00"
	@.str0 = private constant [6 x i8] c"hello\00"
define i32 @main(){
entry:
call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8],[4 x i8]* @print_str_format, i32 0, i32 0), i8* getelementptr ([6 x i8], [6 x i8]* @.str0, i32 0, i32 0))
call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @print_int_format, i32 0, i32 0), i32 42)
 ret i32 0
}
