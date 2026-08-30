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
// It contains basic stats for combat and skill checks.
// Use NewActor to create instances with the fluent builder API.
type Actor struct {
	id              string         // Unique identifier (normalized to lowercase snake_case)
	maxHP           int            // Maximum Hit Points (base HP)
	currentHP       int            // Current Hit Points
	ac              int            // Armor Class (total, including all bonuses)
	initiative      int            // Initiative order (situational)
	combatModifiers []Modifier     // Active offensive modifiers for attack rolls
	attributes      map[string]int // Flexible attribute system (abilities, skills, etc.)
}

// ID returns the actor's normalized ID (lowercase snake_case).
// IDs are read-only after creation.
//
// Example:
//
//	fighter, err := d20.NewActor("ironpants Son of Arathorn").
//		WithHP(45).
//		WithAC(18).
//		Build()
//	fmt.Println(fighter.ID()) // "ironpants_son_of_arathorn"
func (a *Actor) ID() string {
	return a.id
}

// HP returns the actor's current hit points.
func (a *Actor) HP() int {
	return a.currentHP
}

// MaxHP returns the actor's maximum hit points.
func (a *Actor) MaxHP() int {
	return a.maxHP
}

// SetHP sets the actor's current hit points directly.
// HP cannot be set below 0 or above max HP.
func (a *Actor) SetHP(hp int) error {
	if hp < 0 {
		return fmt.Errorf("hp cannot be negative, got %d", hp)
	}
	if hp > a.maxHP {
		return fmt.Errorf("hp cannot exceed max HP (%d), got %d", a.maxHP, hp)
	}
	a.currentHP = hp
	return nil
}

// SetMaxHP sets the actor's maximum hit points.
// Max HP must be greater than 0. Current HP is adjusted if it exceeds the new max.
func (a *Actor) SetMaxHP(maxHP int) error {
	if maxHP <= 0 {
		return fmt.Errorf("max hp must be greater than 0, got %d", maxHP)
	}
	a.maxHP = maxHP
	// Adjust current HP if it exceeds new max
	if a.currentHP > a.maxHP {
		a.currentHP = a.maxHP
	}
	return nil
}

// SubHP reduces the actor's current HP by the specified amount.
// HP will not go below 0.
func (a *Actor) SubHP(damage int) {
	a.currentHP -= damage
	if a.currentHP < 0 {
		a.currentHP = 0
	}
}

// AddHP increases the actor's current HP by the specified amount.
// HP will not exceed max HP.
func (a *Actor) AddHP(amount int) {
	a.currentHP += amount
	if a.currentHP > a.maxHP {
		a.currentHP = a.maxHP
	}
}

// ResetHP restores the actor's current HP to maximum.
func (a *Actor) ResetHP() {
	a.currentHP = a.maxHP
}

// IsKnockedOut returns true if the actor has 0 HP.
func (a *Actor) IsKnockedOut() bool {
	return a.currentHP == 0
}

// AC returns the actor's Armor Class.
func (a *Actor) AC() int {
	return a.ac
}

// SetAC sets the actor's Armor Class. AC must be greater than 0.
func (a *Actor) SetAC(ac int) error {
	if ac <= 0 {
		return fmt.Errorf("ac must be greater than 0, got %d", ac)
	}
	a.ac = ac
	return nil
}

// Initiative returns the actor's initiative modifier.
func (a *Actor) Initiative() int {
	return a.initiative
}

// SetInitiative sets the actor's initiative modifier.
func (a *Actor) SetInitiative(initiative int) {
	a.initiative = initiative
}

// Attribute returns the value of the specified attribute and whether it exists.
// The key is automatically lowercased for consistent lookups.
func (a *Actor) Attribute(key string) (int, bool) {
	value, exists := a.attributes[strings.ToLower(key)]
	return value, exists
}

// SetAttribute sets the value of the specified attribute.
// The key is automatically lowercased for consistency.
func (a *Actor) SetAttribute(key string, value int) {
	a.attributes[strings.ToLower(key)] = value
}

// HasAttribute returns true if the actor has the specified attribute.
// The key is automatically lowercased for consistent lookups.
func (a *Actor) HasAttribute(key string) bool {
	_, exists := a.attributes[strings.ToLower(key)]
	return exists
}

// RemoveAttribute removes the specified attribute from the actor.
// The key is automatically lowercased for consistent lookups.
func (a *Actor) RemoveAttribute(key string) {
	delete(a.attributes, strings.ToLower(key))
}

// IncrementAttribute increases the value of the specified attribute by delta.
// If the attribute doesn't exist, it is created with the delta value.
// The key is automatically lowercased for consistency.
//
// Example:
//
//	actor.SetAttribute("strength", 16)
//	actor.IncrementAttribute("strength", 2) // Now 18 (temporary buff)
func (a *Actor) IncrementAttribute(key string, delta int) {
	key = strings.ToLower(key)
	a.attributes[key] = a.attributes[key] + delta
}

// DecrementAttribute decreases the value of the specified attribute by delta.
// If the attribute doesn't exist, it is created with the negative delta value.
// The key is automatically lowercased for consistency.
//
// Example:
//
//	actor.SetAttribute("hp", 45)
//	actor.DecrementAttribute("hp", 10) // Now 35 (took damage)
func (a *Actor) DecrementAttribute(key string, delta int) {
	key = strings.ToLower(key)
	a.attributes[key] = a.attributes[key] - delta
}

// AddCombatModifier adds a modifier to the actor's combat modifiers.
// Accepts name and value, creates the Modifier internally.
// The modifier name is automatically lowercased for consistency.
//
// Example:
//
//	actor.AddCombatModifier("strength", 3)
//	actor.AddCombatModifier("bless", 2)
func (a *Actor) AddCombatModifier(name string, value int) {
	a.combatModifiers = append(a.combatModifiers, NewModifier(name, value))
}

// RemoveCombatModifier removes all modifiers with the specified reason.
// The reason is automatically lowercased for consistent lookups.
func (a *Actor) RemoveCombatModifier(reason string) {
	reason = strings.ToLower(reason)
	filtered := make([]Modifier, 0, len(a.combatModifiers))
	for _, mod := range a.combatModifiers {
		if mod.Reason != reason {
			filtered = append(filtered, mod)
		}
	}
	a.combatModifiers = filtered
}

// GetCombatModifiers returns a copy of the actor's combat modifiers.
// Returns a copy to prevent external mutations.
func (a *Actor) GetCombatModifiers() []Modifier {
	modifiers := make([]Modifier, len(a.combatModifiers))
	copy(modifiers, a.combatModifiers)
	return modifiers
}

// SkillCheck creates a DiceManager for a skill check using D&D 5e conventions (1d20 + skill modifier).
// The skill value is looked up from the actor's Attributes map.
// Returns a DiceManager pre-configured with the skill modifier. Chain .WithAdvantage() or other
// modifiers as needed, then call .Roll() to execute.
//
// Returns an error if the skill is not found.
//
// Example:
//
//	actor.SetAttribute("athletics", 5)
//	builder, err := actor.SkillCheck("athletics", roller)
//	if err != nil {
//		return err
//	}
//	result, err := builder.WithAdvantage().Roll()
func (a *Actor) SkillCheck(skill string, roller *Roller) (*DiceManager, error) {
	skillValue, exists := a.Attribute(skill)
	if !exists {
		return nil, fmt.Errorf("skill %q not found in actor attributes", skill)
	}

	// Return a DiceManager with skill modifier pre-loaded
	return roller.Dice(1, 20).WithModifier(skill, skillValue), nil
}

// AttackRoll creates a DiceManager for an attack roll using the actor's CombatModifiers.
// Uses D&D 5e conventions (1d20 + combat modifiers).
// Returns a DiceManager pre-configured with all combat modifiers. Chain .WithAdvantage(),
// .WithModifier() for situational bonuses, then call .Roll() to execute.
//
// Example:
//
//	actor.AddCombatModifier("strength", 5)
//	actor.AddCombatModifier("proficiency", 3)
//	result, _ := actor.AttackRoll(roller).WithAdvantage().Roll()
//
//	// With situational modifier
//	result, _ := actor.AttackRoll(roller).WithModifier("flanking", 2).Roll()
func (a *Actor) AttackRoll(roller *Roller) *DiceManager {
	builder := roller.Dice(1, 20)

	// Add all combat modifiers
	for _, mod := range a.combatModifiers {
		builder = builder.WithModifier(mod.Reason, mod.Value)
	}

	return builder
}

// D100SkillCheck returns a 2d10 DiceManager for a roll-under skill check.
// Errors if the skill is missing. Compare outcome.Value to Attribute(skill) for success.
//
//	builder, _ := actor.D100SkillCheck("stealth", roller)
//	outcome, _ := builder.RollPercentile()
func (a *Actor) D100SkillCheck(skill string, roller *Roller) (*DiceManager, error) {
	if _, exists := a.Attribute(skill); !exists {
		return nil, fmt.Errorf("skill %q not found in actor attributes", skill)
	}

	return roller.Dice(2, 10), nil
}
