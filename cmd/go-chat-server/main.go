package main

import (
	"os"

	"github.com/madatsci/go-chat-server/internal"
	"github.com/urfave/cli"
)

var version = "development"

func main() {
	app := cli.NewApp()
	app.Name = "go-chat-server"
	app.Version = version

	app.Commands = []cli.Command{
		{
			Name: "serve",
			Aliases: []string{"s"},
			Usage: "run chat server",
			Action: func(ctx *cli.Context) {
				internal.Run()
			},
		},
		{
			Name: "migrate:up",
			Aliases: []string{"mup"},
			Usage: "apply migrations up",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir",
					Value: "./migrations",
				},
			},
			Action: func(ctx *cli.Context) {
				dir := ctx.String("dir")
				internal.Migrate(dir)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
