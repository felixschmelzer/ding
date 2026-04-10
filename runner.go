package main

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/creack/pty"
	"golang.org/x/term"
)

func runCommand(args []string) (exitCode int, duration time.Duration, err error) {
	cmd := exec.Command(args[0], args[1:]...)

	start := time.Now()

	ptmx, err := pty.Start(cmd)
	if err != nil {
		// Fall back to normal execution if PTY fails (e.g. non-interactive env)
		cmd = exec.Command(args[0], args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		start = time.Now()
		runErr := cmd.Run()
		elapsed := time.Since(start)
		if runErr != nil {
			if exitErr, ok := runErr.(*exec.ExitError); ok {
				return exitErr.ExitCode(), elapsed, nil
			}
			return 1, elapsed, runErr
		}
		return 0, elapsed, nil
	}
	defer ptmx.Close()

	// Mirror terminal resize signals into the PTY
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)
	go func() {
		for range sigCh {
			_ = pty.InheritSize(os.Stdout, ptmx)
		}
	}()
	sigCh <- syscall.SIGWINCH // set initial size

	// Put the host terminal into raw mode so the child gets raw input
	if term.IsTerminal(int(os.Stdin.Fd())) {
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err == nil {
			defer term.Restore(int(os.Stdin.Fd()), oldState)
		}
	}

	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx) // blocks until PTY closes

	signal.Stop(sigCh)
	close(sigCh)

	elapsed := time.Since(start)

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), elapsed, nil
		}
		return 1, elapsed, err
	}

	return 0, elapsed, nil
}
