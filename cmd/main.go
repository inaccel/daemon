package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/inaccel/daemon/internal/driver"
	"github.com/inaccel/daemon/internal/plugins"
	"github.com/inaccel/daemon/pkg"
	"github.com/inaccel/daemon/pkg/env"
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
			&cli.StringFlag{
				Name:  "root",
				Usage: "inacceld root directory",
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

			if context.IsSet("root") {
				config.Root = context.String("root")
			}

			if err := config.Validate(); err != nil {
				return err
			}

			inaccel, err := driver.NewInAccel(config, version)
			if err != nil {
				return err
			}

			var new []plugin.New
			if !env.Disabled("DOCKER") {
				new = append(new, func() plugin.Plugin {
					return plugins.NewDocker(context.Context, inaccel)
				})
			}
			if !env.Disabled("KUBELET") {
				new = append(new, func() plugin.Plugin {
					return plugins.NewKubelet(context.Context, inaccel)
				})
			}

			plugin.Handle(new...)

			return nil
		},
		Commands: []*cli.Command{
			configCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
