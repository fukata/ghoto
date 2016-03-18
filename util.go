package main

import (
	"os"
	"os/exec"
	"bytes"
	"strings"
	"time"
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
