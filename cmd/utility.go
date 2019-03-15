package cmd

import (
	"bufio"
	"encoding/csv"
	"os"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func headerBytes(dst string) (header []byte) {
	r, _ := os.Open(dst)
	br := bufio.NewScanner(r)
	br.Scan()
	header = br.Bytes()
	r.Close()
	return
}

func CopyBytes(source []byte) []byte {
	dst := make([]byte, len(source))
	copy(dst, source)
	return dst
}

func SaveFile(file string, list [][]string) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	CheckErr(err)
	defer func() {
		err = f.Close()
		CheckErr(err)
	}()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	err = writer.WriteAll(list)
	CheckErr(err)
	//fmt.Println("Save to file: ", file)
}
