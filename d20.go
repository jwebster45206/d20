// Package d20 provides dice rolling functionality for tabletop gaming.
// It supports standard dice notation (e.g., "2d6", "1d20+5") and provides
// a seedable random number generator for reproducible results.
package d20

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

// Roller handles dice rolling with a configurable random source.
type Roller struct {
	rng *rand.Rand
}

// NewRoller creates a new Roller with the given seed.
// Use the same seed to get reproducible results.
func NewRoller(seed int64) *Roller {
	return &Roller{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Roll evaluates a dice expression and returns the result.
// Supported formats:
//   - "NdS" - roll N dice with S sides
//   - "NdS+M" - roll N dice with S sides and add modifier M
//   - "NdS-M" - roll N dice with S sides and subtract modifier M
//
// Examples: "2d6", "1d20+5", "3d8-2"
func (r *Roller) Roll(expr string) (Result, error) {
	return roll(r.rng, expr)
}

var diceRegex = regexp.MustCompile(`^(\d+)d(\d+)(([+-])(\d+))?$`)

// roll performs the actual dice rolling logic.
func roll(rng *rand.Rand, expr string) (Result, error) {
	expr = strings.TrimSpace(strings.ToLower(expr))

	matches := diceRegex.FindStringSubmatch(expr)
	if matches == nil {
		return Result{}, fmt.Errorf("invalid dice expression: %s", expr)
	}

	numDice, err := strconv.Atoi(matches[1])
	if err != nil || numDice < 1 {
		return Result{}, fmt.Errorf("invalid number of dice: %s", matches[1])
	}

	numSides, err := strconv.Atoi(matches[2])
	if err != nil || numSides < 1 {
		return Result{}, fmt.Errorf("invalid number of sides: %s", matches[2])
	}

	modifier := 0
	if matches[3] != "" {
		mod, err := strconv.Atoi(matches[5])
		if err != nil {
			return Result{}, fmt.Errorf("invalid modifier: %s", matches[5])
		}
		if matches[4] == "-" {
			modifier = -mod
		} else {
			modifier = mod
		}
	}

	rolls := make([]int, numDice)
	total := 0
	for i := 0; i < numDice; i++ {
		roll := rng.Intn(numSides) + 1
		rolls[i] = roll
		total += roll
	}
	total += modifier

	return Result{
		Total: total,
		Rolls: rolls,
	}, nil
}

// formatResult creates a string representation of the result.
func formatResult(r Result) string {
	if len(r.Rolls) == 0 {
		return fmt.Sprintf("Total: %d", r.Total)
	}

	rollsStr := make([]string, len(r.Rolls))
	for i, roll := range r.Rolls {
		rollsStr[i] = strconv.Itoa(roll)
	}

	sum := 0
	for _, roll := range r.Rolls {
		sum += roll
	}

	if sum == r.Total {
		return fmt.Sprintf("Rolls: [%s], Total: %d", strings.Join(rollsStr, ", "), r.Total)
	}

	modifier := r.Total - sum
	modStr := ""
	if modifier > 0 {
		modStr = fmt.Sprintf("+%d", modifier)
	} else if modifier < 0 {
		modStr = fmt.Sprintf("%d", modifier)
	}

	return fmt.Sprintf("Rolls: [%s]%s, Total: %d", strings.Join(rollsStr, ", "), modStr, r.Total)
}
