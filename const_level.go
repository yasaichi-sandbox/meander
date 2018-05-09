package meander

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
