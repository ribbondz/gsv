package cmd

import (
	"bufio"
	"fmt"
	"os"
)

func Head(path string, n int) {
	r, err := os.Open(path)
	if err != nil {
		fmt.Print("File does not exist.\n")
		return
	}
	defer r.Close()

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
