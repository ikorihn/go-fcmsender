package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

type config struct {
	Services map[string]service `toml:"services"`
}

type service struct {
	CollapseKey string `toml:"collapse_key"`
}

func main() {

	app := &cli.App{
		Name:  "push sender",
		Usage: "",
		Action: func(c *cli.Context) error {
			var conf config
			if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
				return err
			}
			return exec(c, conf)
		},
	}

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "title",
			Aliases:  []string{"t"},
			Value:    "タイトル",
			Usage:    "title for the massage",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "body",
			Aliases:  []string{"b"},
			Value:    "メッセージ",
			Usage:    "body for the massage",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:    "option",
			Aliases: []string{"o"},
			Usage:   "option array key=value",
		},
	}

	app.Flags = flags

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func exec(c *cli.Context, conf config) error {
	title := c.String("title")
	body := c.String("body")
	options := c.StringSlice("option")

	fmt.Printf("TITLE: %v\n", title)
	fmt.Printf("BODY: %v\n", body)
	opt := map[string]string{}
	for _, v := range options {
		keyValue := strings.Split(v, "=")
		key, value := keyValue[0], keyValue[1]
		fmt.Printf("key: %v, value: %v\n", key, value)
		opt[key] = value
	}
	fmt.Printf("OPTIONS: %v\n", opt)

	fmt.Printf("CONFIG: %+v\n", conf)
	/*
		ctx := context.Background()
		client, err := makeClient(ctx)
		if err != nil {
			return err
		}

		orgTokenList := []string{
			"TOKEN1",
			"TOKEN2",
		}

		for _, tokenList := range chunkSlice(orgTokenList, 500) {
			message := makeMessage(title, body, tokenList)

			br, err := client.SendMulticast(context.Background(), message)
			if err != nil {
				log.Fatalln(err)
				return err
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
	*/

	return nil
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
