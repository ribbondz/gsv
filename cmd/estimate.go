package cmd

import (
	"bytes"
	"os"
)

const (
	MBBytes = 1024 * 1024 // 1MB
)

func EstimateRowNumber(path string, header bool, mb int) (estimatedRowN int) {
	r, _ := os.Open(path)
	defer r.Close()

	// guess rows by 10M data
	buf := make([]byte, MBBytes*mb)
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
