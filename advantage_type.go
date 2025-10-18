package d20

// AdvantageType represents the advantage/disadvantage state for a roll.
// This is a core 5e mechanic where advantage means roll twice and take higher,
// disadvantage means roll twice and take lower.
type AdvantageType int

const (
	Disadvantage AdvantageType = iota // Roll twice, take lower
	Normal                            // Roll normally
	Advantage                         // Roll twice, take higher
)
