package d20

import (
	"errors"
	"regexp"
	"testing"
)

func TestRoller_New(t *testing.T) {
	tests := []struct {
		name     string
		seed     int64
		hasError bool
	}{
		{name: "seeded", seed: 42},
		{name: "zero seed", seed: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRoller(tt.seed)
			if r == nil || r.rng == nil {
				t.Fatal("NewRoller returned incomplete roller")
			}
		})
	}
}

func TestRoller_Roll(t *testing.T) {
	tests := []struct {
		name      string
		seed      int64
		notation  string
		hasError  bool
		valueMin  int
		valueMax  int
		diceCount int
		detailRE  string
	}{
		{
			name:      "1d20",
			seed:      42,
			notation:  "1d20",
			valueMin:  1,
			valueMax:  20,
			diceCount: 1,
			detailRE:  `Rolled 1d20\.\.\.`,
		},
		{
			name:      "d20 shorthand",
			seed:      42,
			notation:  "d20",
			valueMin:  1,
			valueMax:  20,
			diceCount: 1,
		},
		{
			name:      "2d6+3",
			seed:      42,
			notation:  "2d6+3",
			valueMin:  5,
			valueMax:  15,
			diceCount: 2,
			detailRE:  `\+3 modifier`,
		},
		{
			name:      "3d8-2",
			seed:      42,
			notation:  "3d8-2",
			valueMin:  1,
			valueMax:  22,
			diceCount: 3,
			detailRE:  `-2 modifier`,
		},
		{
			name:      "1d100 uniform",
			seed:      42,
			notation:  "1d100",
			valueMin:  1,
			valueMax:  100,
			diceCount: 1,
		},
		{
			name:     "invalid empty",
			seed:     42,
			notation: "",
			hasError: true,
		},
		{
			name:     "invalid garbage",
			seed:     42,
			notation: "not-dice",
			hasError: true,
		},
		{
			name:     "invalid zero faces",
			seed:     42,
			notation: "1d0",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := NewRoller(tt.seed).Roll(tt.notation)
			if tt.hasError {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if out.Value < tt.valueMin || out.Value > tt.valueMax {
				t.Errorf("Value %d not in [%d, %d]", out.Value, tt.valueMin, tt.valueMax)
			}
			if tt.diceCount >= 0 && len(out.DiceRolls) != tt.diceCount {
				t.Errorf("DiceRolls len = %d, want %d", len(out.DiceRolls), tt.diceCount)
			}
			if tt.detailRE != "" && !regexp.MustCompile(tt.detailRE).MatchString(out.Detail()) {
				t.Errorf("Detail() = %q, want match %q", out.Detail(), tt.detailRE)
			}
		})
	}
}

func TestParseNotation(t *testing.T) {
	tests := []struct {
		name     string
		notation string
		count    uint
		faces    uint
		hasMod   bool
		modValue int
		hasError bool
	}{
		{name: "1d20", notation: "1d20", count: 1, faces: 20},
		{name: "d20 shorthand", notation: "d20", count: 1, faces: 20},
		{name: "2d6+3", notation: "2d6+3", count: 2, faces: 6, hasMod: true, modValue: 3},
		{name: "3d8-2", notation: "3d8-2", count: 3, faces: 8, hasMod: true, modValue: -2},
		{name: "trim and lowercase", notation: " 2D6+1 ", count: 2, faces: 6, hasMod: true, modValue: 1},
		{name: "empty", notation: "", hasError: true},
		{name: "garbage", notation: "not-dice", hasError: true},
		{name: "zero faces", notation: "1d0", hasError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, faces, mods, err := ParseNotation(tt.notation)
			if tt.hasError {
				if err == nil {
					t.Fatal("expected error")
				}
				if !errors.Is(err, ErrInvalidDiceNotation) {
					t.Errorf("err = %v, want ErrInvalidDiceNotation", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if count != tt.count {
				t.Errorf("count = %d, want %d", count, tt.count)
			}
			if faces != tt.faces {
				t.Errorf("faces = %d, want %d", faces, tt.faces)
			}
			if tt.hasMod {
				if len(mods) != 1 || mods[0].Reason != "modifier" || mods[0].Value != tt.modValue {
					t.Errorf("modifiers = %v, want [{modifier %d}]", mods, tt.modValue)
				}
			} else if len(mods) != 0 {
				t.Errorf("modifiers = %v, want empty", mods)
			}
		})
	}
}

func TestRoller_Concurrent(t *testing.T) {
	roller := NewRandomRoller()
	const n = 50
	errCh := make(chan error, n)
	for range n {
		go func() {
			_, err := roller.Dice(1, 20).Roll()
			errCh <- err
		}()
	}
	for range n {
		if err := <-errCh; err != nil {
			t.Errorf("Roll: %v", err)
		}
	}
}

func TestAdvantageType_ZeroValue(t *testing.T) {
	var a AdvantageType
	if a != Normal {
		t.Errorf("zero AdvantageType = %v, want Normal", a)
	}
}
