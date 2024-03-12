package debug

import (
	"fmt"
	"os"
)

type Condition bool

var IfVerbose Condition = isVerbose()

func If (condition bool) Condition {
	return Condition(condition)
}

func (c Condition) Println(a ...any) {
	if !c {
		return
	}

	fmt.Println(a...)
}

func (c Condition) Printf(format string, a ...any) {
	if !c {
		return
	}
	
	fmt.Printf(format, a...)
}

func (c Condition) PrintAsJSON(data interface{}) {
	if !c {
		return
	}

	PrintAsJSON(data)
}

func isVerbose () Condition	{
	if len(os.Args) < 2 {
		return Condition(false)
	}

	verbose := os.Args[1] == "-v"

	return Condition(verbose)
}