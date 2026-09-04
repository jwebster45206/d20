package d20

import (
	"errors"
	"math/rand/v2"
	"sync"
	"time"
)

var (
	ErrInvalidDiceNotation = errors.New("invalid dice notation format")
)

// Roller executes dice rolls with a seedable random number generator.
// It is safe to share across goroutines.
type Roller struct {
	mu  sync.Mutex
	rng *rand.Rand
}

// NewRoller creates a Roller with the given seed.
// The same seed yields reproducible results.
func NewRoller(seed int64) *Roller {
	return &Roller{
		rng: rand.New(rand.NewPCG(uint64(seed), 0)),
	}
}

// NewRandomRoller creates a Roller seeded with the current time in nanoseconds.
func NewRandomRoller() *Roller {
	return NewRoller(time.Now().UnixNano())
}

// RollExpr parses notation and rolls it. Shorthand for DiceFromExpr + Roll.
// Accepts "1d20", "2d6+3", "3d8-2", or "d20" (count defaults to 1).
func (r *Roller) RollExpr(notation string) (RollOutcome, error) {
	d, err := DiceFromExpr(notation)
	if err != nil {
		return RollOutcome{}, err
	}
	return r.Roll(d)
}
