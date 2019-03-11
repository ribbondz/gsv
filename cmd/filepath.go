package cmd

import (
	"os"
	"strings"
)

func DirToFilename(dir string) string {
	dir = strings.ReplaceAll(dir, string('\\'), "-")
	dir = strings.ReplaceAll(dir, string('/'), "-")
	dir = strings.ReplaceAll(dir, "--", "-")
	if len(dir) > 1 && strings.HasSuffix(dir, "-") {
		dir = dir[0 : len(dir)-1]
	}
	return dir
}

func FileIsExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
