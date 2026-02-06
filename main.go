package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	githubUsername string
	webhookSecret  string
	slackBotToken  string
	slackUserID    string
	debug          bool
)

func main() {
	githubUsername = requireEnv("GITHUB_USERNAME")
	webhookSecret = requireEnv("GITHUB_WEBHOOK_SECRET")
	slackBotToken = os.Getenv("SLACK_BOT_TOKEN")
	slackUserID = os.Getenv("SLACK_USER_ID")
	debug = os.Getenv("DEBUG") == "true"

	log.Printf("Watching PRs for user: %s (debug: %v, slack: %v)", githubUsername, debug, slackBotToken != "")

	http.HandleFunc("/healthz", handleHealthz)
	http.HandleFunc("/webhook", handleWebhook)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s environment variable is required", key)
	}
	return val
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
