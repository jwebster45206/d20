package d20

import (
	"errors"
	"testing"
)

func TestRoller_Roll(t *testing.T) {
	tests := []struct {
		name  string
		dice  Dice
		want  int
		nDice int
		errIs error
	}{
		{"1d20", NewDice(1, 20), 18, 1, nil},
		{"modifier", NewDice(1, 20).WithModifier("strength", 3), 21, 1, nil},
		{"advantage", NewDice(1, 20).WithAdvantage(), 20, 2, nil},
		{"disadvantage", NewDice(1, 20).WithDisadvantage(), 18, 2, nil},
		{"zero count", NewDice(0, 20), 0, 0, ErrRollCountZero},
		{"zero faces", NewDice(1, 0), 0, 0, ErrDieFacesZero},
		{"bad advantage", Dice{Count: 1, Faces: 20, Advantage: 99}, 0, 0, ErrInvalidAdvantage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := NewRoller(42).Roll(tt.dice)
			if tt.errIs != nil {
				if !errors.Is(err, tt.errIs) {
					t.Fatalf("err = %v, want %v", err, tt.errIs)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if out.Value != tt.want {
				t.Errorf("Value = %d, want %d", out.Value, tt.want)
			}
			if len(out.DiceRolls) != tt.nDice {
				t.Errorf("DiceRolls len = %d, want %d", len(out.DiceRolls), tt.nDice)
			}
			if tt.dice.Advantage == Advantage && len(out.DiceRolls) == 2 {
				keep := max(out.DiceRolls[0].Result, out.DiceRolls[1].Result)
				if out.Value != keep {
					t.Errorf("advantage Value = %d, want max %d", out.Value, keep)
				}
			}
			if tt.dice.Advantage == Disadvantage && len(out.DiceRolls) == 2 {
				keep := min(out.DiceRolls[0].Result, out.DiceRolls[1].Result)
				if out.Value != keep {
					t.Errorf("disadvantage Value = %d, want min %d", out.Value, keep)
				}
			}
		})
	}
}

func TestRoller_RollPercentile(t *testing.T) {
	tests := []struct {
		name  string
		seed  int64
		dice  Dice
		want  int // 0 means any 1–100
		errIs error
	}{
		{"2d10", 42, NewDice(2, 10), 0, nil},
		{"00 is 100", 169, NewDice(2, 10), 100, nil},
		{"not 2d10", 42, NewDice(1, 20), 0, ErrPercentileRequires2d10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := NewRoller(tt.seed).RollPercentile(tt.dice)
			if tt.errIs != nil {
				if !errors.Is(err, tt.errIs) {
					t.Fatalf("err = %v, want %v", err, tt.errIs)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(out.DiceRolls) != 2 {
				t.Fatalf("DiceRolls len = %d, want 2", len(out.DiceRolls))
			}
			if tt.want != 0 && out.Value != tt.want {
				t.Errorf("Value = %d, want %d", out.Value, tt.want)
			}
			if out.Value < 1 || out.Value > 100 {
				t.Errorf("Value = %d, want 1–100", out.Value)
			}
			for i, d := range out.DiceRolls {
				if d.Faces != 10 || d.Result < 0 || d.Result > 9 {
					t.Errorf("DiceRolls[%d] = %+v, want d10 digit 0–9", i, d)
				}
			}
		})
	}
}

func TestRoller_RollExpr(t *testing.T) {
	tests := []struct {
		notation string
		want     int
		nDice    int
		err      bool
	}{
		{"2d6+3", 15, 2, false},
		{"not-dice", 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.notation, func(t *testing.T) {
			out, err := NewRoller(42).RollExpr(tt.notation)
			if tt.err {
				if !errors.Is(err, ErrInvalidDiceNotation) {
					t.Fatalf("err = %v, want ErrInvalidDiceNotation", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if out.Value != tt.want || len(out.DiceRolls) != tt.nDice {
				t.Errorf("Value=%d n=%d, want %d n=%d", out.Value, len(out.DiceRolls), tt.want, tt.nDice)
			}
		})
	}
}

func TestRoller_Nil(t *testing.T) {
	var r *Roller
	if _, err := r.Roll(NewDice(1, 20)); !errors.Is(err, ErrNilRoller) {
		t.Fatalf("Roll err = %v, want ErrNilRoller", err)
	}
	if _, err := r.RollPercentile(NewDice(2, 10)); !errors.Is(err, ErrNilRoller) {
		t.Fatalf("RollPercentile err = %v, want ErrNilRoller", err)
	}
}

func TestRoller_Concurrent(t *testing.T) {
	roller := NewRandomRoller()
	const n = 50
	errCh := make(chan error, n)
	for range n {
		go func() {
			_, err := roller.Roll(NewDice(1, 20))
			errCh <- err
		}()
	}
	for range n {
		if err := <-errCh; err != nil {
			t.Errorf("Roll: %v", err)
		}
	}
}

func TestRollOutcome_Detail(t *testing.T) {
	tests := []struct {
		name string
		o    RollOutcome
		want string
	}{
		{
			"modifier",
			NewRollOutcome([]DieRoll{{Faces: 20, Result: 6}}, []Modifier{NewModifier("strength", 3)}, 9),
			"Rolled 1d20... 6; +3 strength; *Result: 9*",
		},
		{
			"advantage",
			NewRollOutcome([]DieRoll{{Faces: 20, Result: 6}, {Faces: 20, Result: 8}}, nil, 8),
			"Rolled 2d20... 6, 8; *Result: 8*",
		},
		{
			"percentile 00",
			NewRollOutcome([]DieRoll{{Faces: 10, Result: 0}, {Faces: 10, Result: 0}}, nil, 100),
			"Rolled 2d10... 0, 0; *Result: 100*",
		},
		{
			"penalty",
			NewRollOutcome([]DieRoll{{Faces: 20, Result: 15}}, []Modifier{NewModifier("cover", -2)}, 13),
			"Rolled 1d20... 15; -2 cover; *Result: 13*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Detail(); got != tt.want {
				t.Errorf("Detail() = %q, want %q", got, tt.want)
			}
		})
	}
}
