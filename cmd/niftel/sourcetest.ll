declare i32 @print(i8*,...)
	@print.str = constant [4 x i8] c"%\d\0A\00"

	define i32 @main(){
entry:
call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8]* @print.str,i32 0,i32 0 42)
 ret i32 0
}
