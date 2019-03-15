package utility

import (
	"os"
	"strings"
)

func FileSize(file string) int {
	f, err := os.Stat(file)
	if os.IsNotExist(err) {
		return 0
	} else {
		return int(f.Size())
	}
}

func FileIsExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func DirToFilename(dir string) string {
	dir = strings.ReplaceAll(dir, string('\\'), "-")
	dir = strings.ReplaceAll(dir, string('/'), "-")
	dir = strings.ReplaceAll(dir, "--", "-")
	if len(dir) > 1 && strings.HasSuffix(dir, "-") {
		dir = dir[0 : len(dir)-1]
	}
	return dir
}
