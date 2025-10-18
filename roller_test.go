package d20

import (
	"strings"
	"testing"
)

func TestNewRoller(t *testing.T) {
	roller := NewRoller(42)
	if roller == nil {
		t.Fatal("NewRoller returned nil")
	}
	if roller.rng == nil {
		t.Fatal("Roller.rng is nil")
	}
}

func TestRoll_Deterministic(t *testing.T) {
	// Test that the same seed produces the same results
	seed := int64(12345)

	roller1 := NewRoller(seed)
	result1, err := roller1.Roll(2, 20, nil)
	if err != nil {
		t.Fatalf("Roll failed: %v", err)
	}

	roller2 := NewRoller(seed)
	result2, err := roller2.Roll(2, 20, nil)
	if err != nil {
		t.Fatalf("Roll failed: %v", err)
	}

	if result1.Value != result2.Value {
		t.Errorf("Expected deterministic results, got %d and %d", result1.Value, result2.Value)
	}

	if len(result1.DiceRolls) != len(result2.DiceRolls) {
		t.Errorf("Different number of dice rolls")
	}

	for i := range result1.DiceRolls {
		if result1.DiceRolls[i] != result2.DiceRolls[i] {
			t.Errorf("Dice roll %d differs: %d vs %d", i, result1.DiceRolls[i], result2.DiceRolls[i])
		}
	}
}

func TestRoll_SingleDie(t *testing.T) {
	roller := NewRoller(42)
	result, err := roller.Roll(1, 20, nil)
	if err != nil {
		t.Fatalf("Roll failed: %v", err)
	}

	if len(result.DiceRolls) != 1 {
		t.Errorf("Expected 1 die roll, got %d", len(result.DiceRolls))
	}

	if result.DiceRolls[0] < 1 || result.DiceRolls[0] > 20 {
		t.Errorf("Die roll %d out of range [1, 20]", result.DiceRolls[0])
	}

	if result.Value != result.DiceRolls[0] {
		t.Errorf("Value %d doesn't match die roll %d (no modifiers)", result.Value, result.DiceRolls[0])
	}
}

func TestRoll_MultipleDice(t *testing.T) {
	roller := NewRoller(42)
	result, err := roller.Roll(3, 6, nil)
	if err != nil {
		t.Fatalf("Roll failed: %v", err)
	}

	if len(result.DiceRolls) != 3 {
		t.Errorf("Expected 3 dice rolls, got %d", len(result.DiceRolls))
	}

	sum := 0
	for i, roll := range result.DiceRolls {
		if roll < 1 || roll > 6 {
			t.Errorf("Die roll %d is %d, out of range [1, 6]", i, roll)
		}
		sum += roll
	}

	if result.Value != sum {
		t.Errorf("Value %d doesn't match sum of dice %d (no modifiers)", result.Value, sum)
	}
}

func TestRoll_WithModifiers(t *testing.T) {
	roller := NewRoller(42)
	modifiers := []Modifier{
		{Value: 3, Reason: "Strength"},
		{Value: 2, Reason: "Proficiency"},
	}

	result, err := roller.Roll(2, 20, modifiers)
	if err != nil {
		t.Fatalf("Roll failed: %v", err)
	}

	diceSum := 0
	for _, roll := range result.DiceRolls {
		diceSum += roll
	}

	expectedValue := diceSum + 3 + 2
	if result.Value != expectedValue {
		t.Errorf("Expected value %d, got %d", expectedValue, result.Value)
	}
}

func TestRoll_WithNegativeModifiers(t *testing.T) {
	roller := NewRoller(42)
	modifiers := []Modifier{
		{Value: 5, Reason: "Bonus"},
		{Value: -2, Reason: "Penalty"},
	}

	result, err := roller.Roll(1, 20, modifiers)
	if err != nil {
		t.Fatalf("Roll failed: %v", err)
	}

	expectedValue := result.DiceRolls[0] + 5 - 2
	if result.Value != expectedValue {
		t.Errorf("Expected value %d, got %d", expectedValue, result.Value)
	}
}

func TestRoll_NoModifiers(t *testing.T) {
	roller := NewRoller(42)
	result, err := roller.Roll(1, 20, nil)
	if err != nil {
		t.Fatalf("Roll failed: %v", err)
	}

	if result.Value != result.DiceRolls[0] {
		t.Errorf("Value should equal die roll with no modifiers")
	}

	// Detail should still be generated
	if result.Detail == "" {
		t.Error("Detail string should not be empty")
	}
}

func TestRoll_EmptyModifiers(t *testing.T) {
	roller := NewRoller(42)
	result, err := roller.Roll(1, 20, []Modifier{})
	if err != nil {
		t.Fatalf("Roll failed: %v", err)
	}

	if result.Value != result.DiceRolls[0] {
		t.Errorf("Value should equal die roll with empty modifiers slice")
	}
}

func TestRoll_InvalidInput_ZeroRollCount(t *testing.T) {
	roller := NewRoller(42)
	_, err := roller.Roll(0, 20, nil)
	if err == nil {
		t.Error("Expected error for rollCount of 0")
	}
}

func TestRoll_InvalidInput_ZeroDieFaces(t *testing.T) {
	roller := NewRoller(42)
	_, err := roller.Roll(1, 0, nil)
	if err == nil {
		t.Error("Expected error for dieFaces of 0")
	}
}

func TestRoll_DifferentDice(t *testing.T) {
	testCases := []struct {
		name     string
		dieFaces uint
	}{
		{"d4", 4},
		{"d6", 6},
		{"d8", 8},
		{"d10", 10},
		{"d12", 12},
		{"d20", 20},
		{"d100", 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			roller := NewRoller(42)
			result, err := roller.Roll(1, tc.dieFaces, nil)
			if err != nil {
				t.Fatalf("Roll failed: %v", err)
			}

			if result.DiceRolls[0] < 1 || result.DiceRolls[0] > int(tc.dieFaces) {
				t.Errorf("Die roll %d out of range [1, %d]", result.DiceRolls[0], tc.dieFaces)
			}
		})
	}
}

func TestFormatRollDetail_Basic(t *testing.T) {
	detail := formatRollDetail(1, 20, []int{15}, nil, 15)

	if !strings.Contains(detail, "Rolled 1d20...") {
		t.Error("Detail should contain rolled dice notation")
	}
	if !strings.Contains(detail, "values 15") {
		t.Error("Detail should contain die value")
	}
	if !strings.Contains(detail, "*Result: 15*") {
		t.Error("Detail should contain final result")
	}
}

func TestFormatRollDetail_WithModifiers(t *testing.T) {
	modifiers := []Modifier{
		{Value: 3, Reason: "Strength"},
		{Value: 2, Reason: "Proficiency"},
	}
	detail := formatRollDetail(2, 20, []int{16, 12}, modifiers, 33)

	if !strings.Contains(detail, "Rolled 2d20...") {
		t.Error("Detail should contain rolled dice notation")
	}
	if !strings.Contains(detail, "values 16, 12") {
		t.Error("Detail should contain die values")
	}
	if !strings.Contains(detail, "+3 strength") {
		t.Error("Detail should contain positive modifier")
	}
	if !strings.Contains(detail, "+2 proficiency") {
		t.Error("Detail should contain proficiency modifier")
	}
	if !strings.Contains(detail, "*Result: 33*") {
		t.Error("Detail should contain final result")
	}
}

func TestFormatRollDetail_WithNegativeModifier(t *testing.T) {
	modifiers := []Modifier{
		{Value: -2, Reason: "Penalty"},
	}
	detail := formatRollDetail(1, 20, []int{15}, modifiers, 13)

	if !strings.Contains(detail, "-2 penalty") {
		t.Errorf("Detail should contain negative modifier: %s", detail)
	}
}

func TestRoll_DetailStringGeneration(t *testing.T) {
	roller := NewRoller(42)
	modifiers := []Modifier{
		{Value: 3, Reason: "Strength"},
		{Value: 2, Reason: "Proficiency"},
	}

	result, err := roller.Roll(2, 20, modifiers)
	if err != nil {
		t.Fatalf("Roll failed: %v", err)
	}

	if result.Detail == "" {
		t.Error("Detail string should not be empty")
	}

	// Detail should contain key components
	if !strings.Contains(result.Detail, "Rolled") {
		t.Error("Detail should contain 'Rolled'")
	}
	if !strings.Contains(result.Detail, "Result:") {
		t.Error("Detail should contain 'Result:'")
	}
}
