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

	fmt.Print("$API: ")
	fmt.Println(a...)
}

func (c Condition) Printf(format string, a ...any) {
	if !c {
		return
	}

	fmt.Print("$API: ")
	fmt.Printf(format, a...)
}

func (c Condition) PrintAsJSON(data interface{}) {
	if !c {
		return
	}

	PrintAsJSON(map[string]interface{}{"$API": data})
}
