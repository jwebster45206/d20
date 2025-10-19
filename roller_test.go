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

func TestRoller_Roll_DiceNotation(t *testing.T) {
	roller := NewRoller(42)

	tests := []struct {
		name     string
		notation string
		wantErr  bool
		wantMin  int // Expected minimum possible value
		wantMax  int // Expected maximum possible value
	}{
		{"Simple d20", "1d20", false, 1, 20},
		{"Shorthand d20", "d20", false, 1, 20},
		{"Multiple dice", "2d6", false, 2, 12},
		{"With positive modifier", "1d20+3", false, 4, 23},
		{"With negative modifier", "3d8-2", false, 1, 22},
		{"Large dice pool", "10d6+5", false, 15, 65},
		{"d100", "1d100", false, 1, 100},
		{"Invalid - no d", "20", true, 0, 0},
		{"Invalid - no faces", "2d", true, 0, 0},
		{"Invalid - letter faces", "1dabc", true, 0, 0},
		{"Invalid - multiple modifiers", "1d20+3+2", true, 0, 0},
		{"Invalid - empty", "", true, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := roller.Roll(tt.notation)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Roll(%q) expected error, got none", tt.notation)
				}
			} else {
				if err != nil {
					t.Errorf("Roll(%q) unexpected error: %v", tt.notation, err)
				}
				if result.Value < tt.wantMin {
					t.Errorf("Roll(%q) got value %d, below minimum %d", tt.notation, result.Value, tt.wantMin)
				}
				if result.Value > tt.wantMax {
					t.Errorf("Roll(%q) got value %d, above maximum %d", tt.notation, result.Value, tt.wantMax)
				}
			}
		})
	}
}

func TestRoller_Roll_Values(t *testing.T) {
	roller := NewRoller(42)

	t.Run("1d20 produces single die", func(t *testing.T) {
		result, err := roller.Roll("1d20")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.DiceRolls) != 1 {
			t.Errorf("expected 1 die roll, got %d", len(result.DiceRolls))
		}
		if result.DiceRolls[0] < 1 || result.DiceRolls[0] > 20 {
			t.Errorf("die roll %d out of range [1, 20]", result.DiceRolls[0])
		}
	})

	t.Run("d20 shorthand works", func(t *testing.T) {
		result, err := roller.Roll("d20")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.DiceRolls) != 1 {
			t.Errorf("expected 1 die roll, got %d", len(result.DiceRolls))
		}
	})

	t.Run("2d6 produces two dice", func(t *testing.T) {
		result, err := roller.Roll("2d6")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.DiceRolls) != 2 {
			t.Errorf("expected 2 die rolls, got %d", len(result.DiceRolls))
		}
	})

	t.Run("1d20+3 applies modifier", func(t *testing.T) {
		result, err := roller.Roll("1d20+3")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := result.DiceRolls[0] + 3
		if result.Value != expected {
			t.Errorf("expected value %d (dice %d + 3), got %d", expected, result.DiceRolls[0], result.Value)
		}
	})

	t.Run("1d20-2 applies negative modifier", func(t *testing.T) {
		result, err := roller.Roll("1d20-2")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := result.DiceRolls[0] - 2
		if result.Value != expected {
			t.Errorf("expected value %d (dice %d - 2), got %d", expected, result.DiceRolls[0], result.Value)
		}
	})
}
