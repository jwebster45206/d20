// Package vocab is a shared vocabulary of well-known attribute and modifier names.
// Package d20 does not import this package; StrikeRoll takes strings.
package vocab

const (
	Strength     = "strength"
	Dexterity    = "dexterity"
	Constitution = "constitution"
	Intelligence = "intelligence"
	Wisdom       = "wisdom"
	Charisma     = "charisma"
	Striking     = "striking" // strike / attack proficiency
	Damage       = "damage"   // damage bonus/penalty; not applied by StrikeRoll
)
