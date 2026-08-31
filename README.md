# D20 - Go Dice Rolling Library

[![CI](https://github.com/jwebster45206/d20/actions/workflows/ci.yml/badge.svg)](https://github.com/jwebster45206/d20/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jwebster45206/d20.svg)](https://pkg.go.dev/github.com/jwebster45206/d20)

A Go library for dice rolling and tabletop RPG helpers: notation and fluent d20-style rolls, named modifiers, advantage/disadvantage, and a simple actor model. It is not a rules engine and does not implement SRD procedures such as ability modifiers or proficiency.

## Features

- **Dice notation**: `"1d20+3"`, `"2d6"`, `"d20"`
- **Dice config**: named modifiers and advantage/disadvantage
- **Actors**: HP/AC plus caller-owned attributes and modifiers
- **RollOutcome.Detail()**: Bioware-style breakdown of dice and modifiers

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/jwebster45206/d20"
)

func main() {
    roller := d20.NewRandomRoller()

    result, _ := roller.RollExpr("1d20+3")
    fmt.Printf("Attack roll: %d\n", result.Value)

    result, err := roller.Roll(d20.NewDice(1, 20).
        WithModifier("strength", 3).
        WithModifier("proficiency", 2))
    if err != nil {
        panic(err)
    }
    fmt.Println(result.Detail())
    // Example: "Rolled 1d20... 15; +3 strength, +2 proficiency; *Result: 20*"

    r, err := roller.Roll(d20.NewDice(1, 20).
        WithAdvantage().
        WithModifier("dexterity", 4))
    if err != nil {
        panic(err)
    }
    fmt.Printf("Roll: %d (dice: %v)\n", r.Value, r.DiceRolls)
}
```

## Core Types

### Roller

The random number generator (RNG). Seed it for tests; `NewRandomRoller` for play. Safe to share across goroutines. `Roll` executes a `Dice` spec.

```go
func NewRoller(seed int64) *Roller
func NewRandomRoller() *Roller

func (r *Roller) Roll(d Dice) (RollOutcome, error)
func (r *Roller) RollPercentile(d Dice) (RollOutcome, error)
func (r *Roller) RollExpr(notation string) (RollOutcome, error) // ParseDiceNotation + Roll
```

### Dice

Roll configuration. Value type: `With*` returns a copy and does not mutate the original.

```go
type Dice struct {
    Count     uint
    Faces     uint
    Modifiers []Modifier
    Advantage AdvantageType
}

func NewDice(count, faces uint) Dice

func (d Dice) WithModifier(name string, value int) Dice
func (d Dice) WithModifiers(modifiers map[string]int) Dice
func (d Dice) WithAdvantage() Dice
func (d Dice) WithDisadvantage() Dice
```

```go
attack := d20.NewDice(1, 20).WithModifier("strength", 4)
roller.Roll(attack)
roller.Roll(attack.WithAdvantage()) // attack unchanged
```

Advantage (ignored by `RollPercentile`): rolls twice per die, uses the higher of each pair; disadvantage uses the lower. All rolls appear in `DiceRolls`.

`RollPercentile` is tens+ones (`00` = 100), distinct from `Roll(NewDice(1, 100))`. Requires 2d10.

### Dice notation

```go
func ParseDiceNotation(notation string) (Dice, error)

var ErrInvalidDiceNotation error
```

Accepted: `"1d20"`, `"d20"`, `"2d6+3"`, `"3d8-2"`. `"1d100"` is a uniform 1–100 die. A trailing `+N`/`-N` becomes a modifier named `"modifier"`. Invalid notation fails immediately.

```go
d, err := d20.ParseDiceNotation("2d6+3")
out, err := roller.Roll(d.WithModifier("bless", 1))
out, err = roller.RollExpr("2d6+3")
```

### RollOutcome

```go
type RollOutcome struct {
    Value     int        // dice total + modifiers
    DiceRolls []DieRoll  // each die: Faces + Result
    Modifiers []Modifier
}

type DieRoll struct {
    Faces  uint
    Result int
}

func (o RollOutcome) Detail() string
```

- Normal: `DiceRolls: [{20, 17}]`, `Detail(): "Rolled 1d20... 17; *Result: 17*"`
- Advantage: `DiceRolls: [{20, 6}, {20, 8}]`, `Value: 8`
- Percentile `00`: `DiceRolls: [{10, 0}, {10, 0}]`, `Detail(): "Rolled 2d10... 0, 0; *Result: 100*"`

## Actor System

Map keys are lowercase by convention. `SkillDice`, `StrikeDice`, and `D100SkillDice` lowercase the *query*; they do not rewrite stored keys. They return a `Dice` spec; the roller executes it.

```go
type Actor struct {
    ID         string
    MaxHP      int
    HP         int
    AC         int
    Initiative int
    Attributes map[string]int // caller-owned numbers (scores or skill bonuses)
    Modifiers  map[string]int // caller-wired roll bonuses
}

func NewActor(id string) *Actor

func (a *Actor) AddHP(amount int)  // no-op if amount <= 0; will not exceed MaxHP
func (a *Actor) SubHP(amount int)  // no-op if amount <= 0; will not go below 0
func (a *Actor) ResetHP()
func (a *Actor) IsKnockedOut() bool // HP <= 0

func (a *Actor) SkillDice(skill string) (Dice, error)
func (a *Actor) StrikeDice(keys ...string) Dice
func (a *Actor) D100SkillDice(skill string) (Dice, error)
```

Well-known names live in `github.com/jwebster45206/d20/vocab`. Custom string keys are allowed.

### Creating actors

```go
fighter := d20.NewActor("Ironpants")
fighter.MaxHP, fighter.HP = 45, 45
fighter.AC = 18

wizard := d20.NewActor("Merlin")
wizard.MaxHP, wizard.HP = 38, 38
wizard.AC = 14
wizard.Attributes = map[string]int{
    "intelligence": 18, // ability score; SkillDice would add +18 if used as a skill
    "athletics":    5,  // skill bonus used by SkillDice
}
wizard.Modifiers = map[string]int{
    "intelligence": 4,
    "striking":     4,
}

roller := d20.NewRandomRoller()
hp, _ := roller.RollExpr("10d12+30")
str, _ := roller.RollExpr("3d6")
barbarian := d20.NewActor("Grog")
barbarian.MaxHP, barbarian.HP = hp.Value, hp.Value
barbarian.AC = 14
barbarian.Attributes["strength"] = str.Value
```

`SkillDice` adds the stored attribute value to 1d20. The library does not derive 5e ability modifiers; store the bonus you want (`athletics: 5`, not `strength: 16`) unless you intend the raw score.

### Modifiers and rolls

`Modifiers` are not derived from `Attributes`. `StrikeDice` applies only the keys you pass; missing keys are skipped. Situational extras go on the returned `Dice`. `RollPercentile` ignores advantage/disadvantage.

```go
import "github.com/jwebster45206/d20/vocab"

actor.Modifiers[vocab.Strength] = 4
actor.Modifiers[vocab.Striking] = 3
actor.Modifiers[vocab.Damage] = 1

result, _ := roller.Roll(actor.StrikeDice(vocab.Strength, vocab.Striking))
result, _ = roller.Roll(actor.StrikeDice(vocab.Strength, vocab.Striking).WithModifier("bless", 1))

d, err := actor.SkillDice("athletics")
if err != nil {
    // skill not found
}
result, _ = roller.Roll(d.WithAdvantage())

d100, err := actor.D100SkillDice("stealth")
outcome, _ := roller.RollPercentile(d100)
success := outcome.Value <= actor.Attributes["stealth"]
```

`vocab.Damage` is for damage rolls you build yourself (`NewDice(1, 8).WithModifier(vocab.Damage, actor.Modifiers[vocab.Damage])`). `StrikeDice` never auto-includes it.

## Game systems

Helpers cover common 5e-style rolls (polyhedral dice, d20 + named modifiers, advantage/disadvantage, HP/AC). Callers supply ability modifiers and proficiency. The [5e SRD](https://dnd.wizards.com/resources/systems-reference-document) is published by Wizards of the Coast under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/). This repository does not include SRD text.

`RollPercentile(NewDice(2, 10))` and `D100SkillDice` implement compatible d100 roll-under mechanics (tens + ones, `00` = 100). They are not a Call of Cthulhu rules implementation. Call of Cthulhu® is a registered trademark of Chaosium Inc.

## References

- [D&D 5th Edition System Reference Document](https://dnd.wizards.com/resources/systems-reference-document)
- [5e SRD CC-BY License](https://creativecommons.org/licenses/by/4.0/)
