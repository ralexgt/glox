package main

import (
	"fmt"
	"os"

	"github.com/ralexgt/glox/global"
)

func main() {
	global.VM = global.Lox{
		HadError: false,
	}

	// It takes 1 argument by default. If no additional argument is provided, open an prompt to execute one line at a time. (REPL)
	// If it takes 2 argument in total, it was provided a script to run.
	// If it takes more than 2 arguments in total it should exit at is not meant to be used like this.
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	}
	if len(os.Args) == 2 {
		global.VM.RunFile(os.Args[1])
	} else {
		global.VM.RunPrompt()
	}
}
