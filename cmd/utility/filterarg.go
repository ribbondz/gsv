package utility

import (
	"strconv"
	"strings"
)

const (
	IsAnd = iota
	IsOr
	IsNull
	IsFloatFilter
	IsStringFilter
)

type ColumnFilter struct {
	Applicable bool
	T          int
	FloatV     []float64 // list of possible values
	StringV    []string  // list of possible values
	allV       []string  // list of all values
}

type Filter struct {
	Op      int
	Filters []ColumnFilter
}

// filter initialization with args
func NewFilter(arg string, columnN int) (*Filter, error) {
	r := &Filter{Op: IsNull}
	for i := 0; i < columnN; i++ {
		r.Filters = append(r.Filters, ColumnFilter{})
	}

	// no filter
	if len(arg) == 0 {
		return r, nil
	}

	if strings.Contains(arg, "&") {
		r.Op = IsAnd
		for _, filter := range strings.Split(arg, "&") {
			err := r.handleOneFilter(filter)
			if err != nil {
				return r, err
			}
		}
	} else if strings.Contains(arg, "|") {
		r.Op = IsOr
		for _, filter := range strings.Split(arg, "|") {
			err := r.handleOneFilter(filter)
			if err != nil {
				return r, err
			}
		}
	} else {
		r.Op = IsAnd
		err := r.handleOneFilter(arg)
		if err != nil {
			return r, err
		}
	}

	return r, nil
}

// handle one filter condition
func (f *Filter) handleOneFilter(arg string) error {
	splits := strings.Split(arg, "=")
	c, err := strconv.Atoi(splits[0])
	if err != nil {
		return err
	} else {
		fc := &f.Filters[c]
		fc.Applicable = true
		fc.allV = append(fc.allV, splits[1])

		if fc.T == IsFloatFilter {
			v, err := strconv.ParseFloat(splits[1], 64)
			if err == nil {
				fc.FloatV = append(fc.FloatV, v)
			} else {
				fc.T = IsStringFilter
				fc.StringV = append([]string{}, fc.allV...) // transform all previous value into string
			}
		} else if fc.T == IsStringFilter {
			fc.StringV = append(fc.StringV, splits[1])
		} else {
			v, err := strconv.ParseFloat(splits[1], 64)
			if err == nil {
				fc.T = IsFloatFilter
				fc.FloatV = append(fc.FloatV, v)
			} else {
				fc.T = IsStringFilter
				fc.StringV = append(fc.StringV, splits[1])
			}
		}
	}

	return nil
}

// apply filter to row
func (f *Filter) FilterOneRowSatisfy(row []string) bool {
	switch f.Op {
	case IsNull: // no filter specified
		return true
	case IsAnd:
		for i, filter := range f.Filters {
			if filter.Applicable {
				if filter.T == IsStringFilter {
					if !SliceContainsString(filter.StringV, row[i]) {
						return false
					}
				} else {
					v, err := strconv.ParseFloat(row[i], 64)
					if err != nil || !SliceContainsFloat(filter.FloatV, v) {
						return false
					}
				}
			}
		}
		return true
	case IsOr:
		for i, filter := range f.Filters {
			if filter.Applicable {
				if filter.T == IsStringFilter && SliceContainsString(filter.StringV, row[i]) {
					return true
				} else {
					v, err := strconv.ParseFloat(row[i], 64)
					if err == nil && SliceContainsFloat(filter.FloatV, v) {
						return true
					}
				}
			}
		}
		return false
	}

	return true
}
