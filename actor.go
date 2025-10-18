package d20

// Actor represents a character, NPC, or monster in the game world.
// It contains basic stats for combat and skill checks.
type Actor struct {
	HP              int            // Hit Points
	AC              int            // Armor Class (total, including all bonuses)
	Initiative      int            // Initiative/speed modifier
	CombatModifiers []Modifier     // Active offensive modifiers for attack rolls
	Attributes      map[string]int // Flexible attribute system (abilities, skills, etc.)
}
