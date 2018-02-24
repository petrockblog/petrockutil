package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/petrockblog/petrockutil/commands"
)

func main() {
	app := cli.NewApp()
	app.Name = "petrockutil"
	app.Version = VERSION
	app.Usage = "Command Line Utility for petrockblock.com gadgets"
	app.Commands = []cli.Command{
		commands.Scan(),
		commands.GamepadBlock(),
	}
	app.Run(os.Args)
}
