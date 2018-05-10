package meander

import (
	"strings"
)

type Cost int8

const (
	Cost1 = Cost(1 + iota)
	// NOTE: When you don't specify a value of constants, the previous one (in this
	// case `Cost(1 + iota)`) is used. And `iota` corresponds to a position of
	// constants in the const block. For these reasons, `Cost2` is equal to 2.
	Cost2
	Cost3
	Cost4
	Cost5
)

var costStrings = map[string]Cost{
	"$":     Cost1,
	"$$":    Cost2,
	"$$$":   Cost3,
	"$$$$":  Cost4,
	"$$$$$": Cost5,
}

func (c Cost) String() string {
	for s, v := range costStrings {
		if v == c {
			return s
		}
	}

	return "不正な値です"
}

// Q. This method returns 0 if `costString` doesn't include `s`. Is this OK?
func ParseCost(s string) Cost {
	return costStrings[s]
}

type CostRange struct {
	From Cost
	To   Cost
}

func (r *CostRange) String() string {
	return r.From.String() + "..." + r.To.String()
}

func ParseCostRange(s string) *CostRange {
	segs := strings.Split(s, "...")

	return &CostRange{
		From: ParseCost(segs[0]),
		To:   ParseCost(segs[1]),
	}
}
