package d20

import (
	"errors"
	"slices"
)

var (
	ErrNilRoller              = errors.New("nil roller")
	ErrRollCountZero          = errors.New("count must be greater than 0")
	ErrDieFacesZero           = errors.New("faces must be greater than 0")
	ErrInvalidAdvantage       = errors.New("invalid advantage type")
	ErrPercentileRequires2d10 = errors.New("RollPercentile requires 2d10")
)

// Dice is a roll configuration: count, faces, modifiers, and advantage.
// It is data. A Roller executes it.
//
//	out, err := roller.Roll(d20.NewDice(1, 20).WithModifier("strength", 3))
type Dice struct {
	Count     uint
	Faces     uint
	Modifiers []Modifier
	Advantage AdvantageType
}

// NewDice returns a Dice with the given count and faces.
func NewDice(count, faces uint) Dice {
	return Dice{Count: count, Faces: faces}
}

// WithModifier returns a copy with an added modifier. The name is lowercased.
// The original Dice is not mutated.
func (d Dice) WithModifier(name string, value int) Dice {
	d.Modifiers = append(slices.Clone(d.Modifiers), NewModifier(name, value))
	return d
}

// WithModifiers returns a copy with the map entries appended as modifiers.
// Names are lowercased. The original Dice is not mutated.
func (d Dice) WithModifiers(modifiers map[string]int) Dice {
	d.Modifiers = slices.Clone(d.Modifiers)
	for name, value := range modifiers {
		d.Modifiers = append(d.Modifiers, NewModifier(name, value))
	}
	return d
}

// WithAdvantage returns a copy that rolls twice per die and uses the higher of each pair.
// Ignored by RollPercentile. The original Dice is not mutated.
func (d Dice) WithAdvantage() Dice {
	d.Modifiers = slices.Clone(d.Modifiers)
	d.Advantage = Advantage
	return d
}

// WithDisadvantage returns a copy that rolls twice per die and uses the lower of each pair.
// Ignored by RollPercentile. The original Dice is not mutated.
func (d Dice) WithDisadvantage() Dice {
	d.Modifiers = slices.Clone(d.Modifiers)
	d.Advantage = Disadvantage
	return d
}

func modifierTotal(mods []Modifier) int {
	n := 0
	for _, m := range mods {
		n += m.Value
	}
	return n
}

// Roll executes d with this roller: standard 1–N faces, plus modifiers.
// Advantage rolls twice per die and keeps the higher of each pair;
// disadvantage keeps the lower. All rolls appear in DiceRolls.
func (r *Roller) Roll(d Dice) (RollOutcome, error) {
	if r == nil {
		return RollOutcome{}, ErrNilRoller
	}
	if d.Count == 0 {
		return RollOutcome{}, ErrRollCountZero
	}
	if d.Faces == 0 {
		return RollOutcome{}, ErrDieFacesZero
	}
	switch d.Advantage {
	case Normal, Advantage, Disadvantage:
	default:
		return RollOutcome{}, ErrInvalidAdvantage
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	var rolls []DieRoll
	var diceTotal int
	faces := d.Faces

	switch d.Advantage {
	case Normal:
		rolls = make([]DieRoll, d.Count)
		for i := range d.Count {
			n := r.rng.IntN(int(faces)) + 1
			rolls[i] = DieRoll{Faces: faces, Result: n}
			diceTotal += n
		}

	case Advantage:
		rolls = make([]DieRoll, d.Count*2)
		for i := range d.Count {
			roll1 := r.rng.IntN(int(faces)) + 1
			roll2 := r.rng.IntN(int(faces)) + 1
			rolls[i*2] = DieRoll{Faces: faces, Result: roll1}
			rolls[i*2+1] = DieRoll{Faces: faces, Result: roll2}
			diceTotal += max(roll1, roll2)
		}

	case Disadvantage:
		rolls = make([]DieRoll, d.Count*2)
		for i := range d.Count {
			roll1 := r.rng.IntN(int(faces)) + 1
			roll2 := r.rng.IntN(int(faces)) + 1
			rolls[i*2] = DieRoll{Faces: faces, Result: roll1}
			rolls[i*2+1] = DieRoll{Faces: faces, Result: roll2}
			diceTotal += min(roll1, roll2)
		}
	}

	return NewRollOutcome(rolls, d.Modifiers, diceTotal+modifierTotal(d.Modifiers)), nil
}

// RollPercentile executes d as Call of Cthulhu-style d100 (tens + ones, 00 = 100).
// Distinct from Roll(NewDice(1, 100)), which is a uniform 1–100 die.
// Requires Count == 2 and Faces == 10. Advantage is ignored.
//
// DiceRolls is two d10 digits [tens, ones] (each Result 0–9).
// Value is the percentile (1–100) plus modifiers.
func (r *Roller) RollPercentile(d Dice) (RollOutcome, error) {
	if r == nil {
		return RollOutcome{}, ErrNilRoller
	}
	if d.Count != 2 || d.Faces != 10 {
		return RollOutcome{}, ErrPercentileRequires2d10
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	tensDigit := r.rng.IntN(10)
	onesDigit := r.rng.IntN(10)

	result := tensDigit*10 + onesDigit
	if result == 0 {
		result = 100
	}

	rolls := []DieRoll{
		{Faces: 10, Result: tensDigit},
		{Faces: 10, Result: onesDigit},
	}
	return NewRollOutcome(rolls, d.Modifiers, result+modifierTotal(d.Modifiers)), nil
}
