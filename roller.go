package d20

import (
	"errors"
	"math/rand"
	"time"
)

// Roller handles dice rolling with a seedable random number generator.
type Roller struct {
	rng *rand.Rand
}

// RollBuilder provides a fluent API for configuring and executing dice rolls.
// Use Dice() to start building a roll, chain configuration methods, then call Roll() to execute.
type RollBuilder struct {
	roller        *Roller
	rollCount     uint
	dieFaces      uint
	modifiers     []Modifier
	advantageType AdvantageType
}

// NewRoller creates a new Roller with the given seed.
// Use the same seed to get reproducible results, or use time.Now().UnixNano()
// for non-deterministic random rolling.
func NewRoller(seed int64) *Roller {
	return &Roller{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// NewRandomRoller is a convenience function that creates a new Roller seeded
// with the current time in nanoseconds.
func NewRandomRoller() *Roller {
	return NewRoller(time.Now().UnixNano())
}

var (
	errRollCountZero = errors.New("rollCount must be greater than 0")
	errDieFacesZero  = errors.New("dieFaces must be greater than 0")
)

// Dice starts building a dice roll with the specified count and faces.
// This is the entry point for the fluent API.
//
// Example:
//
//	result, err := roller.Dice(1, 20).WithModifier("strength", 3).Roll()
func (r *Roller) Dice(rollCount uint, dieFaces uint) *RollBuilder {
	return &RollBuilder{
		roller:        r,
		rollCount:     rollCount,
		dieFaces:      dieFaces,
		modifiers:     []Modifier{},
		advantageType: Normal,
	}
}

// WithModifier adds a single modifier to the roll.
// The modifier name is automatically lowercased for consistency.
//
// Example:
//
//	roller.Dice(1, 20).WithModifier("strength", 3).WithModifier("proficiency", 2).Roll()
func (rb *RollBuilder) WithModifier(name string, value int) *RollBuilder {
	rb.modifiers = append(rb.modifiers, NewModifier(value, name))
	return rb
}

// WithModifiers adds multiple modifiers to the roll at once.
// Accepts a map of name->value pairs. Names are automatically lowercased.
//
// Example:
//
//	mods := map[string]int{"strength": 3, "proficiency": 2}
//	roller.Dice(1, 20).WithModifiers(mods).Roll()
func (rb *RollBuilder) WithModifiers(modifiers map[string]int) *RollBuilder {
	for name, value := range modifiers {
		rb.modifiers = append(rb.modifiers, NewModifier(value, name))
	}
	return rb
}

// WithAdvantage sets the roll to use advantage (roll twice, take higher).
// This is a D&D 5e mechanic.
//
// Example:
//
//	roller.Dice(1, 20).WithAdvantage().Roll()
func (rb *RollBuilder) WithAdvantage() *RollBuilder {
	rb.advantageType = Advantage
	return rb
}

// WithDisadvantage sets the roll to use disadvantage (roll twice, take lower).
// This is a D&D 5e mechanic.
//
// Example:
//
//	roller.Dice(1, 20).WithDisadvantage().Roll()
func (rb *RollBuilder) WithDisadvantage() *RollBuilder {
	rb.advantageType = Disadvantage
	return rb
}

// Normal explicitly sets the roll to normal (no advantage/disadvantage).
// Usually not needed as Normal is the default, but provided for completeness.
//
// Example:
//
//	roller.Dice(1, 20).WithAdvantage().Normal().Roll() // Normal overrides advantage
func (rb *RollBuilder) Normal() *RollBuilder {
	rb.advantageType = Normal
	return rb
}

// Roll executes the configured dice roll and returns the result.
// This is the terminal method that performs the actual roll.
//
// Example:
//
//	result, err := roller.Dice(2, 6).WithModifier("strength", 3).Roll()
func (rb *RollBuilder) Roll() (RollOutcome, error) {
	if rb.rollCount == 0 {
		return RollOutcome{}, errRollCountZero
	}
	if rb.dieFaces == 0 {
		return RollOutcome{}, errDieFacesZero
	}

	// Roll the dice with advantage/disadvantage
	var rolls []int
	var diceTotal int

	switch rb.advantageType {
	case Normal:
		// Roll normally - one roll per die
		rolls = make([]int, rb.rollCount)
		for i := range rb.rollCount {
			rolls[i] = rb.roller.rng.Intn(int(rb.dieFaces)) + 1
			diceTotal += rolls[i]
		}

	case Advantage:
		// Roll twice per die, keep all rolls but use higher values for total
		// For 1d20 with advantage: rolls = [17, 12], used 17
		rolls = make([]int, rb.rollCount*2)
		for i := range rb.rollCount {
			roll1 := rb.roller.rng.Intn(int(rb.dieFaces)) + 1
			roll2 := rb.roller.rng.Intn(int(rb.dieFaces)) + 1
			rolls[i*2] = roll1
			rolls[i*2+1] = roll2
			diceTotal += max(roll1, roll2)
		}

	case Disadvantage:
		// Roll twice per die, keep all rolls but use lower values for total
		// For 1d20 with disadvantage: rolls = [8, 14], used 8
		rolls = make([]int, rb.rollCount*2)
		for i := range rb.rollCount {
			roll1 := rb.roller.rng.Intn(int(rb.dieFaces)) + 1
			roll2 := rb.roller.rng.Intn(int(rb.dieFaces)) + 1
			rolls[i*2] = roll1
			rolls[i*2+1] = roll2
			diceTotal += min(roll1, roll2)
		}
	}

	modifierTotal := 0
	for _, mod := range rb.modifiers {
		modifierTotal += mod.Value
	}

	return NewRollOutcome(rb.rollCount, rb.dieFaces, rolls, rb.modifiers, diceTotal+modifierTotal), nil
}
