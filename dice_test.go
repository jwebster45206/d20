package d20

import (
	"errors"
	"testing"
)

func TestNewDice(t *testing.T) {
	tests := []struct {
		name      string
		count     uint
		faces     uint
		wantCount uint
		wantFaces uint
		errIs     error
	}{
		{"1d20", 1, 20, 1, 20, nil},
		{"zero count", 0, 20, 0, 0, ErrRollCountZero},
		{"zero faces", 1, 0, 0, 0, ErrDieFacesZero},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDice(tt.count, tt.faces)
			if tt.errIs != nil {
				if !errors.Is(err, tt.errIs) {
					t.Fatalf("err = %v, want %v", err, tt.errIs)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if d.Count != tt.wantCount || d.Faces != tt.wantFaces {
				t.Errorf("got %dd%d, want %dd%d", d.Count, d.Faces, tt.wantCount, tt.wantFaces)
			}
		})
	}
}

func TestDice(t *testing.T) {
	orig := MustNewDice(1, 20).WithModifier("strength", 4)
	tests := []struct {
		name string
		got  Dice
		want Dice
	}{
		{"new", MustNewDice(2, 6), Dice{Count: 2, Faces: 6}},
		{
			"modifiers",
			MustNewDice(1, 20).WithModifier("Strength", 3).WithModifier("cover", -2),
			Dice{Count: 1, Faces: 20, Modifiers: []Modifier{
				{Reason: "strength", Value: 3},
				{Reason: "cover", Value: -2},
			}},
		},
		{"advantage", MustNewDice(1, 20).WithAdvantage(), Dice{Count: 1, Faces: 20, Advantage: Advantage}},
		{"disadvantage", MustNewDice(1, 20).WithDisadvantage(), Dice{Count: 1, Faces: 20, Advantage: Disadvantage}},
		{
			"copy does not mutate original",
			orig,
			Dice{Count: 1, Faces: 20, Modifiers: []Modifier{{Reason: "strength", Value: 4}}},
		},
		{
			"copy has extras",
			orig.WithAdvantage().WithModifier("bless", 1),
			Dice{Count: 1, Faces: 20, Advantage: Advantage, Modifiers: []Modifier{
				{Reason: "strength", Value: 4},
				{Reason: "bless", Value: 1},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got.Count != tt.want.Count || tt.got.Faces != tt.want.Faces || tt.got.Advantage != tt.want.Advantage {
				t.Errorf("got %+v, want %+v", tt.got, tt.want)
			}
			if len(tt.got.Modifiers) != len(tt.want.Modifiers) {
				t.Fatalf("Modifiers = %v, want %v", tt.got.Modifiers, tt.want.Modifiers)
			}
			for i, m := range tt.want.Modifiers {
				if tt.got.Modifiers[i] != m {
					t.Errorf("Modifiers[%d] = %+v, want %+v", i, tt.got.Modifiers[i], m)
				}
			}
		})
	}
}

func TestDiceFromExpr(t *testing.T) {
	tests := []struct {
		name     string
		notation string
		count    uint
		faces    uint
		mod      int // 0 means none
		err      bool
	}{
		{"1d20", "1d20", 1, 20, 0, false},
		{"d20", "d20", 1, 20, 0, false},
		{"2d6+3", "2d6+3", 2, 6, 3, false},
		{"3d8-2", "3d8-2", 3, 8, -2, false},
		{"trim and lowercase", " 2D6+1 ", 2, 6, 1, false},
		{"not-dice", "not-dice", 0, 0, 0, true},
		{"1d0", "1d0", 0, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := DiceFromExpr(tt.notation)
			if tt.err {
				if !errors.Is(err, ErrInvalidDiceNotation) {
					t.Fatalf("err = %v, want ErrInvalidDiceNotation", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if d.Count != tt.count || d.Faces != tt.faces {
				t.Errorf("%dd%d, want %dd%d", d.Count, d.Faces, tt.count, tt.faces)
			}
			if tt.mod == 0 {
				if len(d.Modifiers) != 0 {
					t.Errorf("Modifiers = %v, want none", d.Modifiers)
				}
				return
			}
			if len(d.Modifiers) != 1 || d.Modifiers[0].Reason != "modifier" || d.Modifiers[0].Value != tt.mod {
				t.Errorf("Modifiers = %v, want modifier %+d", d.Modifiers, tt.mod)
			}
		})
	}
}

func TestMustNewDice(t *testing.T) {
	d := MustNewDice(2, 6)
	if d.Count != 2 || d.Faces != 6 {
		t.Errorf("got %dd%d, want 2d6", d.Count, d.Faces)
	}
}

func TestMustNewDice_panics(t *testing.T) {
	tests := []struct {
		name         string
		count, faces uint
		errIs        error
	}{
		{"zero count", 0, 20, ErrRollCountZero},
		{"zero faces", 1, 0, ErrDieFacesZero},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				got := recover()
				if got == nil {
					t.Fatal("expected panic")
				}
				err, ok := got.(error)
				if !ok {
					t.Fatalf("panic %v (%T), want error", got, got)
				}
				if !errors.Is(err, tt.errIs) {
					t.Fatalf("panic %v, want %v", err, tt.errIs)
				}
			}()
			MustNewDice(tt.count, tt.faces)
		})
	}
}

func TestMustDiceFromExpr(t *testing.T) {
	d := MustDiceFromExpr("2d6+3")
	if d.Count != 2 || d.Faces != 6 {
		t.Errorf("got %dd%d, want 2d6", d.Count, d.Faces)
	}
	if len(d.Modifiers) != 1 || d.Modifiers[0].Reason != "modifier" || d.Modifiers[0].Value != 3 {
		t.Errorf("Modifiers = %v, want modifier +3", d.Modifiers)
	}
}

func TestMustDiceFromExpr_panics(t *testing.T) {
	defer func() {
		got := recover()
		if got == nil {
			t.Fatal("expected panic")
		}
		err, ok := got.(error)
		if !ok {
			t.Fatalf("panic %v (%T), want error", got, got)
		}
		if !errors.Is(err, ErrInvalidDiceNotation) {
			t.Fatalf("panic %v, want ErrInvalidDiceNotation", err)
		}
	}()
	MustDiceFromExpr("not-dice")
}
