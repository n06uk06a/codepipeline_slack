package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var url = "https://slack.com/api/chat.postMessage"
var username = os.Getenv("USERNAME")
var channel = os.Getenv("CHANNEL")
var token = os.Getenv("TOKEN")
var text = "Pipeline ended."

type codePipelineEvent struct {
	Pipeline string `json:"pipeline"`
	State    string `json:"state"`
}

type slackMessage struct {
	Username    string       `json:"username,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Text        string       `json:"text,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Fallback    string       `json:"fallback,omitempty"`
	Attachments []attachment `json:"attachments,omitempty"`
}

type attachment struct {
	Fields []field `json:"fields,omitempty"`
}

type field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

func lambdaHandler(ctx context.Context, event events.CloudWatchEvent) error {
	logObject(event)
	var detail codePipelineEvent
	if err := json.Unmarshal([]byte(event.Detail), &detail); err != nil {
		panic(err)
	}
	pipeline := detail.Pipeline
	state := detail.State
	var icon string
	if state == "SUCCEEDED" {
		icon = ":sunny:"
	}
	if state == "FAILED" {
		icon = ":umbrella_with_rain_drops:"
	}
	message, err := json.Marshal(
		slackMessage{
			Username:  username,
			IconEmoji: icon,
			Text:      text,
			Channel:   channel,
			Fallback:  text,
			Attachments: []attachment{
				attachment{
					Fields: []field{
						{
							Title: "Pipeline",
							Value: pipeline,
						},
						{

							Title: "Status",
							Value: state,
						},
					},
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	logObject(string(message))
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(message),
	)
	if err != nil {
		return err
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(b))
	defer resp.Body.Close()
	return err
}

func main() {
	lambda.Start(lambdaHandler)
}

func logObject(event interface{}) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	log.Print(string(eventJSON))
}
