package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Option struct {
	From            string
	To              string
	PhotoDir        string
	VideoDir        string
	Recursive       bool
	DryRun          bool
	Excludes        []string
	Concurrency     int
	Force           bool
	SkipInvalidData bool
	Verbose         bool
}

const VERSION = "0.1.1"

func main() {
	app := &cli.App{
		UseShortOptionHandling: true,
		Name:                   "ghoto",
		Version:                VERSION,
		Compiled:               time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "fukata",
				Email: "tatsuya.fukata@gmail.com",
			},
		},
		Usage: "Transfer photo(video)",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "from", Value: "/path/to/src", Usage: "Source directory"},
			&cli.StringFlag{Name: "to", Value: "/path/to/dst", Usage: "Destination directory"},
			&cli.StringFlag{Name: "photo-dir", Aliases: []string{"P"}, Value: "photo", Usage: "Destination photo directory"},
			&cli.StringFlag{Name: "video-dir", Aliases: []string{"V"}, Value: "video", Usage: "Destination video directory"},
			&cli.StringFlag{Name: "exclude", Aliases: []string{"x"}, Value: "", Usage: "Exclude dir/file separate comma."},
			&cli.IntFlag{Name: "concurrency", Aliases: []string{"c"}, Value: runtime.NumCPU(), Usage: "Concurrency num."},
			&cli.BoolFlag{Name: "recursive", Aliases: []string{"r"}, Usage: "Resursive"},
			&cli.BoolFlag{Name: "force", Usage: "Force"},
			&cli.BoolFlag{Name: "skip-invalid-data", Usage: "SkipInvalidData"},
			&cli.BoolFlag{Name: "dry-run", Usage: "Dry Run"},
			&cli.BoolFlag{Name: "verbose, vvv", Usage: "Verbose"},
		},
		Action: func(c *cli.Context) error {
			// options
			option := &Option{
				c.String("from"),
				c.String("to"),
				c.String("photo-dir"),
				c.String("video-dir"),
				c.Bool("recursive"),
				c.Bool("dry-run"),
				strings.Split(c.String("exclude"), ","),
				c.Int("concurrency"),
				c.Bool("force"),
				c.Bool("skip-invalid-data"),
				c.Bool("verbose"),
			}

			// check path
			isDir, err := IsDirectory(option.From)
			if err != nil {
				log.Fatal(err)
			}

			if isDir != false {
				fmt.Errorf("%s is not found.", option.From)
			}

			isDir, err = IsDirectory(option.To)
			if err != nil {
				log.Fatal(err)
			}

			if isDir != false {
				fmt.Errorf("%s is not found.", option.To)
			}

			// move
			var wg sync.WaitGroup
			ch := make(chan int, option.Concurrency)
			Transfer(&wg, ch, option.From, option)
			wg.Wait()
			log.Println("done")

			return err
		},
	}

	app.Run(os.Args)
}
