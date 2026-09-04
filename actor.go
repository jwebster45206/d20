package d20

import (
	"regexp"
	"strings"
)

var (
	// Regex to match any non-alphanumeric character for ID normalization
	nonAlphaNumeric = regexp.MustCompile(`[^a-z0-9]+`)
)

// normalizeID converts a string to lowercase snake_case for consistent IDs.
// Handles spaces, hyphens, special characters, etc.
//
// Examples:
//   - "Ironpants" -> "ironpants"
//   - "Busta the Black" -> "busta_the_black"
//   - "Fighter-1" -> "fighter_1"
func normalizeID(id string) string {
	id = strings.ToLower(id)
	id = nonAlphaNumeric.ReplaceAllString(id, "_")
	id = strings.Trim(id, "_")
	return id
}

// Actor represents a character, NPC, or monster in the game world.
//
// Example:
//
//	fighter := d20.NewActor("ironpants Son of Arathorn")
//	fighter.MaxHP, fighter.HP = 45, 45
//	fighter.AC = 18
//	fmt.Println(fighter.ID) // "ironpants_son_of_arathorn"
type Actor struct {
	ID         string         // Unique identifier (normalized by NewActor)
	MaxHP      int            // Maximum Hit Points
	HP         int            // Current Hit Points
	AC         int            // Armor Class
	Initiative int            // Initiative order (situational)
	Attributes map[string]int // Caller-owned numbers (ability scores or skill bonuses). Keys are lowercase.
	Modifiers  map[string]int // Caller-wired roll bonuses; not derived from Attributes. Keys are lowercase.
}

// NewActor creates an Actor with a normalized ID and initialized maps.
func NewActor(id string) *Actor {
	return &Actor{
		ID:         normalizeID(id),
		Attributes: make(map[string]int),
		Modifiers:  make(map[string]int),
	}
}

// D20Dice returns 1d20 with the named modifier keys applied.
// Missing keys are skipped. Key names are lowercased for lookup.
// Situational extras go on the returned Dice (WithModifier) without mutating this spec.
//
//	d, err := actor.D20Dice("strength", "striking")
//	out, err := roller.Roll(d.WithAdvantage())
func (a *Actor) D20Dice(keys ...string) (Dice, error) {
	d, err := NewDice(1, 20)
	if err != nil {
		return Dice{}, err
	}
	for _, name := range keys {
		name = strings.ToLower(name)
		if v, ok := a.Modifiers[name]; ok {
			d = d.WithModifier(name, v)
		}
	}
	return d, nil
}
