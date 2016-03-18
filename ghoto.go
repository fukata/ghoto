package main

import (
	"os"
	"fmt"
	"log"
	"io/ioutil"
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

		// options
		src := c.String("src")
		dst := c.String("dst")
		recursive := c.Bool("recursive")

		// check path
		isDir, err := IsDirectory(src)
		if err != nil {
			log.Fatal(err)
		}

		if isDir != false {
			fmt.Errorf("%s is not found.", src)
		}

		isDir, err = IsDirectory(dst)
		if err != nil {
			log.Fatal(err)
		}

		if isDir != false {
			fmt.Errorf("%s is not found.", dst)
		}

		fileInfos, readDirErr := ioutil.ReadDir(src + "/")
		if readDirErr != nil {
			log.Fatal(err)
		}

		for _, fileInfo := range fileInfos {
			name := (fileInfo).Name()
			srcPath := src + "/" + name

			log.Printf("file=%s", name)
			isDir, err = IsDirectory(srcPath)
			if err != nil {
				log.Fatal(err)
			}

			if isDir && recursive {
				// recursive
			} else {
				//MoveFile()
			}
		}
	}
	app.Run(os.Args)
}
