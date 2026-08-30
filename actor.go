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
	ID              string         // Unique identifier (normalized by NewActor)
	MaxHP           int            // Maximum Hit Points
	HP              int            // Current Hit Points
	AC              int            // Armor Class
	Initiative      int            // Initiative order (situational)
	CombatModifiers []Modifier     // Offensive modifiers for attack rolls
	Attributes      map[string]int // Flexible attributes (abilities, skills, etc.)
}

// NewActor creates an Actor with a normalized ID and initialized maps.
func NewActor(id string) *Actor {
	return &Actor{
		ID:              normalizeID(id),
		CombatModifiers: []Modifier{},
		Attributes:      make(map[string]int),
	}
}

// SubHP reduces current HP by the specified amount. HP will not go below 0.
func (a *Actor) SubHP(damage int) {
	a.HP -= damage
	if a.HP < 0 {
		a.HP = 0
	}
}

// AddHP increases current HP by the specified amount. HP will not exceed MaxHP.
func (a *Actor) AddHP(amount int) {
	a.HP += amount
	if a.HP > a.MaxHP {
		a.HP = a.MaxHP
	}
}

// ResetHP restores current HP to maximum.
func (a *Actor) ResetHP() {
	a.HP = a.MaxHP
}

// IsKnockedOut returns true if the actor has 0 HP.
func (a *Actor) IsKnockedOut() bool {
	return a.HP == 0
}

// SkillCheck creates a DiceManager for a skill check using D&D 5e conventions (1d20 + skill modifier).
// The skill value is looked up from Attributes (key is lowercased). Returns an error if missing.
//
// Example:
//
//	actor.Attributes["athletics"] = 5
//	builder, err := actor.SkillCheck("athletics", roller)
//	if err != nil {
//		return err
//	}
//	result, err := builder.WithAdvantage().Roll()
func (a *Actor) SkillCheck(skill string, roller *Roller) (*DiceManager, error) {
	skill = strings.ToLower(skill)
	skillValue, exists := a.Attributes[skill]
	if !exists {
		return nil, fmt.Errorf("skill %q not found in actor attributes", skill)
	}

	return roller.Dice(1, 20).WithModifier(skill, skillValue), nil
}

// AttackRoll creates a DiceManager for an attack roll using CombatModifiers.
// Uses D&D 5e conventions (1d20 + combat modifiers).
//
// Example:
//
//	actor.CombatModifiers = []Modifier{
//		NewModifier("strength", 5),
//		NewModifier("proficiency", 3),
//	}
//	result, _ := actor.AttackRoll(roller).WithAdvantage().Roll()
func (a *Actor) AttackRoll(roller *Roller) *DiceManager {
	builder := roller.Dice(1, 20)

	for _, mod := range a.CombatModifiers {
		builder = builder.WithModifier(mod.Reason, mod.Value)
	}

	return builder
}

// D100SkillCheck returns a 2d10 DiceManager for a roll-under skill check.
// Errors if the skill is missing. Compare outcome.Value to Attributes[skill] for success.
//
//	builder, _ := actor.D100SkillCheck("stealth", roller)
//	outcome, _ := builder.RollPercentile()
func (a *Actor) D100SkillCheck(skill string, roller *Roller) (*DiceManager, error) {
	skill = strings.ToLower(skill)
	if _, exists := a.Attributes[skill]; !exists {
		return nil, fmt.Errorf("skill %q not found in actor attributes", skill)
	}

	return roller.Dice(2, 10), nil
}
