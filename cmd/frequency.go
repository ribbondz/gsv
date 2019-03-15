package cmd

import (
	"bufio"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/ribbondz/gsv/cmd/utility"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Frequency(file string, header bool, sep string, colPara utility.ColArgs, out bool, ascending bool, limit int) {
	var et utility.ElapsedTime
	et.Start()

	// check file existence
	if !utility.FileIsExist(file) {
		fmt.Print("File does not exist. Try command 'gsv frequency --help'.")
		return
	}

	// all columns to generate a frequency table
	columnN := ColumnN(file, sep)
	col := utility.AllIncludedCols(colPara, columnN)

	// file processing
	f, _ := os.Open(file)
	defer f.Close()
	br := bufio.NewScanner(f)

	// column names and header drop
	var names []string
	if header {
		br.Scan()
		names = strings.Split(br.Text(), sep)
	} else {
		for i := 0; i < columnN; i++ {
			names = append(names, "col_"+strconv.Itoa(i+1))
		}
	}

	jobs := make(chan []string, 20)
	results := make(chan []map[string]int, 20)
	wg := &sync.WaitGroup{}
	freq := freqMapInit(col)

	// worker
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for job := range jobs {
				results <- processRows(job, sep, col)
			}
		}()
	}

	// collect result
	go func() {
		for result := range results {
			freq = MergeMapList(freq, result)
			wg.Done()
		}
	}()

	N := 0
	n := 0
	batch := []string{}
	for br.Scan() {
		N++

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

	// generate freq table from []map[string]int
	table := GenerateFreqTable(freq, names, ascending)

	// apply limit option
	if limit > 0 && len(table) >= limit {
		table = table[0:limit]
	}

	// apply out option
	if out {
		table = append([][]string{{"col", "value", "count"}}, table...)
		outFile := OutFilename(file)
		utility.SaveFile(outFile, table)
		fmt.Println("Frequency table saved to: ", outFile)
	} else {
		PrintFreqTable(table, N)
		if limit > 0 {
			fmt.Println("Limit: ", limit)
		}
	}

	et.EndAndPrint()
}

func processRows(rows []string, sep string, col []int) []map[string]int {
	r := freqMapInit(col)

	for _, row := range rows {
		for i, field := range strings.Split(row, sep) {
			if utility.SliceContainsInt(col, i) {
				r[i][field]++
			}
		}
	}
	return r
}

func MergeMapList(l1, l2 []map[string]int) []map[string]int {
	for i := range l1 {
		a, b := l1[i], l2[i]
		for k, v := range b {
			a[k] += v
		}
	}
	return l1
}

func freqMapInit(includeColumn []int) []map[string]int {
	var r []map[string]int

	maxColumnIndex := includeColumn[len(includeColumn)-1]
	for i := 0; i <= maxColumnIndex; i++ {
		m := make(map[string]int)
		r = append(r, m)
	}
	return r
}

func GenerateFreqTable(freq []map[string]int, names []string, ascending bool) (r [][]string) {
	for i, l := range freq {
		for k, v := range l {
			r = append(r, []string{names[i], k, strconv.Itoa(v)})
		}
	}

	if ascending {
		sort.Slice(r, func(i, j int) bool {
			a, _ := strconv.Atoi(r[i][2])
			b, _ := strconv.Atoi(r[j][2])
			return a < b
		})
	} else {
		sort.Slice(r, func(i, j int) bool {
			a, _ := strconv.Atoi(r[i][2])
			b, _ := strconv.Atoi(r[j][2])
			return a > b
		})
	}

	return
}

func PrintFreqTable(freq [][]string, totalN int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"col", "value", "count"})
	table.SetBorder(true)
	table.AppendBulk(freq)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetCaption(true, "Total records of file: "+strconv.Itoa(totalN))
	table.Render()
}

func OutFilename(file string) string {
	wd, _ := os.Getwd()
	file = strings.TrimSuffix(file, filepath.Ext(file))
	file = utility.DirToFilename(file)
	timeStr := time.Now().Format("20060102150405")
	return filepath.Join(wd, file+"-frequency-table-"+timeStr+".txt")
}
