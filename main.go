package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Type     string   `json:"type"`
	Fields   []*Field `json:"fields,omitempty"`
	Elements []*Field `json:"elements,omitempty"`
	Text     *Field   `json:"text,omitempty"`
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
	}

	context := []*Field{
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Environment:* %s", environment),
		},
	}

	if operation != "delete" && installedVersion != "" {
		if version == installedVersion {
			os.Exit(0)
		}

		context = append(context, &Field{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Previous:* %s", installedVersion),
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
						Type:     "context",
						Elements: context,
					},
				},
			},
		},
	}

	body, e := json.Marshal(payload)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	request, re := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if re != nil {
		fmt.Println(re)
		os.Exit(1)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, ce := client.Do(request)
	if ce != nil {
		fmt.Println(ce)
		os.Exit(1)
	}
	slackMsg, re := ioutil.ReadAll(resp.Body)
	if re != nil {
		fmt.Println(re)
		os.Exit(1)
	}
	fmt.Println(string(slackMsg))
}
