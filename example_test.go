package d20_test

import (
	"fmt"

	"github.com/jwebster45206/d20"
)

// Example_basicRoll demonstrates a simple dice roll.
func Example_basicRoll() {
	roller := d20.NewRoller(42)
	result, _ := roller.Dice(1, 20).Roll()

	fmt.Printf("Rolled: %d\n", result.Value)
	// Output:
	// Rolled: 6
}

// Example_diceNotation shows using dice notation shorthand.
func Example_diceNotation() {
	roller := d20.NewRoller(42)

	// Simple notation
	result, _ := roller.Roll("1d20")
	fmt.Printf("1d20: %d\n", result.Value)

	// Shorthand (assumes 1d)
	result, _ = roller.Roll("d20")
	fmt.Printf("d20: %d\n", result.Value)

	// With modifier
	result, _ = roller.Roll("1d20+3")
	fmt.Printf("1d20+3: %d\n", result.Value)

	// Multiple dice
	result, _ = roller.Roll("2d6+2")
	fmt.Printf("2d6+2: %d\n", result.Value)

	// Output:
	// 1d20: 6
	// d20: 8
	// 1d20+3: 12
	// 2d6+2: 5
}

// Example_rollWithModifier shows adding a single modifier.
func Example_rollWithModifier() {
	roller := d20.NewRoller(42)
	result, _ := roller.Dice(1, 20).
		WithModifier("strength", 3).
		Roll()

	fmt.Printf("Roll: %d\n", result.Value)
	// Output:
	// Roll: 9
}

// Example_rollWithMultipleModifiers shows adding multiple modifiers.
func Example_rollWithMultipleModifiers() {
	roller := d20.NewRoller(42)
	result, _ := roller.Dice(1, 20).
		WithModifiers(map[string]int{
			"strength": 3,
			"magic":    2,
		}).
		Roll()

	fmt.Printf("Roll: %d\n", result.Value)
	// Output:
	// Roll: 11
}

// Example_withModifiersMap shows adding modifiers from a map.
func Example_withModifiersMap() {
	roller := d20.NewRoller(42)
	mods := map[string]int{
		"strength":    3,
		"proficiency": 2,
	}
	result, _ := roller.Dice(1, 20).WithModifiers(mods).Roll()

	fmt.Printf("Roll: %d\n", result.Value)
	// Output:
	// Roll: 11
}

// Example_rollWithAdvantage shows rolling with advantage.
func Example_rollWithAdvantage() {
	roller := d20.NewRoller(42)
	result, _ := roller.Dice(1, 20).
		WithAdvantage().
		Roll()

	fmt.Printf("Roll: %d, Dice: %v\n", result.Value, result.DiceRolls)
	// Output:
	// Roll: 8, Dice: [6 8]
}

// Example_disadvantage shows rolling with disadvantage (2 dice, take lower).
func Example_disadvantage() {
	roller := d20.NewRoller(42)
	result, _ := roller.Dice(1, 20).WithDisadvantage().Roll()

	fmt.Printf("Rolled: %d (from %v)\n", result.Value, result.DiceRolls)
	// Output:
	// Rolled: 6 (from [6 8])
}

// Example_advantageWithModifier shows combining advantage with modifiers.
func Example_advantageWithModifier() {
	roller := d20.NewRoller(42)
	result, _ := roller.Dice(1, 20).
		WithAdvantage().
		WithModifier("dexterity", 4).
		Roll()

	fmt.Printf("Roll: %d\n", result.Value)
	// Output:
	// Roll: 12
}

// Example_multipleDice shows rolling multiple dice.
func Example_multipleDice() {
	roller := d20.NewRoller(42)
	result, _ := roller.Dice(3, 6).Roll()

	fmt.Printf("Total: %d (rolls: %v)\n", result.Value, result.DiceRolls)
	// Output:
	// Total: 15 (rolls: [6 6 3])
}

// Example_damageDice shows a typical damage roll.
func Example_damageDice() {
	roller := d20.NewRoller(42)
	result, _ := roller.Dice(2, 6).
		WithModifier("strength", 3).
		Roll()

	fmt.Printf("Damage: %d\n", result.Value)
	// Output:
	// Damage: 15
}

// Example_newActor shows creating an actor with the builder pattern.
func Example_newActor() {
	actor, _ := d20.NewActor("Aragorn").
		WithHP(45).
		WithAC(18).
		Build()

	fmt.Printf("ID: %s\n", actor.ID())
	fmt.Printf("HP: %d/%d\n", actor.HP(), actor.MaxHP())
	fmt.Printf("AC: %d\n", actor.AC())
	// Output:
	// ID: aragorn
	// HP: 45/45
	// AC: 18
}

// Example_actorBuilder shows building an actor with attributes and modifiers.
func Example_actorBuilder() {
	actor, _ := d20.NewActor("Fighter").
		WithHP(50).
		WithAC(18).
		WithAttribute("strength", 16).
		WithAttribute("dexterity", 14).
		WithCombatModifier("strength", 3).
		WithCombatModifier("proficiency", 2).
		Build()

	fmt.Printf("HP: %d\n", actor.MaxHP())
	str, _ := actor.Attribute("strength")
	fmt.Printf("Strength: %d\n", str)
	// Output:
	// HP: 50
	// Strength: 16
}

// Example_actorSkillCheck shows making a skill check.
func Example_actorSkillCheck() {
	roller := d20.NewRoller(42)
	actor, _ := d20.NewActor("Rogue").
		WithHP(30).
		WithAC(15).
		WithAttribute("dexterity", 18).
		Build()

	builder, _ := actor.SkillCheck("dexterity", roller)
	result, _ := builder.Roll()

	fmt.Printf("Skill check: %d\n", result.Value)
	// Output:
	// Skill check: 24
}

// Example_actorSkillCheckAdvantage shows a skill check with advantage.
func Example_actorSkillCheckAdvantage() {
	roller := d20.NewRoller(42)
	actor, _ := d20.NewActor("Bard").
		WithHP(38).
		WithAC(14).
		WithAttribute("charisma", 16).
		Build()

	builder, _ := actor.SkillCheck("charisma", roller)
	result, _ := builder.WithAdvantage().Roll()

	fmt.Printf("Check: %d\n", result.Value)
	// Output:
	// Check: 24
}

// Example_actorAttackRoll shows making an attack roll.
func Example_actorAttackRoll() {
	roller := d20.NewRoller(42)
	actor, _ := d20.NewActor("Fighter").
		WithHP(45).
		WithAC(18).
		WithCombatModifier("strength", 4).
		WithCombatModifier("proficiency", 3).
		Build()

	result, _ := actor.AttackRoll(roller).Roll()

	fmt.Printf("Attack: %d\n", result.Value)
	// Output:
	// Attack: 13
}

// Example_actorAttackAdvantage shows an attack with advantage.
func Example_actorAttackAdvantage() {
	roller := d20.NewRoller(42)
	actor, _ := d20.NewActor("Barbarian").
		WithHP(52).
		WithAC(15).
		WithCombatModifier("strength", 5).
		WithCombatModifier("proficiency", 3).
		Build()

	result, _ := actor.AttackRoll(roller).WithAdvantage().Roll()

	fmt.Printf("Attack: %d\n", result.Value)
	// Output:
	// Attack: 16
}

// Example_hpManagement shows managing actor hit points.
func Example_hpManagement() {
	actor, _ := d20.NewActor("Cleric").
		WithHP(38).
		WithAC(16).
		Build()

	fmt.Printf("Start: %d/%d\n", actor.HP(), actor.MaxHP())

	actor.SubHP(15)
	fmt.Printf("After damage: %d/%d\n", actor.HP(), actor.MaxHP())

	actor.AddHP(10)
	fmt.Printf("After healing: %d/%d\n", actor.HP(), actor.MaxHP())

	actor.SubHP(50)
	fmt.Printf("Knocked out: %v\n", actor.IsKnockedOut())

	actor.ResetHP()
	fmt.Printf("After rest: %d/%d\n", actor.HP(), actor.MaxHP())

	// Output:
	// Start: 38/38
	// After damage: 23/38
	// After healing: 33/38
	// Knocked out: true
	// After rest: 38/38
}

// Example_attributes shows managing actor attributes.
func Example_attributes() {
	actor, _ := d20.NewActor("Wizard").
		WithHP(28).
		WithAC(12).
		WithAttribute("intelligence", 18).
		Build()

	intel, exists := actor.Attribute("intelligence")
	fmt.Printf("Intelligence: %d (exists: %v)\n", intel, exists)

	actor.SetAttribute("wisdom", 14)
	wis, _ := actor.Attribute("wisdom")
	fmt.Printf("Wisdom: %d\n", wis)

	actor.IncrementAttribute("intelligence", 2)
	intel, _ = actor.Attribute("intelligence")
	fmt.Printf("Intelligence buffed: %d\n", intel)

	actor.DecrementAttribute("wisdom", 2)
	wis, _ = actor.Attribute("wisdom")
	fmt.Printf("Wisdom debuffed: %d\n", wis)

	// Output:
	// Intelligence: 18 (exists: true)
	// Wisdom: 14
	// Intelligence buffed: 20
	// Wisdom debuffed: 12
}

// Example_combatModifiers shows managing combat modifiers.
func Example_combatModifiers() {
	roller := d20.NewRoller(42)
	actor, _ := d20.NewActor("Paladin").
		WithHP(42).
		WithAC(18).
		WithCombatModifier("strength", 4).
		WithCombatModifier("proficiency", 3).
		Build()

	result, _ := actor.AttackRoll(roller).Roll()
	fmt.Printf("Normal: %d\n", result.Value)

	actor.AddCombatModifier("bless", 1)
	result, _ = actor.AttackRoll(roller).Roll()
	fmt.Printf("With bless: %d\n", result.Value)

	actor.RemoveCombatModifier("bless")
	result, _ = actor.AttackRoll(roller).Roll()
	fmt.Printf("After bless: %d\n", result.Value)

	// Output:
	// Normal: 13
	// With bless: 16
	// After bless: 16
}

// Example_idNormalization shows automatic ID normalization to snake_case.
func Example_idNormalization() {
	actor1, _ := d20.NewActor("Simple").
		WithHP(20).
		WithAC(15).
		Build()
	fmt.Println(actor1.ID())

	actor2, _ := d20.NewActor("UPPERCASE").
		WithHP(20).
		WithAC(15).
		Build()
	fmt.Println(actor2.ID())

	actor3, _ := d20.NewActor("Mixed Case Name").
		WithHP(20).
		WithAC(15).
		Build()
	fmt.Println(actor3.ID())

	actor4, _ := d20.NewActor("Goblin-#3").
		WithHP(20).
		WithAC(15).
		Build()
	fmt.Println(actor4.ID())

	// Output:
	// simple
	// uppercase
	// mixed_case_name
	// goblin_3
}

// Example_normal shows switching back to normal after advantage.
func Example_normal() {
	roller := d20.NewRoller(42)

	builder := roller.Dice(1, 20).WithAdvantage()
	builder = builder.Normal()
	result, _ := builder.Roll()

	fmt.Printf("Roll: %d (dice: %v)\n", result.Value, result.DiceRolls)
	// Output:
	// Roll: 6 (dice: [6])
}

// Example_negativeModifiers shows using negative modifiers as penalties.
func Example_negativeModifiers() {
	roller := d20.NewRoller(42)
	result, _ := roller.Dice(1, 20).
		WithModifier("exhaustion", -2).
		WithModifier("poison", -1).
		Roll()

	fmt.Printf("Roll: %d\n", result.Value)
	// Output:
	// Roll: 3
}

// Example_withAttributesMap shows adding multiple attributes at once.
func Example_withAttributesMap() {
	attrs := map[string]int{
		"strength":     16,
		"dexterity":    14,
		"constitution": 15,
	}
	actor, _ := d20.NewActor("Fighter").
		WithHP(50).
		WithAC(18).
		WithAttributes(attrs).
		Build()

	str, _ := actor.Attribute("strength")
	dex, _ := actor.Attribute("dexterity")
	fmt.Printf("STR: %d, DEX: %d\n", str, dex)
	// Output:
	// STR: 16, DEX: 14
}

// Example_withCombatModifiersMap shows adding multiple combat modifiers at once.
func Example_withCombatModifiersMap() {
	roller := d20.NewRoller(42)
	mods := map[string]int{
		"strength":    4,
		"proficiency": 3,
	}
	actor, _ := d20.NewActor("Fighter").
		WithHP(50).
		WithAC(18).
		WithCombatModifiers(mods).
		Build()

	result, _ := actor.AttackRoll(roller).Roll()
	fmt.Printf("Attack: %d\n", result.Value)
	// Output:
	// Attack: 13
}

// Example_withRoller shows creating an actor with rolled stats using dice notation.
func Example_withRoller() {
	roller := d20.NewRoller(42)
	actor, _ := d20.NewActor("Conan").
		WithRoller(roller).
		WithRolledHP("10d10+30").
		WithRolledAttribute("strength", "3d6+1").
		WithAC(14).
		Build()

	fmt.Printf("HP: %d\n", actor.MaxHP())
	str, _ := actor.Attribute("strength")
	fmt.Printf("Strength: %d\n", str)
	// Output:
	// HP: 92
	// Strength: 13
}

// Example_rolledCharacterCreation shows creating a complete character with rolled stats.
func Example_rolledCharacterCreation() {
	roller := d20.NewRoller(42)

	// Roll all 6 ability scores using 4d6 keep highest 3
	attrs := map[string]string{
		"strength":     "3d6",
		"dexterity":    "3d6",
		"constitution": "3d6",
		"intelligence": "3d6",
		"wisdom":       "3d6",
		"charisma":     "3d6",
	}

	fighter, _ := d20.NewActor("Thorin").
		WithRoller(roller).
		WithRolledHP("10d10+20"). // Level 10 fighter
		WithAC(18).
		WithRolledAttributes(attrs).
		WithAttribute("proficiency", 4). // Mix rolled and fixed
		Build()

	fmt.Printf("HP: %d\n", fighter.MaxHP())
	str, _ := fighter.Attribute("strength")
	fmt.Printf("STR: %d\n", str)
	dex, _ := fighter.Attribute("dexterity")
	fmt.Printf("DEX: %d\n", dex)
	// Output:
	// HP: 82
	// STR: 12
	// DEX: 12
}

// Example_mixedStaticAndRolled shows combining fixed values with rolled stats.
func Example_mixedStaticAndRolled() {
	roller := d20.NewRoller(42)

	// Point buy for abilities, rolled HP
	actor, _ := d20.NewActor("Gimli").
		WithRoller(roller).
		WithRolledHP("8d10+24"). // Level 8, +3 CON
		WithAC(18).
		WithAttribute("strength", 16).          // Fixed point buy
		WithAttribute("constitution", 16).      // Fixed point buy
		WithRolledAttribute("charisma", "3d6"). // Dump stat, rolled
		Build()

	fmt.Printf("HP: %d\n", actor.MaxHP())
	str, _ := actor.Attribute("strength")
	cha, _ := actor.Attribute("charisma")
	fmt.Printf("STR: %d, CHA: %d\n", str, cha)
	// Output:
	// HP: 73
	// STR: 16, CHA: 7
}

// Example_rolledCombatStats shows rolling combat-related values.
func Example_rolledCombatStats() {
	roller := d20.NewRoller(42)

	// Create a barbarian with variable rage damage
	barbarian, _ := d20.NewActor("Grog").
		WithRoller(roller).
		WithHP(95).
		WithAC(14).
		WithAttribute("strength", 18).
		Build()

	// Temporarily check roller usage - in game you'd roll at rage time
	builder := roller.Dice(1, 4)
	rage, _ := builder.Roll()
	barbarian.AddCombatModifier("rage", rage.Value)

	result, _ := barbarian.AttackRoll(roller).Roll()
	fmt.Printf("Raging attack includes +%d rage: %d total\n", rage.Value, result.Value)
	// Output:
	// Raging attack includes +2 rage: 10 total
}
