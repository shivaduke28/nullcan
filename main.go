package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !strings.Contains(s.ChannelName, "無職") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch s.Command {
	case "/nullcan_touch":
		handleTouch(w, &s)
	case "/nullcan_worktime":
		handleWorkTime(w, &s)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleTouch(w http.ResponseWriter, s *slack.SlashCommand) {
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

	go func() {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		delayMs := r.Intn(501) + 1000
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
		api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
		_, _, err := api.PostMessage(
			s.ChannelID,
			slack.MsgOptionText(fmt.Sprintf("<@%s> 打刻しました", s.UserID), false))
		if err != nil {
			log.Printf("Error posting message: %v", err)
		}
	}()
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

	go func() {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		delayMs := r.Intn(501) + 1000
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
		api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
		_, _, err := api.PostMessage(
			s.ChannelID,
			slack.MsgOptionText(fmt.Sprintf("<@%s> 現在の労働時間は00:00(無職)です:smiley:", s.UserID), false))
		if err != nil {
			log.Printf("Error posting message: %v", err)
		}
	}()
}

func main() {
	http.HandleFunc("/slack/commands", handleRequest)
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
