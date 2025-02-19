package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch s.Command {
	case "/nullcan_touch":
		handleToutch(w, &s)
	case "/nullcan_worktime":
		handleWorkTime(w, &s)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleToutch(w http.ResponseWriter, s *slack.SlashCommand) {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	// @shivaduke 打刻します。
	_, _, err := api.PostMessage(
		s.ChannelID,
		slack.MsgOptionText(fmt.Sprintf("<@%s> 打刻します", s.UserID), false))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// @shivaduke 打刻しました。
	_, _, err = api.PostMessage(
		s.ChannelID,
		slack.MsgOptionText(fmt.Sprintf("<@%s> 打刻しました", s.UserID), false))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleWorkTime(w http.ResponseWriter, s *slack.SlashCommand) {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	_, _, err := api.PostMessage(
		s.ChannelID,
		slack.MsgOptionText(fmt.Sprintf("<@%s> 現在の労働時間確認してきます", s.UserID), false))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _, err = api.PostMessage(
		s.ChannelID,
		slack.MsgOptionText(fmt.Sprintf("<@%s> 現在の労働時間は00:00(無職)です:smiley:", s.UserID), false))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/slack/commands", handleRequest)
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
