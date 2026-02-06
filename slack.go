package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type slackMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func sendSlackMessage(text string) error {
	// Skip if Slack not configured
	if slackBotToken == "" || slackUserID == "" {
		log.Printf("[DRY-RUN] %s", text)
		return nil
	}

	// Skip actual API call in tests
	if strings.HasPrefix(slackBotToken, "test") {
		log.Printf("[TEST] Would send Slack message: %s", text)
		return nil
	}

	msg := slackMessage{
		Channel: slackUserID,
		Text:    text,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+slackBotToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API returned status %d", resp.StatusCode)
	}

	log.Printf("Slack message sent: %s", text)
	return nil
}
