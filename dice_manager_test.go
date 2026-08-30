package d20

import (
	"errors"
	"regexp"
	"testing"
)

func TestDiceManager_Roll(t *testing.T) {
	tests := []struct {
		name      string
		seed      int64
		count     uint
		faces     uint
		mods      map[string]int
		adv       AdvantageType
		hasError  bool
		valueMin  int
		valueMax  int
		diceCount int
		detailRE  string
	}{
		{
			name:      "1d20",
			seed:      42,
			count:     1,
			faces:     20,
			valueMin:  1,
			valueMax:  20,
			diceCount: 1,
			detailRE:  `Rolled 1d20\.\.\. \d+; \*Result: \d+\*`,
		},
		{
			name:      "1d20 with strength mod",
			seed:      42,
			count:     1,
			faces:     20,
			mods:      map[string]int{"strength": 3},
			valueMin:  4,
			valueMax:  23,
			diceCount: 1,
			detailRE:  `Rolled 1d20\.\.\. \d+; \+3 strength; \*Result: \d+\*`,
		},
		{
			name:      "advantage two dice",
			seed:      42,
			count:     1,
			faces:     20,
			adv:       Advantage,
			valueMin:  1,
			valueMax:  20,
			diceCount: 2,
			detailRE:  `Rolled 2d20\.\.\. \d+, \d+; \*Result: \d+\*`,
		},
		{
			name:      "disadvantage two dice",
			seed:      42,
			count:     1,
			faces:     20,
			adv:       Disadvantage,
			valueMin:  1,
			valueMax:  20,
			diceCount: 2,
			detailRE:  `Rolled 2d20\.\.\. \d+, \d+; \*Result: \d+\*`,
		},
		{
			name:      "3d6",
			seed:      42,
			count:     3,
			faces:     6,
			valueMin:  3,
			valueMax:  18,
			diceCount: 3,
			detailRE:  `Rolled 3d6\.\.\.`,
		},
		{
			name:     "zero count errors",
			seed:     42,
			count:    0,
			faces:    20,
			hasError: true,
		},
		{
			name:     "zero faces errors",
			seed:     42,
			count:    1,
			faces:    0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := NewRoller(tt.seed).Dice(tt.count, tt.faces)
			if len(tt.mods) > 0 {
				rb = rb.WithModifiers(tt.mods)
			}
			switch tt.adv {
			case Advantage:
				rb = rb.WithAdvantage()
			case Disadvantage:
				rb = rb.WithDisadvantage()
			}
			out, err := rb.Roll()
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
			for _, d := range out.DiceRolls {
				if d.Faces != tt.faces {
					t.Errorf("die Faces = %d, want %d", d.Faces, tt.faces)
				}
				if d.Result < 1 || d.Result > int(tt.faces) {
					t.Errorf("die Result %d out of range [1, %d]", d.Result, tt.faces)
				}
			}
			if tt.detailRE != "" && !regexp.MustCompile(tt.detailRE).MatchString(out.Detail()) {
				t.Errorf("Detail() = %q, want match %q", out.Detail(), tt.detailRE)
			}
			if tt.adv == Advantage && len(out.DiceRolls) == 2 {
				want := max(out.DiceRolls[0].Result, out.DiceRolls[1].Result)
				modSum := 0
				for _, m := range out.Modifiers {
					modSum += m.Value
				}
				if out.Value != want+modSum {
					t.Errorf("advantage Value = %d, want max(%d,%d)+mods=%d",
						out.Value, out.DiceRolls[0].Result, out.DiceRolls[1].Result, want+modSum)
				}
			}
			if tt.adv == Disadvantage && len(out.DiceRolls) == 2 {
				want := min(out.DiceRolls[0].Result, out.DiceRolls[1].Result)
				if out.Value != want {
					t.Errorf("disadvantage Value = %d, want min=%d", out.Value, want)
				}
			}
		})
	}
}

func TestDiceManager_RollPercentile(t *testing.T) {
	tests := []struct {
		name      string
		seed      int64
		count     uint
		faces     uint
		mods      map[string]int
		hasError  bool
		errIs     error
		valueMin  int
		valueMax  int
		diceCount int
		want00    bool
		detailRE  string
	}{
		{
			name:      "2d10 range",
			seed:      42,
			count:     2,
			faces:     10,
			valueMin:  1,
			valueMax:  100,
			diceCount: 2,
			detailRE:  `Rolled 2d10\.\.\. \d+, \d+; \*Result: \d+\*`,
		},
		{
			name:      "unset assumes 2d10",
			seed:      42,
			count:     0,
			faces:     0,
			valueMin:  1,
			valueMax:  100,
			diceCount: 2,
		},
		{
			name:      "00 maps to 100",
			seed:      60,
			count:     2,
			faces:     10,
			valueMin:  100,
			valueMax:  100,
			diceCount: 2,
			want00:    true,
			detailRE:  `Rolled 2d10\.\.\. 0, 0; \*Result: 100\*`,
		},
		{
			name:      "with modifier",
			seed:      42,
			count:     2,
			faces:     10,
			mods:      map[string]int{"penalty": -5},
			valueMin:  -4,
			valueMax:  95,
			diceCount: 2,
			detailRE:  `-5 penalty`,
		},
		{
			name:     "1d20 rejected",
			seed:     42,
			count:    1,
			faces:    20,
			hasError: true,
			errIs:    errPercentileRequires2d10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := NewRoller(tt.seed).Dice(tt.count, tt.faces)
			if len(tt.mods) > 0 {
				rb = rb.WithModifiers(tt.mods)
			}
			out, err := rb.RollPercentile()
			if tt.hasError {
				if err == nil {
					t.Fatal("expected error")
				}
				if tt.errIs != nil && !errors.Is(err, tt.errIs) {
					t.Errorf("err = %v, want %v", err, tt.errIs)
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
			for i, d := range out.DiceRolls {
				if d.Faces != 10 {
					t.Errorf("DiceRolls[%d].Faces = %d, want 10", i, d.Faces)
				}
				if d.Result < 0 || d.Result > 9 {
					t.Errorf("DiceRolls[%d].Result = %d, want 0–9", i, d.Result)
				}
			}
			if tt.want00 && (out.DiceRolls[0].Result != 0 || out.DiceRolls[1].Result != 0) {
				t.Fatalf("DiceRolls = %v, want both 0", out.DiceRolls)
			}
			if len(tt.mods) > 0 {
				if len(out.Modifiers) != len(tt.mods) {
					t.Errorf("Modifiers len = %d, want %d", len(out.Modifiers), len(tt.mods))
				}
				raw := out.DiceRolls[0].Result*10 + out.DiceRolls[1].Result
				if raw == 0 {
					raw = 100
				}
				modSum := 0
				for _, m := range out.Modifiers {
					modSum += m.Value
				}
				if out.Value != raw+modSum {
					t.Errorf("Value = %d, want raw %d + mods %d", out.Value, raw, modSum)
				}
			}
			if tt.detailRE != "" && !regexp.MustCompile(tt.detailRE).MatchString(out.Detail()) {
				t.Errorf("Detail() = %q, want match %q", out.Detail(), tt.detailRE)
			}
		})
	}
}
