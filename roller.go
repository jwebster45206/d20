package d20

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	errRollCountZero          = errors.New("rollCount must be greater than 0")
	errDieFacesZero           = errors.New("dieFaces must be greater than 0")
	ErrInvalidDiceNotation    = errors.New("invalid dice notation format")
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
		rng: rand.New(rand.NewPCG(uint64(seed), 0)),
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
	return r.DiceExpr(notation).Roll()
}

// ParseNotation parses standard dice notation such as "1d20", "d6", "2d6+3", or "3d8-2".
// Count defaults to 1 when omitted ("d20"). A trailing +N or -N becomes a Modifier
// with reason "modifier". The result can be applied to a DiceManager or inspected
// without rolling.
//
// Returns ErrInvalidDiceNotation (via errors.Is) when the string is not valid notation.
func ParseNotation(notation string) (rollCount uint, dieFaces uint, modifiers []Modifier, err error) {
	notation = strings.TrimSpace(strings.ToLower(notation))

	matches := rollNotationFmt.FindStringSubmatch(notation)
	if matches == nil {
		return 0, 0, nil, fmt.Errorf("%w: %s", ErrInvalidDiceNotation, notation)
	}

	rollCount = 1
	if matches[1] != "" {
		count, err := strconv.Atoi(matches[1])
		if err != nil || count <= 0 {
			return 0, 0, nil, fmt.Errorf("%w: invalid roll count", ErrInvalidDiceNotation)
		}
		rollCount = uint(count)
	}

	faces, err := strconv.Atoi(matches[2])
	if err != nil || faces <= 0 {
		return 0, 0, nil, fmt.Errorf("%w: invalid die faces", ErrInvalidDiceNotation)
	}

	modifiers = []Modifier{}
	if matches[3] != "" {
		modValue, err := strconv.Atoi(matches[5])
		if err != nil {
			return 0, 0, nil, fmt.Errorf("%w: invalid modifier value", ErrInvalidDiceNotation)
		}
		if matches[4] == "-" {
			modValue = -modValue
		}
		modifiers = append(modifiers, NewModifier("modifier", modValue))
	}

	return rollCount, uint(faces), modifiers, nil
}
