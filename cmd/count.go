package cmd

import (
	"bytes"
	"fmt"
	"github.com/ribbondz/gsv/cmd/utility"
	"os"
)

func Count(path string, header bool) (nRow int) {
	var et utility.ElapsedTime
	et.Start()
	if !utility.FileIsExist(path) {
		fmt.Println("File doest not exist. Try command 'gsv count --help'.")
		return
	}
	// 1. is directory: count files in directory
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		f, _ := os.Open(path)
		list, _ := f.Readdir(-1)
		fmt.Printf("Total files: %d\n", len(list))

		var t int64 = 0
		for _, i := range list {
			t += i.Size()
		}
		PrintFileSize(int(t))
		f.Close()
		return
	}
	// 2. is file: count lines in file
	r, err := os.Open(path)
	defer r.Close()
	utility.CheckErr(err)
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
	fmt.Printf("%d", nRow)
	et.EndAndPrint()
	return
}

func PrintFileSize(c int) {
	b := float64(c)
	mb := 1024.0 * 1024.0
	gb := mb * 1024.0
	if b < mb { //1MB
		fmt.Printf("Total size:  %.2fKB", b/1024.0)
	} else if b < 1024*mb {
		fmt.Printf("Total size:  %.2fMB", b/mb)
	} else {
		fmt.Printf("Total size:  %.2fGB", b/gb)
	}
}
