package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/slack-go/slack"
)

type UserStatus struct {
	IsWorking bool
	StartTime time.Time
	WorkTime  time.Duration
}

var userStatuses = make(map[string]*UserStatus)

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
	status, exists := userStatuses[s.UserID]
	if !exists {
		status = &UserStatus{}
		userStatuses[s.UserID] = status
	}

	now := time.Now()
	if status.IsWorking {
		status.WorkTime += now.Sub(status.StartTime)
		status.IsWorking = false
	} else {
		status.IsWorking = true
		status.StartTime = now
	}

	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	// @shivaduke 打刻します。
	_, _, err := api.PostMessage(
		s.ChannelID,
		slack.MsgOptionText(fmt.Sprintf("<@%s> 打刻します。", s.UserID), false))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// @shivaduke 打刻しました。
	_, _, err = api.PostMessage(
		s.ChannelID,
		slack.MsgOptionText(fmt.Sprintf("<@%s> 打刻しました。", s.UserID), false))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleWorkTime(w http.ResponseWriter, s *slack.SlashCommand) {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	// @shivaduke 現在の労働時間確認してきます
	_, _, err := api.PostMessage(
		s.ChannelID,
		slack.MsgOptionText("<@%s> 現在の労働時間確認してきます", false))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 未出勤、勤務中、退室中
	status, ok := userStatuses[s.UserID]

	now := time.Now()
	var workingStatus string
	var workTime time.Duration
	if ok && status.IsWorking {
		workingStatus = "勤務中"
		workTime = status.WorkTime + now.Sub(status.StartTime)
	} else if status.StartTime.Day() == now.Day() {
		workingStatus = "退室中"
		workTime = status.WorkTime
	} else {
		workingStatus = "未出勤"
		workTime = 0
	}

	// @shivaduke 現在の労働時間はHH:MM(勤務中)です:スマイリー:
	h := int32(workTime.Hours())
	m := int32(workTime.Minutes()) - h*60
	_, _, err = api.PostMessage(
		s.ChannelID,
		slack.MsgOptionText(fmt.Sprintf("<@%s> 現在の労働時間は%d:%d(%s)です:スマイリー:", s.UserID, h, m, workingStatus), false))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/slack/commands", handleRequest)
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
