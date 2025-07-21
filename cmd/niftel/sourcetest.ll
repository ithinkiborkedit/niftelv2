declare i32 @printf(i8*,...)
	@print_str_format = constant [3 x i8] c"%s\00"
	@print_int_format = constant [3 x i8] c"%d\00"
	@print_float_format = constant [3 x i8] c"%f\00"
	@print_str_open_brace = private constant [2 x i8] c"{\00"
	@print_str_close_brace = private constant [2 x i8] c"}\00"
	@print_str_comma = private constant [3 x i8] c", \00"
	%Person = type { i8*, i64, i8* }
@.str0 = private constant [5 x i8] c"test\00"
@.str1 = private constant [8 x i8] c"finland\00"
@.str2 = private constant [6 x i8] c"billy\00"
@.str3 = private constant [6 x i8] c"paris\00"
define i32 @main(){
entry:
 %t0 = alloca %Person
 %t1 = alloca %Person
%t2 = getelementptr %Person, %Person* %t1, i32 0, i32 0
store i8* getelementptr ([5 x i8], [5 x i8]* @.str0, i32 0, i32 0), i8** %t2
%t3 = getelementptr %Person, %Person* %t1, i32 0, i32 1
store i64 2, i64* %t3
%t4 = getelementptr %Person, %Person* %t1, i32 0, i32 2
store i8* getelementptr ([8 x i8], [8 x i8]* @.str1, i32 0, i32 0), i8** %t4
 %t5 = getelementptr %Person, %Person* %t0, i32 0, i32 0
 %t6 = getelementptr %Person, %Person* %t1, i32 0, i32 0
 %t7 = load i8*, i8** %t6
 store i8* %t7, i8** %t5
 %t8 = getelementptr %Person, %Person* %t0, i32 0, i32 1
 %t9 = getelementptr %Person, %Person* %t1, i32 0, i32 1
 %t10 = load i64, i64* %t9
 store i64 %t10, i64* %t8
 %t11 = getelementptr %Person, %Person* %t0, i32 0, i32 2
 %t12 = getelementptr %Person, %Person* %t1, i32 0, i32 2
 %t13 = load i8*, i8** %t12
 store i8* %t13, i8** %t11
 %t14 = alloca %Person
 %t15 = alloca %Person
%t16 = getelementptr %Person, %Person* %t15, i32 0, i32 0
store i8* getelementptr ([6 x i8], [6 x i8]* @.str2, i32 0, i32 0), i8** %t16
%t17 = getelementptr %Person, %Person* %t15, i32 0, i32 1
store i64 234, i64* %t17
%t18 = getelementptr %Person, %Person* %t15, i32 0, i32 2
store i8* getelementptr ([6 x i8], [6 x i8]* @.str3, i32 0, i32 0), i8** %t18
 %t19 = getelementptr %Person, %Person* %t14, i32 0, i32 0
 %t20 = getelementptr %Person, %Person* %t15, i32 0, i32 0
 %t21 = load i8*, i8** %t20
 store i8* %t21, i8** %t19
 %t22 = getelementptr %Person, %Person* %t14, i32 0, i32 1
 %t23 = getelementptr %Person, %Person* %t15, i32 0, i32 1
 %t24 = load i64, i64* %t23
 store i64 %t24, i64* %t22
 %t25 = getelementptr %Person, %Person* %t14, i32 0, i32 2
 %t26 = getelementptr %Person, %Person* %t15, i32 0, i32 2
 %t27 = load i8*, i8** %t26
 store i8* %t27, i8** %t25
call i32 (i8*,...) @printf(i8* getelementptr ([2 x i8], [2 x i8]* @print_str_open_brace, i32 0, i32 0))
 %t28 = getelementptr %Person, %Person* %t0, i32 0, i32 0
%t29 = load i8*, i8** %t28
call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @print_str_format, i32 0, i32 0), i8* %t29)
call i32 (i8*,...) @printf(i8* getelementptr ([2 x i8], [2 x i8]* @print_str_comma, i32 0, i32 0))
 %t30 = getelementptr %Person, %Person* %t0, i32 0, i32 1
%t31 = load i64, i64* %t30
call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @print_int_format, i32 0, i32 0), i64 %t31)
call i32 (i8*,...) @printf(i8* getelementptr ([2 x i8], [2 x i8]* @print_str_comma, i32 0, i32 0))
 %t32 = getelementptr %Person, %Person* %t0, i32 0, i32 2
%t33 = load i8*, i8** %t32
call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @print_str_format, i32 0, i32 0), i8* %t33)
call i32 (i8*,...) @printf(i8* getelementptr ([2 x i8], [2 x i8]* @print_str_close_brace, i32 0, i32 0))

call i32 (i8*,...) @printf(i8* getelementptr ([2 x i8], [2 x i8]* @print_str_open_brace, i32 0, i32 0))
 %t34 = getelementptr %Person, %Person* %t14, i32 0, i32 0
%t35 = load i8*, i8** %t34
call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @print_str_format, i32 0, i32 0), i8* %t35)
call i32 (i8*,...) @printf(i8* getelementptr ([2 x i8], [2 x i8]* @print_str_comma, i32 0, i32 0))
 %t36 = getelementptr %Person, %Person* %t14, i32 0, i32 1
%t37 = load i64, i64* %t36
call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @print_int_format, i32 0, i32 0), i64 %t37)
call i32 (i8*,...) @printf(i8* getelementptr ([2 x i8], [2 x i8]* @print_str_comma, i32 0, i32 0))
 %t38 = getelementptr %Person, %Person* %t14, i32 0, i32 2
%t39 = load i8*, i8** %t38
call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @print_str_format, i32 0, i32 0), i8* %t39)
call i32 (i8*,...) @printf(i8* getelementptr ([2 x i8], [2 x i8]* @print_str_close_brace, i32 0, i32 0))

 ret i32 0
}
