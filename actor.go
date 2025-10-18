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

// ActorBuilder provides a fluent API for creating Actors.
// Use NewActor() to start building, chain configuration methods, then call Build().
type ActorBuilder struct {
	id              string
	maxHP           int
	ac              int
	initiative      int
	combatModifiers []Modifier
	attributes      map[string]int
}

// NewActor starts building a new Actor with required ID, HP, and AC.
// The ID is automatically normalized to lowercase snake_case.
// Initiative defaults to 0 if not set via WithInitiative().
//
// Example:
//
//	fighter := d20.NewActor("ironpants", 45, 18).
//	    WithAttribute("strength", 16).
//	    WithCombatModifier("strength", 3).
//	    Build()
func NewActor(id string, hp, ac int) *ActorBuilder {
	return &ActorBuilder{
		id:              normalizeID(id),
		maxHP:           hp,
		ac:              ac,
		initiative:      0, // Default to 0
		combatModifiers: []Modifier{},
		attributes:      make(map[string]int),
	}
}

// WithInitiative sets a custom initiative value for the actor.
// Optional - defaults to 0 if not called.
//
// Example:
//
//	rogue := d20.NewActor("tom", 30, 15).WithInitiative(4).Build()
func (ab *ActorBuilder) WithInitiative(initiative int) *ActorBuilder {
	ab.initiative = initiative
	return ab
}

// WithAttribute adds a single attribute to the actor.
// The key is automatically lowercased for consistency.
//
// Example:
//
//	actor := d20.NewActor("fighter", 45, 18).WithAttribute("strength", 16).Build()
func (ab *ActorBuilder) WithAttribute(key string, value int) *ActorBuilder {
	ab.attributes[strings.ToLower(key)] = value
	return ab
}

// WithAttributes adds multiple attributes at once from a map.
// All keys are automatically lowercased for consistency.
//
// Example:
//
//	attrs := map[string]int{"strength": 16, "dexterity": 14}
//	actor := d20.NewActor("fighter", 45, 18).WithAttributes(attrs).Build()
func (ab *ActorBuilder) WithAttributes(attrs map[string]int) *ActorBuilder {
	for key, value := range attrs {
		ab.attributes[strings.ToLower(key)] = value
	}
	return ab
}

// WithCombatModifier adds a single combat modifier to the actor.
// The modifier name is automatically lowercased for consistency.
// Creates the Modifier struct internally.
//
// Example:
//
//	actor := d20.NewActor("fighter", 45, 18).WithCombatModifier("strength", 3).Build()
func (ab *ActorBuilder) WithCombatModifier(name string, value int) *ActorBuilder {
	ab.combatModifiers = append(ab.combatModifiers, NewModifier(value, name))
	return ab
}

// WithCombatModifiers adds multiple combat modifiers at once from a map.
// All modifier names are automatically lowercased for consistency.
// Creates Modifier structs internally.
//
// Example:
//
//	mods := map[string]int{"strength": 3, "proficiency": 2}
//	actor := d20.NewActor("fighter", 45, 18).WithCombatModifiers(mods).Build()
func (ab *ActorBuilder) WithCombatModifiers(mods map[string]int) *ActorBuilder {
	for name, value := range mods {
		ab.combatModifiers = append(ab.combatModifiers, NewModifier(value, name))
	}
	return ab
}

// Build creates and returns the final Actor instance.
// Validates that HP and AC are greater than 0.
//
// Example:
//
//	actor := d20.NewActor("ironpants", 45, 18).
//	    WithAttribute("strength", 16).
//	    Build()
func (ab *ActorBuilder) Build() (*Actor, error) {
	// Validate
	if ab.maxHP <= 0 {
		return nil, fmt.Errorf("hp must be greater than 0, got %d", ab.maxHP)
	}
	if ab.ac <= 0 {
		return nil, fmt.Errorf("ac must be greater than 0, got %d", ab.ac)
	}

	return &Actor{
		id:              ab.id,
		maxHP:           ab.maxHP,
		currentHP:       ab.maxHP, // Start at full HP
		ac:              ab.ac,
		initiative:      ab.initiative,
		combatModifiers: ab.combatModifiers,
		attributes:      ab.attributes,
	}, nil
}

// ID returns the actor's normalized ID (lowercase snake_case).
// IDs are read-only after creation.
//
// Example:
//
//	fighter := d20.NewActor("ironpants Son of Arathorn", 45, 18).Build()
//	fmt.Println(fighter.ID()) // "aragorn_son_of_arathorn"
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
	a.combatModifiers = append(a.combatModifiers, NewModifier(value, name))
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

// SkillCheck creates a RollBuilder for a skill check using D&D 5e conventions (1d20 + skill modifier).
// The skill value is looked up from the actor's Attributes map.
// Returns a RollBuilder pre-configured with the skill modifier. Chain .WithAdvantage() or other
// modifiers as needed, then call .Roll() to execute.
//
// Returns an error if the skill is not found.
//
// Example:
//
//	actor.SetAttribute("athletics", 5)
//	result, _ := actor.SkillCheck("athletics", roller).WithAdvantage().Roll()
func (a *Actor) SkillCheck(skill string, roller *Roller) (*RollBuilder, error) {
	skillValue, exists := a.Attribute(skill)
	if !exists {
		return nil, fmt.Errorf("skill %q not found in actor attributes", skill)
	}

	// Return a RollBuilder with skill modifier pre-loaded
	return roller.Dice(1, 20).WithModifier(skill, skillValue), nil
}

// AttackRoll creates a RollBuilder for an attack roll using the actor's CombatModifiers.
// Uses D&D 5e conventions (1d20 + combat modifiers).
// Returns a RollBuilder pre-configured with all combat modifiers. Chain .WithAdvantage(),
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
func (a *Actor) AttackRoll(roller *Roller) *RollBuilder {
	builder := roller.Dice(1, 20)

	// Add all combat modifiers
	for _, mod := range a.combatModifiers {
		builder = builder.WithModifier(mod.Reason, mod.Value)
	}

	return builder
}

// D100SkillCheck performs a percentile skill check for d100 systems like Call of Cthulhu.
// Returns success (rolled <= skill value), the roll outcome, and any error.
//
// The bonus parameter implements Call of Cthulhu's bonus/penalty die mechanic:
//   - bonus > 0: Roll multiple d10s for tens digit, take the LOWEST (better chance)
//   - bonus < 0: Roll multiple d10s for tens digit, take the HIGHEST (worse chance)
//   - bonus = 0: Normal d100 roll (1d10 for tens, 1d10 for ones)
//
// Example:
//
//	actor.SetAttribute("stealth", 45)  // 45% skill
//	success, roll, _ := actor.D100SkillCheck("stealth", roller, 0)  // Normal roll
//	bonusSuccess, roll, _ := actor.D100SkillCheck("stealth", roller, 1)  // Bonus die
//	penaltySuccess, roll, _ := actor.D100SkillCheck("stealth", roller, -1) // Penalty die
func (a *Actor) D100SkillCheck(skill string, roller *Roller, bonus int) (bool, RollOutcome, error) {
	skillValue, exists := a.Attribute(skill)
	if !exists {
		return false, RollOutcome{}, fmt.Errorf("skill %q not found in actor attributes", skill)
	}

	var tensDigit, onesDigit int

	if bonus == 0 {
		// Normal d100: 1d10 for tens, 1d10 for ones
		tensDigit = (roller.rng.Intn(10)) * 10
		onesDigit = roller.rng.Intn(10)
	} else if bonus > 0 {
		// Bonus die: Roll (1 + bonus) d10s for tens, take LOWEST
		rolls := bonus + 1
		bestTens := 9 // Start with worst
		for i := 0; i < rolls; i++ {
			roll := roller.rng.Intn(10)
			if roll < bestTens {
				bestTens = roll
			}
		}
		tensDigit = bestTens * 10
		onesDigit = roller.rng.Intn(10)
	} else { // bonus < 0
		// Penalty die: Roll (1 + |bonus|) d10s for tens, take HIGHEST
		rolls := -bonus + 1
		worstTens := 0 // Start with best
		for i := 0; i < rolls; i++ {
			roll := roller.rng.Intn(10)
			if roll > worstTens {
				worstTens = roll
			}
		}
		tensDigit = worstTens * 10
		onesDigit = roller.rng.Intn(10)
	}

	// Calculate final d100 result (00 = 100)
	result := tensDigit + onesDigit
	if result == 0 {
		result = 100
	}

	// Success if rolled <= skill value
	success := result <= skillValue

	// Create roll outcome with skill modifier shown
	rolls := []int{result}
	modifiers := []Modifier{
		{Value: skillValue, Reason: fmt.Sprintf("%s (target)", strings.ToLower(skill))},
	}

	outcome := NewRollOutcome(1, 100, rolls, modifiers, result)

	return success, outcome, nil
}
