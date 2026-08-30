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

// Example_diceExpr shows building a roll from notation, then chaining fluent modifiers.
func Example_diceExpr() {
	roller := d20.NewRoller(42)
	result, _ := roller.DiceExpr("2d6+3").
		WithModifier("strength", 2).
		Roll()

	fmt.Printf("Roll: %d\n", result.Value)
	fmt.Println(result.Detail())
	// Output:
	// Roll: 17
	// Rolled 2d6... 6, 6; +3 modifier, +2 strength; *Result: 17*
}

// Example_parseNotation shows inspecting dice notation without rolling.
func Example_parseNotation() {
	count, faces, mods, err := d20.ParseNotation("2d6+3")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%dd%d", count, faces)
	if len(mods) > 0 {
		fmt.Printf(" %+d", mods[0].Value)
	}
	fmt.Println()
	// Output:
	// 2d6 +3
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
	// Roll: 8, Dice: [{20 6} {20 8}]
}

// Example_disadvantage shows rolling with disadvantage (2 dice, take lower).
func Example_disadvantage() {
	roller := d20.NewRoller(42)
	result, _ := roller.Dice(1, 20).WithDisadvantage().Roll()

	fmt.Printf("Rolled: %d (from %v)\n", result.Value, result.DiceRolls)
	// Output:
	// Rolled: 6 (from [{20 6} {20 8}])
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
	// Total: 15 (rolls: [{6 6} {6 6} {6 3}])
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

// Example_newActor shows creating an actor and setting fields.
func Example_newActor() {
	actor := d20.NewActor("Aragorn")
	actor.MaxHP, actor.HP = 45, 45
	actor.AC = 18

	fmt.Printf("ID: %s\n", actor.ID)
	fmt.Printf("HP: %d/%d\n", actor.HP, actor.MaxHP)
	fmt.Printf("AC: %d\n", actor.AC)
	// Output:
	// ID: aragorn
	// HP: 45/45
	// AC: 18
}

// Example_actorFields shows setting attributes and combat modifiers.
func Example_actorFields() {
	actor := d20.NewActor("Fighter")
	actor.MaxHP, actor.HP = 50, 50
	actor.AC = 18
	actor.Attributes["strength"] = 16
	actor.Attributes["dexterity"] = 14
	actor.CombatModifiers = []d20.Modifier{
		d20.NewModifier("strength", 3),
		d20.NewModifier("proficiency", 2),
	}

	fmt.Printf("HP: %d\n", actor.MaxHP)
	fmt.Printf("Strength: %d\n", actor.Attributes["strength"])
	// Output:
	// HP: 50
	// Strength: 16
}

// Example_actorSkillCheck shows making a skill check.
func Example_actorSkillCheck() {
	roller := d20.NewRoller(42)
	actor := d20.NewActor("Rogue")
	actor.MaxHP, actor.HP = 30, 30
	actor.AC = 15
	actor.Attributes["dexterity"] = 18

	builder, _ := actor.SkillCheck("dexterity", roller)
	result, _ := builder.Roll()

	fmt.Printf("Skill check: %d\n", result.Value)
	// Output:
	// Skill check: 24
}

// Example_actorSkillCheckAdvantage shows a skill check with advantage.
func Example_actorSkillCheckAdvantage() {
	roller := d20.NewRoller(42)
	actor := d20.NewActor("Bard")
	actor.MaxHP, actor.HP = 38, 38
	actor.AC = 14
	actor.Attributes["charisma"] = 16

	builder, _ := actor.SkillCheck("charisma", roller)
	result, _ := builder.WithAdvantage().Roll()

	fmt.Printf("Check: %d\n", result.Value)
	// Output:
	// Check: 24
}

// Example_actorAttackRoll shows making an attack roll.
func Example_actorAttackRoll() {
	roller := d20.NewRoller(42)
	actor := d20.NewActor("Fighter")
	actor.MaxHP, actor.HP = 45, 45
	actor.AC = 18
	actor.CombatModifiers = []d20.Modifier{
		d20.NewModifier("strength", 4),
		d20.NewModifier("proficiency", 3),
	}

	result, _ := actor.AttackRoll(roller).Roll()

	fmt.Printf("Attack: %d\n", result.Value)
	// Output:
	// Attack: 13
}

// Example_actorAttackAdvantage shows an attack with advantage.
func Example_actorAttackAdvantage() {
	roller := d20.NewRoller(42)
	actor := d20.NewActor("Barbarian")
	actor.MaxHP, actor.HP = 52, 52
	actor.AC = 15
	actor.CombatModifiers = []d20.Modifier{
		d20.NewModifier("strength", 5),
		d20.NewModifier("proficiency", 3),
	}

	result, _ := actor.AttackRoll(roller).WithAdvantage().Roll()

	fmt.Printf("Attack: %d\n", result.Value)
	// Output:
	// Attack: 16
}

// Example_hpManagement shows managing actor hit points.
func Example_hpManagement() {
	actor := d20.NewActor("Cleric")
	actor.MaxHP, actor.HP = 38, 38
	actor.AC = 16

	fmt.Printf("Start: %d/%d\n", actor.HP, actor.MaxHP)

	actor.SubHP(15)
	fmt.Printf("After damage: %d/%d\n", actor.HP, actor.MaxHP)

	actor.AddHP(10)
	fmt.Printf("After healing: %d/%d\n", actor.HP, actor.MaxHP)

	actor.SubHP(50)
	fmt.Printf("Knocked out: %v\n", actor.IsKnockedOut())

	actor.ResetHP()
	fmt.Printf("After rest: %d/%d\n", actor.HP, actor.MaxHP)

	// Output:
	// Start: 38/38
	// After damage: 23/38
	// After healing: 33/38
	// Knocked out: true
	// After rest: 38/38
}

// Example_attributes shows managing actor attributes.
func Example_attributes() {
	actor := d20.NewActor("Wizard")
	actor.MaxHP, actor.HP = 28, 28
	actor.AC = 12
	actor.Attributes["intelligence"] = 18

	intel, exists := actor.Attributes["intelligence"]
	fmt.Printf("Intelligence: %d (exists: %v)\n", intel, exists)

	actor.Attributes["wisdom"] = 14
	fmt.Printf("Wisdom: %d\n", actor.Attributes["wisdom"])

	actor.Attributes["intelligence"] += 2
	fmt.Printf("Intelligence buffed: %d\n", actor.Attributes["intelligence"])

	actor.Attributes["wisdom"] -= 2
	fmt.Printf("Wisdom debuffed: %d\n", actor.Attributes["wisdom"])

	// Output:
	// Intelligence: 18 (exists: true)
	// Wisdom: 14
	// Intelligence buffed: 20
	// Wisdom debuffed: 12
}

// Example_combatModifiers shows managing combat modifiers.
func Example_combatModifiers() {
	roller := d20.NewRoller(42)
	actor := d20.NewActor("Paladin")
	actor.MaxHP, actor.HP = 42, 42
	actor.AC = 18
	actor.CombatModifiers = []d20.Modifier{
		d20.NewModifier("strength", 4),
		d20.NewModifier("proficiency", 3),
	}

	result, _ := actor.AttackRoll(roller).Roll()
	fmt.Printf("Normal: %d\n", result.Value)

	actor.CombatModifiers = append(actor.CombatModifiers, d20.NewModifier("bless", 1))
	result, _ = actor.AttackRoll(roller).Roll()
	fmt.Printf("With bless: %d\n", result.Value)

	actor.CombatModifiers = actor.CombatModifiers[:2]
	result, _ = actor.AttackRoll(roller).Roll()
	fmt.Printf("After bless: %d\n", result.Value)

	// Output:
	// Normal: 13
	// With bless: 16
	// After bless: 16
}

// Example_idNormalization shows automatic ID normalization to snake_case.
func Example_idNormalization() {
	fmt.Println(d20.NewActor("Simple").ID)
	fmt.Println(d20.NewActor("UPPERCASE").ID)
	fmt.Println(d20.NewActor("Mixed Case Name").ID)
	fmt.Println(d20.NewActor("Goblin-#3").ID)

	// Output:
	// simple
	// uppercase
	// mixed_case_name
	// goblin_3
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

// Example_withAttributesMap shows assigning attributes from a map.
func Example_withAttributesMap() {
	actor := d20.NewActor("Fighter")
	actor.MaxHP, actor.HP = 50, 50
	actor.AC = 18
	actor.Attributes = map[string]int{
		"strength":     16,
		"dexterity":    14,
		"constitution": 15,
	}

	fmt.Printf("STR: %d, DEX: %d\n", actor.Attributes["strength"], actor.Attributes["dexterity"])
	// Output:
	// STR: 16, DEX: 14
}

// Example_combatModifiersSlice shows setting combat modifiers as a slice.
func Example_combatModifiersSlice() {
	roller := d20.NewRoller(42)
	actor := d20.NewActor("Fighter")
	actor.MaxHP, actor.HP = 50, 50
	actor.AC = 18
	actor.CombatModifiers = []d20.Modifier{
		d20.NewModifier("strength", 4),
		d20.NewModifier("proficiency", 3),
	}

	result, _ := actor.AttackRoll(roller).Roll()
	fmt.Printf("Attack: %d\n", result.Value)
	// Output:
	// Attack: 13
}

// Example_rolledStats shows creating an actor with rolled stats.
func Example_rolledStats() {
	roller := d20.NewRoller(42)
	hp, _ := roller.Roll("10d10+30")
	str, _ := roller.Roll("3d6+1")

	actor := d20.NewActor("Conan")
	actor.MaxHP, actor.HP = hp.Value, hp.Value
	actor.AC = 14
	actor.Attributes["strength"] = str.Value

	fmt.Printf("HP: %d\n", actor.MaxHP)
	fmt.Printf("Strength: %d\n", actor.Attributes["strength"])
	// Output:
	// HP: 92
	// Strength: 13
}

// Example_rolledCharacterCreation shows creating a complete character with rolled stats.
func Example_rolledCharacterCreation() {
	roller := d20.NewRoller(42)

	hp, _ := roller.Roll("10d10+20")
	actor := d20.NewActor("Thorin")
	actor.MaxHP, actor.HP = hp.Value, hp.Value
	actor.AC = 18
	actor.Attributes["proficiency"] = 4

	// Roll abilities in a stable order (sorted keys)
	for _, key := range []string{"charisma", "constitution", "dexterity", "intelligence", "strength", "wisdom"} {
		out, _ := roller.Roll("3d6")
		actor.Attributes[key] = out.Value
	}

	fmt.Printf("HP: %d\n", actor.MaxHP)
	fmt.Printf("STR: %d\n", actor.Attributes["strength"])
	fmt.Printf("DEX: %d\n", actor.Attributes["dexterity"])
	// Output:
	// HP: 82
	// STR: 14
	// DEX: 15
}

// Example_mixedStaticAndRolled shows combining fixed values with rolled stats.
func Example_mixedStaticAndRolled() {
	roller := d20.NewRoller(42)

	hp, _ := roller.Roll("8d10+24")
	cha, _ := roller.Roll("3d6")

	actor := d20.NewActor("Gimli")
	actor.MaxHP, actor.HP = hp.Value, hp.Value
	actor.AC = 18
	actor.Attributes["strength"] = 16
	actor.Attributes["constitution"] = 16
	actor.Attributes["charisma"] = cha.Value

	fmt.Printf("HP: %d\n", actor.MaxHP)
	fmt.Printf("STR: %d, CHA: %d\n", actor.Attributes["strength"], actor.Attributes["charisma"])
	// Output:
	// HP: 73
	// STR: 16, CHA: 7
}

// Example_rolledCombatStats shows rolling combat-related values.
func Example_rolledCombatStats() {
	roller := d20.NewRoller(42)

	barbarian := d20.NewActor("Grog")
	barbarian.MaxHP, barbarian.HP = 95, 95
	barbarian.AC = 14
	barbarian.Attributes["strength"] = 18

	rage, _ := roller.Dice(1, 4).Roll()
	barbarian.CombatModifiers = append(barbarian.CombatModifiers, d20.NewModifier("rage", rage.Value))

	result, _ := barbarian.AttackRoll(roller).Roll()
	fmt.Printf("Raging attack includes +%d rage: %d total\n", rage.Value, result.Value)
	// Output:
	// Raging attack includes +2 rage: 10 total
}
