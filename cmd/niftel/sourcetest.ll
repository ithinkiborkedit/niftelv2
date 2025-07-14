declare i32 @printf(i8*,...)
	@print_str_format = constant [4 x i8] c"%s\0A\00"
	@print_int_format = constant [4 x i8] c"%d\0A\00"
	define i32 @main(){
entry:
@.str7 = private constant [8 x i8] c"s%!\(string=\"hello\")00"
call i32 (i8*,...) @printf(i8* getelementptr ([8 x i8], [8 x i8]* @print_str_format, i32 0, i32 0), i32 @.str7)
 ret i32 0
}
