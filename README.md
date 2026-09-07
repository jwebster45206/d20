# D20 - Go Dice Rolling Library

[![CI](https://github.com/jwebster45206/d20/actions/workflows/ci.yml/badge.svg)](https://github.com/jwebster45206/d20/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jwebster45206/d20.svg)](https://pkg.go.dev/github.com/jwebster45206/d20)

A Go library for dice rolling and tabletop RPG helpers: dice notation, named modifiers, advantage/disadvantage, and a simple actor model. It is not a rules engine and does not implement SRD procedures such as ability modifiers or proficiency.

## Features

- **Dice notation**: `"1d20+3"`, `"2d6"`, `"d20"`
- **Dice config**: named modifiers and advantage/disadvantage
- **Actors**: HP/AC, name, and caller-owned attributes, modifiers, and actions
- **RollOutcome.Detail()**: Bioware-style breakdown of dice and modifiers

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/jwebster45206/d20"
)

func main() {
    roller := d20.NewRandomRoller()

    result, _ := roller.RollExpr("1d20+3")
    fmt.Printf("Attack roll: %d\n", result.Value)

    d := d20.MustNewDice(1, 20)
    result, err := roller.Roll(d.
        WithModifier("strength", 3).
        WithModifier("proficiency", 2))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result.Detail())
    // Example: "Rolled 1d20... 15; +3 strength, +2 proficiency; *Result: 20*"

    d = d20.MustNewDice(1, 20)
    r, err := roller.Roll(d.
        WithAdvantage().
        WithModifier("dexterity", 4))
    if err != nil {
        log.Fatal(err)
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
func (r *Roller) RollExpr(notation string) (RollOutcome, error) // DiceFromExpr + Roll

var (
    ErrNilRoller              error
    ErrRollCountZero          error
    ErrDieFacesZero           error
    ErrInvalidAdvantage       error
    ErrPercentileRequires2d10 error
)
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

func NewDice(count, faces uint) (Dice, error)
func MustNewDice(count, faces uint) Dice

func (d Dice) WithModifier(name string, value int) Dice
func (d Dice) WithModifiers(modifiers map[string]int) Dice
func (d Dice) WithAdvantage() Dice
func (d Dice) WithDisadvantage() Dice
```

```go
attack := d20.MustNewDice(1, 20)
attack = attack.WithModifier("strength", 4)
roller.Roll(attack)
roller.Roll(attack.WithAdvantage()) // attack unchanged
```

Advantage (ignored by `RollPercentile`): rolls twice per die, uses the higher of each pair; disadvantage uses the lower. All rolls appear in `DiceRolls`.

`RollPercentile` is tens+ones (`00` = 100), distinct from rolling a uniform 1d100. Requires 2d10.

### Dice notation

```go
func DiceFromExpr(expr string) (Dice, error)
func MustDiceFromExpr(expr string) Dice

var ErrInvalidDiceNotation error
```

Accepted: `"1d20"`, `"d20"`, `"2d6+3"`, `"3d8-2"`. `"1d100"` is a uniform 1–100 die. A trailing `+N`/`-N` becomes a modifier named `"modifier"`. Invalid notation fails.

```go
d := d20.MustDiceFromExpr("2d6+3")
out, err := roller.Roll(d.WithModifier("bless", 1))
out, err = roller.RollExpr("2d6+3")
d, err = d20.DiceFromExpr(userNotation) // when the string is data
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

Map keys are lowercase by convention. `D20Dice` and `Dice` lowercase the *query*; they do not rewrite stored keys. They return a `Dice` spec; the roller executes it.

```go
type Actor struct {
    ID         string
    Name       string
    Alignment  string
    MaxHP      int
    HP         int
    AC         int
    Attributes map[string]int // caller-owned numbers (scores or skill bonuses)
    Modifiers  map[string]int // caller-wired roll bonuses
    Actions    map[string]Action
}

type Action struct {
    ID      string
    Name    string
    Type    string
    Attempt Dice
    Effect  Dice
    Charges *uint // nil means unlimited
}

func NewActor(id string) *Actor

func (a *Actor) Normalize() error
func (a *Action) Normalize()

func (a *Actor) Action(id string) (Action, bool)
func (a *Actor) Dice(d Dice, keys ...string) Dice
func (a *Actor) D20Dice(keys ...string) Dice
func (a *Actor) DiceFromExpr(expr string, keys ...string) (Dice, error)
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
    "intelligence": 18,
    "athletics":    5,
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

HP, AC, name, alignment, attributes, and actions are caller-owned fields. Mutate them directly. `Normalize` rewrites actor ID and map keys to lowercase snake_case, and copies each action map key onto `Action.ID`.

### Modifiers and rolls

`Modifiers` are not derived from `Attributes`. `D20Dice` applies only the keys you pass; missing keys are skipped. Situational extras go on the returned `Dice`. `RollPercentile` ignores advantage/disadvantage.

```go
import "github.com/jwebster45206/d20/vocab"

actor.Modifiers[vocab.Strength] = 4
actor.Modifiers[vocab.Striking] = 3
actor.Modifiers[vocab.Damage] = 1

d := actor.D20Dice(vocab.Strength, vocab.Striking)
result, _ := roller.Roll(d)
result, _ = roller.Roll(d.WithModifier("bless", 1))
result, _ = roller.Roll(d.WithAdvantage())

dmg := actor.Dice(d20.MustDiceFromExpr("1d8"), vocab.Damage, vocab.Strength)
```


`RollPercentile` of 2d10 (`NewDice(2, 10)`) implements compatible d100 roll-under mechanics (tens + ones, `00` = 100). 
