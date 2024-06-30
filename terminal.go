package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Terminal struct {
	in      io.Reader
	out     io.Writer
	noColor bool
}

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

func (t *Terminal) Write(message string) {
	_, _ = fmt.Fprintf(t.out, message)
}

func (t *Terminal) WriteError(message string) {
	t.writeColor(Red, message)
}

func (t *Terminal) WriteInfo(message string) {
	t.writeColor(Green, message)
}

func (t *Terminal) WriteWarning(message string) {
	t.writeColor(Yellow, message)
}

func (t *Terminal) WriteDebug(message string) {
	t.writeColor(Blue, message)
}

func (t *Terminal) writeColor(color string, message string) {
	if t.noColor {
		t.Write(message)
		return
	}
	t.Write(color + message + Reset)
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
