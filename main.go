package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
)

func main() {
	ctx := context.Background()
	client, err := makeClient(ctx)
	if err != nil {
		return
	}

	title := "title"
	body := "body"

	orgTokenList := []string{
		"TOKEN1",
		"TOKEN2",
	}

	for _, tokenList := range chunkSlice(orgTokenList, 500) {
		message := makeMessage(title, body, tokenList)

		br, err := client.SendMulticast(context.Background(), message)
		if err != nil {
			log.Fatalln(err)
		}

		if br.FailureCount > 0 {
			var failedTokens []string
			for idx, resp := range br.Responses {
				if !resp.Success {
					// The order of responses corresponds to the order of the registration tokens.
					failedTokens = append(failedTokens, tokenList[idx])
				}
			}
			fmt.Printf("List of tokens that caused failures: %v\n", failedTokens)
		}

		// See the BatchResponse reference documentation
		// for the contents of response.
		fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
	}
}

func makeClient(ctx context.Context) (*messaging.Client, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return nil, err
	}

	return app.Messaging(ctx)
}

func makeMessage(title, body string, tokenList []string) *messaging.MulticastMessage {
	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: "",
		},
		Tokens: tokenList,
	}

	return message
}

func chunkSlice(slice []string, size int) [][]string {
	chunkedTokenList := make([][]string, 0)
	sliceSize := len(slice)
	for i := 0; i < sliceSize; i += size {
		end := i + size
		if sliceSize < end {
			end = sliceSize
		}
		chunkedTokenList = append(chunkedTokenList, slice[i:end])
	}

	return chunkedTokenList
}