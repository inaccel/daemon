package main

import (
	"os"

	"github.com/inaccel/daemon/pkg"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var configCommand = &cli.Command{
	Name:      "config",
	Aliases:   []string{"c"},
	Usage:     "Information on the inacceld config",
	ArgsUsage: " ",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "defaults",
			Usage: "see the output of the default config",
		},
	},
	Action: func(context *cli.Context) error {
		config := pkg.NewConfig()

		if !context.Bool("defaults") {
			if err := config.Read(context.String("config-file")); err != nil {
				if context.IsSet("config-file") {
					return err
				}

				logrus.Warn(err)
			}

			if context.IsSet("root") {
				config.Root = context.String("root")
			}
		}

		return config.Write(os.Stdout.Name(), os.ModePerm)
	},
}
