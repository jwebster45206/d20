package d20

import "strings"

// Modifier is a bonus or penalty applied to a dice roll.
// Fields are public for flexibility, but it's recommended to use lowercase
// for the Reason field to maintain consistency in formatted output.
type Modifier struct {
	Value  int    // Positive for bonus, negative for penalty
	Reason string // Description of the modifier source (e.g., "strength", "proficiency")
}

// NewModifier creates a new Modifier with the reason automatically lowercased
// for consistent formatting. This is the recommended way to create modifiers.
func NewModifier(reason string, value int) Modifier {
	return Modifier{
		Value:  value,
		Reason: strings.ToLower(reason),
	}
}
