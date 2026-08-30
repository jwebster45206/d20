package d20

import "testing"

func TestRollOutcome_Detail(t *testing.T) {
	t.Run("normal with modifiers", func(t *testing.T) {
		o := NewRollOutcome(
			[]DieRoll{{Faces: 20, Result: 6}},
			[]Modifier{NewModifier("strength", 3)},
			9,
		)
		want := "Rolled 1d20... 6; +3 strength; *Result: 9*"
		if got := o.Detail(); got != want {
			t.Errorf("Detail() = %q, want %q", got, want)
		}
		if len(o.Modifiers) != 1 || o.Modifiers[0].Value != 3 {
			t.Errorf("Modifiers = %+v, want strength +3", o.Modifiers)
		}
	})

	t.Run("advantage shows 2d20", func(t *testing.T) {
		o := NewRollOutcome(
			[]DieRoll{{Faces: 20, Result: 6}, {Faces: 20, Result: 8}},
			nil,
			8,
		)
		want := "Rolled 2d20... 6, 8; *Result: 8*"
		if got := o.Detail(); got != want {
			t.Errorf("Detail() = %q, want %q", got, want)
		}
	})

	t.Run("percentile 2d10 and 00 to 100", func(t *testing.T) {
		o := NewRollOutcome(
			[]DieRoll{{Faces: 10, Result: 0}, {Faces: 10, Result: 0}},
			nil,
			100,
		)
		want := "Rolled 2d10... 0, 0; *Result: 100*"
		if got := o.Detail(); got != want {
			t.Errorf("Detail() = %q, want %q", got, want)
		}
	})

	t.Run("mixed faces", func(t *testing.T) {
		o := NewRollOutcome(
			[]DieRoll{{Faces: 20, Result: 10}, {Faces: 6, Result: 4}},
			nil,
			14,
		)
		want := "Rolled 1d20 + 1d6... 10, 4; *Result: 14*"
		if got := o.Detail(); got != want {
			t.Errorf("Detail() = %q, want %q", got, want)
		}
	})
}
