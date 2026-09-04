package d20_test

import (
	"fmt"
	"log"

	"github.com/jwebster45206/d20"
)

func Example_basicRoll() {
	roller := d20.NewRoller(42)
	d, err := d20.NewDice(1, 20)
	if err != nil {
		log.Fatal(err)
	}
	result, err := roller.Roll(d)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Rolled: %d\n", result.Value)
	// Output:
	// Rolled: 18
}

func Example_diceExpr() {
	roller := d20.NewRoller(42)
	d, err := d20.DiceFromExpr("2d6+3")
	if err != nil {
		log.Fatal(err)
	}
	result, err := roller.Roll(d.WithModifier("strength", 2))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Roll: %d\n", result.Value)
	fmt.Println(result.Detail())
	// Output:
	// Roll: 17
	// Rolled 2d6... 6, 6; +3 modifier, +2 strength; *Result: 17*
}

func Example_parseDiceNotation() {
	d, err := d20.DiceFromExpr("2d6+3")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%dd%d", d.Count, d.Faces)
	if len(d.Modifiers) > 0 {
		fmt.Printf(" %+d", d.Modifiers[0].Value)
	}
	fmt.Println()
	// Output:
	// 2d6 +3
}

func Example_diceNotation() {
	roller := d20.NewRoller(42)

	result, err := roller.RollExpr("1d20")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("1d20: %d\n", result.Value)

	result, err = roller.RollExpr("d20")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("d20: %d\n", result.Value)

	result, err = roller.RollExpr("1d20+3")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("1d20+3: %d\n", result.Value)

	result, err = roller.RollExpr("2d6+2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("2d6+2: %d\n", result.Value)

	// Output:
	// 1d20: 18
	// d20: 20
	// 1d20+3: 6
	// 2d6+2: 5
}

func Example_rollWithModifier() {
	roller := d20.NewRoller(42)
	d, err := d20.NewDice(1, 20)
	if err != nil {
		log.Fatal(err)
	}
	result, err := roller.Roll(d.WithModifier("strength", 3))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Roll: %d\n", result.Value)
	// Output:
	// Roll: 21
}

func Example_rollWithMultipleModifiers() {
	roller := d20.NewRoller(42)
	d, err := d20.NewDice(1, 20)
	if err != nil {
		log.Fatal(err)
	}
	result, err := roller.Roll(d.
		WithModifier("strength", 3).
		WithModifier("cover", -2))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Roll: %d\n", result.Value)
	fmt.Println(result.Detail())
	// Output:
	// Roll: 19
	// Rolled 1d20... 18; +3 strength, -2 cover; *Result: 19*
}

func Example_rollWithAdvantage() {
	roller := d20.NewRoller(42)
	d, err := d20.NewDice(1, 20)
	if err != nil {
		log.Fatal(err)
	}
	result, err := roller.Roll(d.WithAdvantage())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Roll: %d, Dice: %v\n", result.Value, result.DiceRolls)
	// Output:
	// Roll: 20, Dice: [{20 18} {20 20}]
}

func Example_disadvantage() {
	roller := d20.NewRoller(42)
	d, err := d20.NewDice(1, 20)
	if err != nil {
		log.Fatal(err)
	}
	result, err := roller.Roll(d.WithDisadvantage())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Rolled: %d (from %v)\n", result.Value, result.DiceRolls)
	// Output:
	// Rolled: 18 (from [{20 18} {20 20}])
}

func Example_advantageWithModifier() {
	roller := d20.NewRoller(42)
	d, err := d20.NewDice(1, 20)
	if err != nil {
		log.Fatal(err)
	}
	result, err := roller.Roll(d.
		WithAdvantage().
		WithModifier("dexterity", 4))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Roll: %d\n", result.Value)
	// Output:
	// Roll: 24
}

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

func Example_actorSkillDice() {
	roller := d20.NewRoller(42)
	actor := d20.NewActor("Rogue")
	actor.MaxHP, actor.HP = 30, 30
	actor.AC = 15
	actor.Attributes["athletics"] = 5

	d, err := actor.SkillDice("athletics")
	if err != nil {
		log.Fatal(err)
	}
	result, err := roller.Roll(d)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Skill check: %d\n", result.Value)
	fmt.Println(result.Detail())
	// Output:
	// Skill check: 23
	// Rolled 1d20... 18; +5 athletics; *Result: 23*
}

func Example_actorStrikeDice() {
	roller := d20.NewRoller(42)
	actor := d20.NewActor("Fighter")
	actor.MaxHP, actor.HP = 45, 45
	actor.AC = 18
	actor.Modifiers["strength"] = 4
	actor.Modifiers["striking"] = 3

	d, err := actor.StrikeDice("strength", "striking")
	if err != nil {
		log.Fatal(err)
	}
	result, err := roller.Roll(d)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Strike: %d\n", result.Value)
	// Output:
	// Strike: 25
}

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

// Bless is on the returned Dice only; StrikeDice itself is unchanged.
func Example_strikeModifiers() {
	roller := d20.NewRoller(42)
	actor := d20.NewActor("Paladin")
	actor.MaxHP, actor.HP = 42, 42
	actor.AC = 18
	actor.Modifiers["strength"] = 4
	actor.Modifiers["striking"] = 3
	actor.Modifiers["damage"] = 2

	strike, err := actor.StrikeDice("strength", "striking")
	if err != nil {
		log.Fatal(err)
	}
	result, err := roller.Roll(strike)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Detail())

	result, err = roller.Roll(strike.WithModifier("bless", 1))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Detail())

	result, err = roller.Roll(strike)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Detail())

	// Output:
	// Rolled 1d20... 18; +4 strength, +3 striking; *Result: 25*
	// Rolled 1d20... 20; +4 strength, +3 striking, +1 bless; *Result: 28*
	// Rolled 1d20... 3; +4 strength, +3 striking; *Result: 10*
}

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

func Example_rolledStats() {
	roller := d20.NewRoller(42)
	hp, err := roller.RollExpr("10d10+30")
	if err != nil {
		log.Fatal(err)
	}
	str, err := roller.RollExpr("3d6+1")
	if err != nil {
		log.Fatal(err)
	}

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
