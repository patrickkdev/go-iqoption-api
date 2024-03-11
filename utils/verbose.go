package utils

import (
	"fmt"
	"os"
)

func PrintlnIfVerbose(a ...any) {
	if len(os.Args) < 2 {
		return
	}
	
	verbose := os.Args[1] == "-v"
	if verbose {
		fmt.Println(a...)
	}
}