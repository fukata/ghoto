package main

import (
	"os"
	"os/exec"
	"bytes"
	"strings"
	"time"
	"io/ioutil"
	"regexp"
	"log"
	"sync"
)

var (
	re = regexp.MustCompile(`(?i)^[^\.].*\.(dng|cr2|jpg|jpeg|arw|orf)$`)
)
func IsDirectory(name string) (isDir bool, err error) {
	info, err := os.Stat(name)
	if err != nil {
			return false, err
	}
	return info.IsDir(), nil
}

func GetExifData(file string) (map[string]string, error) {
	cmd := exec.Command("exiftool", file)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	tags := make(map[string]string)

	data := strings.Trim(out.String(), " \r\n")
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		k, v := strings.Replace(strings.TrimSpace(line[0:32]), " ", "", -1), strings.TrimSpace(line[33:])
		tags[k] = v
	}

	return tags, nil
}

func GetDateDirPath(date string) (string, error) {
	t, err := time.Parse("2006:01:02 15:04:05", date)
	if err != nil {
		return "", err
	}
	return t.Format("2006/01/02"), nil
}

func MoveFile(src, dst string) (error) {
	cmd := exec.Command("mv", src, dst)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
	//return os.Rename(src, dst)
}

func IsIgnoreFile(name string, option *Option) (bool) {
	return name == "." || name == ".." || name == "lightroom"
}

func GetFileNum(from string, option *Option) (int, error) {
	fileInfos, readErr := ioutil.ReadDir(from + "/")
	if readErr != nil {
		return 0, readErr
	}

	num := 0
	for _, fileInfo := range fileInfos {
		name := (fileInfo).Name()
		if IsIgnoreFile(name, option) {
			continue
		}

		filePath := from + "/" + name
		isDir, err := IsDirectory(filePath)
		if err != nil {
			return 0, err
		}

		if isDir {
			if option.Recursive {
				subNum, subErr := GetFileNum(from + "/" + name, option)
				if subErr != nil {
					return 0, subErr
				}
				num += subNum
			}
		} else {
			if re.MatchString(name) {
				num += 1
			}
		}
	}

	return num, nil
}

func MoveFiles(wg *sync.WaitGroup, ch chan int, from string, option *Option) {
	fileInfos, readErr := ioutil.ReadDir(from + "/")
	if readErr != nil {
		log.Fatal(readErr)
	}

	for _, fileInfo := range fileInfos {
		wg.Add(1)
		go func(fileInfo os.FileInfo) {
			defer wg.Done()
			ch <- 1

			name := (fileInfo).Name()
			if IsIgnoreFile(name, option) {
				<-ch
				return
			}

			filePath := from + "/" + name
			isDir, err := IsDirectory(filePath)
			if err != nil {
				log.Fatal(err)
			}

			if isDir {
				if option.Recursive {
					go MoveFiles(wg, ch, filePath, option)
				}
			} else {
				if re.MatchString(name) == false {
					<-ch
					return
				}

				exif, _ := GetExifData(filePath)
				dateDirPath, dateDirPathErr := GetDateDirPath(exif["Date/TimeOriginal"])
				//TODO don't have exif
				if dateDirPathErr != nil {
					log.Fatal(dateDirPathErr)
				}

				dstDir := option.To + "/" + option.PhotoDir + "/" + dateDirPath + "/"
				if option.DryRun == false {
					mkdirErr := os.MkdirAll(dstDir, 0755)
					if mkdirErr != nil {
						log.Fatal(mkdirErr)
					}
				}

				dstPath := dstDir + name
				log.Printf("%s -> %s", filePath, dstPath)

				if option.DryRun == false {
					moveErr := MoveFile(filePath, dstPath)
					if moveErr != nil {
						log.Fatal(moveErr)
					}
				}
			}

			<-ch
		}(fileInfo)
	}
}
