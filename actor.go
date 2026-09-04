package d20

import (
	"fmt"
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
//   - "Goblin#3" -> "goblin_3"
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

// SubHP reduces current HP by the specified amount. HP will not go below 0.
// Non-positive amounts are ignored.
func (a *Actor) SubHP(amount int) {
	if amount <= 0 {
		return
	}
	a.HP -= amount
	if a.HP < 0 {
		a.HP = 0
	}
}

// AddHP increases current HP by the specified amount. HP will not exceed MaxHP.
// Non-positive amounts are ignored.
func (a *Actor) AddHP(amount int) {
	if amount <= 0 {
		return
	}
	a.HP += amount
	if a.HP > a.MaxHP {
		a.HP = a.MaxHP
	}
}

// ResetHP restores current HP to maximum.
func (a *Actor) ResetHP() {
	a.HP = a.MaxHP
}

// IsKnockedOut returns true if the actor has 0 or fewer HP.
func (a *Actor) IsKnockedOut() bool {
	return a.HP <= 0
}

// SkillDice returns 1d20 plus the stored attribute value as a modifier.
// The library does not derive 5e ability modifiers; callers store the bonus they want added.
// The skill name is lowercased for lookup; map keys are lowercase by convention.
// Returns an error if the skill is missing.
//
//	d, err := actor.SkillDice("athletics")
//	out, err := roller.Roll(d.WithAdvantage())
func (a *Actor) SkillDice(skill string) (Dice, error) {
	skill = strings.ToLower(skill)
	skillValue, exists := a.Attributes[skill]
	if !exists {
		return Dice{}, fmt.Errorf("skill %q not found in actor attributes", skill)
	}
	d, err := NewDice(1, 20)
	if err != nil {
		return Dice{}, err
	}
	return d.WithModifier(skill, skillValue), nil
}

// StrikeDice returns 1d20 with the named modifier keys applied.
// Missing keys are skipped. Key names are lowercased for lookup.
// Situational extras go on the returned Dice (WithModifier) without mutating this spec.
//
//	d, err := actor.StrikeDice("strength", "striking")
//	out, err := roller.Roll(d.WithAdvantage())
func (a *Actor) StrikeDice(keys ...string) (Dice, error) {
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

// D100SkillDice returns 2d10 for a roll-under skill check.
// Errors if the skill is missing. Compare outcome.Value to Attributes[skill] for success.
//
//	d, err := actor.D100SkillDice("stealth")
//	out, err := roller.RollPercentile(d)
func (a *Actor) D100SkillDice(skill string) (Dice, error) {
	skill = strings.ToLower(skill)
	if _, exists := a.Attributes[skill]; !exists {
		return Dice{}, fmt.Errorf("skill %q not found in actor attributes", skill)
	}
	return NewDice(2, 10)
}
