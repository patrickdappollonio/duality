package main

import (
	"fmt"
	"os"
)

const interpreter = "/bin/bash -c"

func main() {
	shell := os.Getenv("DUALITY_SHELL_COMMAND")
	if shell == "" {
		shell = interpreter
	}

	if err := run(shell, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
