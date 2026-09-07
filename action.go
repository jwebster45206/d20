package d20

type Action struct {
	ID      string
	Name    string
	Type    string
	Attempt Dice
	Effect  Dice
	Charges *uint
}

func NewAction(id, name, actionType string, attempt, effect Dice, charges *uint) Action {
	return Action{
		ID:      id,
		Name:    name,
		Type:    actionType,
		Attempt: attempt,
		Effect:  effect,
		Charges: charges,
	}
}

// Normalize standardizes ID and type to lower snake_case.
func (a *Action) Normalize() {
	a.ID = normalizeID(a.ID)
	a.Type = normalizeID(a.Type)
}
