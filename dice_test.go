package d20

import (
	"errors"
	"regexp"
	"testing"
)

func TestRollOutcome_Detail(t *testing.T) {
	tests := []struct {
		name        string
		rolls       []DieRoll
		modifiers   []Modifier
		value       int
		modCount    int
		detailRE    string
		detailExact string
	}{
		{
			name:        "normal with modifiers",
			rolls:       []DieRoll{{Faces: 20, Result: 6}},
			modifiers:   []Modifier{NewModifier("strength", 3)},
			value:       9,
			modCount:    1,
			detailExact: "Rolled 1d20... 6; +3 strength; *Result: 9*",
		},
		{
			name:        "advantage shows 2d20",
			rolls:       []DieRoll{{Faces: 20, Result: 6}, {Faces: 20, Result: 8}},
			value:       8,
			detailExact: "Rolled 2d20... 6, 8; *Result: 8*",
		},
		{
			name:        "percentile 00 is 100",
			rolls:       []DieRoll{{Faces: 10, Result: 0}, {Faces: 10, Result: 0}},
			value:       100,
			detailExact: "Rolled 2d10... 0, 0; *Result: 100*",
		},
		{
			name:        "mixed faces",
			rolls:       []DieRoll{{Faces: 20, Result: 10}, {Faces: 6, Result: 4}},
			value:       14,
			detailExact: "Rolled 1d20 + 1d6... 10, 4; *Result: 14*",
		},
		{
			name:      "negative modifier in detail",
			rolls:     []DieRoll{{Faces: 20, Result: 15}},
			modifiers: []Modifier{NewModifier("cover", -2)},
			value:     13,
			modCount:  1,
			detailRE:  `Rolled 1d20\.\.\. 15; -2 cover; \*Result: 13\*`,
		},
		{
			name:        "empty rolls",
			rolls:       nil,
			value:       0,
			detailExact: "Rolled 0d0...; *Result: 0*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewRollOutcome(tt.rolls, tt.modifiers, tt.value)
			if tt.modCount > 0 && len(o.Modifiers) != tt.modCount {
				t.Errorf("Modifiers len = %d, want %d", len(o.Modifiers), tt.modCount)
			}
			if tt.modCount > 0 && len(tt.modifiers) > 0 {
				if o.Modifiers[0].Value != tt.modifiers[0].Value {
					t.Errorf("Modifiers[0].Value = %d, want %d", o.Modifiers[0].Value, tt.modifiers[0].Value)
				}
			}
			got := o.Detail()
			if tt.detailExact != "" && got != tt.detailExact {
				t.Errorf("Detail() = %q, want %q", got, tt.detailExact)
			}
			if tt.detailRE != "" && !regexp.MustCompile(tt.detailRE).MatchString(got) {
				t.Errorf("Detail() = %q, want match %q", got, tt.detailRE)
			}
			if o.Value != tt.value {
				t.Errorf("Value = %d, want %d", o.Value, tt.value)
			}
		})
	}
}

func setupDice(seed int64, count, faces uint, expr string) *DiceManager {
	r := NewRoller(seed)
	if expr != "" {
		return r.DiceExpr(expr)
	}
	return r.Dice(count, faces)
}

func TestDiceManager_Roll(t *testing.T) {
	tests := []struct {
		name           string
		seed           int64
		count          uint
		faces          uint
		expr           string
		mods           map[string]int
		adv            AdvantageType
		hasError       bool
		errIs          error
		hasNotationMod bool
		notationMod    int
		valueMin       int
		valueMax       int
		diceCount      int
		detailRE       string
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
			errIs:    errRollCountZero,
		},
		{
			name:     "zero faces errors",
			seed:     42,
			count:    1,
			faces:    0,
			hasError: true,
			errIs:    errDieFacesZero,
		},
		{
			name:           "expr 2d6+3",
			seed:           42,
			expr:           "2d6+3",
			count:          2,
			faces:          6,
			hasNotationMod: true,
			notationMod:    3,
			valueMin:       5,
			valueMax:       15,
			diceCount:      2,
			detailRE:       `\+3 modifier`,
		},
		{
			name:      "expr d20",
			seed:      42,
			expr:      "d20",
			count:     1,
			faces:     20,
			valueMin:  1,
			valueMax:  20,
			diceCount: 1,
			detailRE:  `Rolled 1d20\.\.\.`,
		},
		{
			name:           "expr 2d20+1 with bless",
			seed:           42,
			expr:           "2d20+1",
			count:          2,
			faces:          20,
			mods:           map[string]int{"bless": 1},
			hasNotationMod: true,
			notationMod:    1,
			valueMin:       4,
			valueMax:       42,
			diceCount:      2,
			detailRE:       `\+1 modifier, \+1 bless|\+1 bless, \+1 modifier`,
		},
		{
			name:     "expr not-dice",
			seed:     42,
			expr:     "not-dice",
			hasError: true,
			errIs:    ErrInvalidDiceNotation,
		},
		{
			name:     "expr empty",
			seed:     42,
			expr:     " ",
			hasError: true,
			errIs:    ErrInvalidDiceNotation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm := setupDice(tt.seed, tt.count, tt.faces, tt.expr)
			if dm.RollCount != tt.count {
				t.Errorf("RollCount = %d, want %d", dm.RollCount, tt.count)
			}
			if dm.DieFaces != tt.faces {
				t.Errorf("DieFaces = %d, want %d", dm.DieFaces, tt.faces)
			}
			if len(tt.mods) > 0 {
				dm = dm.WithModifiers(tt.mods)
			}
			if tt.hasNotationMod {
				found := false
				for _, m := range dm.Modifiers {
					if m.Reason == "modifier" && m.Value == tt.notationMod {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Modifiers = %v, want notation modifier %d", dm.Modifiers, tt.notationMod)
				}
			}
			switch tt.adv {
			case Advantage:
				dm = dm.WithAdvantage()
			case Disadvantage:
				dm = dm.WithDisadvantage()
			}
			out, err := dm.Roll()
			if tt.hasError {
				if err == nil {
					t.Fatal("expected error")
				}
				if tt.errIs != nil && !errors.Is(err, tt.errIs) {
					t.Errorf("err = %v, want %v", err, tt.errIs)
				}
				if dm.Error() == nil {
					t.Fatal("Error() = nil, want error")
				}
				if tt.errIs != nil && !errors.Is(dm.Error(), tt.errIs) {
					t.Errorf("Error() = %v, want %v", dm.Error(), tt.errIs)
				}
				return
			}
			if err != nil || dm.Error() != nil {
				t.Fatalf("unexpected error: %v (Error=%v)", err, dm.Error())
			}
			if out.Value < tt.valueMin || out.Value > tt.valueMax {
				t.Errorf("Value %d not in [%d, %d]", out.Value, tt.valueMin, tt.valueMax)
			}
			if len(out.DiceRolls) != tt.diceCount {
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
		expr      string
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
			seed:      169,
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
		{
			name:      "expr 2d10",
			seed:      42,
			expr:      "2d10",
			count:     2,
			faces:     10,
			valueMin:  1,
			valueMax:  100,
			diceCount: 2,
			detailRE:  `Rolled 2d10\.\.\.`,
		},
		{
			name:     "expr 1d20 rejected",
			seed:     42,
			expr:     "1d20",
			count:    1,
			faces:    20,
			hasError: true,
			errIs:    errPercentileRequires2d10,
		},
		{
			name:     "expr not-dice",
			seed:     42,
			expr:     "not-dice",
			hasError: true,
			errIs:    ErrInvalidDiceNotation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm := setupDice(tt.seed, tt.count, tt.faces, tt.expr)
			if dm.RollCount != tt.count {
				t.Errorf("RollCount = %d, want %d", dm.RollCount, tt.count)
			}
			if dm.DieFaces != tt.faces {
				t.Errorf("DieFaces = %d, want %d", dm.DieFaces, tt.faces)
			}
			if len(tt.mods) > 0 {
				dm = dm.WithModifiers(tt.mods)
			}
			out, err := dm.RollPercentile()
			if tt.hasError {
				if err == nil {
					t.Fatal("expected error")
				}
				if tt.errIs != nil && !errors.Is(err, tt.errIs) {
					t.Errorf("err = %v, want %v", err, tt.errIs)
				}
				if dm.Error() == nil {
					t.Fatal("Error() = nil, want error")
				}
				if tt.errIs != nil && !errors.Is(dm.Error(), tt.errIs) {
					t.Errorf("Error() = %v, want %v", dm.Error(), tt.errIs)
				}
				return
			}
			if err != nil || dm.Error() != nil {
				t.Fatalf("unexpected error: %v (Error=%v)", err, dm.Error())
			}
			if dm.RollCount != 2 || dm.DieFaces != 10 {
				t.Errorf("after percentile RollCount, DieFaces = %d, %d, want 2, 10", dm.RollCount, dm.DieFaces)
			}
			if out.Value < tt.valueMin || out.Value > tt.valueMax {
				t.Errorf("Value %d not in [%d, %d]", out.Value, tt.valueMin, tt.valueMax)
			}
			if len(out.DiceRolls) != tt.diceCount {
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
