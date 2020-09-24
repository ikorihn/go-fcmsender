package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

type config struct {
	Services map[string]service `toml:"services"`
}

type service struct {
	CollapseKey string `toml:"collapse_key"`
}

type action interface {
	Run(c *cli.Context, conf config) error
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

			var action action

			action = &debug{}
			return action.Run(c, conf)
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
		&cli.StringFlag{
			Name:     "service",
			Aliases:  []string{"s"},
			Usage:    "service to send",
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

type debug struct{}

func (d *debug) Run(c *cli.Context, conf config) error {
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

	fmt.Printf("CONFIG: %+v\n", conf.Services[c.String("service")].CollapseKey)

	return nil
}
