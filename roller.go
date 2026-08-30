package d20

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	errRollCountZero          = errors.New("rollCount must be greater than 0")
	errDieFacesZero           = errors.New("dieFaces must be greater than 0")
	errInvalidDiceNotation    = errors.New("invalid dice notation format")
	errPercentileRequires2d10 = errors.New("RollPercentile requires 2d10 (or unset count/faces to assume 2d10)")
)

// rollNotationFmt matches patterns like: 1d20, 2d6+3, 3d8-2, d20+5
var rollNotationFmt = regexp.MustCompile(`^(\d*)d(\d+)(([+-])(\d+))?$`)

// Roller executes dice rolls with a seedable random number generator.
type Roller struct {
	mu  sync.Mutex
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

// Roll provides a simple shorthand API for rolling dice using standard dice notation.
// Accepts strings like "1d20", "2d6+3", "3d8-2", or "d20" (assumes 1d20).
// This is a convenience method that doesn't use the fluent API.
//
// Examples:
//   - "1d20" - Roll one 20-sided die
//   - "2d6+3" - Roll two 6-sided dice and add 3
//   - "3d8-2" - Roll three 8-sided dice and subtract 2
//   - "d20" - Roll one 20-sided die (shorthand)
//
// Returns a RollOutcome with the result, or an error if the notation is invalid.
func (r *Roller) Roll(notation string) (RollOutcome, error) {
	notation = strings.TrimSpace(strings.ToLower(notation))

	matches := rollNotationFmt.FindStringSubmatch(notation)
	if matches == nil {
		return RollOutcome{}, fmt.Errorf("%w: %s", errInvalidDiceNotation, notation)
	}

	// Parse roll count (default to 1 if not specified)
	rollCount := uint(1)
	if matches[1] != "" {
		count, err := strconv.Atoi(matches[1])
		if err != nil || count <= 0 {
			return RollOutcome{}, fmt.Errorf("%w: invalid roll count", errInvalidDiceNotation)
		}
		rollCount = uint(count)
	}

	// Parse die faces
	dieFaces, err := strconv.Atoi(matches[2])
	if err != nil || dieFaces <= 0 {
		return RollOutcome{}, fmt.Errorf("%w: invalid die faces", errInvalidDiceNotation)
	}

	// Parse modifier (if present)
	var modifiers []Modifier
	if matches[3] != "" {
		modValue, err := strconv.Atoi(matches[5])
		if err != nil {
			return RollOutcome{}, fmt.Errorf("%w: invalid modifier value", errInvalidDiceNotation)
		}
		if matches[4] == "-" {
			modValue = -modValue
		}
		modifiers = append(modifiers, NewModifier("modifier", modValue))
	}

	// Use the fluent API internally
	builder := r.Dice(rollCount, uint(dieFaces))
	for _, mod := range modifiers {
		builder = builder.WithModifier(mod.Reason, mod.Value)
	}

	return builder.Roll()
}
