# D20 - Go Dice Rolling Library

[![CI](https://github.com/jwebster45206/d20/actions/workflows/ci.yml/badge.svg)](https://github.com/jwebster45206/d20/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jwebster45206/d20.svg)](https://pkg.go.dev/github.com/jwebster45206/d20)

A Go library for dice rolling and D20 system mechanics, designed with 5e SRD compatibility in mind. This library provides a foundation for tabletop RPG applications with clean interfaces for dice rolling, modifiers, and basic actor statistics.

## Features

- **Flexible Dice Rolling**: Support for any dice combination with modifiers
- **5e SRD Compatible**: Follows D&D 5th Edition System Reference Document conventions
- **Actor System**: Basic character/creature representation for combat and skill checks
- **Detailed Roll Results**: Bioware-style roll result formatting with full breakdowns
- **Extensible Design**: Clean interfaces for future expansion

## Installation

```bash
go get github.com/jwebster45206/d20
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/jwebster45206/d20"
)

func main() {
    // Create a roller
    roller := d20.NewRoller(42) // seed for reproducible results
    
    // Roll some dice
    result, err := roller.Roll(2, 20, []d20.Modifier{
        {Value: 3, Reason: "Strength"},
        {Value: 2, Reason: "Proficiency"},
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Println(result.Detail)
    // Rolled 2d20...; values 16, 12; +3 strength, +2 proficiency; *Result: 33*
}
```

## Core Types

### Roller

The `Roller` struct provides the core dice rolling functionality:

```go
type Roller struct {
    // Internal random number generator (seeded for deterministic results)
}

func NewRoller(seed int64) *Roller
func NewRandomRoller() *Roller
func (r *Roller) Roll(rollCount uint, dieFaces uint, modifiers []Modifier) (RollValue, error)
```

- **rollCount**: Number of dice to roll
- **dieFaces**: Number of faces on each die (d4, d6, d8, d10, d12, d20, d100)
- **modifiers**: Array of modifiers to apply to the roll

### Modifier

Modifiers represent buffs or debuffs applied to dice rolls:

```go
type Modifier struct {
    Value  int    // Positive for bonus, negative for penalty
    Reason string // Description of the modifier source
}
```

### RollValue

The result of a dice roll operation:

```go
type RollValue struct {
    Value     int    // Final calculated result
    DiceRolls []int  // Raw values from each die rolled
    Detail    string // Formatted roll description
}
```

The `Detail` field provides a Bioware-style formatted string showing the complete roll breakdown:
- `Rolled 2d20...; values 16, 12; +3 strength, +2 proficiency; *Result: 33*`

## Actor System

Actors represent characters, NPCs, and monsters in the game world:

```go
type Actor struct {
    // Private fields - use constructor and accessors
}

// Constructor
func NewActor(hp, ac, initiative int) (*Actor, error)
func NewActorWithAttributes(hp, ac, init int, attrs map[string]int) (*Actor, error)

// HP Management
func (a *Actor) HP() int                  // Current HP
func (a *Actor) MaxHP() int               // Maximum HP
func (a *Actor) SetHP(hp int) error       // Set current HP (0 to max)
func (a *Actor) SetMaxHP(maxHP int) error // Set maximum HP (auto-adjusts current if needed)
func (a *Actor) TakeDamage(damage int)    // Reduce HP (won't go below 0)
func (a *Actor) Heal(amount int)          // Increase HP (won't exceed max)
func (a *Actor) ResetHP()                 // Restore to max HP
func (a *Actor) IsAlive() bool            // Returns true if HP > 0

// AC and Initiative
func (a *Actor) AC() int
func (a *Actor) SetAC(ac int) error
func (a *Actor) Initiative() int
func (a *Actor) SetInitiative(init int)

// Attribute Management (keys automatically lowercased)
func (a *Actor) Attribute(key string) (int, bool)
func (a *Actor) SetAttribute(key string, value int)
func (a *Actor) HasAttribute(key string) bool
func (a *Actor) RemoveAttribute(key string)

// Combat Modifier Management (reasons automatically lowercased)
func (a *Actor) AddCombatModifier(m Modifier)
func (a *Actor) RemoveCombatModifier(reason string)
func (a *Actor) GetCombatModifiers() []Modifier
```

**Design Notes**: Actor uses private fields with accessor methods to:
- Track both maximum and current HP separately (damage/healing affect current, leveling affects max)
- Enforce data validation (HP and AC must be positive, current HP can't exceed max)
- Automatically lowercase all attribute keys and modifier reasons for consistency
- Prevent direct slice/map mutations that could cause bugs
- Provide a clean, discoverable API

### Creating Actors

Use the constructor functions to create actors with validated initial values:

```go
// Basic constructor with validation
fighter, err := d20.NewActor(45, 18, 2)
if err != nil {
    panic(err) // HP and AC must be > 0
}

// Or create with initial attributes
wizard, err := d20.NewActorWithAttributes(22, 12, 3, map[string]int{
    "intelligence": 18,
    "wisdom":       14,
    "arcana":       8,
    "history":      6,
})
```

### Combat Modifiers

The `CombatModifiers` field contains all active offensive bonuses and penalties that apply to the actor's attack rolls. This includes:

- **Ability Modifiers**: Strength for melee, Dexterity for ranged/finesse weapons
- **Proficiency Bonus**: If proficient with the weapon being used
- **Equipment Bonuses**: Magic weapon bonuses (+1, +2, +3 weapons)
- **Spell Effects**: Bless, Guidance, or other temporary bonuses
- **Class Features**: Fighting styles, rage bonuses, etc.

Note: AC represents the actor's total Armor Class including all static bonuses (armor, shields, Dex modifier, natural armor, etc.). Situational AC modifiers (cover, spells) are handled at the time of attack resolution.

### Attributes

The flexible attribute system supports standard D&D 5e ability scores and derived statistics:

- **Core Abilities**: `strength`, `dexterity`, `constitution`, `intelligence`, `wisdom`, `charisma`
- **Skills**: `athletics`, `stealth`, `perception`, `insight`, etc.
- **Custom Attributes**: Any string key with integer value

### Actor Methods

Actors have methods for common 5e mechanics:

```go
// AdvantageType represents advantage, normal, or disadvantage rolls
type AdvantageType int

const (
    Disadvantage AdvantageType = iota
    Normal
    Advantage
)

// SkillCheck performs a skill check using 5e conventions (d20 + modifiers)
func (a *Actor) SkillCheck(skill string, roller *Roller, advantage AdvantageType) (RollValue, error)

// D100SkillCheck performs a percentile skill check (roll under skill value)
func (a *Actor) D100SkillCheck(skill string, roller *Roller, bonus int) (bool, RollValue, error)

// AttackRoll makes an attack roll using the actor's CombatModifiers
func (a *Actor) AttackRoll(roller *Roller, advantage AdvantageType) (RollValue, error)

// AttackRollWithModifiers makes an attack roll with additional situational modifiers
func (a *Actor) AttackRollWithModifiers(roller *Roller, advantage AdvantageType, extraModifiers []Modifier) (RollValue, error)
```

#### Advantage/Disadvantage Mechanics

- **Advantage**: Roll 2d20, take the higher result
- **Disadvantage**: Roll 2d20, take the lower result  
- **Normal**: Roll 1d20

This system is core to 5e and applies to attack rolls, skill checks, and saving throws.

#### D100 System Support

The library also supports d100/percentile systems like Call of Cthulhu:

- **D100SkillCheck**: Roll d100, succeed if result ≤ skill value
- **Bonus Die**: Roll 2d10 for tens digit, take the better result (equivalent to advantage)
- **Penalty Die**: Roll 2d10 for tens digit, take the worse result (equivalent to disadvantage)
- **Combat**: Uses skill checks (Fighting, Firearms, etc.) rather than separate attack rolls

#### Custom Dice Systems

#### D100 System Support

The library supports d100/percentile systems like Call of Cthulhu:

- **D100SkillCheck**: Roll d100, succeed if result ≤ skill value
- **Bonus Die** (bonus > 0): Roll multiple d10s for tens digit, take the LOWEST (better chance)
- **Penalty Die** (bonus < 0): Roll multiple d10s for tens digit, take the HIGHEST (worse chance)
- **Combat**: Uses skill checks (Fighting, Firearms, etc.) rather than separate attack rolls

### Skill Checks

Skill checks use the actor's attribute values from the `Attributes` map:

```go
// D&D 5e skill checks (d20 + modifiers)
stealthCheck, _ := actor.SkillCheck("stealth", roller, d20.Normal)
athleticsCheck, _ := actor.SkillCheck("athletics", roller, d20.Advantage)
perceptionCheck, _ := actor.SkillCheck("perception", roller, d20.Disadvantage)

// Call of Cthulhu skill checks (d100, roll under skill value)
success, roll, _ := investigator.D100SkillCheck("stealth", roller, false)
bonusSuccess, roll, _ := investigator.D100SkillCheck("fighting", roller, true) // bonus die

// Check success for d100 systems
if success {
    fmt.Printf("Skill check succeeded with %d (needed ≤ %d)", roll.Value, investigator.Attributes["stealth"])
}
```

## Game System Compliance

### D&D 5th Edition SRD

This library follows the D&D 5th Edition System Reference Document (SRD 5.2) under CC-BY licensing:

- **Dice Mechanics**: Standard polyhedral dice (d4, d6, d8, d10, d12, d20, d100)
- **Ability Scores**: Six core abilities with standard modifiers
- **Advantage/Disadvantage**: Roll twice, take higher/lower
- **Proficiency Bonus**: Scalable bonus system
- **Combat Stats**: AC, HP, Initiative as per SRD

### Call of Cthulhu®

This library implements d100 percentile mechanics compatible with Call of Cthulhu® by Chaosium Inc:

- **Percentile Skills**: Roll d100, succeed if result ≤ skill value
- **Bonus/Penalty Dice**: Multiple d10s for tens digit, take best/worst
- **Roll-Under System**: Success determined by rolling under skill percentage

Call of Cthulhu® is a registered trademark of Chaosium Inc. This library implements compatible game mechanics but does not include copyrighted content from Call of Cthulhu sourcebooks.

## Examples

### Basic Dice Rolling

```go
roller := d20.NewRoller(time.Now().UnixNano())

// Simple d20 roll
result, _ := roller.Roll(1, 20, nil)

// Attack roll with modifiers
attackRoll, _ := roller.Roll(1, 20, []d20.Modifier{
    {Value: 5, Reason: "Strength"},
    {Value: 3, Reason: "Proficiency"},
})

// Damage roll
damageRoll, _ := roller.Roll(1, 8, []d20.Modifier{
    {Value: 3, Reason: "Strength"},
})
```

### Actor Usage

```go
// Create a character using the constructor
fighter, _ := d20.NewActor(45, 18, 2)

// Set up attributes
fighter.SetAttribute("strength", 16)
fighter.SetAttribute("dexterity", 14)
fighter.SetAttribute("constitution", 15)
fighter.SetAttribute("athletics", 5)  // includes proficiency
fighter.SetAttribute("stealth", 2)    // dex modifier only

// Add combat modifiers
fighter.AddCombatModifier(d20.Modifier{Value: 3, Reason: "Strength"})
fighter.AddCombatModifier(d20.Modifier{Value: 3, Reason: "Proficiency"})
fighter.AddCombatModifier(d20.Modifier{Value: 1, Reason: "Magic Weapon"})

// Or create with attributes in one call
wizard, _ := d20.NewActorWithAttributes(22, 12, 3, map[string]int{
    "intelligence": 18,
    "wisdom":       14,
    "arcana":       8,
    "history":      6,
})

// Make an attack roll using the actor's combat modifiers
attackRoll, _ := fighter.AttackRoll(roller, d20.Normal)

// Attack with advantage (flanking, help action, etc.)
advantageAttack, _ := fighter.AttackRoll(roller, d20.Advantage)

// Attack with additional situational modifiers
situationalAttack, _ := fighter.AttackRollWithModifiers(roller, d20.Normal, []d20.Modifier{
    d20.NewModifier(2, "Flanking"),
    d20.NewModifier(-2, "Partial Cover"),
})

// Perform skill checks
stealthCheck, _ := fighter.SkillCheck("stealth", roller, d20.Normal)
athleticsCheck, _ := fighter.SkillCheck("athletics", d20.Advantage)

// HP management
fighter.TakeDamage(15)
if fighter.IsAlive() {
    fmt.Printf("Fighter has %d/%d HP remaining\n", fighter.HP(), fighter.MaxHP())
}
fighter.Heal(8)

// Level up - increase max HP
fighter.SetMaxHP(50)
fighter.ResetHP() // Full heal after rest
```

### D100 System Usage

```go
// Create a Call of Cthulhu investigator
investigator, _ := d20.NewActorWithAttributes(12, 10, 0, map[string]int{
    "stealth":      45,  // 45% skill
    "fighting":     60,  // 60% skill  
    "firearms":     25,  // 25% skill
    "spot_hidden":  70,  // 70% skill
})

// Perform d100 skill checks
stealthSuccess, stealthRoll, _ := investigator.D100SkillCheck("stealth", roller, false)
fightingSuccess, fightingRoll, _ := investigator.D100SkillCheck("fighting", roller, true) // bonus die

// Check results
if stealthSuccess {
    fmt.Printf("Stealth succeeded: rolled %d ≤ %d", stealthRoll.Value, investigator.Attributes["stealth"])
}
```


## References

- [D&D 5th Edition System Reference Document](https://dnd.wizards.com/resources/systems-reference-document)
- [5e SRD CC-BY License](https://creativecommons.org/licenses/by/4.0/)

## Implementation Plan - v0 (Code Agent friend, this section is just for you!)

This is the complete implementation plan for v0 of the D20 library. We're building everything from scratch with the new design.

### Phase 1: Core Foundation ⚡ PRIORITY
**Goal**: Get basic dice rolling working with the new API
**Status**: Complete.

### Phase 2: Actor Methods 🎯 CORE FEATURES
**Goal**: Implement the main Actor functionality for D&D 5e
**Status**: Complete

### Phase 3: Extended Dice Systems 🌍 MULTI-SYSTEM
**Goal**: Support other RPG systems beyond D&D 5e
**Status**: Complete

#### 3.1 D100 System Support
```go
func (a *Actor) D100SkillCheck(skill string, roller *Roller, bonus bool) (bool, RollValue, error)
```
- ✅ Roll d100 vs skill percentage
- ✅ Bonus/penalty die mechanics (2d10, take better/worse tens digit)
- ✅ Return success boolean + roll details
- ✅ Support Call of Cthulhu mechanics and mention by name (if this is allowed)
- ✅ Add Chaosium to attributions. 

### Phase 4: Testing & Polish 🧪 QUALITY
**Goal**: Ensure reliability and usability

Add examples and shore up unit tests if needed.

### Success Criteria for v0

- ✅ All types compile and work as specified in README
- ✅ Can create characters and perform basic D&D 5e skill checks
- ✅ Can make attack rolls with advantage/disadvantage
- ✅ Bioware-style roll formatting works correctly
- ✅ Deterministic testing with seeds works
- ✅ Basic support for D100
- ✅ Clean, well-documented API ready for story-engine integration

