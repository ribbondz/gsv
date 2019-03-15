package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/ribbondz/gsv/cmd/utility"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	IsInt = iota
	IsFloat
	IsString
	IsNull

	BatchRowsPerStat = 2000 //rows per batch
)

type ColStats struct {
	cType     int
	nulls     int
	minLength int
	maxLength int

	intStats   IntColStats
	floatStats FloatColStats
	strStats   StringColStats
}

type StringColStats struct {
	min       string
	max       string
	uniqueMap map[string]int
}

type IntColStats struct {
	min       int
	max       int
	total     int // can overflow, not handle it yet
	uniqueMap map[int]int
}

type FloatColStats struct {
	min   float64
	max   float64
	total float64
}

func Stats(file string, header bool, sep string) {
	var et utility.ElapsedTime
	et.Start()

	// check file existence
	if !utility.FileIsExist(file) {
		fmt.Print("File does not exist.")
		return
	}

	// column types
	colTypes, firstValue, err := GuessColType(file, header, sep) // "" is string
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// stats initial
	stat := statsInit(colTypes, firstValue)

	// stats processing
	f, _ := os.Open(file)
	defer f.Close()
	br := bufio.NewScanner(f)

	// column names and header drop
	var names []string
	if header {
		br.Scan()
		names = strings.Split(br.Text(), sep)
	} else {
		for i := range colTypes {
			names = append(names, "col"+strconv.Itoa(i+1))
		}
	}

	jobs := make(chan []string, 20)
	results := make(chan []ColStats, 20)
	wg := &sync.WaitGroup{}

	// worker
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for job := range jobs {
				results <- processRow(job, colTypes, firstValue, sep)
			}
		}()
	}

	// collect result
	go func() {
		for result := range results {
			stat = mergeStats(stat, result)
			wg.Done()
		}
	}()

	// reading file in main thread
	var batch []string
	var totalN = 0
	var n = 0
	for br.Scan() {
		totalN++

		batch = append(batch, br.Text())

		n++
		if n > BatchRowsPerStat {
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
	PrintStats(stat, names, totalN)

	et.EndAndPrint()
}

// len(lines) > 0
func processRow(lines []string, colTypes []int, firstValue []string, sep string) []ColStats {
	stats := statsInit(colTypes, firstValue)

	for _, line := range lines {
		fields := strings.Split(line, sep)
		for i, field := range fields {
			cs := &stats[i]

			l := len(field)
			if l < cs.minLength {
				cs.minLength = l
			}
			if l > cs.maxLength {
				cs.maxLength = l
			}

			// null
			if field == "" || field == "NA" || field == "Na" || field == "na" || field == "Null" || field == "NULL" {
				cs.nulls++
			} else {
				switch cs.cType {
				case IsString:
					if field < cs.strStats.min {
						cs.strStats.min = field
					}
					if field > cs.strStats.max {
						cs.strStats.max = field
					}
					cs.strStats.uniqueMap[field] = 1
				case IsInt:
					if v, err := strconv.ParseInt(field, 10, 64); err == nil {
						b := int(v)
						if b < cs.intStats.min {
							cs.intStats.min = b
						}
						if b > cs.intStats.max {
							cs.intStats.max = b
						}
						cs.intStats.total += b
						cs.intStats.uniqueMap[b] = 1
					} else {
						fmt.Println("parsing error happen to a row.")
					}
				case IsFloat:
					if b, err := strconv.ParseFloat(field, 64); err == nil {
						if b < cs.floatStats.min {
							cs.floatStats.min = b
						}
						if b > cs.floatStats.max {
							cs.floatStats.max = b
						}
						cs.floatStats.total += b
					} else {
						fmt.Println("parsing error happen to a row.")
					}
				}
			}
		}
	}

	return stats
}

// merge two stats
func mergeStats(dst []ColStats, s []ColStats) []ColStats {
	for i := range dst {
		a, b := &dst[i], s[i]
		a.nulls += b.nulls

		if a.minLength > b.minLength {
			a.minLength = b.minLength
		}
		if a.maxLength < b.maxLength {
			a.maxLength = b.maxLength
		}

		switch a.cType {
		case IsString:
			if a.strStats.min == "" {
				a.strStats.min = b.strStats.min
			} else if a.strStats.min > b.strStats.min {
				a.strStats.min = b.strStats.min
			}

			if a.strStats.max == "" {
				a.strStats.max = b.strStats.max
			} else if a.strStats.max < b.strStats.max {
				a.strStats.max = b.strStats.max
			}

			for k, _ := range b.strStats.uniqueMap {
				a.strStats.uniqueMap[k] = 1
			}
		case IsInt:
			if a.intStats.min > b.intStats.min {
				a.intStats.min = b.intStats.min
			}

			if a.intStats.max < b.intStats.max {
				a.intStats.max = b.intStats.max
			}

			a.intStats.total += b.intStats.total
			for k, _ := range b.intStats.uniqueMap {
				a.intStats.uniqueMap[k] = 1
			}
		case IsFloat:
			if a.floatStats.min > b.floatStats.min {
				a.floatStats.min = b.floatStats.min
			}

			if a.floatStats.max < b.floatStats.max {
				a.floatStats.max = b.floatStats.max
			}

			a.floatStats.total += b.floatStats.total
		}
	}
	return dst
}

func PrintStats(stat []ColStats, names []string, totalN int) {
	// avoid zero division
	if totalN == 0 {
		totalN++
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"col", "type", "null", "unique", "min", "max", "mean", "min_length", "max_length"})
	table.SetBorder(true)
	for i, s := range stat {
		switch s.cType {
		case IsString:
			table.Append([]string{
				names[i],
				"string",
				strconv.Itoa(s.nulls),
				strconv.Itoa(len(s.strStats.uniqueMap)),
				s.strStats.min,
				s.strStats.max,
				"-",
				strconv.Itoa(s.minLength),
				strconv.Itoa(s.maxLength),
			})
		case IsInt:
			table.Append([]string{
				names[i],
				"int",
				strconv.Itoa(s.nulls),
				strconv.Itoa(len(s.intStats.uniqueMap)),
				strconv.Itoa(s.intStats.min),
				strconv.Itoa(s.intStats.max),
				strconv.FormatFloat(float64(s.intStats.total)/float64(totalN), 'f', 4, 64),
				strconv.Itoa(s.minLength),
				strconv.Itoa(s.maxLength),
			})
		case IsFloat:
			table.Append([]string{
				names[i],
				"float",
				strconv.Itoa(s.nulls),
				"-",
				strconv.FormatFloat(s.floatStats.min, 'f', 4, 64),
				strconv.FormatFloat(s.floatStats.max, 'f', 4, 64),
				strconv.FormatFloat(s.floatStats.total/float64(totalN), 'f', 4, 64),
				strconv.Itoa(s.minLength),
				strconv.Itoa(s.maxLength),
			})
		}
	}

	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetCaption(true, "Total records: "+strconv.Itoa(totalN))
	table.Render()
}

func statsInit(colTypes []int, firstValue []string) (stat []ColStats) {
	for i, ct := range colTypes {
		var cs ColStats
		cs.strStats.uniqueMap = make(map[string]int)
		cs.intStats.uniqueMap = make(map[int]int)
		cs.cType = ct

		// initial min and max to the first value of file
		if ct == IsString {
			cs.strStats.min = firstValue[i]
			cs.strStats.max = firstValue[i]
		} else if ct == IsInt {
			v, _ := strconv.Atoi(firstValue[i])
			cs.intStats.min = v
			cs.intStats.max = v
		} else {
			v, _ := strconv.ParseFloat(firstValue[i], 64)
			cs.floatStats.min = v
			cs.floatStats.max = v
		}

		// initial max_length min_length
		cs.minLength = 9999999
		cs.maxLength = -9999999
		stat = append(stat, cs)
	}
	return
}

func GuessColType(file string, header bool, sep string) ([]int, []string, error) {
	var (
		guessN     = 2000
		line       = ""
		cType      []int
		firstValue []string
		fields     []string
	)

	// type initialization
	cn := ColumnN(file, sep)
	for i := 0; i < cn; i++ {
		cType = append(cType, IsNull)
		firstValue = append(firstValue, "")
	}

	f, _ := os.Open(file)
	defer f.Close()
	br := bufio.NewScanner(f)

	if header {
		br.Scan()
		br.Bytes()
	}

	for br.Scan() && guessN > 0 {
		guessN--

		line = br.Text()
		fields = strings.Split(line, sep)
		if len(fields) != len(cType) {
			return []int{}, []string{}, errors.New("rows have unequal length")
		}

		for i, field := range fields {
			// skip null values
			if field == "" || field == "NA" || field == "Na" || field == "na" || field == "Null" || field == "NULL" {
				continue
			}

			// obtain first not-null value
			if firstValue[i] == "" && len(strings.TrimSpace(field)) > 0 {
				firstValue[i] = field
			}

			// string is always string
			if cType[i] == IsString {
				continue
			}

			// is int
			if _, err := strconv.ParseInt(field, 10, 64); err == nil {
				if cType[i] == IsNull {
					cType[i] = IsInt
				}
				continue
			}

			// is float
			if _, err := strconv.ParseFloat(field, 64); err == nil {
				if cType[i] == IsNull || cType[i] == IsInt {
					cType[i] = IsFloat
				}
				continue
			}

			cType[i] = IsString
		}
	}

	return cType, firstValue, nil
}

func ColumnN(file string, sep string) int {
	f, _ := os.Open(file)
	defer f.Close()
	br := bufio.NewScanner(f)

	br.Scan()
	line := br.Text()
	return len(strings.Split(line, sep))
}
