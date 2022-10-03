package main

import (
	"context"
	"os"

	"github.com/flytrap/gin-base/internal/app"
	"github.com/flytrap/gin-base/pkg/logger"
	"github.com/urfave/cli/v2"
)

var VERSION = "0.0.1"

// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @scheme bearer
func main() {
	ctx := context.Background()
	app := cli.NewApp()
	app.Name = "gin-base"
	app.Version = VERSION
	app.Usage = "gin based on GIN + GORM + WIRE."
	app.Commands = []*cli.Command{
		newWebCmd(ctx),
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
		},
		Action: func(c *cli.Context) error {
			return app.Run(ctx,
				app.SetConfigFile(c.String("conf")),
				app.SetVersion(VERSION))
		},
	}
}
