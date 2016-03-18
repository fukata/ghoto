package main

import (
	"os"
	"fmt"
	"log"
	"io/ioutil"
	"regexp"
	"sync"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ghoto"
	app.Usage = "Transfer photo(video)"
	app.Flags = []cli.Flag {
		cli.StringFlag {
			Name: "from",
			Value: "/path/to/src",
			Usage: "Source directory",
		},
		cli.StringFlag {
			Name: "to",
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
		cli.BoolFlag {
			Name: "dry-run",
			Usage: "Dry Run",
		},
	}
	app.Action = func(c *cli.Context) {
		fmt.Printf("Hello ghoto from=%s, to=%s, photo-dir=%s, video-dir=%s, recursive=%s, dry-run=%s \n", c.String("from"), c.String("to"), c.String("photo-dir"), c.String("video-dir"), c.Bool("recursive"), c.Bool("dry-run"))

		re := regexp.MustCompile(`(?i)^[^\.].*\.(dng|cr2|jpg|jpeg|arw|orf)$`)

		// options
		from := c.String("from")
		to := c.String("to")
		photoDir := c.String("photo-dir")
		//videoDir := c.String("video-dir")
		recursive := c.Bool("recursive")
		dryRun := c.Bool("dry-run")

		// check path
		isDir, err := IsDirectory(from)
		if err != nil {
			log.Fatal(err)
		}

		if isDir != false {
			fmt.Errorf("%s is not found.", from)
		}

		isDir, err = IsDirectory(to)
		if err != nil {
			log.Fatal(err)
		}

		if isDir != false {
			fmt.Errorf("%s is not found.", to)
		}

		fileInfos, readDirErr := ioutil.ReadDir(from + "/")
		if readDirErr != nil {
			log.Fatal(err)
		}

		var wg sync.WaitGroup

		//fileNames := []string{} 

		for _, fi := range fileInfos {
			wg.Add(1)
			go func(fileInfo os.FileInfo) {
				defer wg.Done()
				name := (fileInfo).Name()
				filePath := from + "/" + name

				//log.Printf("file=%s", filePath)
				isDir, err = IsDirectory(filePath)
				if err != nil {
					log.Fatal(err)
				}

				if isDir && recursive {
					// recursive
				} else {
					if re.MatchString(name) == false {
						return
					}

					//log.Printf("match file=%s", filePath)
					exif, _ := GetExifData(filePath)
					dateDirPath, dateDirPathErr := GetDateDirPath(exif["Date/TimeOriginal"])
					if dateDirPathErr != nil {
						log.Fatal(dateDirPathErr)
					}

					dstDir := to + "/" + photoDir + "/" + dateDirPath + "/"
					if dryRun == false {
						mkdirErr := os.MkdirAll(dstDir, 0755)
						if mkdirErr != nil {
							log.Fatal(mkdirErr)
						}
					}

					dstPath := dstDir + name
					log.Printf("%s -> %s", filePath, dstPath)

					if dryRun == false {
						moveErr := MoveFile(filePath, dstPath)
						if moveErr != nil {
							log.Fatal(moveErr)
						}
					}
				}
			}(fi)
		}

		wg.Wait()
	}
	app.Run(os.Args)
}
