package main

import (
	"os"

	"github.com/photoprism/photoprism/internal/commands"
	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/urfave/cli"
)

var version = "development"
var log = event.Log

func main() {
	app := cli.NewApp()
	app.Name = "PhotoPrism-Album"
	app.Usage = "Manipulate your Albums"
	app.Version = version
	app.Copyright = "(c) 2018-2020 PhotoPrism.org <hello@photoprism.org>"
	app.EnableBashCompletion = true
	app.Flags = config.GlobalFlags

	app.Commands = []cli.Command{
		commands.AlbumCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}
