# D20 - Go Dice Rolling Library

[![CI](https://github.com/jwebster45206/d20/actions/workflows/ci.yml/badge.svg)](https://github.com/jwebster45206/d20/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jwebster45206/d20.svg)](https://pkg.go.dev/github.com/jwebster45206/d20)

A Go library for dice rolling and tabletop RPG helpers: notation and fluent d20-style rolls, named modifiers, advantage/disadvantage, and a simple actor model. It is not a rules engine and does not implement SRD procedures such as ability modifiers or proficiency.

## Features

- **Dice Shorthand**: Parse standard string notation like "1d20+3" or "2d6"
- **Dice Longhand**: Fluent builder API with named modifiers and advantage/disadvantage
- **5e-style rolls**: d20 checks, advantage/disadvantage, and named modifiers (you supply the numbers)
- **Actor System**: Basic character/creature representation for combat and skill checks
- **Detailed Roll Results**: Bioware-inspired roll result formatting with full breakdowns

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/jwebster45206/d20"
)

func main() {
    // Create a roller
    roller := d20.NewRandomRoller()
    
    // Simple dice notation shorthand
    result, _ := roller.Roll("1d20+3")
    fmt.Printf("Attack roll: %d\n", result.Value)
    
    // With named modifiers
    result, err := roller.Dice(1, 20).
        WithModifier("strength", 3).
        WithModifier("proficiency", 2).
        Roll()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Attack roll: %d\n", result.Value)
    
    // Roll with advantage
    advResult, _ := roller.Dice(1, 20).
        WithAdvantage().
        WithModifier("dexterity", 4).
        Roll()
    fmt.Printf("Roll: %d (dice: %v)\n", advResult.Value, advResult.DiceRolls)
    
    // Roll and display detailed results
    detailResult, _ := roller.Dice(1, 20).
        WithModifier("strength", 3).
        WithModifier("proficiency", 2).
        Roll()
    fmt.Println(detailResult.Detail())
    // Example: "Rolled 1d20... 15; +3 strength, +2 proficiency; *Result: 20*"
}
```

## Core Types

### Roller & DiceManager

The `Roller` provides dice rolling functionality through both a fluent builder API and simple dice notation:

```go
// Create a roller
func NewRoller(seed int64) *Roller
func NewRandomRoller() *Roller

// Dice notation shorthand - simple and fast
func (r *Roller) Roll(notation string) (RollOutcome, error)

// Start building a roll - for complex scenarios
func (r *Roller) Dice(rollCount, dieFaces uint) *DiceManager
func (r *Roller) DiceExpr(notation string) *DiceManager

// Parse notation without rolling
func ParseNotation(notation string) (rollCount uint, dieFaces uint, modifiers []Modifier, err error)

var ErrInvalidDiceNotation error

// DiceManager - fluent API for configuring rolls
type DiceManager struct {
    RollCount     uint
    DieFaces      uint
    Modifiers     []Modifier
    AdvantageType AdvantageType
    Results       []RollOutcome // last successful roll; empty after a failed roll
}

func (rb *DiceManager) WithModifier(name string, value int) *DiceManager
func (rb *DiceManager) WithModifiers(modifiers map[string]int) *DiceManager
func (rb *DiceManager) WithAdvantage() *DiceManager
func (rb *DiceManager) WithDisadvantage() *DiceManager
func (rb *DiceManager) Error() error
func (rb *DiceManager) Roll() (RollOutcome, error)
func (rb *DiceManager) RollPercentile() (RollOutcome, error)
```

**Dice Notation Shorthand:**

The `Roll()` method and `ParseNotation` / `DiceExpr` accept standard dice notation strings:
- `"1d20"` - Roll one 20-sided die
- `"d20"` - Shorthand for 1d20
- `"2d6+3"` - Roll two 6-sided dice and add 3
- `"3d8-2"` - Roll three 8-sided dice and subtract 2
- `"1d100"` - Uniform 1–100 (notation). For Call of Cthulhu-style tens+ones (`00` = 100), use `Dice(2, 10).RollPercentile()`.

`ParseNotation` returns count, faces, and modifiers without a roller. `DiceExpr` applies that parse to a `DiceManager` so you can still chain `WithModifier`, `WithAdvantage`, and `Roll`. A notation bonus becomes a modifier named `"modifier"` and stacks with fluent modifiers. Invalid notation is stored on the manager: `Roll()` / `RollPercentile()` return it, and `Error()` is non-nil. `Results` holds the last successful outcome (one element) and is emptied on failure.

**Advantage/Disadvantage Mechanics:**
- **Advantage**: Rolls 2 dice, uses the higher value, returns both in `DiceRolls`
- **Disadvantage**: Rolls 2 dice, uses the lower value, returns both in `DiceRolls`
- **Normal**: Rolls 1 die per count, returns all in `DiceRolls`

This transparency allows you to see all dice rolled, even when using advantage/disadvantage.

### RollOutcome

The result of a dice roll operation:

```go
type RollOutcome struct {
    Value     int            // Final calculated result (dice + modifiers)
    DiceRolls []DieRoll      // Each die: Faces + Result (2 dice for adv/dis)
    Modifiers []Modifier      // Modifiers applied to the roll
}

type DieRoll struct {
    Faces  uint // Die size
    Result int  // Face showing
}

func (o RollOutcome) Detail() string  // Bioware-style description
```

**Examples:**
- Normal roll: `DiceRolls: [{20, 17}]`, `Detail(): "Rolled 1d20... 17; *Result: 17*"`
- With advantage: `DiceRolls: [{20, 6}, {20, 8}]`, `Value: 8`, `Detail(): "Rolled 2d20... 6, 8; *Result: 8*"`
- With modifiers: `Detail(): "Rolled 1d20... 6; +3 strength; *Result: 9*"`
- Percentile: `DiceRolls: [{10, 0}, {10, 0}]`, `Detail(): "Rolled 2d10... 0, 0; *Result: 100*"`

## Actor System

Actors represent characters, NPCs, and monsters. Attribute keys are caller-owned (lowercase by convention). `SkillCheck` and `D100SkillCheck` lowercase the skill name on lookup.

```go
type Actor struct {
    ID         string
    MaxHP      int
    HP         int
    AC         int
    Initiative int
    Attributes map[string]int // scores, skills
    Modifiers  map[string]int // caller-wired bonuses
}

func NewActor(id string) *Actor

func (a *Actor) AddHP(amount int)         // Increase HP (won't exceed MaxHP)
func (a *Actor) SubHP(amount int)         // Reduce HP (won't go below 0)
func (a *Actor) ResetHP()                 // Restore to MaxHP
func (a *Actor) IsKnockedOut() bool       // Returns true if HP == 0

func (a *Actor) SkillCheck(skill string, roller *Roller) (*DiceManager, error)
func (a *Actor) StrikeRoll(roller *Roller, keys ...string) *DiceManager
func (a *Actor) D100SkillCheck(skill string, roller *Roller) (*DiceManager, error)
```

Well-known names live in `github.com/jwebster45206/d20/vocab` (`vocab.Strength`, `vocab.Striking`, `vocab.Damage`). Custom string keys are allowed. 

### Creating Actors

```go
fighter := d20.NewActor("Ironpants")
fighter.MaxHP, fighter.HP = 45, 45
fighter.AC = 18

wizard := d20.NewActor("Merlin")
wizard.MaxHP, wizard.HP = 38, 38
wizard.AC = 14
wizard.Attributes = map[string]int{
    "intelligence": 18,
    "wisdom":       16,
}
wizard.Modifiers = map[string]int{
    "intelligence": 4,
    "striking":     4,
}

// Rolled stats: roll, then assign
roller := d20.NewRandomRoller()
hp, _ := roller.Roll("10d12+30")
str, _ := roller.Roll("3d6")
barbarian := d20.NewActor("Grog")
barbarian.MaxHP, barbarian.HP = hp.Value, hp.Value
barbarian.AC = 14
barbarian.Attributes["strength"] = str.Value
barbarian.Attributes["proficiency"] = 4
```

### Modifiers

`Modifiers` are bonuses the caller wires. They are not derived from `Attributes`. `StrikeRoll` applies only the keys you pass; missing keys are skipped. Situational extras go on the `DiceManager`.

```go
import "github.com/jwebster45206/d20/vocab"

actor.Modifiers[vocab.Strength] = 4
actor.Modifiers[vocab.Striking] = 3
actor.Modifiers[vocab.Damage] = 1

result, _ := actor.StrikeRoll(roller, vocab.Strength, vocab.Striking).Roll()
result, _ = actor.StrikeRoll(roller, vocab.Strength, vocab.Striking).WithModifier("bless", 1).Roll()
```

`vocab.Damage` is for damage rolls you build yourself (`Dice(1, 8).WithModifier(vocab.Damage, actor.Modifiers[vocab.Damage])`). `StrikeRoll` never auto-includes it.

### Attributes

`Attributes` is a `map[string]int`. Typical keys:

- **Core Abilities**: `strength`, `dexterity`, `constitution`, `intelligence`, `wisdom`, `charisma`
- **Skills**: `athletics`, `stealth`, `perception`, `insight`, etc.
- **Custom**: any string key with an integer value

```go
actor.Attributes["strength"] = 16
actor.Attributes["strength"] += 2  // temporary buff
delete(actor.Attributes, "stealth")
```

### Actor Roll Methods

Actor roll methods return `*DiceManager` for further configuration:

```go
builder, err := actor.SkillCheck("stealth", roller)
if err != nil {
    // skill not found in attributes
}
result, _ := builder.Roll()
result, _ = builder.WithAdvantage().Roll()

builder = actor.StrikeRoll(roller, vocab.Strength, vocab.Striking)
result, _ = builder.Roll()
result, _ = builder.WithAdvantage().Roll()
result, _ = builder.WithModifier("bless", 1).Roll()

builder, err = actor.D100SkillCheck("stealth", roller)
outcome, _ := builder.RollPercentile()
success := outcome.Value <= actor.Attributes["stealth"]
```

Advantage and disadvantage are configured on the `DiceManager`:

- **Advantage**: Rolls 2 dice, uses higher, shows both in `DiceRolls`
- **Disadvantage**: Rolls 2 dice, uses lower, shows both in `DiceRolls`
- **Normal**: Rolls the standard number of dice

d100/percentile rolls:

- **`Roll("1d100")`**: Uniform 1–100 via dice notation.
- **`Dice(2, 10).RollPercentile()` / `D100SkillCheck`**: Call of Cthulhu-style tens + ones (`00` = 100).

## Game systems

### 5e-style helpers

This library provides helpers for common 5e-style rolls (polyhedral dice, d20 + named modifiers, advantage/disadvantage, HP/AC on an actor). It does not compute ability modifiers, proficiency bonuses, or other SRD procedures — callers supply those numbers.

The [5e SRD](https://dnd.wizards.com/resources/systems-reference-document) is published by Wizards of the Coast under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/). This repository does not include SRD text.

### Call of Cthulhu-style d100

`Dice(2, 10).RollPercentile()` and `D100SkillCheck` implement compatible d100 roll-under mechanics (tens + ones, `00` = 100). They are not a Call of Cthulhu rules implementation.

Call of Cthulhu® is a registered trademark of Chaosium Inc. This library does not include copyrighted content from Call of Cthulhu sourcebooks.

## Examples

Both APIs return a structured `RollOutcome` (`Value`, `DiceRolls`, `Modifiers`, `Detail()`). The difference is only how you *request* the roll.

### Simple input — dice notation

Use `Roll(notation)` when a string is enough (optional `+`/`-` modifier in the notation):

```go
roller := d20.NewRoller(time.Now().UnixNano())

result, _ := roller.Roll("1d20")
fmt.Printf("Rolled: %d\n", result.Value)

result, _ = roller.Roll("d20+5")
fmt.Printf("With modifier: %d\n", result.Value)

result, _ = roller.Roll("2d6+3")
fmt.Printf("Damage: %d\n", result.Value)

fmt.Println(result.Detail()) // still a full RollOutcome
```

### Structured input — DiceManager

Use `Dice(count, faces)` or `DiceExpr(notation)` when you need named modifiers, advantage/disadvantage, or percentile:

```go
roller := d20.NewRoller(time.Now().UnixNano())

result, _ := roller.Dice(1, 20).
    WithModifier("strength", 5).
    WithModifier("proficiency", 3).
    Roll()
fmt.Printf("Attack: %d\n", result.Value)
fmt.Println(result.Detail())

result, _ = roller.DiceExpr("2d20+1").
    WithAdvantage().
    WithModifier("bless", 1).
    Roll()
fmt.Printf("Roll: %d (from dice: %v)\n", result.Value, result.DiceRolls)

result, _ = roller.Dice(2, 10).RollPercentile()
fmt.Printf("Percentile: %d (digits: %v)\n", result.Value, result.DiceRolls)
```

### Actor Usage

```go
import (
    "github.com/jwebster45206/d20"
    "github.com/jwebster45206/d20/vocab"
)

roller := d20.NewRoller(time.Now().UnixNano())

fighter := d20.NewActor("Ironpants")
fighter.MaxHP, fighter.HP = 45, 45
fighter.AC = 18
fighter.Attributes = map[string]int{
    "strength":     16,
    "dexterity":    14,
    "constitution": 15,
    "athletics":    5,
    "stealth":      2,
}
fighter.Modifiers = map[string]int{
    vocab.Strength: 3,
    vocab.Striking: 3,
}

wizard := d20.NewActor("Merlin")
wizard.MaxHP, wizard.HP = 22, 22
wizard.AC = 12
wizard.Attributes[vocab.Intelligence] = 18
wizard.Attributes[vocab.Wisdom] = 14
wizard.Modifiers[vocab.Intelligence] = 4

hp, _ := roller.Roll("12d12+48")
barbarian := d20.NewActor("Grog")
barbarian.MaxHP, barbarian.HP = hp.Value, hp.Value
barbarian.AC = 14
barbarian.Attributes["proficiency"] = 5
for _, key := range []string{"strength", "dexterity", "constitution", "intelligence", "wisdom", "charisma"} {
    out, _ := roller.Roll("3d6")
    barbarian.Attributes[key] = out.Value
}

result, _ := fighter.StrikeRoll(roller, vocab.Strength, vocab.Striking).Roll()
fmt.Printf("Strike: %d\n", result.Value)

result, _ = fighter.StrikeRoll(roller, vocab.Strength, vocab.Striking).WithAdvantage().Roll()
fmt.Printf("Strike with advantage: %d (dice: %v)\n", result.Value, result.DiceRolls)

result, _ = fighter.StrikeRoll(roller, vocab.Strength, vocab.Striking).
    WithModifier("bless", 1).
    WithModifier("cover_penalty", -2).
    Roll()

builder, _ := fighter.SkillCheck("stealth", roller)
result, _ = builder.Roll()

builder, _ = fighter.SkillCheck("athletics", roller)
result, _ = builder.WithAdvantage().Roll()

fighter.SubHP(15)
if !fighter.IsKnockedOut() {
    fmt.Printf("Fighter has %d/%d HP remaining\n", fighter.HP, fighter.MaxHP)
}
fighter.AddHP(8)

fighter.MaxHP = 50
fighter.ResetHP()

fighter.Attributes["strength"] += 2
fighter.Attributes["dexterity"] -= 1
```

### D100 System Usage

```go
investigator := d20.NewActor("Detective Morgan")
investigator.MaxHP, investigator.HP = 12, 12
investigator.AC = 10
investigator.Attributes = map[string]int{
    "stealth":     45,
    "fighting":    60,
    "firearms":    25,
    "spot_hidden": 70,
    "sanity":      65,
}

investigator.Attributes["sanity"] += 1
investigator.Attributes["sanity"] -= 3

builder, _ := investigator.D100SkillCheck("stealth", roller)
outcome, _ := builder.RollPercentile()
skillValue := investigator.Attributes["stealth"]
if outcome.Value <= skillValue {
    fmt.Printf("Stealth succeeded: rolled %d ≤ %d\n", outcome.Value, skillValue)
} else {
    fmt.Printf("Stealth failed: rolled %d\n", outcome.Value)
}

builder, _ = investigator.D100SkillCheck("fighting", roller)
outcome, _ = builder.RollPercentile()
```

## Future Enhancements

### Advanced Dice Notation
Currently, the dice notation parser supports basic formats like `"1d20"`, `"2d6+3"`, and `"3d8-2"`. 

**Planned additions:**
- `kh` (keep highest): `"4d6kh3"` - Roll 4d6, keep highest 3 (common for D&D ability scores)
- `kl` (keep lowest): `"4d6kl3"` - Roll 4d6, keep lowest 3
- `dh` (drop highest): `"4d6dh1"` - Roll 4d6, drop highest 1
- `dl` (drop lowest): `"4d6dl1"` - Roll 4d6, drop lowest 1 (equivalent to kh3)

## References

- [D&D 5th Edition System Reference Document](https://dnd.wizards.com/resources/systems-reference-document)
- [5e SRD CC-BY License](https://creativecommons.org/licenses/by/4.0/)


