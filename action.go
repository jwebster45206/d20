package d20

// Action is a named roll pair: Attempt to see if it happens, Effect if it does.
// Charges is optional remaining uses (nil means unlimited).
type Action struct {
	ID      string
	Name    string
	Type    string
	Attempt Dice
	Effect  Dice
	Charges *uint
}

// Normalize standardizes ID and type to lower snake_case.
func (a *Action) Normalize() {
	a.ID = normalizeID(a.ID)
	a.Type = normalizeID(a.Type)
}
