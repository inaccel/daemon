package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/inaccel/daemon/internal/driver"
	"github.com/inaccel/daemon/internal/plugins"
	"github.com/inaccel/daemon/pkg"
	"github.com/inaccel/daemon/pkg/plugin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var version string

func main() {
	app := &cli.App{
		Name:    "inacceld",
		Version: version,
		Usage:   "A self-sufficient runtime for accelerators.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config-file",
				Usage: "daemon configuration file",
				Value: "/etc/inaccel/daemon.json",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "enable debug output",
			},
		},
		Before: func(context *cli.Context) error {
			log.SetOutput(ioutil.Discard)

			logrus.SetFormatter(new(logrus.JSONFormatter))

			if context.Bool("debug") {
				logrus.SetLevel(logrus.DebugLevel)
			}

			return nil
		},
		Action: func(context *cli.Context) error {
			config := pkg.NewConfig()

			if err := config.Read(context.String("config-file")); err != nil {
				if context.IsSet("config-file") {
					return err
				}

				logrus.Warn(err)
			}

			if err := config.Validate(); err != nil {
				return err
			}

			inaccel, err := driver.NewInAccel(config, version)
			if err != nil {
				return err
			}

			var new []plugin.New
			if os.Getenv("DOCKER") != "disabled" {
				new = append(new, func() plugin.Plugin {
					return plugins.NewDocker(context.Context, inaccel)
				})
			}
			if os.Getenv("KUBELET") != "disabled" {
				new = append(new, func() plugin.Plugin {
					return plugins.NewKubelet(context.Context, inaccel)
				})
			}

			plugin.Handle(new...)

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Information on the inacceld config",
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
					}

					return config.Write(os.Stdout.Name(), os.ModePerm)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
