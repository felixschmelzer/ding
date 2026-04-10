package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	BotToken string `toml:"bot_token"`
	ChatID   string `toml:"chat_id"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "done", "config.toml"), nil
}

func loadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	if cfg.BotToken == "" || cfg.ChatID == "" {
		return nil, fmt.Errorf("config is incomplete")
	}
	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}

func runConfig() error {
	existing, _ := loadConfig()

	fmt.Println("done — Telegram notification setup")
	fmt.Println("────────────────────────────────────")
	fmt.Println()
	fmt.Println("You need a Telegram bot token and your chat ID.")
	fmt.Println("  1. Message @BotFather on Telegram → /newbot → copy the token")
	fmt.Println("  2. Message your new bot, then open:")
	fmt.Println("     https://api.telegram.org/bot<TOKEN>/getUpdates")
	fmt.Println("     Look for \"id\" inside the \"chat\" object — that's your chat ID.")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	prompt := func(label, current string) (string, error) {
		if current != "" {
			fmt.Printf("%s [%s]: ", label, current)
		} else {
			fmt.Printf("%s: ", label)
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			return current, nil
		}
		return line, nil
	}

	var cur Config
	if existing != nil {
		cur = *existing
	}

	cfg := &Config{}
	var err error

	cfg.BotToken, err = prompt("Bot token", cur.BotToken)
	if err != nil {
		return err
	}
	cfg.ChatID, err = prompt("Chat ID", cur.ChatID)
	if err != nil {
		return err
	}

	if cfg.BotToken == "" || cfg.ChatID == "" {
		return fmt.Errorf("bot token and chat ID are required")
	}

	if err := saveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	path, _ := configPath()
	fmt.Printf("\nConfig saved to %s\n", path)
	return nil
}
