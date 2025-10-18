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

// Roll performs a dice roll with the specified parameters.
// Returns a RollOutcome containing the final result, individual dice rolls,
// and a formatted detail string in Bioware style.
//
// Example: Roll(2, 20, []Modifier{{Value: 3, Reason: "strength"}, {Value: 2, Reason: "proficiency"}})
func (r *Roller) Roll(rollCount uint, dieFaces uint, modifiers []Modifier) (RollOutcome, error) {
	// Validate input
	if rollCount == 0 {
		return RollOutcome{}, errRollCountZero
	}
	if dieFaces == 0 {
		return RollOutcome{}, errDieFacesZero
	}

	// Roll the dice
	rolls := make([]int, rollCount)
	diceTotal := 0
	for i := range rollCount {
		roll := r.rng.Intn(int(dieFaces)) + 1
		rolls[i] = roll
		diceTotal += roll
	}

	// Apply modifiers
	modifierTotal := 0
	for _, mod := range modifiers {
		modifierTotal += mod.Value
	}

	return NewRollOutcome(rollCount, dieFaces, rolls, modifiers, diceTotal+modifierTotal), nil
}
