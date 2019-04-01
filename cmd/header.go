package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/ribbondz/gsv/cmd/utility"
)

func Header(file string, sep string) {
	// check file existence
	if !utility.FileIsExist(file) {
		fmt.Print("File does not exist. Try command 'gsv header --help'.")
		return
	}

	// open file
	f, _ := os.Open(file)
	defer f.Close()
	br := bufio.NewScanner(f)

	// first and second rows
	br.Scan()
	row1 := strings.Split(br.Text(), sep)
	br.Scan()
	row2 := strings.Split(br.Text(), sep)

	var result [][]string
	for i := range row1 {
		result = append(result, []string{
			strconv.Itoa(i),
			row1[i],
			row2[i],
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "header", "row example"})
	table.SetBorder(true)
	table.AppendBulk(result)
	table.SetCaption(true, "Total columns: "+strconv.Itoa(len(result)))
	table.Render()
}
