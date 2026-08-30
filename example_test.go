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
	// Rolled: 18
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

	result, _ := roller.Roll("1d20")
	fmt.Printf("1d20: %d\n", result.Value)

	result, _ = roller.Roll("d20")
	fmt.Printf("d20: %d\n", result.Value)

	result, _ = roller.Roll("1d20+3")
	fmt.Printf("1d20+3: %d\n", result.Value)

	result, _ = roller.Roll("2d6+2")
	fmt.Printf("2d6+2: %d\n", result.Value)

	// Output:
	// 1d20: 18
	// d20: 20
	// 1d20+3: 6
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
	// Roll: 21
}

// Example_rollWithMultipleModifiers shows adding multiple modifiers, including a penalty.
func Example_rollWithMultipleModifiers() {
	roller := d20.NewRoller(42)
	result, _ := roller.Dice(1, 20).
		WithModifier("strength", 3).
		WithModifier("cover", -2).
		Roll()

	fmt.Printf("Roll: %d\n", result.Value)
	fmt.Println(result.Detail())
	// Output:
	// Roll: 19
	// Rolled 1d20... 18; +3 strength, -2 cover; *Result: 19*
}

// Example_rollWithAdvantage shows rolling with advantage.
func Example_rollWithAdvantage() {
	roller := d20.NewRoller(42)
	result, _ := roller.Dice(1, 20).
		WithAdvantage().
		Roll()

	fmt.Printf("Roll: %d, Dice: %v\n", result.Value, result.DiceRolls)
	// Output:
	// Roll: 20, Dice: [{20 18} {20 20}]
}

// Example_disadvantage shows rolling with disadvantage (2 dice, take lower).
func Example_disadvantage() {
	roller := d20.NewRoller(42)
	result, _ := roller.Dice(1, 20).WithDisadvantage().Roll()

	fmt.Printf("Rolled: %d (from %v)\n", result.Value, result.DiceRolls)
	// Output:
	// Rolled: 18 (from [{20 18} {20 20}])
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
	// Roll: 24
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

// Example_actorSkillCheck shows making a skill check with a stored skill bonus.
func Example_actorSkillCheck() {
	roller := d20.NewRoller(42)
	actor := d20.NewActor("Rogue")
	actor.MaxHP, actor.HP = 30, 30
	actor.AC = 15
	actor.Attributes["athletics"] = 5

	builder, _ := actor.SkillCheck("athletics", roller)
	result, _ := builder.Roll()

	fmt.Printf("Skill check: %d\n", result.Value)
	fmt.Println(result.Detail())
	// Output:
	// Skill check: 23
	// Rolled 1d20... 18; +5 athletics; *Result: 23*
}

// Example_actorStrikeRoll shows making a strike roll.
func Example_actorStrikeRoll() {
	roller := d20.NewRoller(42)
	actor := d20.NewActor("Fighter")
	actor.MaxHP, actor.HP = 45, 45
	actor.AC = 18
	actor.Modifiers["strength"] = 4
	actor.Modifiers["striking"] = 3

	result, _ := actor.StrikeRoll(roller, "strength", "striking").Roll()

	fmt.Printf("Strike: %d\n", result.Value)
	// Output:
	// Strike: 25
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

// Example_strikeModifiers shows wiring modifiers and applying a subset on a strike.
// Bless is on the DiceManager only; a new StrikeRoll does not keep it.
func Example_strikeModifiers() {
	roller := d20.NewRoller(42)
	actor := d20.NewActor("Paladin")
	actor.MaxHP, actor.HP = 42, 42
	actor.AC = 18
	actor.Modifiers["strength"] = 4
	actor.Modifiers["striking"] = 3
	actor.Modifiers["damage"] = 2

	result, _ := actor.StrikeRoll(roller, "strength", "striking").Roll()
	fmt.Println(result.Detail())

	result, _ = actor.StrikeRoll(roller, "strength", "striking").WithModifier("bless", 1).Roll()
	fmt.Println(result.Detail())

	result, _ = actor.StrikeRoll(roller, "strength", "striking").Roll()
	fmt.Println(result.Detail())

	// Output:
	// Rolled 1d20... 18; +4 strength, +3 striking; *Result: 25*
	// Rolled 1d20... 20; +4 strength, +3 striking, +1 bless; *Result: 28*
	// Rolled 1d20... 3; +4 strength, +3 striking; *Result: 10*
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
	// HP: 95
	// Strength: 11
}
