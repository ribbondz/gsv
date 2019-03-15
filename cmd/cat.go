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
	"sync"
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
	defer dstW.Close()

	if header {
		headerContent := utility.HeaderBytes(files[0])
		WriteBytes(dstW, headerContent)
	}

	// progress bar
	bar := progressbar.New(len(files))
	bar.RenderBlank()

	jobs := make(chan string, 100)
	results := make(chan []byte, 100)

	for w := 1; w <= runtime.NumCPU(); w++ {
		go worker(w, jobs, results, header)
	}

	go func() {
		for _, file := range files {
			jobs <- file
		}
		close(jobs)
	}()

	wg := &sync.WaitGroup{}
	for range files {
		content := <-results
		WriteBytes(dstW, content)
		wg.Add(1)
		go func() {
			bar.Add(1)
			wg.Done()
		}()
	}

	wg.Wait()
	bar.Finish()
	dstW.Sync()
	fmt.Printf("\n\nSave to file: %s\n", dst)
	et.EndAndPrint()
}

func worker(id int, jobs <-chan string, results chan<- []byte, header bool) {
	for path := range jobs {
		results <- ReadOneFile(path, header)
	}
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
