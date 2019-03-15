package cmd

import (
	"bytes"
	"fmt"
	"os"
)

func Count(path string, header bool) (nRow int) {
	var et ElapsedTime
	et.Start()

	if !FileIsExist(path) {
		fmt.Println("File doest not exist.")
		return
	}

	r, err := os.Open(path)
	defer r.Close()
	CheckErr(err)

	var bufSize = MBBytes * 10 // 10MB
	buf := make([]byte, bufSize)

	for {
		n, _ := r.Read(buf)

		// count char '\n'
		nRow += bytes.Count(buf[0:n], []byte{'\n'})

		// read finished.
		if n < bufSize {
			break
		}
	}

	if header {
		nRow--
	}

	fmt.Printf("%d\n", nRow)
	et.EndAndPrint()
	return
}
