package main

import (
	"os"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ghoto"
	app.Usage = "Transfer photo(video)"
	app.Action = func(c *cli.Context) {
		println("Hello ghoto")
	}
	app.Run(os.Args)
}
