package utility

import (
	"fmt"
	"time"
)

type ElapsedTime struct {
	startTime time.Time
	elapsed   int
}

func (t *ElapsedTime) Start() {
	t.startTime = time.Now()
}

func (t *ElapsedTime) EndAndPrint() {
	a := time.Since(t.startTime) // ns
	a = a / 1e9                  // s

	if a > 60 {
		fmt.Printf("Time consumed: %dm %ds\n", a/60, a%60)
	} else {
		fmt.Printf("Time consumed: %ds\n", a)
	}
}
