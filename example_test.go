package d20_test

import (
	"fmt"

	"github.com/jwebster45206/d20"
)

// ExampleNewRoller demonstrates creating a roller with a seed for deterministic results.
func ExampleNewRoller() {
	roller := d20.NewRoller(42)
	result, _ := roller.Roll(1, 20, nil)
	fmt.Printf("Rolled: %d\n", result.Value)
	// Output: Rolled: 6
}

// ExampleNewRandomRoller demonstrates creating a roller with random seed.
func ExampleNewRandomRoller() {
	roller := d20.NewRandomRoller()
	result, _ := roller.Roll(1, 6, nil)
	// Result will vary each time
	fmt.Printf("Random d6 roll: %v\n", result.Value >= 1 && result.Value <= 6)
	// Output: Random d6 roll: true
}

// ExampleRoller_Roll demonstrates a basic dice roll.
func ExampleRoller_Roll() {
	roller := d20.NewRoller(42)
	result, _ := roller.Roll(1, 20, nil)
	fmt.Println(result.Detail)
	// Output: Rolled 1d20...; values 6; *Result: 6*
}

// ExampleRoller_Roll_withModifiers demonstrates rolling with multiple modifiers.
func ExampleRoller_Roll_withModifiers() {
	roller := d20.NewRoller(42)
	result, _ := roller.Roll(2, 6, []d20.Modifier{
		d20.NewModifier(3, "Strength"),
		d20.NewModifier(2, "Proficiency"),
	})
	fmt.Println(result.Detail)
	// Output: Rolled 2d6...; values 6, 6; +3 strength, +2 proficiency; *Result: 17*
}

// ExampleRoller_Roll_damage demonstrates a damage roll.
func ExampleRoller_Roll_damage() {
	roller := d20.NewRoller(123)
	result, _ := roller.Roll(2, 8, []d20.Modifier{
		d20.NewModifier(4, "Strength"),
	})
	fmt.Println(result.Detail)
	// Output: Rolled 2d8...; values 4, 2; +4 strength; *Result: 10*
}

// ExampleNewActor demonstrates creating a basic actor.
func ExampleNewActor() {
	fighter, err := d20.NewActor(45, 18, 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("HP: %d, AC: %d, Initiative: %d\n", fighter.HP(), fighter.AC(), fighter.Initiative())
	// Output: HP: 45, AC: 18, Initiative: 2
}

// ExampleNewActorWithAttributes demonstrates creating an actor with initial attributes.
func ExampleNewActorWithAttributes() {
	wizard, err := d20.NewActorWithAttributes(22, 12, 3, map[string]int{
		"intelligence": 18,
		"wisdom":       14,
		"arcana":       8,
		"history":      6,
	})
	if err != nil {
		panic(err)
	}

	arcana, _ := wizard.Attribute("arcana")
	fmt.Printf("Wizard HP: %d, Arcana: %d\n", wizard.HP(), arcana)
	// Output: Wizard HP: 22, Arcana: 8
}

// ExampleActor_SkillCheck demonstrates a D&D 5e skill check.
func ExampleActor_SkillCheck() {
	roller := d20.NewRoller(42)
	rogue, _ := d20.NewActor(30, 15, 3)
	rogue.SetAttribute("stealth", 9) // +4 Dex, +5 Proficiency

	result, _ := rogue.SkillCheck("stealth", roller, d20.Normal)
	fmt.Println(result.Detail)
	// Output: Rolled 1d20...; values 6; +9 stealth; *Result: 15*
}

// ExampleActor_SkillCheck_advantage demonstrates skill check with advantage.
func ExampleActor_SkillCheck_advantage() {
	roller := d20.NewRoller(100)
	barbarian, _ := d20.NewActor(60, 14, 1)
	barbarian.SetAttribute("athletics", 7) // +4 Str, +3 Proficiency

	result, _ := barbarian.SkillCheck("athletics", roller, d20.Advantage)
	fmt.Println(result.Detail)
	// Output: Rolled 1d20...; values 9; +7 athletics; *Result: 16*
}

// ExampleActor_SkillCheck_disadvantage demonstrates skill check with disadvantage.
func ExampleActor_SkillCheck_disadvantage() {
	roller := d20.NewRoller(100)
	barbarian, _ := d20.NewActor(60, 14, 1)
	barbarian.SetAttribute("stealth", 0) // +0 Dex, not proficient, heavy armor

	result, _ := barbarian.SkillCheck("stealth", roller, d20.Disadvantage)
	fmt.Println(result.Detail)
	// Output: Rolled 1d20...; values 4; +0 stealth; *Result: 4*
}

// ExampleActor_D100SkillCheck demonstrates a Call of Cthulhu percentile skill check.
func ExampleActor_D100SkillCheck() {
	roller := d20.NewRoller(10)
	investigator, _ := d20.NewActorWithAttributes(12, 10, 0, map[string]int{
		"stealth": 45,
	})

	success, result, _ := investigator.D100SkillCheck("stealth", roller, 0)
	fmt.Printf("Success: %v, Roll: %d\n", success, result.Value)
	// Output: Success: false, Roll: 48
}

// ExampleActor_D100SkillCheck_bonusDie demonstrates Call of Cthulhu bonus die.
func ExampleActor_D100SkillCheck_bonusDie() {
	roller := d20.NewRoller(50)
	investigator, _ := d20.NewActorWithAttributes(12, 10, 0, map[string]int{
		"fighting": 60,
	})

	success, result, _ := investigator.D100SkillCheck("fighting", roller, 1)
	fmt.Printf("Success: %v, Roll: %d\n", success, result.Value)
	// Output: Success: true, Roll: 17
}

// ExampleActor_D100SkillCheck_penaltyDie demonstrates Call of Cthulhu penalty die.
func ExampleActor_D100SkillCheck_penaltyDie() {
	roller := d20.NewRoller(50)
	investigator, _ := d20.NewActorWithAttributes(12, 10, 0, map[string]int{
		"spot_hidden": 70,
	})

	success, result, _ := investigator.D100SkillCheck("spot_hidden", roller, -1)
	fmt.Printf("Success: %v, Roll: %d\n", success, result.Value)
	// Output: Success: true, Roll: 47
}

// ExampleActor_AttackRoll demonstrates a basic attack roll.
func ExampleActor_AttackRoll() {
	roller := d20.NewRoller(42)
	fighter, _ := d20.NewActor(45, 18, 2)
	fighter.AddCombatModifier(d20.NewModifier(3, "Strength"))
	fighter.AddCombatModifier(d20.NewModifier(3, "Proficiency"))

	result, _ := fighter.AttackRoll(roller, d20.Normal)
	fmt.Println(result.Detail)
	// Output: Rolled 1d20...; values 6; +3 strength, +3 proficiency; *Result: 12*
}

// ExampleActor_AttackRoll_advantage demonstrates attack with advantage.
func ExampleActor_AttackRoll_advantage() {
	roller := d20.NewRoller(100)
	ranger, _ := d20.NewActor(40, 16, 4)
	ranger.AddCombatModifier(d20.NewModifier(4, "Dexterity"))
	ranger.AddCombatModifier(d20.NewModifier(3, "Proficiency"))

	result, _ := ranger.AttackRoll(roller, d20.Advantage)
	fmt.Println(result.Detail)
	// Output: Rolled 1d20...; values 9; +4 dexterity, +3 proficiency; *Result: 16*
}

// ExampleActor_AttackRollWithModifiers demonstrates attack with situational modifiers.
func ExampleActor_AttackRollWithModifiers() {
	roller := d20.NewRoller(42)
	paladin, _ := d20.NewActor(50, 20, 1)
	paladin.AddCombatModifier(d20.NewModifier(4, "Strength"))
	paladin.AddCombatModifier(d20.NewModifier(4, "Proficiency"))

	result, _ := paladin.AttackRollWithModifiers(roller, d20.Normal, []d20.Modifier{
		d20.NewModifier(2, "Bless"),
		d20.NewModifier(1, "Magic Weapon"),
	})
	fmt.Println(result.Detail)
	// Output: Rolled 1d20...; values 6; +4 strength, +4 proficiency, +2 bless, +1 magic weapon; *Result: 17*
}

// ExampleActor_SubHP demonstrates HP tracking with damage and healing.
func ExampleActor_SubHP() {
	fighter, _ := d20.NewActor(45, 18, 2)

	fmt.Printf("Initial HP: %d/%d\n", fighter.HP(), fighter.MaxHP())

	fighter.SubHP(15)
	fmt.Printf("After damage: %d/%d\n", fighter.HP(), fighter.MaxHP())

	fighter.AddHP(8)
	fmt.Printf("After healing: %d/%d\n", fighter.HP(), fighter.MaxHP())

	fighter.ResetHP()
	fmt.Printf("After rest: %d/%d\n", fighter.HP(), fighter.MaxHP())

	// Output:
	// Initial HP: 45/45
	// After damage: 30/45
	// After healing: 38/45
	// After rest: 45/45
}

// ExampleActor_SetMaxHP demonstrates leveling up with max HP increase.
func ExampleActor_SetMaxHP() {
	wizard, _ := d20.NewActor(22, 12, 3)
	fmt.Printf("Level 4: %d/%d HP\n", wizard.HP(), wizard.MaxHP())

	// Level up to level 5 - current HP stays at 22
	wizard.SetMaxHP(27)
	fmt.Printf("Level 5: %d/%d HP\n", wizard.HP(), wizard.MaxHP())

	// After a long rest
	wizard.ResetHP()
	fmt.Printf("After rest: %d/%d HP\n", wizard.HP(), wizard.MaxHP())

	// Output:
	// Level 4: 22/22 HP
	// Level 5: 22/27 HP
	// After rest: 27/27 HP
}

// ExampleActor_AddCombatModifier demonstrates managing combat modifiers.
func ExampleActor_AddCombatModifier() {
	fighter, _ := d20.NewActor(45, 18, 2)

	// Add permanent modifiers
	fighter.AddCombatModifier(d20.NewModifier(3, "Strength"))
	fighter.AddCombatModifier(d20.NewModifier(3, "Proficiency"))

	// Add temporary modifier
	fighter.AddCombatModifier(d20.NewModifier(2, "Bless"))

	mods := fighter.GetCombatModifiers()
	fmt.Printf("Active modifiers: %d\n", len(mods))

	// Remove temporary modifier
	fighter.RemoveCombatModifier("bless")

	mods = fighter.GetCombatModifiers()
	fmt.Printf("After removing Bless: %d\n", len(mods))

	// Output:
	// Active modifiers: 3
	// After removing Bless: 2
}

// ExampleActor_SetAttribute demonstrates attribute management.
func ExampleActor_SetAttribute() {
	rogue, _ := d20.NewActor(30, 15, 3)

	// Set attributes
	rogue.SetAttribute("dexterity", 18)
	rogue.SetAttribute("stealth", 9)
	rogue.SetAttribute("sleight_of_hand", 9)

	// Case-insensitive lookup
	if dex, ok := rogue.Attribute("DEXTERITY"); ok {
		fmt.Printf("Dexterity: %d\n", dex)
	}

	// Check if attribute exists
	if rogue.HasAttribute("stealth") {
		fmt.Println("Has stealth proficiency")
	}

	// Remove attribute
	rogue.RemoveAttribute("sleight_of_hand")
	fmt.Printf("Has sleight of hand: %v\n", rogue.HasAttribute("sleight_of_hand"))

	// Output:
	// Dexterity: 18
	// Has stealth proficiency
	// Has sleight of hand: false
}

// ExampleRollWithAdvantage demonstrates the advantage/disadvantage utility function.
func ExampleRollWithAdvantage() {
	roller := d20.NewRoller(100)

	// Normal roll
	normal, _ := d20.RollWithAdvantage(roller, 1, 20, d20.Normal)
	fmt.Printf("Normal: %v\n", normal)

	// Advantage (rolls twice, takes higher)
	adv, _ := d20.RollWithAdvantage(roller, 1, 20, d20.Advantage)
	fmt.Printf("Advantage: %v\n", adv)

	// Disadvantage (rolls twice, takes lower)
	disadv, _ := d20.RollWithAdvantage(roller, 1, 20, d20.Disadvantage)
	fmt.Printf("Disadvantage: %v\n", disadv)

	// Output:
	// Normal: [4]
	// Advantage: [9]
	// Disadvantage: [1]
}

// ExampleNewModifier demonstrates creating modifiers.
func ExampleNewModifier() {
	bonus := d20.NewModifier(5, "Strength")
	penalty := d20.NewModifier(-2, "Exhaustion")

	fmt.Printf("Bonus: +%d (%s)\n", bonus.Value, bonus.Reason)
	fmt.Printf("Penalty: %d (%s)\n", penalty.Value, penalty.Reason)

	// Output:
	// Bonus: +5 (strength)
	// Penalty: -2 (exhaustion)
}
