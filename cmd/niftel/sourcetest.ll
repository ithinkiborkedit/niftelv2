declare i32 @printf(i8*,...)
	@print_str_format = constant [4 x i8] c"%s\0A\00"
	@print_int_format = constant [4 x i8] c"%d\0A\00"
	define i32 @main(){
entry:
@.str4 = private constant [5 x i8] c"ello\00"
call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @print_str_format, i32 0, i32 0), i8* getelementptr ([5 x i8], [5 x i8]* @.str4, i32 0, i32 0))
 ret i32 0
}
