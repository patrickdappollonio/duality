//go:build !linux

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func runProgram(ctx context.Context, stdout, stderr io.Writer, id, program, shell string, shellArgs []string) error {
	writelnf(os.Stderr, "[%s] running: %s %s %q", id, shell, strings.Join(shellArgs, " "), program)

	cmd := exec.Command(shell, append(shellArgs, program)...)
	cmd.Env = os.Environ()
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	go func() {
		<-ctx.Done()
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Signal(syscall.SIGTERM)
		}
	}()

	if err := cmd.Run(); err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			return nil
		}

		return fmt.Errorf("error running %q: %w", program, err)
	}

	return nil
}
