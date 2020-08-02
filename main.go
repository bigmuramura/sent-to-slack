package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/slack-go/slack"
)

// REGION ...
const REGION = "ap-northeast-1"

func main() {
	// The name you saved in Secrets Manager
	secretName := "SlackInfo"

	// Fetch a value from the Secrets Manager
	secretChannelID, secretSlackToken, err := getSecret(secretName)
	if err != nil {
		log.Fatal(err)
	}

	// Send to Slack channel
	api := slack.New(secretSlackToken)
	sendMessage := "Test Message"
	channelID, timestamp, err := api.PostMessage(secretChannelID, slack.MsgOptionText(sendMessage, false))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func getSecret(sec string) (string, string, error) {
	secretName := sec

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New(),
		aws.NewConfig().WithRegion(REGION))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", "", err
	}
	fmt.Printf("%T型\n%[1]v\n", result) // 說明ポイント1

	secretString := *result.SecretString
	fmt.Printf("%T型\n%[1]v\n", secretString) // 說明ポイント2

	res := make(map[string]interface{})
	if err := json.Unmarshal([]byte(secretString), &res); err != nil {
		return "", "", err
	}
	fmt.Printf("%T型\n%[1]v\n", res) // 說明ポイント3

	return res["channel_id"].(string), res["slack_token"].(string), nil
}
