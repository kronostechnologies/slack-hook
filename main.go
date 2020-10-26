package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Field struct {
	Type  string `json:"type"`
	Text  string `json:"text,omitempty"`
	Emoji bool   `json:"emoji,omitempty"`
}

type Block struct {
	Type   string   `json:"type"`
	Fields []*Field `json:"fields,omitempty"`
	Text   *Field   `json:"text,omitempty"`
}

type Attachment struct {
	Blocks []*Block `json:"blocks,omitempty"`
	Color  string   `json:"color,omitempty"`
}

type Payload struct {
	IconEmoji   string        `json:"icon_emoji,omitempty"`
	Username    string        `json:"username,omitempty"`
	Channel     string        `json:"channel,omitempty"`
	Blocks      []*Block      `json:"blocks,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

func main() {
	url := os.Getenv("SLACK_URL")
	channel := os.Getenv("SLACK_CHANNEL")
	iconEmoji := os.Getenv("SLACK_ICON_EMOJI")
	username := os.Getenv("SLACK_USERNAME")
	operation := os.Getenv("RELEASE_OPERATION")
	version := os.Getenv("RELEASE_VERSION")
	installedVersion := os.Getenv("RELEASE_INSTALLED_VERSION")
	application := os.Getenv("RELEASE_APPLICATION")
	environment := os.Getenv("RELEASE_ENVIRONMENT")

	var msgAction string
	var color []byte

	switch operation {
	case "install":
		msgAction = "installed"
		color = []byte{46, 182, 125}
	case "upgrade":
		msgAction = "upgraded"
		color = []byte{54, 197, 240}
	case "delete":
		msgAction = "deleted"
		color = []byte{224, 30, 90}
	case "rollback":
		msgAction = "rolled back"
		color = []byte{255, 161, 0}
	}

	payload := Payload{
		Channel:   channel,
		IconEmoji: iconEmoji,
		Username:  username,
		Attachments: []*Attachment{
			{
				Color: fmt.Sprintf("#%x", color),
				Blocks: []*Block{
					{
						Type: "section",
						Text: &Field{
							Type: "mrkdwn",
							Text: fmt.Sprintf("%s is being %s", application, msgAction),
						},
					},
					{
						Type: "section",
						Fields: []*Field{
							{
								Type: "mrkdwn",
								Text: fmt.Sprintf("*Environment*\n%s", environment),
							},
							{
								Type: "mrkdwn",
								Text: fmt.Sprintf("*Application*\n%s", application),
							},
						},
					},
					{
						Type: "section",
						Fields: []*Field{
							{
								Type: "mrkdwn",
								Text: fmt.Sprintf("*Previous version*\n%s", installedVersion),
							}, {
								Type: "mrkdwn",
								Text: fmt.Sprintf("*Version*\n%s", version),
							},
						},
					},
				},
			},
		},
	}

	body, e := json.Marshal(payload)
	if e != nil {
		log.Panicln(e)
	}

	request, re := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if re != nil {
		log.Panicln(re)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	_, ce := client.Do(request)
	if ce != nil {
		log.Panicln(ce)
	}
}
