package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func sendTelegram(cfg *Config, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.BotToken)
	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id":    {cfg.ChatID},
		"text":       {message},
		"parse_mode": {"HTML"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}
	return nil
}

func buildRunningMessage(cmd, elapsed string) string {
	return fmt.Sprintf(
		"⏳ <b>Still running</b>\n<code>%s</code>\nRunning for: %s",
		cmd, elapsed,
	)
}

func buildMessage(cmd string, exitCode int, duration, finishTime string) string {
	if exitCode == 0 {
		return fmt.Sprintf(
			"✅ <b>Done</b>\n<code>%s</code>\nExit: 0 | Duration: %s\nFinished: %s",
			cmd, duration, finishTime,
		)
	}
	return fmt.Sprintf(
		"❌ <b>Failed</b>\n<code>%s</code>\nExit: %d | Duration: %s\nFinished: %s",
		cmd, exitCode, duration, finishTime,
	)
}
