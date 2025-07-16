declare i32 @printf(i8*,...)
	@print_str_format = constant [4 x i8] c"%s\0A\00"
	@print_int_format = constant [4 x i8] c"%d\0A\00"
	@print_float_format = constant [4 x i8] c"%f\0A\00"
	define i32 @main(){
entry:
 ret i32 0
}
