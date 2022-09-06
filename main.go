package main

import (
	"context"
	"os"

	"github.com/flytrap/gin_template/internal/app"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var VERSION = "0.0.1"

func main() {
	ctx := context.Background()
	app := cli.NewApp()
	app.Name = "gosearch"
	app.Version = VERSION
	app.Usage = "gosearch based on GIN + GORM + WIRE."
	app.Commands = []*cli.Command{
		newWebCmd(ctx),
		newImportCmd(ctx),
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Error(err.Error())
	}
}

func newWebCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "Run http server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{"c"},
				Usage:   "App configuration file(.json,.yaml,.toml)",
			},
			&cli.StringFlag{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "init database data(.json)",
			},
			&cli.StringFlag{
				Name:  "www",
				Usage: "Static site directory",
			},
		},
		Action: func(c *cli.Context) error {
			return app.Run(ctx,
				app.SetConfigFile(c.String("conf")),
				app.SetInitFile(c.String("init")),
				app.SetVersion(VERSION))
		},
	}
}

func newImportCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "import",
		Usage: "import data",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{"c"},
				Usage:   "App configuration file(.json,.yaml,.toml)",
			},
			&cli.StringFlag{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "init database data(.json)",
			},
		},
		Action: func(c *cli.Context) error {
			return app.Import(ctx,
				app.SetConfigFile(c.String("conf")),
				app.SetInitFile(c.String("init")),
				app.SetVersion(VERSION))
		},
	}
}
