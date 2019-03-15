package utility

import (
	"testing"
)

func TestColArgParse(t *testing.T) {
	testArgs := []string{"1", "1,2", "1,2,", "1,3-4,6", "-1", "!0"}
	trueArgs := []ColArgs{
		{false, []int{1}, []int{}},
		{false, []int{1, 2}, []int{}},
		{false, []int{1, 2}, []int{}},
		{false, []int{1, 3, 4, 6}, []int{}},
		{true, []int{}, []int{}},
		{true, []int{}, []int{0}},
	}

	for i := range testArgs {
		p, _ := ParseColArg(testArgs[i])
		a := trueArgs[i]
		if p.all != a.all || !SliceIntEqual(p.include, a.include) || !SliceIntEqual(p.exclude, a.exclude) {
			t.Error("col arg parse error.")
		}
	}
}

func SliceIntEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
