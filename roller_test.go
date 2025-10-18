package d20

import (
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

func TestRollBuilder_SimpleDiceRoll(t *testing.T) {
	roller := NewRoller(42)
	result, err := roller.Dice(1, 20).Roll()
	if err != nil {
		t.Fatalf("Roll() error: %v", err)
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

func TestRollBuilder_WithModifier(t *testing.T) {
	roller := NewRoller(42)
	result, err := roller.Dice(1, 20).WithModifier("strength", 3).Roll()
	if err != nil {
		t.Fatalf("Roll() error: %v", err)
	}

	// Value should be dice + modifier
	expected := result.DiceRolls[0] + 3
	if result.Value != expected {
		t.Errorf("Value %d doesn't match expected %d", result.Value, expected)
	}
}

func TestRollBuilder_WithAdvantage(t *testing.T) {
	roller := NewRoller(42)
	result, err := roller.Dice(1, 20).WithAdvantage().Roll()
	if err != nil {
		t.Fatalf("Roll() error: %v", err)
	}

	// With advantage, should have 2 dice rolls (both rolls visible)
	if len(result.DiceRolls) != 2 {
		t.Errorf("Expected 2 dice rolls, got %d", len(result.DiceRolls))
	}

	// Value should be the higher roll (no modifiers)
	higherRoll := max(result.DiceRolls[0], result.DiceRolls[1])
	if result.Value != higherRoll {
		t.Errorf("Expected value %d (higher roll), got %d", higherRoll, result.Value)
	}
}

func TestRollBuilder_WithDisadvantage(t *testing.T) {
	roller := NewRoller(42)
	result, err := roller.Dice(1, 20).WithDisadvantage().Roll()
	if err != nil {
		t.Fatalf("Roll() error: %v", err)
	}

	// With disadvantage, should have 2 dice rolls (both rolls visible)
	if len(result.DiceRolls) != 2 {
		t.Errorf("Expected 2 dice rolls, got %d", len(result.DiceRolls))
	}

	// Value should be the lower roll (no modifiers)
	lowerRoll := min(result.DiceRolls[0], result.DiceRolls[1])
	if result.Value != lowerRoll {
		t.Errorf("Expected value %d (lower roll), got %d", lowerRoll, result.Value)
	}
}

func TestRollBuilder_InvalidInput(t *testing.T) {
	roller := NewRoller(42)

	t.Run("Zero roll count", func(t *testing.T) {
		_, err := roller.Dice(0, 20).Roll()
		if err == nil {
			t.Error("Expected error for rollCount = 0")
		}
	})

	t.Run("Zero die faces", func(t *testing.T) {
		_, err := roller.Dice(1, 0).Roll()
		if err == nil {
			t.Error("Expected error for dieFaces = 0")
		}
	})
}
