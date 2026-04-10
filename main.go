package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: done <command> [args...]\n       done --config")
		os.Exit(1)
	}

	if os.Args[1] == "--config" {
		if err := runConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "done: %v\n", err)
			os.Exit(1)
		}
		return
	}

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "done: no config found, run 'done --config' to set up")
		os.Exit(1)
	}

	args := os.Args[1:]
	exitCode, duration, err := runCommand(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "done: %v\n", err)
		os.Exit(1)
	}

	msg := buildMessage(
		strings.Join(args, " "),
		exitCode,
		formatDuration(duration),
		time.Now().Format("15:04:05"),
	)

	if err := sendTelegram(cfg, msg); err != nil {
		fmt.Fprintf(os.Stderr, "done: notification failed: %v\n", err)
	}

	os.Exit(exitCode)
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
