package main

import (
	"os"
	"fmt"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ghoto"
	app.Usage = "Transfer photo(video)"
	app.Flags = []cli.Flag {
		cli.StringFlag {
			Name: "src, s",
			Value: "/path/to/src",
			Usage: "Source directory",
		},
		cli.StringFlag {
			Name: "dst, d",
			Value: "/path/to/dst",
			Usage: "Destination directory",
		},
		cli.StringFlag {
			Name: "photo-dir, P",
			Value: "photo",
			Usage: "Destination photo directory",
		},
		cli.StringFlag {
			Name: "video-dir, V",
			Value: "video",
			Usage: "Destination video directory",
		},
		cli.BoolFlag {
			Name: "recursive, r",
			Usage: "Resursive",
		},

	}
	app.Action = func(c *cli.Context) {
		fmt.Printf("Hello ghoto src=%s, dst=%s, photo-dir=%s, video-dir=%s, recursive=%s \n", c.String("src"), c.String("dst"), c.String("photo-dir"), c.String("video-dir"), c.Bool("recursive"))
	}
	app.Run(os.Args)
}
