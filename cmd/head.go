package cmd

import (
	"bufio"
	"fmt"
	"os"
)

func Head(path string, n int) {
	// check file existence
	if !FileIsExist(path) {
		fmt.Print("File does not exist.")
		return
	}

	r, err := os.Open(path)
	defer r.Close()
	CheckErr(err)

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
