package global

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ralexgt/glox/scanner"
)

type Lox struct {
	HadError bool
}

var VM = Lox{}

// runs an entire file as a script all at once
func (l *Lox) RunFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	l.run(string(bytes))

	if l.HadError {
		os.Exit(65)
	}

	return nil
}

// takes in the prompt in the repl and runs the line or errors (func run)
func (l *Lox) RunPrompt() error {
	input := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if ok := input.Scan(); !ok {
			break
		}
		line := input.Text()
		l.run(line)
		l.HadError = false
	}

	return nil
}

// Error logic - the client decides how to use the implemented logic
func (l *Lox) ReportError(line int, err error) {
	l.report(line, "", err)
}

func (l *Lox) report(line int, where string, err error) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, err)
	l.HadError = true
}

func (l *Lox) run(source string) {
	scanner := scanner.NewScanner(source)
	scanner.ScanTokens()

	for _, token := range scanner.Tokens {
		fmt.Println(token)
	}
}
