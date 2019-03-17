package cmd

import (
	"bytes"
	"fmt"
	"github.com/ribbondz/gsv/cmd/utility"
	"github.com/schollz/progressbar/v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	MBBytes = 1024 * 1024 // 1MB
)

func Cat(dir string, header bool, pattern string) {
	var et utility.ElapsedTime
	et.Start()

	// all files
	files := fileList(dir, pattern)
	if len(files) == 0 {
		fmt.Print("No files matched.")
		return
	}
	fmt.Printf("Total files: %d\n\n", len(files))

	// dst file
	dst := dstFile(dir)
	dstW, _ := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if header {
		headerContent := utility.HeaderBytes(files[0])
		WriteBytes(dstW, headerContent)
	}

	// progress bar
	bar := progressbar.NewOptions(len(files), progressbar.OptionSetRenderBlankState(true))

	jobs := make(chan string, 100)    // file pool, can only open 100 files at a time
	results := make(chan []byte, 100) // file results

	// put file into the pool
	go func() {
		for _, file := range files {
			jobs <- file
		}
		close(jobs)
	}()

	// worker, read file
	for w := 1; w <= runtime.NumCPU(); w++ {
		go func() {
			for path := range jobs {
				results <- ReadOneFile(path, header)
			}
		}()
	}

	n := 0 // update progress bar every 5 files
	for range files {
		n++
		content := <-results
		WriteBytes(dstW, content)

		if n > 4 {
			n = 0
			go func() {
				bar.Add(1)
			}()
		}
	}

	bar.Finish()
	dstW.Sync()
	dstW.Close()
	fmt.Printf("\n\nSaved to file: %s\n", dst)
	et.EndAndPrint()
}

func fileList(dir string, pattern string) (files []string) {
	wd, err := os.Getwd()
	utility.CheckErr(err)
	p1 := filepath.Join(wd, dir, pattern)

	fmt.Printf("Match pattern: %s\n", p1)

	files, err = filepath.Glob(p1)
	utility.CheckErr(err)
	return
}

func dstFile(dir string) string {
	wd, _ := os.Getwd()
	timeStr := time.Now().Format("20060102150405")

	// clean dir to prevent save output to sub directories.
	dir = utility.DirToFilename(dir)
	return filepath.Join(wd, dir+"-"+timeStr+".txt")
}

func ReadOneFile(path string, header bool) (byteContent []byte) {
	byteContent, err := ioutil.ReadFile(path)
	utility.CheckErr(err)

	// header
	if header {
		n := bytes.IndexByte(byteContent, '\n')
		if n > -1 {
			byteContent = byteContent[n+1:]
		}
	}

	// unify all new line '\r\n' to '\n'
	byteContent = bytes.ReplaceAll(byteContent, []byte{'\r', '\n'}, []byte{'\n'})
	byteContent = bytes.TrimSpace(byteContent)
	return
}

func WriteBytes(w *os.File, content []byte) (n int, err error) {
	// avoid adding new empty lines if content is empty
	if len(content) == 0 {
		return
	}
	content = append(content, '\n')
	n, err = w.Write(content)
	return
}
