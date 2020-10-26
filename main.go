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
	color := os.Getenv("SLACK_COLOR")
	operation := os.Getenv("RELEASE_OPERATION")
	version := os.Getenv("RELEASE_VERSION")
	installedVersion := os.Getenv("RELEASE_INSTALLED_VERSION")
	application := os.Getenv("RELEASE_APPLICATION")
	environment := os.Getenv("RELEASE_ENVIRONMENT")

	var msgAction string
	var msgVersion string
	var msgUsername string

	switch operation {
	case "install":
		msgAction = "Installing"
		msgVersion = version
	case "upgrade":
		msgAction = "Installing"
		msgVersion = version
	case "delete":
		msgAction = "Deleting"
		msgVersion = installedVersion
	case "rollback":
		msgAction = "Rolling back to"
		msgVersion = version
	}

	if username != "" && environment != "" {
		msgUsername = fmt.Sprintf("%s (%s)", username, environment)
	} else if username != "" {
		msgUsername = username
	} else if environment != "" {
		msgUsername = environment
	} else {
		msgUsername = "slack-hook"
	}

	context := []*Field{
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Environment: *\n%s", environment),
		},
	}

	if operation != "delete" && installedVersion != "" {
		context = append(context, &Field{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Previous: *\n%s", installedVersion),
		})
	}

	payload := Payload{
		Channel:   channel,
		IconEmoji: iconEmoji,
		Username:  msgUsername,
		Attachments: []*Attachment{
			{
				Color: color,
				Blocks: []*Block{
					{
						Type: "section",
						Text: &Field{
							Type: "mrkdwn",
							Text: fmt.Sprintf("%s *%s* %s", msgAction, application, msgVersion),
						},
					},
					{
						Type:   "context",
						Fields: context,
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
