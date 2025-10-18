package d20

import (
	"fmt"
	"strings"
)

// Actor represents a character, NPC, or monster in the game world.
// It contains basic stats for combat and skill checks.
// Use NewActor or NewActorWithAttributes to create instances.
type Actor struct {
	maxHP           int            // Maximum Hit Points (base HP)
	currentHP       int            // Current Hit Points
	ac              int            // Armor Class (total, including all bonuses)
	initiative      int            // Initiative/speed modifier
	combatModifiers []Modifier     // Active offensive modifiers for attack rolls
	attributes      map[string]int // Flexible attribute system (abilities, skills, etc.)
}

// NewActor creates a new Actor with the specified HP, AC, and Initiative.
// HP and AC must be greater than 0. Current HP is initialized to max HP.
func NewActor(hp, ac, initiative int) (*Actor, error) {
	if hp <= 0 {
		return nil, fmt.Errorf("hp must be greater than 0, got %d", hp)
	}
	if ac <= 0 {
		return nil, fmt.Errorf("ac must be greater than 0, got %d", ac)
	}

	return &Actor{
		maxHP:           hp,
		currentHP:       hp,
		ac:              ac,
		initiative:      initiative,
		combatModifiers: []Modifier{},
		attributes:      make(map[string]int),
	}, nil
}

// NewActorWithAttributes creates a new Actor with initial attributes.
// All attribute keys are automatically lowercased for consistency.
// HP and AC must be greater than 0.
func NewActorWithAttributes(hp, ac, initiative int, attrs map[string]int) (*Actor, error) {
	actor, err := NewActor(hp, ac, initiative)
	if err != nil {
		return nil, err
	}
	for key, value := range attrs {
		actor.SetAttribute(key, value)
	}
	return actor, nil
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

// AddCombatModifier adds a modifier to the actor's combat modifiers.
// The modifier's reason is automatically lowercased for consistency.
func (a *Actor) AddCombatModifier(m Modifier) {
	m.Reason = strings.ToLower(m.Reason)
	a.combatModifiers = append(a.combatModifiers, m)
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

// SkillCheck performs a skill check using D&D 5e conventions (1d20 + skill modifier).
// The skill value is looked up from the actor's Attributes map.
// Returns an error if the skill is not found.
//
// Example:
//
//	actor.SetAttribute("athletics", 5)
//	result, _ := actor.SkillCheck("athletics", roller, d20.Advantage)
func (a *Actor) SkillCheck(skill string, roller *Roller, advantage AdvantageType) (RollOutcome, error) {
	skillValue, exists := a.Attribute(skill)
	if !exists {
		return RollOutcome{}, fmt.Errorf("skill %q not found in actor attributes", skill)
	}

	// Roll 1d20 with advantage/disadvantage
	rolls, err := RollWithAdvantage(roller, 1, 20, advantage)
	if err != nil {
		return RollOutcome{}, err
	}

	// Create modifier for the skill
	modifiers := []Modifier{
		{Value: skillValue, Reason: strings.ToLower(skill)},
	}

	// Calculate total
	total := rolls[0] + skillValue

	return NewRollOutcome(1, 20, rolls, modifiers, total), nil
}

// SkillCheckWithDice performs a skill check with custom dice (for other RPG systems).
// This allows flexibility for non-D20 systems like GURPS (3d6), Savage Worlds, etc.
//
// Example:
//
//	actor.SetAttribute("fighting", 12)
//	result, _ := actor.SkillCheckWithDice("fighting", roller, 3, 6, d20.Normal) // GURPS 3d6
func (a *Actor) SkillCheckWithDice(skill string, roller *Roller, rollCount uint, dieFaces uint, advantage AdvantageType) (RollOutcome, error) {
	skillValue, exists := a.Attribute(skill)
	if !exists {
		return RollOutcome{}, fmt.Errorf("skill %q not found in actor attributes", skill)
	}

	// Roll dice with advantage/disadvantage
	rolls, err := RollWithAdvantage(roller, rollCount, dieFaces, advantage)
	if err != nil {
		return RollOutcome{}, err
	}

	// Create modifier for the skill
	modifiers := []Modifier{
		{Value: skillValue, Reason: strings.ToLower(skill)},
	}

	// Calculate total
	total := skillValue
	for _, roll := range rolls {
		total += roll
	}

	return NewRollOutcome(rollCount, dieFaces, rolls, modifiers, total), nil
}

// AttackRoll makes an attack roll using the actor's CombatModifiers.
// Uses D&D 5e conventions (1d20 + combat modifiers).
//
// Example:
//
//	actor.AddCombatModifier(d20.Modifier{Value: 5, Reason: "strength"})
//	actor.AddCombatModifier(d20.Modifier{Value: 3, Reason: "proficiency"})
//	result, _ := actor.AttackRoll(roller, d20.Advantage)
func (a *Actor) AttackRoll(roller *Roller, advantage AdvantageType) (RollOutcome, error) {
	return a.AttackRollWithModifiers(roller, advantage, nil)
}

// AttackRollWithModifiers makes an attack roll with additional situational modifiers.
// Combines the actor's base CombatModifiers with extra situational modifiers.
//
// Example:
//
//	situationalMods := []d20.Modifier{
//	    {Value: 2, Reason: "flanking"},
//	    {Value: -2, Reason: "partial cover"},
//	}
//	result, _ := actor.AttackRollWithModifiers(roller, d20.Normal, situationalMods)
func (a *Actor) AttackRollWithModifiers(roller *Roller, advantage AdvantageType, extraModifiers []Modifier) (RollOutcome, error) {
	// Roll 1d20 with advantage/disadvantage
	rolls, err := RollWithAdvantage(roller, 1, 20, advantage)
	if err != nil {
		return RollOutcome{}, err
	}

	// Combine base combat modifiers with extra modifiers
	allModifiers := make([]Modifier, 0, len(a.combatModifiers)+len(extraModifiers))
	allModifiers = append(allModifiers, a.combatModifiers...)

	// Lowercase extra modifier reasons
	for _, mod := range extraModifiers {
		mod.Reason = strings.ToLower(mod.Reason)
		allModifiers = append(allModifiers, mod)
	}

	// Calculate total
	total := rolls[0]
	for _, mod := range allModifiers {
		total += mod.Value
	}

	return NewRollOutcome(1, 20, rolls, allModifiers, total), nil
}
