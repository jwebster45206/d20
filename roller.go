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
	ErrInvalidDiceNotation = errors.New("invalid dice notation format")
)

// rollNotationFmt matches patterns like: 1d20, 2d6+3, 3d8-2, d20+5
var rollNotationFmt = regexp.MustCompile(`^(\d*)d(\d+)(([+-])(\d+))?$`)

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

// RollExpr parses notation and rolls it. Shorthand for ParseDiceNotation + Roll.
// Accepts "1d20", "2d6+3", "3d8-2", or "d20" (count defaults to 1).
func (r *Roller) RollExpr(notation string) (RollOutcome, error) {
	d, err := ParseDiceNotation(notation)
	if err != nil {
		return RollOutcome{}, err
	}
	return r.Roll(d)
}

// ParseDiceNotation parses standard dice notation such as "1d20", "d6", "2d6+3", or "3d8-2".
// Count defaults to 1 when omitted ("d20"). A trailing +N or -N becomes a Modifier
// named "modifier". Inspect the returned Dice without rolling, or pass it to Roller.Roll.
//
// Returns ErrInvalidDiceNotation (via errors.Is) when the string is not valid notation.
func ParseDiceNotation(notation string) (Dice, error) {
	notation = strings.TrimSpace(strings.ToLower(notation))

	matches := rollNotationFmt.FindStringSubmatch(notation)
	if matches == nil {
		return Dice{}, fmt.Errorf("%w: %s", ErrInvalidDiceNotation, notation)
	}

	count := 1
	if matches[1] != "" {
		n, err := strconv.Atoi(matches[1])
		if err != nil || n <= 0 {
			return Dice{}, fmt.Errorf("%w: invalid roll count", ErrInvalidDiceNotation)
		}
		count = n
	}

	faces, err := strconv.Atoi(matches[2])
	if err != nil || faces <= 0 {
		return Dice{}, fmt.Errorf("%w: invalid die faces", ErrInvalidDiceNotation)
	}

	d := NewDice(uint(count), uint(faces))
	if matches[3] != "" {
		modValue, err := strconv.Atoi(matches[5])
		if err != nil {
			return Dice{}, fmt.Errorf("%w: invalid modifier value", ErrInvalidDiceNotation)
		}
		if matches[4] == "-" {
			modValue = -modValue
		}
		d = d.WithModifier("modifier", modValue)
	}
	return d, nil
}
