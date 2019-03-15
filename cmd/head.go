package cmd

import (
	"bufio"
	"fmt"
	"github.com/ribbondz/gsv/cmd/utility"
	"os"
)

func Head(path string, n int) {
	// check file existence
	if !utility.FileIsExist(path) {
		fmt.Print("File does not exist.")
		return
	}

	r, err := os.Open(path)
	defer r.Close()
	utility.CheckErr(err)

	br := bufio.NewScanner(r)
	i := 0
	for br.Scan() {
		i++
		fmt.Println(br.Text())
		if i >= n {
			break
		}
	}
}
