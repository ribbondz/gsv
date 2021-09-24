package utility

import (
	"errors"
	"strconv"
	"strings"
)

type ColArgs struct {
	All     bool
	Include []int
	Exclude []int
}

// ParseColArg
// examples:
// 1,2
// 1,2-4,6
// !1
// -1
func ParseColArg(col string) (r ColArgs, err error) {
	// all columns
	if col == "-1" || col == "" {
		r.All = true
		return
	}
	for _, c := range strings.Split(col, ",") {
		// avoid    1,2-4,  =>   ['1', '2', '']
		if len(c) == 0 {
			continue
		}
		if strings.Contains(c, "-") {
			s := strings.Split(c, "-")
			min, err1 := strconv.Atoi(s[0])
			max, err2 := strconv.Atoi(s[1])
			// if any syntax error exist, abort the program
			if err1 != nil || err2 != nil {
				return r, errors.New("error column select")
			}
			for ; min <= max; min++ {
				r.Include = append(r.Include, min)
			}
		} else if strings.Contains(c, "!") {
			c := strings.ReplaceAll(c, "!", "")
			v, err := strconv.Atoi(c)
			if err != nil {
				return r, err
			}
			r.Exclude = append(r.Exclude, v)
		} else {
			if v, err := strconv.Atoi(c); err == nil {
				r.Include = append(r.Include, v)
			} else {
				return r, errors.New("error column select")
			}
		}
	}
	if len(r.Include) == 0 && len(r.Include) > 0 {
		r.All = true
	}
	return
}

func AllIncludedCols(col ColArgs, totalColumn int) (r []int) {
	if col.All {
		for i := 0; i < totalColumn; i++ {
			if !SliceContainsInt(col.Exclude, i) {
				r = append(r, i)
			}
		}
	} else {
		for _, v := range col.Include {
			if !SliceContainsInt(col.Exclude, v) {
				r = append(r, v)
			}
		}
	}
	return
}
