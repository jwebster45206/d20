package d20

// Result represents the outcome of a dice roll.
type Result struct {
	Total int
	Rolls []int
}

// String returns a string representation of the roll result.
func (r Result) String() string {
	return formatResult(r)
}
