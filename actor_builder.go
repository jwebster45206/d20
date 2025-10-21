package d20

import (
	"errors"
	"fmt"
	"strings"
)

// ActorBuilder provides a fluent API for creating Actors.
// Use NewActor() to start building, chain configuration methods, then call Build().
type ActorBuilder struct {
	id              string
	maxHP           int
	ac              int
	initiative      int
	combatModifiers []Modifier
	attributes      map[string]int
	roller          *Roller
	errors          []error
}

// NewActor starts building a new Actor with required ID.
// The ID is automatically normalized to lowercase snake_case.
// Initiative defaults to 0.
//
// Example:
//
//	fighter := d20.NewActor("ironpants").
//	    WithHP(45).
//	    WithAC(16).
//	    WithAttribute("strength", 16).
//	    WithCombatModifier("strength", 3).
//	    Build()
func NewActor(id string) *ActorBuilder {
	return &ActorBuilder{
		id:              normalizeID(id),
		initiative:      0, // Default to 0
		combatModifiers: []Modifier{},
		attributes:      make(map[string]int),
	}
}

// WithRoller adds a Roller to the builder, enabling rolled stat methods.
// Must be called before using WithRolledAttribute() or WithRolledHP().
//
// Example:
//
//	roller := d20.NewRoller()
//	actor := d20.NewActor("fighter", 0, 0).
//	    WithRoller(roller).
//	    WithRolledHP("10d10+30").
//	    WithRolledAttribute("strength", "3d6").
//	    Build()
func (ab *ActorBuilder) WithRoller(roller *Roller) *ActorBuilder {
	ab.roller = roller
	return ab
}

func (ab *ActorBuilder) WithHP(maxHP int) *ActorBuilder {
	ab.maxHP = maxHP
	return ab
}

func (ab *ActorBuilder) WithRolledHP(roll string) *ActorBuilder {
	if ab.roller == nil {
		ab.errors = append(ab.errors, fmt.Errorf("roller not set, call WithRoller() first"))
		return ab
	}
	outcome, err := ab.roller.Roll(roll)
	if err != nil {
		ab.errors = append(ab.errors, fmt.Errorf("failed to roll HP: %w", err))
		return ab
	}
	ab.maxHP = outcome.Value
	return ab
}

func (ab *ActorBuilder) WithAC(ac int) *ActorBuilder {
	ab.ac = ac
	return ab
}

func (ab *ActorBuilder) WithAttribute(key string, value int) *ActorBuilder {
	ab.attributes[strings.ToLower(key)] = value
	return ab
}

func (ab *ActorBuilder) WithAttributes(attrs map[string]int) *ActorBuilder {
	for key, value := range attrs {
		ab.attributes[strings.ToLower(key)] = value
	}
	return ab
}

func (ab *ActorBuilder) WithRolledAttribute(key string, roll string) *ActorBuilder {
	if ab.roller == nil {
		ab.errors = append(ab.errors, fmt.Errorf("roller not set, call WithRoller() first"))
		return ab
	}
	outcome, err := ab.roller.Roll(roll)
	if err != nil {
		ab.errors = append(ab.errors, fmt.Errorf("failed to roll %s for attribute %s: %w", roll, key, err))
		return ab
	}
	ab.attributes[strings.ToLower(key)] = outcome.Value
	return ab
}

func (ab *ActorBuilder) WithRolledAttributes(rolls map[string]string) *ActorBuilder {
	for key, roll := range rolls {
		ab.WithRolledAttribute(key, roll)
	}
	return ab
}

func (ab *ActorBuilder) WithCombatModifier(name string, value int) *ActorBuilder {
	ab.combatModifiers = append(ab.combatModifiers, NewModifier(name, value))
	return ab
}

func (ab *ActorBuilder) WithCombatModifiers(mods map[string]int) *ActorBuilder {
	for name, value := range mods {
		ab.combatModifiers = append(ab.combatModifiers, NewModifier(name, value))
	}
	return ab
}

func (ab *ActorBuilder) Build() (*Actor, error) {
	if ab.maxHP <= 0 {
		ab.errors = append(ab.errors, fmt.Errorf("hp must be greater than 0, got %d", ab.maxHP))
	}

	if len(ab.errors) > 0 {
		return nil, errors.Join(ab.errors...)
	}

	return &Actor{
		id:              ab.id,
		maxHP:           ab.maxHP,
		currentHP:       ab.maxHP,
		ac:              ab.ac,
		initiative:      ab.initiative,
		combatModifiers: ab.combatModifiers,
		attributes:      ab.attributes,
	}, nil
}
