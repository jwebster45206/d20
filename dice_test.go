package d20

import (
	"regexp"
	"testing"
)

func TestRollOutcome_Detail(t *testing.T) {
	tests := []struct {
		name        string
		rolls       []DieRoll
		modifiers   []Modifier
		value       int
		hasError    bool // unused; outcomes have no error path
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
			if len(o.Modifiers) != tt.modCount && tt.modCount > 0 {
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
