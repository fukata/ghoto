package main

import (
	"fmt"
	"io"
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
	photoRe = regexp.MustCompile(`(?i)^[^\.].*\.(dng|cr2|jpg|jpeg|arw|orf)$`)
	videoRe = regexp.MustCompile(`(?i)^[^\.].*\.(mov|mp4|wmv|avi)$`)
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

func GetDateDirPath(exif map[string]string) (string, error) {
	t, err := time.Parse("2006:01:02 15:04:05", exif["Date/TimeOriginal"])
	if err == nil {
		return t.Format("2006/01/02"), nil
	}

	t, err = time.Parse("2006:01:02 15:04:05", exif["CreateDate"])
	if err == nil {
		return t.Format("2006/01/02"), nil
	}

	return "", err
}

func MoveFile(src, dst string, option *Option) (error) {
	if !option.Force {
		_, err := os.Stat(dst)
		if err != nil {
			return err
		}
	}

	inputFile, err := os.Open(src)
	if err != nil {
			return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(dst)
	if err != nil {
			inputFile.Close()
			return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
			return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(src)
	if err != nil {
			return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func IsIgnoreFile(name string, option *Option) (bool) {
	if name == "." || name == ".." {
		return true
	}

	for _, v := range option.Excludes {
		if name == v {
			return true
		}
	}

	return false
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
			if photoRe.MatchString(name) || videoRe.MatchString(name) {
				num += 1
			}
		}
	}

	return num, nil
}

func TransferFile(name, filePath, dir string, option *Option) {
	exif, _ := GetExifData(filePath)
	if option.Verbose {
		log.Println(exif)
	}
	dateDirPath, dateDirPathErr := GetDateDirPath(exif)
	if dateDirPathErr != nil {
		if option.SkipInvalidData {
			log.Printf("skip %s. Can't get date dir path.", filePath)
			return
		} else {
			log.Fatal(dateDirPathErr)
		}
	}

	dstDir := option.To + "/" + dir + "/" + dateDirPath + "/"
	if option.DryRun == false {
		mkdirErr := os.MkdirAll(dstDir, 0755)
		if mkdirErr != nil {
			log.Fatal(mkdirErr)
		}
	}

	dstPath := dstDir + name
	log.Printf("%s -> %s", filePath, dstPath)

	if option.DryRun == false {
		moveErr := MoveFile(filePath, dstPath, option)
		if moveErr != nil {
			log.Fatal(moveErr)
		}
	}
}

func Transfer(wg *sync.WaitGroup, ch chan int, from string, option *Option) {
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
					Transfer(wg, ch, filePath, option)
				}
			} else {
				if photoRe.MatchString(name) {
					TransferFile(name, filePath, option.PhotoDir, option)
				} else if videoRe.MatchString(name) {
					TransferFile(name, filePath, option.VideoDir, option)
				}
			}

			<-ch
		}(fileInfo)
	}
}
