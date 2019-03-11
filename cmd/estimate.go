package cmd

import (
	"bytes"
	"os"
)

func EstimateRowNumber(path string, header bool, mb int) (estimatedRowN int) {
	r, err := os.Open(path)
	CheckErr(err)
	defer func() {
		err = r.Close()
		CheckErr(err)
	}()

	// guess rows by 10M data
	buf := make([]byte, BufSize*mb)
	r.Read(buf)

	// delete header
	if header {
		n := bytes.IndexByte(buf, '\n')
		if n > -1 {
			buf = buf[n+1:]
		}
	}

	// delete last incomplete row
	n := bytes.LastIndexByte(buf, '\n')
	if n > -1 {
		buf = buf[:n]
	}
	buf = bytes.TrimSpace(buf)

	bytesPerRow := len(buf) / (bytes.Count(buf, []byte{'\n'}) + 1)
	if bytesPerRow == 0 {
		return 0
	}

	// total file size
	total := FileSize(path)
	estimatedRowN = total / bytesPerRow
	return
}

func FileSize(file string) int {
	f, err := os.Stat(file)
	if os.IsNotExist(err) {
		return 0
	} else {
		return int(f.Size())
	}
}
