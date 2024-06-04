package debug

import (
	"fmt"
)

type Condition bool

var IfVerbose Condition = false

func If(condition bool) Condition {
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
