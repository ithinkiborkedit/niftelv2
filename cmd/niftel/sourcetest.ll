declare i32 @printf(i8*,...)
	@print_str_format = constant [4 x i8] c"%s\0A\00"
	@print_int_format = constant [4 x i8] c"%d\0A\00"
	@print_float_format = constant [4 x i8] c"%f\0A\00"
	@print_str_open_brace = private constant [2 x i8] c"{\00"
	@print_str_close_brace = private constant [2 x i8] c"}\00"
	@print_str_comma = private constant [3 x i8] c", \00"
	%Person = type { i8*, i64 }
@.str0 = private constant [5 x i8] c"test\00"
define i32 @main(){
entry:
 %t0 = alloca %Person
 %t1 = alloca %Person
%t2 = getelementptr %Person, %Person* %t1, i32 0, i32 0
store i8* getelementptr ([5 x i8], [5 x i8]* @.str0, i32 0, i32 0), i8** %t2
%t3 = getelementptr %Person, %Person* %t1, i32 0, i32 1
store i64 2, i64* %t3
 %t4 = getelementptr %Person, %Person* %t0, i32 0, i32 0
 %t5 = getelementptr %Person, %Person* %t1, i32 0, i32 0
 %t6 = load i8*, i8** %t5
 store i8* %t6, i8** %t4
 %t7 = getelementptr %Person, %Person* %t0, i32 0, i32 1
 %t8 = getelementptr %Person, %Person* %t1, i32 0, i32 1
 %t9 = load i64, i64* %t8
 store i64 %t9, i64* %t7
call i32 (i8*,...) @printf(i8* getelementptr ([2 x i8], [2 x i8]* @print_str_open_brace, i32 0, i32 0))
 %t10 = getelementptr %Person, %Person* %t0, i32 0, i32 0
%t11 = load i8*, i8** %t10
call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @print_str_format, i32 0, i32 0), i8* %t11)
call i32 (i8*,...) @printf(i8* getelementptr ([2 x i8], [2 x i8]* @print_str_comma, i32 0, i32 0))
 %t12 = getelementptr %Person, %Person* %t0, i32 0, i32 1
%t13 = load i64, i64* %t12
call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @print_int_format, i32 0, i32 0), i64 %t13)
call i32 (i8*,...) @printf(i8* getelementptr ([2 x i8], [2 x i8]* @print_str_close_brace, i32 0, i32 0))

 ret i32 0
}
