package output

import "fmt"

type IOutput interface {
	Verbose(a ...any)
	FVerbose(format string, a ...any)
	Print(a ...any)
	Printf(format string, a ...any)
}

type Output struct {
}

func (*Output) Print(a ...any) {
	fmt.Println(a...)
}

func (*Output) Printf(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Println()
}
