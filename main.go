package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "os"

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
	res := struct {
		ResponseType string `json:"response_type"`
		Text         string `json:"text"`
	}{
		ResponseType: slack.ResponseTypeInChannel,
		Text:         fmt.Sprintf("<@%s> 打刻します", s.UserID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

	// api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	// api.PostMessage(
	// 	s.ChannelID,
	// 	slack.MsgOptionText(fmt.Sprintf("<@%s> 打刻しました", s.UserID), false))
}

func handleWorkTime(w http.ResponseWriter, s *slack.SlashCommand) {
	res := struct {
		ResponseType string `json:"response_type"`
		Text         string `json:"text"`
	}{
		ResponseType: slack.ResponseTypeInChannel,
		Text:         fmt.Sprintf("<@%s> 現在の労働時間確認してきます", s.UserID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

	// api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	// api.PostMessage(
	// 	s.ChannelID,
	// 	slack.MsgOptionText(fmt.Sprintf("<@%s> 現在の労働時間は00:00(無職)です:smiley:", s.UserID), false))
}

func main() {
	http.HandleFunc("/slack/commands", handleRequest)
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
