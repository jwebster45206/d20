package d20

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// Regex to match any non-alphanumeric character for ID normalization
	nonAlphaNumeric = regexp.MustCompile(`[^a-z0-9]+`)

	ErrDuplicateKey  = errors.New("duplicate key after normalization")
	ErrEmptyActionID = errors.New("empty action id")
)

// normalizeID converts a string to lowercase snake_case for consistent IDs and map keys.
// Handles spaces, hyphens, special characters, etc.
//
// Examples:
//   - "Ironpants" -> "ironpants"
//   - "Busta the Black" -> "busta_the_black"
//   - "Fighter-1" -> "fighter_1"
func normalizeID(id string) string {
	id = strings.ToLower(id)
	id = nonAlphaNumeric.ReplaceAllString(id, "_")
	id = strings.Trim(id, "_")
	return id
}

func normalizeIntMap(m map[string]int) (map[string]int, error) {
	if m == nil {
		return nil, nil
	}
	out := make(map[string]int, len(m))
	for k, v := range m {
		nk := normalizeID(k)
		if _, exists := out[nk]; exists {
			return nil, fmt.Errorf("%w: %q", ErrDuplicateKey, nk)
		}
		out[nk] = v
	}
	return out, nil
}

func normalizeActions(actions map[string]Action) (map[string]Action, error) {
	if actions == nil {
		return nil, nil
	}
	out := make(map[string]Action, len(actions))
	for k, action := range actions {
		nk := normalizeID(k)
		if nk == "" {
			return nil, ErrEmptyActionID
		}
		action.Normalize()
		action.ID = nk
		if _, exists := out[nk]; exists {
			return nil, fmt.Errorf("%w: %q", ErrDuplicateKey, nk)
		}
		out[nk] = action
	}
	return out, nil
}

// Actor represents a character, NPC, or monster in the game world.
//
// Example:
//
//	fighter := d20.NewActor("ironpants Son of Arathorn")
//	fighter.MaxHP, fighter.HP = 45, 45
//	fighter.AC = 18
//	fmt.Println(fighter.ID) // "ironpants_son_of_arathorn"
type Actor struct {
	ID         string
	Name       string            // Display name; not normalized
	Alignment  string            // Caller-owned label
	MaxHP      int               // Maximum Hit Points
	HP         int               // Current Hit Points
	AC         int               // Armor Class
	Attributes map[string]int    // Caller-owned numbers (ability scores or skill bonuses)
	Modifiers  map[string]int    // Caller-wired roll bonuses; not derived from Attributes
	Actions    map[string]Action // Attempt/effect roll pairs for quick access (like a combat menu)
}

// NewActor creates an Actor with a normalized ID and initialized maps.
func NewActor(id string) *Actor {
	return &Actor{
		ID:         normalizeID(id),
		Attributes: make(map[string]int),
		Modifiers:  make(map[string]int),
		Actions:    make(map[string]Action),
	}
}

// Normalize rewrites ID, Attributes keys, Modifiers keys, and Action map keys to lowercase snake_case.
// Each action map key is copied onto Action.ID. Type is snake_cased.
// Two keys that collapse to the same name return ErrDuplicateKey and leave the actor unchanged.
// An empty action key after normalization returns ErrEmptyActionID.
// If HP is greater than MaxHP, MaxHP is set to HP.
func (a *Actor) Normalize() error {
	id := normalizeID(a.ID)
	attrs, err := normalizeIntMap(a.Attributes)
	if err != nil {
		return err
	}
	mods, err := normalizeIntMap(a.Modifiers)
	if err != nil {
		return err
	}
	acts, err := normalizeActions(a.Actions)
	if err != nil {
		return err
	}
	a.ID = id
	a.Attributes = attrs
	a.Modifiers = mods
	a.Actions = acts
	if a.HP > a.MaxHP {
		a.MaxHP = a.HP
	}
	return nil
}

// Action returns the action stored under id. The query is normalized for lookup.
// Missing IDs return false. The returned value is a copy.
func (a *Actor) Action(id string) (Action, bool) {
	if a == nil {
		return Action{}, false
	}
	act, ok := a.Actions[normalizeID(id)]
	return act, ok
}

// Dice returns a copy of d with the named modifier keys applied.
// Missing keys are skipped. Key names are normalized for lookup.
// Situational extras go on the returned Dice (WithModifier) without mutating this spec.
//
//	d := actor.Dice(d20.MustDiceFromExpr("1d6"), "damage", "strength")
//	out, err := roller.Roll(d)
func (a *Actor) Dice(d Dice, keys ...string) Dice {
	for _, name := range keys {
		name = normalizeID(name)
		if v, ok := a.Modifiers[name]; ok {
			d = d.WithModifier(name, v)
		}
	}
	return d
}

// D20Dice returns 1d20 with the named modifier keys applied.
// Missing keys are skipped. Key names are normalized for lookup.
//
//	d := actor.D20Dice("strength", "striking")
//	out, err := roller.Roll(d.WithAdvantage())
func (a *Actor) D20Dice(keys ...string) Dice {
	return a.Dice(MustNewDice(1, 20), keys...)
}

// DiceFromExpr parses notation and applies the named modifier keys.
// Missing keys are skipped. Invalid notation returns ErrInvalidDiceNotation.
func (a *Actor) DiceFromExpr(expr string, keys ...string) (Dice, error) {
	d, err := DiceFromExpr(expr)
	if err != nil {
		return Dice{}, err
	}
	return a.Dice(d, keys...), nil
}
