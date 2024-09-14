package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/patrickdappollonio/duality/prefixer"
	squids "github.com/sqids/sqids-go"
	"golang.org/x/sync/errgroup"
)

func run(interpreter string, args []string) error {
	emptyCount := 0
	for _, arg := range args {
		if strings.TrimSpace(arg) == "" {
			emptyCount++
		}
	}

	if emptyCount == len(args) {
		return errors.New("no commands to run")
	}

	shell, shellArgs := "", []string{}
	interpreterFields := strings.Fields(interpreter)
	switch len(interpreterFields) {
	case 0:
		return fmt.Errorf("invalid shell interpreter %q", interpreter)
	case 1:
		shell = interpreterFields[0]
	default:
		shell = interpreterFields[0]
		shellArgs = interpreterFields[1:]
	}

	parent, done := context.WithCancel(context.Background())
	defer done()

	eg, ctx := errgroup.WithContext(parent)

	sq, _ := squids.New(squids.Options{
		MinLength: 8,
		Alphabet:  "abcdefghijklmnopqrstuvwxyz0123456789",
	})

	go func() {
		chExit := make(chan os.Signal, 1)
		signal.Notify(chExit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-chExit
		done()
	}()

	for pos, command := range args {
		eg.Go(func() error {
			if strings.TrimSpace(command) == "" {
				return nil
			}

			if fields := strings.Fields(command); len(fields) == 0 {
				return fmt.Errorf("no fields found in %q", command)
			}

			id, _ := sq.Encode(unicode(pos+1, command))

			stdout := prefixer.NewPrefixWriter(os.Stdout, fmt.Sprintf("[%s][stdout] ", id))
			stderr := prefixer.NewPrefixWriter(os.Stderr, fmt.Sprintf("[%s][stderr] ", id))

			if err := runProgram(ctx, stdout, stderr, id, command, shell, shellArgs); err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "[%s] %q executed successfully (exit code: 0)\n", id, command)
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func writelnf(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format+"\n", args...)
}

func unicode(idx int, s string) []uint64 {
	total := 0
	for _, char := range s {
		total += int(char)
	}
	return []uint64{uint64(idx), uint64(total)}
}
