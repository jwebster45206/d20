package d20

import (
	"fmt"
	"strings"
)

// RollOutcome is the complete result of a dice roll operation.
type RollOutcome struct {
	Value     int    // Final calculated result (dice total + modifiers)
	DiceRolls []int  // Raw values from each die rolled
	Detail    string // Formatted roll description in Bioware style
}

// NewRollOutcome creates a new RollOutcome with formatted detail string.
// The detail string follows Bioware-style formatting:
// "Rolled 2d20... values 16, 12; +3 strength, +2 proficiency; *Result: 33*"
func NewRollOutcome(rollCount uint, dieFaces uint, rolls []int, modifiers []Modifier, finalValue int) RollOutcome {
	return RollOutcome{
		Value:     finalValue,
		DiceRolls: rolls,
		Detail:    formatRollDetail(rollCount, dieFaces, rolls, modifiers, finalValue),
	}
}

// formatRollDetail creates a display-formatted string for a roll result.
func formatRollDetail(rollCount uint, dieFaces uint, rolls []int, modifiers []Modifier, finalValue int) string {
	var parts []string

	// Dice notation (e.g., "Rolled 2d20...")
	parts = append(parts, fmt.Sprintf("Rolled %dd%d...", rollCount, dieFaces))

	// Individual die values
	if len(rolls) > 0 {
		rollStrs := make([]string, len(rolls))
		for i, r := range rolls {
			rollStrs[i] = fmt.Sprintf("%d", r)
		}
		parts = append(parts, fmt.Sprintf("values %s", strings.Join(rollStrs, ", ")))
	}

	// Modifiers
	if len(modifiers) > 0 {
		modStrs := make([]string, len(modifiers))
		for i, mod := range modifiers {
			sign := "+"
			val := mod.Value
			if val < 0 {
				sign = "" // val is already negative, will display as -X
			} else {
				sign = "+"
			}
			modStrs[i] = fmt.Sprintf("%s%d %s", sign, val, strings.ToLower(mod.Reason))
		}
		parts = append(parts, strings.Join(modStrs, ", "))
	}

	// Final result
	parts = append(parts, fmt.Sprintf("*Result: %d*", finalValue))
	return strings.Join(parts, "; ")
}
