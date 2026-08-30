package d20

// AdvantageType represents the advantage/disadvantage state for a roll.
// 5e: advantage rolls twice and takes higher, disadvantage rolls twice and takes lower.
type AdvantageType int

const (
	Normal       AdvantageType = iota // Roll normally (zero value)
	Advantage                         // Roll twice, take higher
	Disadvantage                      // Roll twice, take lower
)
