package d20

import (
	"fmt"
	"strings"
)

// DieRoll is one physical die result: its face count and the number showing.
type DieRoll struct {
	Faces  uint // Die size (e.g. 20 for a d20, 10 for a percentile digit)
	Result int  // Face showing (1–Faces for standard rolls; 0–9 for percentile digits)
}

// RollOutcome is the complete result of a dice roll operation.
// It contains the final result, raw values, and modifiers used in the calculation.
type RollOutcome struct {
	Value     int        // Final calculated result (dice total + modifiers)
	DiceRolls []DieRoll  // Each die rolled, with faces and result
	Modifiers []Modifier // Modifiers applied to the roll
}

// NewRollOutcome creates a new RollOutcome from die results and modifiers.
func NewRollOutcome(rolls []DieRoll, modifiers []Modifier, finalValue int) RollOutcome {
	rollsCopy := make([]DieRoll, len(rolls))
	copy(rollsCopy, rolls)
	modsCopy := make([]Modifier, len(modifiers))
	copy(modsCopy, modifiers)
	return RollOutcome{
		Value:     finalValue,
		DiceRolls: rollsCopy,
		Modifiers: modsCopy,
	}
}

// Detail returns a Bioware-style formatted description of the roll, for example:
// "Rolled 2d20... 16, 12; +3 strength, +2 proficiency; *Result: 33*"
func (o RollOutcome) Detail() string {
	result := "Rolled " + diceNotation(o.DiceRolls) + "..."

	if len(o.DiceRolls) > 0 {
		rollStrs := make([]string, len(o.DiceRolls))
		for i, r := range o.DiceRolls {
			rollStrs[i] = fmt.Sprintf("%d", r.Result)
		}
		result += " " + strings.Join(rollStrs, ", ")
	}

	if len(o.Modifiers) > 0 {
		modStrs := make([]string, len(o.Modifiers))
		for i, mod := range o.Modifiers {
			sign := "+"
			val := mod.Value
			if val < 0 {
				sign = "" // val is already negative, will display as -X
			}
			modStrs[i] = fmt.Sprintf("%s%d %s", sign, val, strings.ToLower(mod.Reason))
		}
		result += "; " + strings.Join(modStrs, ", ")
	}

	result += "; *Result: " + fmt.Sprintf("%d*", o.Value)
	return result
}

// diceNotation derives NdF (or mixed 1dA + 1dB) from the dice actually rolled.
func diceNotation(rolls []DieRoll) string {
	if len(rolls) == 0 {
		return "0d0"
	}
	faces := rolls[0].Faces
	same := true
	for _, r := range rolls[1:] {
		if r.Faces != faces {
			same = false
			break
		}
	}
	if same {
		return fmt.Sprintf("%dd%d", len(rolls), faces)
	}
	parts := make([]string, len(rolls))
	for i, r := range rolls {
		parts[i] = fmt.Sprintf("1d%d", r.Faces)
	}
	return strings.Join(parts, " + ")
}
