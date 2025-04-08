package filter_types

type FilterType int

const (
	None FilterType = iota + 0
	Block
	BoxFilter
	DoubleBoxFilter
)
