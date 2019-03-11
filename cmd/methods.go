package cmd

import (
	"bufio"
	"os"
)

func headerBytes(dst string) (header []byte) {
	r, _ := os.Open(dst)
	br := bufio.NewScanner(r)
	br.Scan()
	header = br.Bytes()
	r.Close()
	return
}

func CopyBytes(source []byte) (dst []byte) {
	copy(dst, source)
	return
}