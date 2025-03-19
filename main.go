package main

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
	hadError bool
}

var vm = Lox{}

func main() {
	// It takes 1 argument by default. If no additional argument is provided, open an prompt to execute one line at a time. (REPL)
	// If it takes 2 argument in total, it was provided a script to run.
	// If it takes more than 2 arguments in total it should exit at is not meant to be used like this.
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	}
	if len(os.Args) == 2 {
		vm.runFile(os.Args[1])
	} else {
		vm.runPrompt()
	}
}

// runs an entire file as a script all at once
func (l *Lox) runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	l.run(string(bytes))

	if l.hadError {
		os.Exit(65)
	}

	return nil
}

// takes in the prompt in the repl and runs the line or errors (func run)
func (l *Lox) runPrompt() error {
	input := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if ok := input.Scan(); !ok {
			break
		}
		line := input.Text()
		l.run(line)
		l.hadError = false
	}

	return nil
}

func (l *Lox) run(source string) {
	scanner := NewScanner(source)
	scanner.scanTokens()

	for _, token := range scanner.tokens {
		fmt.Println(token)
	}
}

// Error logic - the client decides how to use the implemented logic
func (l *Lox) reportError(line int, err error) {
	l.report(line, "", err)
}

func (l *Lox) report(line int, where string, err error) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, err)
	l.hadError = true
}
