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

// Roll is shorthand for DiceExpr(notation).Roll().
// Accepts "1d20", "2d6+3", "3d8-2", or "d20" (count defaults to 1).
func (r *Roller) Roll(notation string) (RollOutcome, error) {
	return r.DiceExpr(notation).Roll()
}

// Dice starts a fluent roll with the given count and faces.
//
//	result, err := roller.Dice(1, 20).WithModifier("strength", 3).Roll()
func (r *Roller) Dice(rollCount uint, dieFaces uint) *DiceManager {
	return &DiceManager{
		roller:        r,
		RollCount:     rollCount,
		DieFaces:      dieFaces,
		Modifiers:     []Modifier{},
		AdvantageType: Normal,
	}
}

// DiceExpr starts a fluent roll from standard notation (e.g. "2d20+1", "d6", "3d8-2").
// Invalid notation is stored on the manager and returned from Roll, RollPercentile, and Error.
// A trailing +N or -N becomes a Modifier named "modifier" and stacks with later WithModifier calls.
//
//	result, err := roller.DiceExpr("2d20+1").WithModifier("bless", 1).Roll()
func (r *Roller) DiceExpr(notation string) *DiceManager {
	dm := &DiceManager{
		roller:        r,
		Modifiers:     []Modifier{},
		AdvantageType: Normal,
	}
	count, faces, mods, err := ParseNotation(notation)
	if err != nil {
		dm.err = err
		return dm
	}
	dm.RollCount = count
	dm.DieFaces = faces
	dm.Modifiers = mods
	return dm
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
