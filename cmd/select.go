package cmd

import (
	"bufio"
	"fmt"
	"github.com/ribbondz/gsv/cmd/utility"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type RowContent struct {
	n       int
	content string
}

func Select(file string, header bool, sep string, filterPara string, colPara utility.ColArgs, out bool) {
	var et utility.ElapsedTime
	et.Start()

	// check file existence
	if !utility.FileIsExist(file) {
		fmt.Print("File does not exist. Try command 'gsv select --help'.")
		return
	}

	// saved columns
	columnN := ColumnN(file, sep)                    // how many columns
	col := utility.AllIncludedCols(colPara, columnN) // all included columns []int

	// filters
	filter, err := utility.NewFilter(filterPara, columnN)
	if err != nil {
		fmt.Print("Filter syntax error. Try command 'gsv select --help'.")
		return
	}

	// file processing
	f, _ := os.Open(file)
	defer f.Close()
	br := bufio.NewScanner(f)

	// writer
	dstFilename := OutFilenameFilter(file)
	r, err := os.OpenFile(dstFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	bw := bufio.NewWriter(r)

	// header with saved column
	if header {
		br.Scan()
		s := keepSavedCol(strings.Split(br.Text(), sep), col, columnN)
		ss := strings.Join(s, sep)
		if out {
			bw.WriteString(ss)
			bw.Write([]byte{'\n'})
		} else {
			fmt.Println(ss)
		}
	}

	jobs := make(chan []string, 20)      // batch rows
	results := make(chan RowContent, 20) // batch results, filtered out rows joined in a string
	wg := &sync.WaitGroup{}              // wait all batch processing
	total := 0                           // filter out rows count

	// worker, process batch rows
	// the number of worker defaults to cpu cores
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for job := range jobs {
				results <- FilterProcessRows(filter, job, sep, col, columnN)
			}
		}()
	}

	// collect batch result, and merge it into main result
	go func() {
		for result := range results {
			// avoid no content
			if len(result.content) > 0 {
				total += result.n
				if out {
					bw.WriteString(result.content)
				} else {
					fmt.Print(result.content) // has \n, so use fmt.Print
				}
			}
			wg.Done() // indicate work done
		}
	}()

	n := 0             // batch number of rows
	var batch []string //batch holder
	for br.Scan() {
		batch = append(batch, br.Text())
		n++
		if n > BatchRowsPerStat { // 2000 rows per batch
			wg.Add(1)
			jobs <- batch
			n = 0
			batch = []string{}
		}
	}

	if len(batch) > 0 {
		wg.Add(1)
		jobs <- batch
	}
	close(jobs)

	wg.Wait()

	// delete out file if not saving
	bw.Flush()
	r.Close()
	if !out {
		os.Remove(dstFilename)
	}

	fmt.Println("Total filtered rows: ", total)
	et.EndAndPrint()
}

func keepSavedCol(fields []string, col []int, columnN int) []string {
	if len(col) == columnN {
		return fields
	}

	var dst []string
	for i, v := range fields {
		if utility.SliceContainsInt(col, i) {
			dst = append(dst, v)
		}
	}
	return dst
}

func FilterProcessRows(f *utility.Filter, rows []string, sep string, col []int, columnN int) RowContent {
	var r [][]string
	for _, row := range rows {
		splits := strings.Split(row, sep)
		if f.FilterOneRowSatisfy(splits) {
			r = append(r, keepSavedCol(splits, col, columnN))
		}
	}

	// quick return
	if len(r) == 0 {
		return RowContent{0, ""}
	}

	var sb strings.Builder
	for _, s := range r {
		sb.WriteString(strings.Join(s, sep))
		sb.WriteByte('\n')
	}

	return RowContent{len(r), sb.String()}
}

// output filename, data.txt has the default out filename data-current-time.txt
func OutFilenameFilter(file string) string {
	wd, _ := os.Getwd()
	file = strings.TrimSuffix(file, filepath.Ext(file))
	file = utility.DirToFilename(file)
	timeStr := time.Now().Format("20060102150405")
	return filepath.Join(wd, file+"-select-"+timeStr+".txt")
}
