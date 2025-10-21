# D20 - Go Dice Rolling Library

[![CI](https://github.com/jwebster45206/d20/actions/workflows/ci.yml/badge.svg)](https://github.com/jwebster45206/d20/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jwebster45206/d20.svg)](https://pkg.go.dev/github.com/jwebster45206/d20)

A Go library for dice rolling and D20 system mechanics, designed with 5e SRD compatibility in mind. This library provides a foundation for tabletop RPG applications with clean interfaces for dice rolling, modifiers, and basic actor statistics.

## Features

- **Dice Shorthand**: Parse standard string notation like "1d20+3" or "2d6"
- **Dice Longhand**: Fluent builder API with support for all dice combinations, named modifiers, and advantage/disadvantage
- **5e SRD Compatible**: Follows D&D 5th Edition System Reference Document conventions
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
    
    // Or use the fluent API for more control and named modifiers
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
    fmt.Println(detailResult.Detail)
    // Example: "Rolled 1d20... 15; +3 strength, +2 proficiency; *Result: 20*"
}
```

## Core Types

### Roller & RollBuilder

The `Roller` provides dice rolling functionality through both a fluent builder API and simple dice notation:

```go
// Create a roller
func NewRoller(seed int64) *Roller
func NewRandomRoller() *Roller

// Dice notation shorthand - simple and fast
func (r *Roller) Roll(notation string) (RollOutcome, error)

// Start building a roll - for complex scenarios
func (r *Roller) Dice(rollCount, dieFaces int) *RollBuilder

// RollBuilder - fluent API for configuring rolls
type RollBuilder struct { /* private fields */ }

func (rb *RollBuilder) WithModifier(name string, value int) *RollBuilder
func (rb *RollBuilder) WithModifiers(modifiers map[string]int) *RollBuilder
func (rb *RollBuilder) WithAdvantage() *RollBuilder
func (rb *RollBuilder) WithDisadvantage() *RollBuilder
func (rb *RollBuilder) Normal() *RollBuilder
func (rb *RollBuilder) Roll() (*RollOutcome, error)
```

**Dice Notation Shorthand:**

The `Roll()` method accepts standard dice notation strings:
- `"1d20"` - Roll one 20-sided die
- `"d20"` - Shorthand for 1d20
- `"2d6+3"` - Roll two 6-sided dice and add 3
- `"3d8-2"` - Roll three 8-sided dice and subtract 2
- `"1d100"` - Percentile dice

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
    DiceRolls []int          // Raw die values (2 dice for adv/dis, 1+ for normal)
    Detail    string         // Human-readable description
}
```

**Examples:**
- Normal roll: `DiceRolls: [17]`, `Detail: "Rolled 1d20... 17; *Result: 17*"`
- With advantage: `DiceRolls: [6, 8]`, `Value: 8`, `Detail: "Rolled 1d20... 6, 8; *Result: 8*"`
- With modifiers: `Detail: "Rolled 1d20... 6; +3 strength; *Result: 9*"`

## Actor System

Actors represent characters, NPCs, and monsters in the game world. The library uses a fluent builder pattern for creating actors:

```go
type Actor struct {
    // Private fields - use builder and accessors
}

// ActorBuilder - fluent API for creating actors
type ActorBuilder struct { /* private fields */ }

func NewActor(id string) *ActorBuilder
func (ab *ActorBuilder) WithHP(hp int) *ActorBuilder
func (ab *ActorBuilder) WithAC(ac int) *ActorBuilder
func (ab *ActorBuilder) WithAttribute(name string, value int) *ActorBuilder
func (ab *ActorBuilder) WithAttributes(attrs map[string]int) *ActorBuilder
func (ab *ActorBuilder) WithCombatModifier(name string, value int) *ActorBuilder
func (ab *ActorBuilder) WithCombatModifiers(mods map[string]int) *ActorBuilder
func (ab *ActorBuilder) Build() (*Actor, error)

// Rolled stat methods - require WithRoller() first
func (ab *ActorBuilder) WithRoller(roller *Roller) *ActorBuilder
func (ab *ActorBuilder) WithRolledHP(roll string) *ActorBuilder
func (ab *ActorBuilder) WithRolledAttribute(key string, roll string) *ActorBuilder
func (ab *ActorBuilder) WithRolledAttributes(attrs map[string]string) *ActorBuilder

// HP Management
func (a *Actor) ID() string               // Normalized identifier
func (a *Actor) HP() int                  // Current HP
func (a *Actor) MaxHP() int               // Maximum HP
func (a *Actor) SetHP(hp int)             // Set current HP (0 to max)
func (a *Actor) SetMaxHP(maxHP int)       // Set maximum HP (auto-adjusts current if needed)
func (a *Actor) AddHP(amount int)         // Increase HP (won't exceed max)
func (a *Actor) SubHP(amount int)         // Reduce HP (won't go below 0)
func (a *Actor) ResetHP()                 // Restore to max HP
func (a *Actor) IsKnockedOut() bool       // Returns true if HP <= 0

// AC and Initiative
func (a *Actor) AC() int
func (a *Actor) SetAC(ac int)
func (a *Actor) Initiative() int
func (a *Actor) SetInitiative(init int)

// Attribute Management 
func (a *Actor) Attribute(key string) (int, bool)
func (a *Actor) SetAttribute(key string, value int)
func (a *Actor) HasAttribute(key string) bool
func (a *Actor) RemoveAttribute(key string)
func (a *Actor) IncrementAttribute(key string, amount int)
func (a *Actor) DecrementAttribute(key string, amount int)

// Combat Modifier Management 
func (a *Actor) AddCombatModifier(name string, value int)
func (a *Actor) RemoveCombatModifier(name string)

// Roll Methods
func (a *Actor) SkillCheck(skill string, roller *Roller) (*RollBuilder, error)
func (a *Actor) AttackRoll(roller *Roller) *RollBuilder
func (a *Actor) D100SkillCheck(skill string, roller *Roller) (bool, *RollOutcome, error)
```

### Creating Actors

Use the builder pattern to create actors with optional configuration:

```go
// Basic actor with fixed stats
fighter, _ := d20.NewActor("Ironpants").
    WithHP(45).
    WithAC(18).
    Build()

// Actor with attributes and modifiers
wizard, _ := d20.NewActor("Merlin").
    WithHP(38).
    WithAC(14).
    WithAttributes(map[string]int{
        "intelligence": 18,
        "wisdom":       16,
    }).
    WithCombatModifiers(map[string]int{
        "intelligence": 4,
        "proficiency":  4,
    }).
    Build()

// Actor with rolled stats using dice notation
roller := d20.NewRandomRoller()
barbarian, _ := d20.NewActor("Grog").
    WithRoller(roller).
    WithRolledHP("10d12+30").        // Roll hit points
    WithRolledAttribute("strength", "3d6").    // Roll ability score
    WithRolledAttributes(map[string]string{       // Roll multiple abilities
        "dexterity":    "3d6",
        "constitution": "3d6",
    }).
    WithAC(14).
    WithAttribute("proficiency", 4).  // Mix rolled and fixed values
    Build()
```

**Rolled Stats**: Use `WithRoller()` to enable dice rolling during character creation. This is perfect for:
- Traditional ability score rolling (3d6, or 4d6 keep highest 3 as a variant)
- Random HP generation
- Variable starting stats
- Quick NPC generation

All rolled methods accept dice notation strings like `"3d6"`, `"2d8+3"`, `"1d20+5"`, etc.

### Combat Modifiers

Combat modifiers apply to an actor's attack rolls. Add them during construction or modify at runtime:

```go
actor.AddCombatModifier("strength", 4)
actor.AddCombatModifier("proficiency", 3)
actor.AddCombatModifier("magic_weapon", 1)
actor.RemoveCombatModifier("magic_weapon")
```

Common combat modifiers include:
- **Ability Modifiers**: Strength for melee, Dexterity for ranged/finesse weapons
- **Proficiency Bonus**: If proficient with the weapon being used
- **Equipment Bonuses**: Magic weapon bonuses (+1, +2, +3 weapons)
- **Spell Effects**: Bless, Guidance, or other temporary bonuses
- **Class Features**: Fighting styles, rage bonuses, etc.

### Attributes

The flexible attribute system supports standard D&D 5e ability scores and derived statistics:

- **Core Abilities**: `strength`, `dexterity`, `constitution`, `intelligence`, `wisdom`, `charisma`
- **Skills**: `athletics`, `stealth`, `perception`, `insight`, etc.
- **Custom Attributes**: Any string key with integer value

### Actor Roll Methods

Actor roll methods return `*RollBuilder` for flexible configuration:

```go
// SkillCheck - Returns a RollBuilder configured with the skill modifier
builder, err := actor.SkillCheck("stealth", roller)
if err != nil {
    // skill not found in attributes
}
result, _ := builder.Roll()                    // Normal roll
result, _ := builder.WithAdvantage().Roll()    // With advantage
result, _ := builder.WithDisadvantage().Roll() // With disadvantage

// AttackRoll - Returns a RollBuilder with all combat modifiers applied
builder := actor.AttackRoll(roller)
result, _ := builder.Roll()                             // Normal attack
result, _ := builder.WithAdvantage().Roll()             // Attack with advantage
result, _ := builder.WithModifier("bless", 1).Roll()    // Add temporary modifier

// D100SkillCheck - Percentile system (Call of Cthulhu, etc.)
success, outcome, err := actor.D100SkillCheck("stealth", roller)
```

#### Advantage/Disadvantage Mechanics

Advantage and disadvantage are configured on the `RollBuilder`:

- **Advantage**: Rolls 2 dice, uses higher, shows both in `DiceRolls: [6, 8]`
- **Disadvantage**: Rolls 2 dice, uses lower, shows both in `DiceRolls: [6, 8]`
- **Normal**: Rolls standard number of dice

This system is core to 5e and applies to attack rolls, skill checks, and saving throws. The library returns all dice rolled for transparency.

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

Skill checks use the actor's attribute values and return configurable `RollBuilder`:

```go
// D&D 5e skill checks (d20 + modifiers)
builder, _ := actor.SkillCheck("stealth", roller)
result, _ := builder.Roll()  // Normal check

builder, _ = actor.SkillCheck("athletics", roller)
result, _ = builder.WithAdvantage().Roll()  // With advantage

builder, _ = actor.SkillCheck("perception", roller)
result, _ = builder.WithDisadvantage().Roll()  // With disadvantage

// Call of Cthulhu skill checks (d100, roll under skill value)
success, outcome, _ := investigator.D100SkillCheck("stealth", roller)

// Check success for d100 systems
if success {
    skillValue, _ := investigator.Attribute("stealth")
    fmt.Printf("Skill check succeeded with %d (needed ≤ %d)", outcome.Value, skillValue)
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

// Dice notation shorthand - quick and simple
result, _ := roller.Roll("1d20")
fmt.Printf("Rolled: %d\n", result.Value)

result, _ = roller.Roll("d20+5")
fmt.Printf("With modifier: %d\n", result.Value)

result, _ = roller.Roll("2d6+3")
fmt.Printf("Damage: %d\n", result.Value)

// Fluent API - more control and options
result, _ = roller.Dice(1, 20).Roll()
fmt.Printf("Rolled: %d\n", result.Value)

// Attack roll with modifiers
result, _ = roller.Dice(1, 20).
    WithModifier("strength", 5).
    WithModifier("proficiency", 3).
    Roll()
fmt.Printf("Attack: %d\n", result.Value)

// Using a map for multiple modifiers
result, _ = roller.Dice(1, 20).
    WithModifiers(map[string]int{
        "strength":    5,
        "proficiency": 3,
    }).
    Roll()

// Damage roll
result, _ = roller.Dice(1, 8).
    WithModifier("strength", 3).
    Roll()
fmt.Printf("Damage: %d\n", result.Value)

// Roll with advantage (shows all dice)
result, _ = roller.Dice(1, 20).
    WithAdvantage().
    WithModifier("dexterity", 4).
    Roll()
fmt.Printf("Roll: %d (from dice: %v)\n", result.Value, result.DiceRolls)
// Output: Roll: 18 (from dice: [15, 18])
```

### Actor Usage

```go
roller := d20.NewRoller(time.Now().UnixNano())

// Create a character using the builder pattern
fighter, _ := d20.NewActor("Ironpants").
    WithHP(45).
    WithAC(18).
    WithAttributes(map[string]int{
        "strength":     16,
        "dexterity":    14,
        "constitution": 15,
        "athletics":    5,  // includes proficiency
        "stealth":      2,  // dex modifier only
    }).
    WithCombatModifiers(map[string]int{
        "strength":     3,
        "proficiency":  3,
        "magic_weapon": 1,
    }).
    Build()

// Or create a simple actor and modify it
wizard, _ := d20.NewActor("Merlin").
    WithHP(22).
    WithAC(12).
    Build()

wizard.SetAttribute("intelligence", 18)
wizard.SetAttribute("wisdom", 14)
wizard.AddCombatModifier("intelligence", 4)

// Create a character with rolled stats
barbarian, _ := d20.NewActor("Grog").
    WithRoller(roller).
    WithRolledHP("12d12+48").        // Roll for hit points
    WithRolledAttributes(map[string]string{
        "strength":     "3d6",       // Roll ability scores
        "dexterity":    "3d6",
        "constitution": "3d6",
        "intelligence": "3d6",       // Dump stat
        "wisdom":       "3d6",
        "charisma":     "3d6",
    }).
    WithAC(14).
    WithAttribute("proficiency", 5). // Mix rolled and fixed
    Build()

// Make an attack roll using the actor's combat modifiers
result, _ := fighter.AttackRoll(roller).Roll()
fmt.Printf("Attack: %d\n", result.Value)

// Attack with advantage (flanking, help action, etc.)
result, _ = fighter.AttackRoll(roller).WithAdvantage().Roll()
fmt.Printf("Attack with advantage: %d (dice: %v)\n", result.Value, result.DiceRolls)

// Attack with additional situational modifiers
result, _ = fighter.AttackRoll(roller).
    WithModifier("bless", 1).
    WithModifier("cover_penalty", -2).
    Roll()

// Perform skill checks
builder, _ := fighter.SkillCheck("stealth", roller)
result, _ = builder.Roll()

builder, _ = fighter.SkillCheck("athletics", roller)
result, _ = builder.WithAdvantage().Roll()

// HP management
fighter.SubHP(15)
if !fighter.IsKnockedOut() {
    fmt.Printf("Fighter has %d/%d HP remaining\n", fighter.HP(), fighter.MaxHP())
}
fighter.AddHP(8)

// Level up - increase max HP
fighter.SetMaxHP(50)
fighter.ResetHP() // Full heal after rest

// Attribute changes
fighter.IncrementAttribute("strength", 2) 
fighter.DecrementAttribute("dexterity", 1) 
```

### D100 System Usage

```go
// Create a Call of Cthulhu investigator
investigator, _ := d20.NewActor("Detective Morgan").
    WithHP(12).
    WithAC(10).
    WithAttributes(map[string]int{
        "stealth":     45,  // 45% skill
        "fighting":    60,  // 60% skill  
        "firearms":    25,  // 25% skill
        "spot_hidden": 70,  // 70% skill
        "sanity":      65,  // Current sanity points
    }).
    Build()

// Attribute changes - sanity loss and recovery
investigator.IncrementAttribute("sanity", 1)  // Therapy or rest
investigator.DecrementAttribute("sanity", 3)  // Witnessed something horrifying     

// Perform d100 skill checks
success, outcome, _ := investigator.D100SkillCheck("stealth", roller)
if success {
    skillValue, _ := investigator.Attribute("stealth")
    fmt.Printf("Stealth succeeded: rolled %d ≤ %d\n", outcome.Value, skillValue)
} else {
    fmt.Printf("Stealth failed: rolled %d\n", outcome.Value)
}

// Combat using Fighting skill
success, outcome, _ = investigator.D100SkillCheck("fighting", roller)
```

## Character Creation Workflows

The library supports multiple character creation methods:

### Traditional Rolled Stats (3d6, etc.)

```go
roller := d20.NewRandomRoller()

// Classic 3d6 in order
fighter, _ := d20.NewActor("Hrothgar").
    WithRoller(roller).
    WithRolledAttributes(map[string]string{
        "strength":     "3d6",
        "dexterity":    "3d6",
        "constitution": "3d6",
        "intelligence": "3d6",
        "wisdom":       "3d6",
        "charisma":     "3d6",
    }).
    WithRolledHP("1d10").
    WithAC(16).
    Build()

// Alternative: Roll 4d6 per stat (simulates rolling 4 and dropping lowest mentally)
wizard, _ := d20.NewActor("Gandalf").
    WithRoller(roller).
    WithRolledAttributes(map[string]string{
        "strength":     "3d6",
        "dexterity":    "3d6",
        "constitution": "3d6",
        "intelligence": "3d6",
        "wisdom":       "3d6",
        "charisma":     "3d6",
    }).
    WithRolledHP("1d6+1").
    WithAC(12).
    Build()
```

### Point Buy / Standard Array (Fixed Stats)

```go
// Point buy or standard array (15, 14, 13, 12, 10, 8)
paladin, _ := d20.NewActor("Arthas").
    WithHP(44).  // Fixed HP (average per level)
    WithAC(18).
    WithAttributes(map[string]int{
        "strength":     15,
        "dexterity":    10,
        "constitution": 14,
        "intelligence": 8,
        "wisdom":       12,
        "charisma":     13,
    }).
    Build()
```

### Hybrid Approach (Mix Rolled and Fixed)

```go
// Roll HP but use point buy for stats
ranger, _ := d20.NewActor("Aragorn").
    WithRoller(roller).
    WithRolledHP("10d10+20").  // Rolled HP for excitement
    WithAC(17).
    WithAttributes(map[string]int{
        "strength":     15,     // Point buy stats
        "dexterity":    16,
        "constitution": 14,
        "intelligence": 10,
        "wisdom":       14,
        "charisma":     12,
    }).
    Build()

// Or roll only dump stats
minmaxer, _ := d20.NewActor("That Guy").
    WithRoller(roller).
    WithHP(45).
    WithAC(18).
    WithAttribute("strength", 18).     // Fixed primary stat
    WithAttribute("constitution", 16). // Fixed important stat
    WithRolledAttribute("intelligence", "3d6"). // Roll dump stat for fun
    WithRolledAttribute("charisma", "3d6").     // Roll dump stat for fun
    Build()
```

## Future Enhancements

### Advanced Dice Notation
Currently, the dice notation parser supports basic formats like `"1d20"`, `"2d6+3"`, and `"3d8-2"`. 

**Planned additions:**
- `kh` (keep highest): `"4d6kh3"` - Roll 4d6, keep highest 3 (common for D&D ability scores)
- `kl` (keep lowest): `"4d6kl3"` - Roll 4d6, keep lowest 3
- `dh` (drop highest): `"4d6dh1"` - Roll 4d6, drop highest 1
- `dl` (drop lowest): `"4d6dl1"` - Roll 4d6, drop lowest 1 (equivalent to kh3)

This would allow more expressive character generation methods while maintaining compatibility with the existing API.

## References

- [D&D 5th Edition System Reference Document](https://dnd.wizards.com/resources/systems-reference-document)
- [5e SRD CC-BY License](https://creativecommons.org/licenses/by/4.0/)


