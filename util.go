package main

import (
	"os"
)

func IsDirectory(name string) (isDir bool, err error) {
	info, err := os.Stat(name)
	if err != nil {
			return false, err
	}
	return info.IsDir(), nil
}

func MoveFile(src, dst string) (error) {
	return os.Rename(src, dst)
}
