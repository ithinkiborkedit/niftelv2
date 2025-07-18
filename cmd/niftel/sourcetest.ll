declare i32 @printf(i8*,...)
	@print_str_format = constant [4 x i8] c"%s\0A\00"
	@print_int_format = constant [4 x i8] c"%d\0A\00"
	@print_float_format = constant [4 x i8] c"%f\0A\00"
	%Person = type { i8*, i64 }
define i32 @main(){
entry:
 %t0 = alloca %Person
 %t1 = alloca %Person
%t2 = getelementptr %Person, %Person* %t1 i32 0, i32 0
store i8* getelementptr ([5 x i8], [5 x i8]* @.str0, i32 0, i32 0), %t2* i8*
%t3 = getelementptr %Person, %Person* %t1 i32 0, i32 1
store i64 2, %t3* i64
 store %Person* %t1, %Person** %t0
 ret i32 0
}
