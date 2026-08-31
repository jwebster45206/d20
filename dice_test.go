package d20

import (
	"errors"
	"testing"
)

func TestDice(t *testing.T) {
	orig := NewDice(1, 20).WithModifier("strength", 4)
	tests := []struct {
		name string
		got  Dice
		want Dice
	}{
		{"new", NewDice(2, 6), Dice{Count: 2, Faces: 6}},
		{
			"modifiers",
			NewDice(1, 20).WithModifier("Strength", 3).WithModifier("cover", -2),
			Dice{Count: 1, Faces: 20, Modifiers: []Modifier{
				{Reason: "strength", Value: 3},
				{Reason: "cover", Value: -2},
			}},
		},
		{"advantage", NewDice(1, 20).WithAdvantage(), Dice{Count: 1, Faces: 20, Advantage: Advantage}},
		{"disadvantage", NewDice(1, 20).WithDisadvantage(), Dice{Count: 1, Faces: 20, Advantage: Disadvantage}},
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

func TestParseDiceNotation(t *testing.T) {
	tests := []struct {
		notation string
		count    uint
		faces    uint
		mod      int // 0 means none
		err      bool
	}{
		{"1d20", 1, 20, 0, false},
		{"d20", 1, 20, 0, false},
		{"2d6+3", 2, 6, 3, false},
		{"3d8-2", 3, 8, -2, false},
		{"not-dice", 0, 0, 0, true},
		{"1d0", 0, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.notation, func(t *testing.T) {
			d, err := ParseDiceNotation(tt.notation)
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
