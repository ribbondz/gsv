package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ribbondz/gsv/cmd/utility"
	"github.com/schollz/progressbar/v2"
	"hash/fnv"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	BarUpdateThreshold = 1024 * 1024 * 20   // 20MB
	Batch              = 1024 * 1024 * 1000 // 1000MB
)

type BufHandler struct {
	summary     map[string]int // summary
	dstDir      string
	sep         []byte
	column      int
	headerBytes []byte
	header      bool
	lineN       int
}

func Partition(file string, header bool, column int, sep string, summary bool) {
	var et utility.ElapsedTime
	et.Start()

	// check file existence
	if !utility.FileIsExist(file) {
		fmt.Print("File does not exist. Try command 'gsv partition --help'.")
		return
	}
	r, _ := os.Open(file)
	defer r.Close()
	br := bufio.NewScanner(r)

	// estimate number of rows
	estimatedTotalN := utility.EstimateRowNumber(file, header, 20) //20MB
	fmt.Printf("Estimated row number: %d\n\n", estimatedTotalN)

	// struct to hold all options
	var handler BufHandler
	handler.summary = make(map[string]int)
	handler.dstDir = dstDirectory(file) // mkdir and return the path
	handler.sep = []byte(sep)
	handler.column = column
	handler.header = header

	if header && br.Scan() {
		handler.headerBytes = utility.CopyBytes(br.Bytes()) // must copy because s.token change under the hood
	}

	// progress bar
	size := utility.FileSize(file)
	bar := progressbar.NewOptions(size, progressbar.OptionSetBytes(size), progressbar.OptionSetRenderBlankState(true))

	type task struct {
		m     map[string][]byte // cached content
		byteN int               // number of processed bytes
	}
	jobs := make(chan task)

	// goroutine reading
	go func() {
		var (
			byteN = 0
			m     = make(map[string][]byte)
			line  []byte
		)

		for br.Scan() {
			handler.lineN++
			// submit to write if currentN is no less than conf.batch
			if byteN > Batch {
				jobs <- task{m, byteN}
				byteN = 0
				m = make(map[string][]byte)
			}

			line = br.Bytes()
			byteN += len(line) + 2 // 2 is for line terminator
			fields := bytes.Split(line, handler.sep)
			if len(fields) > handler.column {
				f := string(fields[handler.column])
				a := append(m[f], line...)
				a = append(a, '\n')
				m[f] = a
			}
		}

		// submit final content in the batch
		if len(m) > 0 {
			jobs <- task{m, byteN}
		}

		// close chan
		// so that the write function knows that there will be no content
		close(jobs)

		if err := br.Err(); err != nil {
			fmt.Println(err)
		}
	}()

	// write
	for t := range jobs {
		handler.SaveContent(t.m, bar)
	}
	bar.Finish()

	// print summary info
	fmt.Printf("\n\nLine count: %d, unique column value: %d\n", handler.lineN, len(handler.summary))

	// summary
	if summary {
		summaryFile := summaryFilename(file)
		WriteSummary(summaryFile, handler.summary)
	}

	et.EndAndPrint()
}

func (handler *BufHandler) SaveContent(content map[string][]byte, bar *progressbar.ProgressBar) {
	result := make(chan int, 100) // bytes saved
	done := make(chan int)        // wait progress bar update

	go func() {
		t := 0
		for b := range result {
			t += b
			if t > BarUpdateThreshold { // update every 20MB
				bar.Add(t)
				t = 0
			}
		}
		bar.Add(t)
		done <- 1
	}()

	for k, v := range content {
		result <- len(v)

		// append a header for first time write
		if handler.header {
			_, ok := handler.summary[k]
			if !ok {
				lh := len(handler.headerBytes) // copy is more efficient than append
				t := make([]byte, lh+1+len(v))
				copy(t, handler.headerBytes)
				copy(t[lh:], []byte{'\n'})
				copy(t[lh+1:], v)
				v = t
			}
		}

		handler.summary[k] += bytes.Count(v, []byte{'\n'})
		AppendToFile(handler.dstDir, k, v)
	}

	close(result)
	<-done
}

func dstDirectory(file string) string {
	wd, _ := os.Getwd()
	file = strings.TrimSuffix(file, filepath.Ext(file))
	file = utility.DirToFilename(file)
	timeStr := time.Now().Format("20060102150405")

	dir := filepath.Join(wd, file+"-"+timeStr)
	err := os.MkdirAll(dir, os.ModePerm)
	utility.CheckErr(err)
	return dir
}

func summaryFilename(file string) string {
	wd, _ := os.Getwd()
	file = strings.TrimSuffix(file, filepath.Ext(file))
	file = utility.DirToFilename(file)
	timeStr := time.Now().Format("20060102150405")
	return filepath.Join(wd, file+"-split-summary-"+timeStr+".txt")
}

func AppendToFile(dir string, col string, content []byte) {
	name := HashedFileName(col)
	file := filepath.Join(dir, name)
	f, _ := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	bf := bufio.NewWriter(f)
	bf.Write(content)
	bf.Flush()
	f.Close()
}

func HashedFileName(name string) (filename string) {
	h := fnv.New64a()
	h.Write([]byte(name))
	filename = strconv.FormatUint(h.Sum64(), 10) + ".txt"
	return
}

func WriteSummary(path string, summary map[string]int) {
	var result [][]string

	for k, v := range summary {
		result = append(result, []string{
			k,
			strconv.Itoa(v),
		})
	}

	// sort by col
	sort.Slice(result, func(i, j int) bool {
		a, _ := strconv.ParseInt(result[i][1], 10, 64)
		b, _ := strconv.ParseInt(result[j][1], 10, 64)
		return a > b
	})

	// add header
	result = utility.PrependStringSlice(result, []string{"col", "count"})

	utility.SaveFile(path, result)
	fmt.Printf("Summary file saved to: %s\n", path)
}
