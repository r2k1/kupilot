package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Terminal struct {
	in  io.Reader
	out io.Writer
}

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

func NewTerminal(in io.Reader, out io.Writer) *Terminal {
	return &Terminal{
		in:  in,
		out: out,
	}
}

func (t *Terminal) Write(message string) {
	_, _ = fmt.Fprintf(t.out, message)
}

func (t *Terminal) WriteError(message string) {
	_, _ = fmt.Fprintf(t.out, Red+message+Reset)
}

func (t *Terminal) WriteInfo(message string) {
	_, _ = fmt.Fprintf(t.out, Green+message+Reset)
}

func (t *Terminal) WriteWarning(message string) {
	_, _ = fmt.Fprintf(t.out, Yellow+message+Reset)
}

func (t *Terminal) WriteDebug(message string) {
	_, _ = fmt.Fprintf(t.out, Blue+message+Reset)
}

func (t *Terminal) Read() (string, error) {
	t.Write("\n> ")
	scanner := bufio.NewScanner(t.in)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read user input: %w", err)
	}
	userInput := scanner.Text()
	if userInput == "exit" || userInput == "quit" {
		os.Exit(0)
	}
	return userInput, nil

}
