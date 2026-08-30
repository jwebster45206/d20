package d20

// DiceManager provides a fluent API for configuring and executing dice rolls.
// Use Dice() or DiceExpr() to start building a roll, chain configuration
// methods, then call Roll() to execute.
type DiceManager struct {
	RollCount     uint
	DieFaces      uint
	Modifiers     []Modifier
	AdvantageType AdvantageType

	roller *Roller
	err    error
}

// Error returns the error from the last failed parse or terminal roll.
// It is nil after a successful Roll or RollPercentile, and after construction
// via Dice with no subsequent failure.
func (dm *DiceManager) Error() error {
	return dm.err
}

func (dm *DiceManager) fail(err error) (RollOutcome, error) {
	dm.err = err
	return RollOutcome{}, err
}

func (dm *DiceManager) succeed(outcome RollOutcome) (RollOutcome, error) {
	dm.err = nil
	return outcome, nil
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
		RollCount:     rollCount,
		DieFaces:      dieFaces,
		Modifiers:     []Modifier{},
		AdvantageType: Normal,
	}
}

// DiceExpr starts building a dice roll from standard notation (e.g. "2d20+1", "d6", "3d8-2").
// Invalid notation is stored on the manager and returned from Roll, RollPercentile, and Error.
// Notation bonuses are a Modifier with reason "modifier" and stack with later WithModifier calls.
//
// Example:
//
//	result, err := roller.DiceExpr("2d20+1").WithModifier("bless", 1).Roll()
func (r *Roller) DiceExpr(notation string) *DiceManager {
	dm := &DiceManager{
		roller:        r,
		Modifiers:     []Modifier{},
		AdvantageType: Normal,
	}
	count, faces, mods, err := ParseNotation(notation)
	if err != nil {
		dm.err = err
		return dm
	}
	dm.RollCount = count
	dm.DieFaces = faces
	dm.Modifiers = mods
	return dm
}

// WithModifier adds a single modifier to the roll.
// The modifier name is automatically lowercased for consistency.
//
// Example:
//
//	roller.Dice(1, 20).WithModifier("strength", 3).WithModifier("proficiency", 2).Roll()
func (dm *DiceManager) WithModifier(name string, value int) *DiceManager {
	dm.Modifiers = append(dm.Modifiers, NewModifier(name, value))
	return dm
}

// WithModifiers adds multiple modifiers to the roll at once.
// Accepts a map of name->value pairs. Names are automatically lowercased.
//
// Example:
//
//	mods := map[string]int{"strength": 3, "proficiency": 2}
//	roller.Dice(1, 20).WithModifiers(mods).Roll()
func (dm *DiceManager) WithModifiers(modifiers map[string]int) *DiceManager {
	for name, value := range modifiers {
		dm.Modifiers = append(dm.Modifiers, NewModifier(name, value))
	}
	return dm
}

// WithAdvantage sets the roll to use advantage (roll twice, take higher).
// 5e mechanic. Ignored by RollPercentile.
//
// Example:
//
//	roller.Dice(1, 20).WithAdvantage().Roll()
func (dm *DiceManager) WithAdvantage() *DiceManager {
	dm.AdvantageType = Advantage
	return dm
}

// WithDisadvantage sets the roll to use disadvantage (roll twice, take lower).
// 5e mechanic. Ignored by RollPercentile.
//
// Example:
//
//	roller.Dice(1, 20).WithDisadvantage().Roll()
func (dm *DiceManager) WithDisadvantage() *DiceManager {
	dm.AdvantageType = Disadvantage
	return dm
}

// RollPercentile rolls a Call of Cthulhu-style d100 (tens digit + ones digit, 00 = 100).
// This is distinct from Roll() / Roll("1d100"), which use standard 1–N die faces.
// AdvantageType is ignored; this always rolls tens + ones.
//
// Requires 2d10. If both RollCount and DieFaces are unset (0), 2d10 is assumed.
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
func (dm *DiceManager) RollPercentile() (RollOutcome, error) {
	if err := dm.Error(); err != nil {
		return dm.fail(err)
	}
	if dm.RollCount == 0 && dm.DieFaces == 0 {
		dm.RollCount = 2
		dm.DieFaces = 10
	}
	if dm.RollCount != 2 || dm.DieFaces != 10 {
		return dm.fail(errPercentileRequires2d10)
	}

	dm.roller.mu.Lock()
	defer dm.roller.mu.Unlock()

	tensDigit := dm.roller.rng.IntN(10)
	onesDigit := dm.roller.rng.IntN(10)

	result := tensDigit*10 + onesDigit
	if result == 0 {
		result = 100
	}

	modifierTotal := 0
	for _, mod := range dm.Modifiers {
		modifierTotal += mod.Value
	}

	rolls := []DieRoll{
		{Faces: 10, Result: tensDigit},
		{Faces: 10, Result: onesDigit},
	}
	return dm.succeed(NewRollOutcome(rolls, dm.Modifiers, result+modifierTotal))
}

// Roll executes the configured dice roll and returns the result.
// This is the terminal method that performs the actual roll.
//
// Example:
//
//	result, err  := roller.Dice(2, 6).WithModifier("strength", 3).Roll()
func (dm *DiceManager) Roll() (RollOutcome, error) {
	if err := dm.Error(); err != nil {
		return dm.fail(err)
	}
	if dm.RollCount == 0 {
		return dm.fail(errRollCountZero)
	}
	if dm.DieFaces == 0 {
		return dm.fail(errDieFacesZero)
	}

	dm.roller.mu.Lock()
	defer dm.roller.mu.Unlock()

	var rolls []DieRoll
	var diceTotal int
	faces := dm.DieFaces

	switch dm.AdvantageType {
	case Normal:
		rolls = make([]DieRoll, dm.RollCount)
		for i := range dm.RollCount {
			n := dm.roller.rng.IntN(int(faces)) + 1
			rolls[i] = DieRoll{Faces: faces, Result: n}
			diceTotal += n
		}

	case Advantage:
		// Roll twice per die, keep all rolls but use higher values for total
		rolls = make([]DieRoll, dm.RollCount*2)
		for i := range dm.RollCount {
			roll1 := dm.roller.rng.IntN(int(faces)) + 1
			roll2 := dm.roller.rng.IntN(int(faces)) + 1
			rolls[i*2] = DieRoll{Faces: faces, Result: roll1}
			rolls[i*2+1] = DieRoll{Faces: faces, Result: roll2}
			diceTotal += max(roll1, roll2)
		}

	case Disadvantage:
		rolls = make([]DieRoll, dm.RollCount*2)
		for i := range dm.RollCount {
			roll1 := dm.roller.rng.IntN(int(faces)) + 1
			roll2 := dm.roller.rng.IntN(int(faces)) + 1
			rolls[i*2] = DieRoll{Faces: faces, Result: roll1}
			rolls[i*2+1] = DieRoll{Faces: faces, Result: roll2}
			diceTotal += min(roll1, roll2)
		}
	}

	modifierTotal := 0
	for _, mod := range dm.Modifiers {
		modifierTotal += mod.Value
	}

	return dm.succeed(NewRollOutcome(rolls, dm.Modifiers, diceTotal+modifierTotal))
}
