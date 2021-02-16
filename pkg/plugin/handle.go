package plugin

import (
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

func Handle(new ...New) {
	plugin := make([]Plugin, len(new))

	for i := range plugin {
		plugin[i] = new[i]()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, unix.SIGHUP, unix.SIGINT, unix.SIGTERM)

loop:
	for sig := range c {
		switch sig {
		case unix.SIGHUP:
			for i := range plugin {
				if plugin[i].IsClosed() {
					plugin[i] = new[i]()
				}
			}
		case unix.SIGINT, unix.SIGTERM:
			for i := range plugin {
				plugin[i].Close()
			}

			break loop
		}
	}
}
