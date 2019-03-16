package utility

import (
	"bufio"
	"encoding/csv"
	"os"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func HeaderBytes(dst string) (header []byte) {
	r, _ := os.Open(dst)
	br := bufio.NewScanner(r)
	br.Scan()
	header = br.Bytes()
	r.Close()
	return
}

func CopyBytes(source []byte) []byte {
	dst := make([]byte, len(source))
	copy(dst, source)
	return dst
}

func SaveFile(file string, list [][]string) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	CheckErr(err)

	writer := csv.NewWriter(f)
	writer.WriteAll(list)
	writer.Flush()
	f.Close()
}

func SliceContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SliceContainsFloat(s []float64, e float64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SliceContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
