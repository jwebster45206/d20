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
    // Rolled 2d20... values 16, 12; +3 strength, +2 proficiency; *Result: 33*
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
- `Rolled 2d20... values 16, 12; +3 strength, +2 proficiency; *Result: 33*`

## Actor System

Actors represent characters, NPCs, and monsters in the game world:

```go
type Actor struct {
    HP              int               // Hit Points
    AC              int               // Armor Class (total, including all bonuses)
    Initiative      int               // Initiative/speed modifier
    CombatModifiers []Modifier        // Active offensive modifiers (weapon, strength, proficiency, etc.)
    Attributes      map[string]int    // Flexible attribute system
}
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

// SkillCheckWithDice performs a skill check with custom dice (for other RPG systems)
func (a *Actor) SkillCheckWithDice(skill string, roller *Roller, rollCount uint, dieFaces uint, advantage AdvantageType) (RollValue, error)

// D100SkillCheck performs a percentile skill check (roll under skill value)
func (a *Actor) D100SkillCheck(skill string, roller *Roller, bonus bool) (bool, RollValue, error)

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

For maximum flexibility, use `SkillCheckWithDice` to support any RPG system:

- **GURPS**: 3d6 roll-under system
- **Savage Worlds**: Exploding dice mechanics
- **Custom Games**: Any combination of dice count and faces

### Skill Checks

Skill checks use the actor's attribute values from the `Attributes` map:

```go
// D&D 5e skill checks (d20 + modifiers)
stealthCheck, _ := actor.SkillCheck("stealth", roller, d20.Normal)
athleticsCheck, _ := actor.SkillCheck("athletics", roller, d20.Advantage)
perceptionCheck, _ := actor.SkillCheck("perception", roller, d20.Disadvantage)

// Custom dice systems using SkillCheckWithDice
gurpsCheck, _ := actor.SkillCheckWithDice("stealth", roller, 3, 6, d20.Normal)    // GURPS 3d6
savageCheck, _ := actor.SkillCheckWithDice("fighting", roller, 1, 8, d20.Normal)  // Savage Worlds d8
customCheck, _ := actor.SkillCheckWithDice("magic", roller, 2, 12, d20.Advantage) // Custom 2d12 system

// Call of Cthulhu skill checks (d100, roll under skill value)
success, roll, _ := investigator.D100SkillCheck("stealth", roller, false)
bonusSuccess, roll, _ := investigator.D100SkillCheck("fighting", roller, true) // bonus die

// Check success for d100 systems
if success {
    fmt.Printf("Skill check succeeded with %d (needed ≤ %d)", roll.Value, investigator.Attributes["stealth"])
}
```

## 5e SRD Compliance

This library follows the D&D 5th Edition System Reference Document (SRD 5.2) under CC-BY licensing:

- **Dice Mechanics**: Standard polyhedral dice (d4, d6, d8, d10, d12, d20, d100)
- **Ability Scores**: Six core abilities with standard modifiers
- **Advantage/Disadvantage**: Roll twice, take higher/lower
- **Proficiency Bonus**: Scalable bonus system
- **Combat Stats**: AC, HP, Initiative as per SRD

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
// Create a character
fighter := &d20.Actor{
    HP:         45,
    AC:         18, // Total AC including armor, shield, dex modifier
    Initiative: 2,
    CombatModifiers: []d20.Modifier{
        {Value: 3, Reason: "Strength"},
        {Value: 3, Reason: "Proficiency"},
        {Value: 1, Reason: "Magic Weapon"},
    },
    Attributes: map[string]int{
        "strength":     16,
        "dexterity":    14, 
        "constitution": 15,
        "athletics":    5,  // includes proficiency
        "stealth":      2,  // dex modifier only
    },
}

// Make an attack roll using the actor's combat modifiers
attackRoll, _ := fighter.AttackRoll(roller, d20.Normal)

// Attack with advantage (flanking, help action, etc.)
advantageAttack, _ := fighter.AttackRoll(roller, d20.Advantage)

// Attack with additional situational modifiers
situationalAttack, _ := fighter.AttackRollWithModifiers(roller, d20.Normal, []d20.Modifier{
    {Value: 2, Reason: "Flanking"},
    {Value: -2, Reason: "Partial Cover"},
})

// Perform skill checks
stealthCheck, _ := fighter.SkillCheck("stealth", roller, d20.Normal)
athleticsCheck, _ := fighter.SkillCheck("athletics", roller, d20.Advantage)
```

### D100 System Usage

```go
// Create a Call of Cthulhu investigator
investigator := &d20.Actor{
    HP: 12,
    Attributes: map[string]int{
        "stealth":   45,  // 45% skill
        "fighting":  60,  // 60% skill  
        "firearms":  25,  // 25% skill
        "spot_hidden": 70, // 70% skill
    },
}

// Perform d100 skill checks
stealthSuccess, stealthRoll, _ := investigator.D100SkillCheck("stealth", roller, false)
fightingSuccess, fightingRoll, _ := investigator.D100SkillCheck("fighting", roller, true) // bonus die

// Check results
if stealthSuccess {
    fmt.Printf("Stealth succeeded: rolled %d ≤ %d", stealthRoll.Value, investigator.Attributes["stealth"])
}
```

## Contributing

This library aims to provide a solid foundation for D20-based gaming applications while maintaining compatibility with 5e SRD conventions. Contributions should:

1. Follow 5e SRD guidelines and terminology
2. Maintain clean, testable interfaces
3. Include comprehensive documentation
4. Add appropriate test coverage

## License

This project is licensed under the MIT License. Content inspired by or derived from the D&D 5th Edition System Reference Document is used under the Creative Commons Attribution 4.0 International License.

## References

- [D&D 5th Edition System Reference Document](https://dnd.wizards.com/resources/systems-reference-document)
- [5e SRD CC-BY License](https://creativecommons.org/licenses/by/4.0/)


## Implementation Plan - v0 (Code Agent friend, this section is just for you!)

This is the complete implementation plan for v0 of the D20 library. We're building everything from scratch with the new design.

### Phase 1: Core Foundation ⚡ PRIORITY
**Goal**: Get basic dice rolling working with the new API

#### 1.1 Core Types & Structs
```go
// types.go
type Modifier struct {
    Value  int
    Reason string
}

type RollValue struct {
    Value     int
    DiceRolls []int
    Detail    string
}

type AdvantageType int
const (
    Disadvantage AdvantageType = iota
    Normal
    Advantage
)

type Roller struct {
    rng *rand.Rand
}

type Actor struct {
    HP              int
    AC              int
    Initiative      int
    CombatModifiers []Modifier
    Attributes      map[string]int
}
```

#### 1.2 Core Roller Implementation
```go
// roller.go
func NewRoller(seed int64) *Roller
func (r *Roller) Roll(rollCount uint, dieFaces uint, modifiers []Modifier) (RollValue, error)
```
- ✅ Implement seeded random number generation
- ✅ Roll multiple dice and sum results
- ✅ Apply modifiers to final result
- ✅ Generate Bioware-style Detail string
- ✅ Input validation (rollCount > 0, dieFaces > 0)

#### 1.3 Detail Formatting
```go
// Format: "*Rolled 2d20*; Values 16, 12; Modifier +3 Strength, +2 Proficiency; *Result: 33*"
func formatRollDetail(rollCount uint, dieFaces uint, rolls []int, modifiers []Modifier, finalValue int) string
```

#### Unit Tests
Test `Roll` with deterministic seeds.

### Phase 2: Actor Methods 🎯 CORE FEATURES
**Goal**: Implement the main Actor functionality for D&D 5e

#### 2.1 Basic Skill Checks
```go
// actor.go
func (a *Actor) SkillCheck(skill string, roller *Roller, advantage AdvantageType) (RollValue, error)
```
- Look up skill value from Attributes map
- Handle Normal/Advantage/Disadvantage dice rolling
- Add skill modifier to d20 roll
- Calls RollWithAdvantage
- Return formatted RollValue

#### 2.2 Attack Rolls
```go
func (a *Actor) AttackRoll(roller *Roller, advantage AdvantageType) (RollValue, error)
func (a *Actor) AttackRollWithModifiers(roller *Roller, advantage AdvantageType, extraModifiers []Modifier) (RollValue, error)
```
- Use actor's CombatModifiers for base attack bonus
- Handle advantage/disadvantage mechanics
- Combine base and situational modifiers
- Calls RollWithAdvantage
- Return formatted attack roll result

#### 2.3 Advantage/Disadvantage Logic

```go
func RollWithAdvantage(roller *Roller, rollCount uint, dieFaces uint, advantage AdvantageType) ([]int, error)
```
- Normal: Roll normally
- Advantage: Roll twice, take higher
- Disadvantage: Roll twice, take lower
- Handle multiple dice with advantage (each die gets advantage)
- **Called internally by SkillCheck, AttackRoll, and SkillCheckWithDice**
- **Public utility** - also available for direct use (saving throws, initiative, etc.)

#### Unit Tests
Test each public method with deterministic seeds.

### Phase 3: Extended Dice Systems 🌍 MULTI-SYSTEM
**Goal**: Support other RPG systems beyond D&D 5e

#### 3.1 Flexible Skill Checks
```go
func (a *Actor) SkillCheckWithDice(skill string, roller *Roller, rollCount uint, dieFaces uint, advantage AdvantageType) (RollValue, error)
```
- Allow custom dice for any RPG system
- GURPS 3d6, Savage Worlds variable dice, etc.
- Reuse advantage/disadvantage logic

#### 3.2 D100 System Support
```go
func (a *Actor) D100SkillCheck(skill string, roller *Roller, bonus bool) (bool, RollValue, error)
```
- ✅ Roll d100 vs skill percentage
- ✅ Bonus/penalty die mechanics (2d10, take better/worse tens digit)
- ✅ Return success boolean + roll details
- ✅ Support Call of Cthulhu mechanics

### Phase 4: Testing & Polish 🧪 QUALITY
**Goal**: Ensure reliability and usability

#### 4.1 Unit Tests
```go
// roller_test.go, actor_test.go
```
- ✅ Test deterministic rolling with seeds
- ✅ Test advantage/disadvantage mechanics
- ✅ Test all Actor methods
- ✅ Test edge cases (invalid skills, zero modifiers, etc.)
- ✅ Test Bioware-style formattingulhu, GURPS)
- ✅ Validate roll distributions are correct

### Implementation Priority Order

1. **Start Here**: Core types, NewRoller, basic Roll method
2. **Essential**: SkillCheck and AttackRoll (basic D&D 5e functionality)
3. **Important**: Advantage/disadvantage mechanics
4. **Nice to Have**: SkillCheckWithDice (flexible systems)
5. **Advanced**: D100SkillCheck (Call of Cthulhu support)
6. **Final**: Comprehensive testing and polish

### Success Criteria for v0

- ✅ All types compile and work as specified in README
- ✅ Can create characters and perform basic D&D 5e skill checks
- ✅ Can make attack rolls with advantage/disadvantage
- ✅ Bioware-style roll formatting works correctly
- ✅ Deterministic testing with seeds works
- ✅ Basic support for D100
- ✅ Clean, well-documented API ready for story-engine integration

