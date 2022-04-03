package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	photoRe = regexp.MustCompile(`(?i)^[^\.].*\.(dng|cr2|jpg|jpeg|arw|orf)$`)
	videoRe = regexp.MustCompile(`(?i)^[^\.].*\.(mov|mp4|wmv|avi)$`)
)

func isDirectory(name string) (isDir bool, err error) {
	info, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func isFile(name string) (isFile bool, err error) {
	info, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

func readExifData(file string) (map[string]string, error) {
	cmdPath := "exiftool" // for linux or mac
	if runtime.GOOS == "windows" {
		cmdPath = "exiftool.exe"
	}

	cmd := exec.Command(cmdPath, file)

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

func pathFromExif(exif map[string]string) (string, error) {
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

func moveFile(src, dst string, option *Option) error {
	log.Printf("moveFile. src=%s, dst=%s", src, dst)
	if !option.Force {
		// 既にファイルが存在するなら書き込まない
		isFile, err := isFile(dst)
		if isFile || err != nil {
			return errors.New(fmt.Sprintf("dst=%s is already exists. if you want to overwrite please use --force option", dst))
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

func isIgnoreFile(name string, option *Option) bool {
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

func transferFile(name, filePath, dir string, option *Option) {
	if option.Verbose {
		log.Printf("transferFile. name=%s, filePath=%s, dir=%s", name, filePath, dir)
	}

	exif, _ := readExifData(filePath)
	if option.Verbose {
		log.Printf("exif=%v", exif)
	}
	dateDirPath, dateDirPathErr := pathFromExif(exif)
	if dateDirPathErr != nil {
		if option.SkipInvalidData {
			log.Printf("skip %s. Can't get date dir path.", filePath)
			return
		} else {
			log.Fatal(dateDirPathErr)
		}
	}

	dstDir := filepath.Join(option.To, dir, dateDirPath)
	if option.DryRun == false {
		mkdirErr := os.MkdirAll(dstDir, 0755)
		if mkdirErr != nil {
			log.Fatal(mkdirErr)
		}
	}

	dstPath := filepath.Join(dstDir, name)
	log.Printf("%s -> %s", filePath, dstPath)

	if option.DryRun == false {
		moveErr := moveFile(filePath, dstPath, option)
		if moveErr != nil {
			log.Fatal(moveErr)
		}
	}
}

func Transfer(wg *sync.WaitGroup, ch chan int, from string, option *Option) {
	if option.Verbose {
		log.Printf("Transfer. from=%s", from)
	}

	fileInfos, readErr := ioutil.ReadDir(from)
	if readErr != nil {
		log.Fatal(readErr)
	}

	for _, fileInfo := range fileInfos {
		wg.Add(1)
		go func(fileInfo os.FileInfo) {
			defer wg.Done()
			ch <- 1

			name := (fileInfo).Name()
			if option.Verbose {
				log.Printf("dir=%s, file=%s", from, name)
			}
			if isIgnoreFile(name, option) {
				if option.Verbose {
					log.Printf("ignored. dir=%s, file=%s", from, name)
				}
				<-ch
				return
			}

			filePath := filepath.Join(from, name)
			isDir, err := isDirectory(filePath)
			if err != nil {
				log.Fatal(err)
			}

			if isDir {
				if option.Recursive {
					Transfer(wg, ch, filePath, option)
				}
			} else {
				if photoRe.MatchString(name) {
					transferFile(name, filePath, option.PhotoDir, option)
				} else if videoRe.MatchString(name) {
					transferFile(name, filePath, option.VideoDir, option)
				}
			}

			<-ch
		}(fileInfo)
	}
}
