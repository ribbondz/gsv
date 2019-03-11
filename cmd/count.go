package cmd

import (
	"bytes"
	"fmt"
	"os"
)

func Count(path string, header bool) (nRow int) {
	if !FileIsExist(path) {
		fmt.Println("File doest not exist.")
		return
	}

	r, err := os.Open(path)
	CheckErr(err)

	buf := make([]byte, BufSize)

	for {
		n, _ := r.Read(buf)

		// count char '\n'
		nRow += bytes.Count(buf[0:n], []byte{'\n'})

		// read finished.
		if n < BufSize {
			break
		}
	}

	if header {
		nRow--
	}

	r.Close()
	fmt.Printf("%d\n", nRow)
	return
}
