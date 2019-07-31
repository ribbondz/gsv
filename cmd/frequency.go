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

	columnN := ColumnN(file, sep)                    // how many columns
	col := utility.AllIncludedCols(colPara, columnN) // all included columns []int

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

	jobs := make(chan []string, 20)            // batch rows
	results := make(chan []map[string]int, 20) // batch processed result
	wg := &sync.WaitGroup{}                    // wait for all batches to be processed
	freq := freqMapInit(col)                   // data structure to save frequency table, []map[string]int

	// worker, process batch rows
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for job := range jobs {
				results <- processRows(job, sep, col)
			}
		}()
	}

	// collect batch result, and merge it into main result
	go func() {
		for result := range results {
			freq = MergeMapList(freq, result)
			wg.Done()
		}
	}()

	N := 0              // total number of rows
	n := 0              // batch number of rows
	batch := []string{} //batch holder
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

	// wait all batch result to be processed
	wg.Wait()

	// generate freq table ([][]string) from results ([]map[string]int)
	// apply ascending option
	// apply limit option
	table := GenerateFreqTable(freq, names, ascending, limit)

	// apply out option
	if out {
		table = utility.PrependStringSlice(table, []string{"col", "value", "count"})
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

// process batch rows
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

// batch frequency tables are merged into main table sequentially
func MergeMapList(l1, l2 []map[string]int) []map[string]int {
	for i := range l1 {
		a, b := l1[i], l2[i]
		for k, v := range b {
			a[k] += v
		}
	}
	return l1
}

// initial frequency table to a list of N maps,
// each map for a column,
// the map records column value (key) and counts (value)
func freqMapInit(includeColumn []int) []map[string]int {
	var r []map[string]int

	maxColumnIndex := includeColumn[len(includeColumn)-1]
	for i := 0; i <= maxColumnIndex; i++ {
		m := make(map[string]int)
		r = append(r, m)
	}
	return r
}

// transform list of N maps to a frequency table,
// the table has structure [][]string
// the table can be feed into a tablewriter to print in stdout, or to be saved into a file
// the function also sort records according to -ascending flag,
// default to descending order
func GenerateFreqTable(freq []map[string]int, names []string, ascending bool, limit int) (r [][]string) {
	for i, l := range freq {
		var one [][]string
		for k, v := range l {
			one = append(one, []string{names[i], k, strconv.Itoa(v)})
		}

		// apply ascending option
		if ascending {
			sort.Slice(one, func(i, j int) bool {
				a, _ := strconv.Atoi(one[i][2])
				b, _ := strconv.Atoi(one[j][2])
				return a < b
			})
		} else {
			sort.Slice(one, func(i, j int) bool {
				a, _ := strconv.Atoi(one[i][2])
				b, _ := strconv.Atoi(one[j][2])
				return a > b
			})
		}

		// apply limit option
		if limit > 0 && len(one) >= limit {
			one = one[0:limit]
		}

		t := make([][]string, len(r)+len(one))
		copy(t, r)
		copy(t[len(r):], one)
		r = t
	}

	return
}

// use tablewriter to print frequency table to stdout
func PrintFreqTable(freq [][]string, totalN int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"col", "value", "count"})
	table.SetBorder(true)
	table.AppendBulk(freq)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetCaption(true, "Total records of file:"+strconv.Itoa(totalN))
	table.Render()
}

// output filename, data.txt has the default out filename data-current-time.txt
func OutFilename(file string) string {
	wd, _ := os.Getwd()
	file = strings.TrimSuffix(file, filepath.Ext(file))
	file = utility.DirToFilename(file)
	timeStr := time.Now().Format("20060102150405")
	return filepath.Join(wd, file+"-frequency-table-"+timeStr+".txt")
}
