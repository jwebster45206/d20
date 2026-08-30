package d20

// DiceManager provides a fluent API for configuring and executing dice rolls.
// Use Dice() to start building a roll, chain configuration methods, then call Roll() to execute.
type DiceManager struct {
	roller        *Roller
	rollCount     uint
	dieFaces      uint
	modifiers     []Modifier
	advantageType AdvantageType
}

// Dice starts building a dice roll with the specified count and faces.
// This is the entry point for the fluent API.
//
// Example:
//
//	result, err := roller.Dice(1, 20).WithModifier("strength", 3).Roll()
func (r *Roller) Dice(rollCount uint, dieFaces uint) *DiceManager {
	return &DiceManager{
		roller:        r,
		rollCount:     rollCount,
		dieFaces:      dieFaces,
		modifiers:     []Modifier{},
		advantageType: Normal,
	}
}

// WithModifier adds a single modifier to the roll.
// The modifier name is automatically lowercased for consistency.
//
// Example:
//
//	roller.Dice(1, 20).WithModifier("strength", 3).WithModifier("proficiency", 2).Roll()
func (rb *DiceManager) WithModifier(name string, value int) *DiceManager {
	rb.modifiers = append(rb.modifiers, NewModifier(name, value))
	return rb
}

// WithModifiers adds multiple modifiers to the roll at once.
// Accepts a map of name->value pairs. Names are automatically lowercased.
//
// Example:
//
//	mods := map[string]int{"strength": 3, "proficiency": 2}
//	roller.Dice(1, 20).WithModifiers(mods).Roll()
func (rb *DiceManager) WithModifiers(modifiers map[string]int) *DiceManager {
	for name, value := range modifiers {
		rb.modifiers = append(rb.modifiers, NewModifier(name, value))
	}
	return rb
}

// WithAdvantage sets the roll to use advantage (roll twice, take higher).
// 5e mechanic.
//
// Example:
//
//	roller.Dice(1, 20).WithAdvantage().Roll()
func (rb *DiceManager) WithAdvantage() *DiceManager {
	rb.advantageType = Advantage
	return rb
}

// WithDisadvantage sets the roll to use disadvantage (roll twice, take lower).
// 5e mechanic.
//
// Example:
//
//	roller.Dice(1, 20).WithDisadvantage().Roll()
func (rb *DiceManager) WithDisadvantage() *DiceManager {
	rb.advantageType = Disadvantage
	return rb
}

// RollPercentile rolls a Call of Cthulhu-style d100 (tens digit + ones digit, 00 = 100).
// This is distinct from Roll() / Roll("1d100"), which use standard 1–N die faces.
//
// Requires 2d10. If both rollCount and dieFaces are unset (0), 2d10 is assumed.
// Any other configuration returns an error.
//
// DiceRolls is always two d10 digits [tens, ones] (each Result 0–9).
// Value is the percentile result (1–100) plus any configured modifiers.
// Modifiers on the builder are included in the outcome and applied to Value.
//
// Example:
//
//	result, err := roller.Dice(2, 10).RollPercentile()
//	result, err := roller.Dice(0, 0).RollPercentile() // assumes 2d10
func (rb *DiceManager) RollPercentile() (RollOutcome, error) {
	if rb.rollCount == 0 && rb.dieFaces == 0 {
		rb.rollCount = 2
		rb.dieFaces = 10
	}
	if rb.rollCount != 2 || rb.dieFaces != 10 {
		return RollOutcome{}, errPercentileRequires2d10
	}

	rb.roller.mu.Lock()
	defer rb.roller.mu.Unlock()

	tensDigit := rb.roller.rng.Intn(10)
	onesDigit := rb.roller.rng.Intn(10)

	result := tensDigit*10 + onesDigit
	if result == 0 {
		result = 100
	}

	modifierTotal := 0
	for _, mod := range rb.modifiers {
		modifierTotal += mod.Value
	}

	rolls := []DieRoll{
		{Faces: 10, Result: tensDigit},
		{Faces: 10, Result: onesDigit},
	}
	return NewRollOutcome(rolls, rb.modifiers, result+modifierTotal), nil
}

// Roll executes the configured dice roll and returns the result.
// This is the terminal method that performs the actual roll.
//
// Example:
//
//	result, err  := roller.Dice(2, 6).WithModifier("strength", 3).Roll()
func (rb *DiceManager) Roll() (RollOutcome, error) {
	if rb.rollCount == 0 {
		return RollOutcome{}, errRollCountZero
	}
	if rb.dieFaces == 0 {
		return RollOutcome{}, errDieFacesZero
	}

	rb.roller.mu.Lock()
	defer rb.roller.mu.Unlock()

	var rolls []DieRoll
	var diceTotal int
	faces := rb.dieFaces

	switch rb.advantageType {
	case Normal:
		rolls = make([]DieRoll, rb.rollCount)
		for i := range rb.rollCount {
			n := rb.roller.rng.Intn(int(faces)) + 1
			rolls[i] = DieRoll{Faces: faces, Result: n}
			diceTotal += n
		}

	case Advantage:
		// Roll twice per die, keep all rolls but use higher values for total
		rolls = make([]DieRoll, rb.rollCount*2)
		for i := range rb.rollCount {
			roll1 := rb.roller.rng.Intn(int(faces)) + 1
			roll2 := rb.roller.rng.Intn(int(faces)) + 1
			rolls[i*2] = DieRoll{Faces: faces, Result: roll1}
			rolls[i*2+1] = DieRoll{Faces: faces, Result: roll2}
			diceTotal += max(roll1, roll2)
		}

	case Disadvantage:
		rolls = make([]DieRoll, rb.rollCount*2)
		for i := range rb.rollCount {
			roll1 := rb.roller.rng.Intn(int(faces)) + 1
			roll2 := rb.roller.rng.Intn(int(faces)) + 1
			rolls[i*2] = DieRoll{Faces: faces, Result: roll1}
			rolls[i*2+1] = DieRoll{Faces: faces, Result: roll2}
			diceTotal += min(roll1, roll2)
		}
	}

	modifierTotal := 0
	for _, mod := range rb.modifiers {
		modifierTotal += mod.Value
	}

	return NewRollOutcome(rolls, rb.modifiers, diceTotal+modifierTotal), nil
}
