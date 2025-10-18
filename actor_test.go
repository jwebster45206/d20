package d20

import (
	"testing"
)

// Test ActorBuilder - NewActor and Build
func TestActorBuilder_NewActorAndBuild(t *testing.T) {
	actor, err := NewActor("TEST-ACTOR", 20, 15).Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	// ID should be normalized to lowercase snake_case
	if actor.ID() != "test_actor" {
		t.Errorf("Expected ID 'test_actor', got '%s'", actor.ID())
	}
	if actor.MaxHP() != 20 {
		t.Errorf("Expected MaxHP 20, got %d", actor.MaxHP())
	}
	if actor.HP() != 20 {
		t.Errorf("Expected HP 20 (starts at max), got %d", actor.HP())
	}
	if actor.AC() != 15 {
		t.Errorf("Expected AC 15, got %d", actor.AC())
	}
	if actor.Initiative() != 0 {
		t.Errorf("Expected Initiative 0 (default), got %d", actor.Initiative())
	}
}

// Test ActorBuilder validation
func TestActorBuilder_BuildValidation(t *testing.T) {
	// HP must be > 0
	_, err := NewActor("test", 0, 15).Build()
	if err == nil {
		t.Error("Expected error for HP <= 0, got nil")
	}

	// AC must be > 0
	_, err = NewActor("test", 20, 0).Build()
	if err == nil {
		t.Error("Expected error for AC <= 0, got nil")
	}
}

// Test ActorBuilder.WithInitiative
func TestActorBuilder_WithInitiative(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).
		WithInitiative(3).
		Build()

	if actor.Initiative() != 3 {
		t.Errorf("Expected initiative 3, got %d", actor.Initiative())
	}
}

// Test ActorBuilder.WithAttribute
func TestActorBuilder_WithAttribute(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).
		WithAttribute("Strength", 16).
		WithAttribute("DEX", 14).
		Build()

	// Keys should be lowercased
	str, exists := actor.Attribute("strength")
	if !exists || str != 16 {
		t.Errorf("Expected strength 16, got %d (exists: %v)", str, exists)
	}

	dex, exists := actor.Attribute("dex")
	if !exists || dex != 14 {
		t.Errorf("Expected dex 14, got %d (exists: %v)", dex, exists)
	}
}

// Test ActorBuilder.WithAttributes
func TestActorBuilder_WithAttributes(t *testing.T) {
	attrs := map[string]int{
		"Strength":     16,
		"DEXTERITY":    14,
		"constitution": 15,
	}
	actor, _ := NewActor("hero", 20, 15).
		WithAttributes(attrs).
		Build()

	// All keys should be lowercased
	str, _ := actor.Attribute("strength")
	if str != 16 {
		t.Errorf("Expected strength 16, got %d", str)
	}

	dex, _ := actor.Attribute("dexterity")
	if dex != 14 {
		t.Errorf("Expected dexterity 14, got %d", dex)
	}

	con, _ := actor.Attribute("constitution")
	if con != 15 {
		t.Errorf("Expected constitution 15, got %d", con)
	}
}

// Test ActorBuilder.WithCombatModifier
func TestActorBuilder_WithCombatModifier(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).
		WithCombatModifier("flanking", 2).
		WithCombatModifier("bless", 1).
		Build()

	mods := actor.GetCombatModifiers()
	if len(mods) != 2 {
		t.Errorf("Expected 2 combat modifiers, got %d", len(mods))
	}

	// Verify via attack roll
	roller := NewRoller(42)
	builder := actor.AttackRoll(roller)
	result, _ := builder.Roll()

	// Expect dice + flanking(2) + bless(1) = dice + 3
	expected := result.DiceRolls[0] + 3
	if result.Value != expected {
		t.Errorf("Expected value %d, got %d", expected, result.Value)
	}
}

// Test ActorBuilder.WithCombatModifiers
func TestActorBuilder_WithCombatModifiers(t *testing.T) {
	mods := map[string]int{
		"Flanking": 2,
		"BLESS":    1,
	}
	actor, _ := NewActor("hero", 20, 15).
		WithCombatModifiers(mods).
		Build()

	combatMods := actor.GetCombatModifiers()
	if len(combatMods) != 2 {
		t.Errorf("Expected 2 combat modifiers, got %d", len(combatMods))
	}
}

// Test Actor.ID normalization
func TestActor_IDNormalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Simple", "simple"},
		{"UPPERCASE", "uppercase"},
		{"Mixed Case", "mixed_case"},
		{"With-Dashes", "with_dashes"},
		{"Multiple   Spaces", "multiple_spaces"},
		{"Special!@#Characters", "special_characters"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actor, _ := NewActor(tt.input, 20, 15).Build()
			if actor.ID() != tt.expected {
				t.Errorf("Expected ID '%s', got '%s'", tt.expected, actor.ID())
			}
		})
	}
}

// Test Actor.SetHP
func TestActor_SetHP(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()

	// Valid HP change
	err := actor.SetHP(10)
	if err != nil {
		t.Errorf("SetHP(10) error: %v", err)
	}
	if actor.HP() != 10 {
		t.Errorf("Expected HP 10, got %d", actor.HP())
	}

	// HP cannot be negative
	err = actor.SetHP(-5)
	if err == nil {
		t.Error("Expected error for negative HP, got nil")
	}

	// HP cannot exceed max
	err = actor.SetHP(25)
	if err == nil {
		t.Error("Expected error for HP > MaxHP, got nil")
	}
}

// Test Actor.SetMaxHP
func TestActor_SetMaxHP(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()

	// Valid max HP change
	err := actor.SetMaxHP(30)
	if err != nil {
		t.Errorf("SetMaxHP(30) error: %v", err)
	}
	if actor.MaxHP() != 30 {
		t.Errorf("Expected MaxHP 30, got %d", actor.MaxHP())
	}

	// Max HP must be > 0
	err = actor.SetMaxHP(0)
	if err == nil {
		t.Error("Expected error for MaxHP <= 0, got nil")
	}

	// Current HP adjusted if exceeds new max
	_ = actor.SetHP(30)
	_ = actor.SetMaxHP(15)
	if actor.HP() != 15 {
		t.Errorf("Expected HP adjusted to 15, got %d", actor.HP())
	}
}

// Test Actor.SubHP
func TestActor_SubHP(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()

	actor.SubHP(5)
	if actor.HP() != 15 {
		t.Errorf("Expected HP 15, got %d", actor.HP())
	}

	// HP cannot go below 0
	actor.SubHP(20)
	if actor.HP() != 0 {
		t.Errorf("Expected HP 0 (floor), got %d", actor.HP())
	}
}

// Test Actor.AddHP
func TestActor_AddHP(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()
	_ = actor.SetHP(10)

	actor.AddHP(5)
	if actor.HP() != 15 {
		t.Errorf("Expected HP 15, got %d", actor.HP())
	}

	// HP cannot exceed max
	actor.AddHP(20)
	if actor.HP() != 20 {
		t.Errorf("Expected HP 20 (ceiling at max), got %d", actor.HP())
	}
}

// Test Actor.ResetHP
func TestActor_ResetHP(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()
	_ = actor.SetHP(5)

	actor.ResetHP()
	if actor.HP() != 20 {
		t.Errorf("Expected HP reset to 20, got %d", actor.HP())
	}
}

// Test Actor.IsKnockedOut
func TestActor_IsKnockedOut(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()

	if actor.IsKnockedOut() {
		t.Error("Expected not knocked out at full HP")
	}

	_ = actor.SetHP(0)
	if !actor.IsKnockedOut() {
		t.Error("Expected knocked out at 0 HP")
	}
}

// Test Actor.SetAC
func TestActor_SetAC(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()

	err := actor.SetAC(18)
	if err != nil {
		t.Errorf("SetAC(18) error: %v", err)
	}
	if actor.AC() != 18 {
		t.Errorf("Expected AC 18, got %d", actor.AC())
	}

	// AC must be > 0
	err = actor.SetAC(0)
	if err == nil {
		t.Error("Expected error for AC <= 0, got nil")
	}
}

// Test Actor.SetInitiative
func TestActor_SetInitiative(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()

	actor.SetInitiative(5)
	if actor.Initiative() != 5 {
		t.Errorf("Expected Initiative 5, got %d", actor.Initiative())
	}
}

// Test Actor.Attribute and SetAttribute
func TestActor_AttributeAndSetAttribute(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()

	// Initially no attributes
	_, exists := actor.Attribute("strength")
	if exists {
		t.Error("Expected no strength attribute initially")
	}

	// Set attribute
	actor.SetAttribute("Strength", 16)
	str, exists := actor.Attribute("strength")
	if !exists || str != 16 {
		t.Errorf("Expected strength 16, got %d (exists: %v)", str, exists)
	}
}

// Test Actor.HasAttribute
func TestActor_HasAttribute(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()

	if actor.HasAttribute("strength") {
		t.Error("Expected no strength attribute initially")
	}

	actor.SetAttribute("strength", 16)
	if !actor.HasAttribute("strength") {
		t.Error("Expected to have strength attribute")
	}
}

// Test Actor.RemoveAttribute
func TestActor_RemoveAttribute(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()
	actor.SetAttribute("strength", 16)

	actor.RemoveAttribute("strength")
	if actor.HasAttribute("strength") {
		t.Error("Expected strength attribute to be removed")
	}
}

// Test Actor.IncrementAttribute
func TestActor_IncrementAttribute(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()
	actor.SetAttribute("strength", 16)

	actor.IncrementAttribute("strength", 2)
	str, _ := actor.Attribute("strength")
	if str != 18 {
		t.Errorf("Expected strength 18, got %d", str)
	}

	// Create new attribute if doesn't exist
	actor.IncrementAttribute("newstat", 5)
	newstat, _ := actor.Attribute("newstat")
	if newstat != 5 {
		t.Errorf("Expected newstat 5, got %d", newstat)
	}
}

// Test Actor.DecrementAttribute
func TestActor_DecrementAttribute(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()
	actor.SetAttribute("strength", 16)

	actor.DecrementAttribute("strength", 2)
	str, _ := actor.Attribute("strength")
	if str != 14 {
		t.Errorf("Expected strength 14, got %d", str)
	}

	// Create new attribute if doesn't exist
	actor.DecrementAttribute("newstat", 5)
	newstat, _ := actor.Attribute("newstat")
	if newstat != -5 {
		t.Errorf("Expected newstat -5, got %d", newstat)
	}
}

// Test Actor.AddCombatModifier and GetCombatModifiers
func TestActor_AddCombatModifier(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).Build()

	actor.AddCombatModifier("flanking", 2)
	actor.AddCombatModifier("bless", 1)

	mods := actor.GetCombatModifiers()
	if len(mods) != 2 {
		t.Errorf("Expected 2 combat modifiers, got %d", len(mods))
	}
}

// Test Actor.RemoveCombatModifier
func TestActor_RemoveCombatModifier(t *testing.T) {
	actor, _ := NewActor("hero", 20, 15).
		WithCombatModifier("flanking", 2).
		WithCombatModifier("bless", 1).
		Build()

	actor.RemoveCombatModifier("flanking")

	mods := actor.GetCombatModifiers()
	if len(mods) != 1 {
		t.Errorf("Expected 1 combat modifier, got %d", len(mods))
	}

	// Verify bless is still there
	if mods[0].Reason != "bless" {
		t.Errorf("Expected remaining modifier to be 'bless', got '%s'", mods[0].Reason)
	}
}

// Test Actor.SkillCheck - returns RollBuilder
func TestActor_SkillCheck(t *testing.T) {
	roller := NewRoller(42)
	actor, _ := NewActor("hero", 20, 15).
		WithAttribute("dexterity", 16).
		Build()

	builder, err := actor.SkillCheck("dexterity", roller)
	if err != nil {
		t.Fatalf("SkillCheck() error: %v", err)
	}

	result, err := builder.Roll()
	if err != nil {
		t.Fatalf("Roll() error: %v", err)
	}

	// Value should be dice + dexterity (16)
	expected := result.DiceRolls[0] + 16
	if result.Value != expected {
		t.Errorf("Expected value %d, got %d", expected, result.Value)
	}
}

// Test Actor.SkillCheck - missing skill
func TestActor_SkillCheck_MissingSkill(t *testing.T) {
	roller := NewRoller(42)
	actor, _ := NewActor("hero", 20, 15).Build()

	_, err := actor.SkillCheck("nonexistent", roller)
	if err == nil {
		t.Error("Expected error for missing skill, got nil")
	}
}

// Test Actor.SkillCheck with advantage
func TestActor_SkillCheck_WithAdvantage(t *testing.T) {
	roller := NewRoller(42)
	actor, _ := NewActor("hero", 20, 15).
		WithAttribute("stealth", 5).
		Build()

	builder, _ := actor.SkillCheck("stealth", roller)
	result, _ := builder.WithAdvantage().Roll()

	// With advantage, should have 2 dice rolls (both rolls visible)
	if len(result.DiceRolls) != 2 {
		t.Errorf("Expected 2 dice rolls (advantage), got %d", len(result.DiceRolls))
	}

	// Value should be higher roll + stealth (5)
	higherRoll := max(result.DiceRolls[0], result.DiceRolls[1])
	expected := higherRoll + 5
	if result.Value != expected {
		t.Errorf("Expected value %d, got %d", expected, result.Value)
	}
}

// Test Actor.AttackRoll - returns RollBuilder
func TestActor_AttackRoll(t *testing.T) {
	roller := NewRoller(42)
	actor, _ := NewActor("hero", 20, 15).
		WithCombatModifier("strength", 3).
		WithCombatModifier("proficiency", 2).
		Build()

	builder := actor.AttackRoll(roller)
	result, err := builder.Roll()
	if err != nil {
		t.Fatalf("Roll() error: %v", err)
	}

	// Value should be dice + 5 (3+2)
	expected := result.DiceRolls[0] + 5
	if result.Value != expected {
		t.Errorf("Expected value %d, got %d", expected, result.Value)
	}
}

// Test Actor.AttackRoll with no modifiers
func TestActor_AttackRoll_NoModifiers(t *testing.T) {
	roller := NewRoller(42)
	actor, _ := NewActor("hero", 20, 15).Build()

	builder := actor.AttackRoll(roller)
	result, _ := builder.Roll()

	// Just the dice roll, no modifiers
	if result.Value != result.DiceRolls[0] {
		t.Errorf("Expected value %d (dice only), got %d", result.DiceRolls[0], result.Value)
	}
}

// Test Actor.AttackRoll with advantage
func TestActor_AttackRoll_WithAdvantage(t *testing.T) {
	roller := NewRoller(42)
	actor, _ := NewActor("hero", 20, 15).
		WithCombatModifier("strength", 3).
		Build()

	builder := actor.AttackRoll(roller)
	result, _ := builder.WithAdvantage().Roll()

	// With advantage, should have 2 dice rolls (both rolls visible)
	if len(result.DiceRolls) != 2 {
		t.Errorf("Expected 2 dice rolls (advantage), got %d", len(result.DiceRolls))
	}

	// Value should be higher roll + strength (3)
	higherRoll := max(result.DiceRolls[0], result.DiceRolls[1])
	expected := higherRoll + 3
	if result.Value != expected {
		t.Errorf("Expected value %d, got %d", expected, result.Value)
	}
}

// Test Actor.D100SkillCheck
func TestActor_D100SkillCheck(t *testing.T) {
	roller := NewRoller(42)
	actor, _ := NewActor("hero", 20, 15).
		WithAttribute("stealth", 45).
		Build()

	success, outcome, err := actor.D100SkillCheck("stealth", roller, 0)
	if err != nil {
		t.Fatalf("D100SkillCheck() error: %v", err)
	}

	// Outcome should have a value between 1 and 100
	if outcome.Value < 1 || outcome.Value > 100 {
		t.Errorf("Expected value 1-100, got %d", outcome.Value)
	}

	// Success should be true if value <= 45
	expectedSuccess := outcome.Value <= 45
	if success != expectedSuccess {
		t.Errorf("Expected success %v, got %v (rolled %d vs 45)", expectedSuccess, success, outcome.Value)
	}
}

// Test Actor.D100SkillCheck - missing skill
func TestActor_D100SkillCheck_MissingSkill(t *testing.T) {
	roller := NewRoller(42)
	actor, _ := NewActor("hero", 20, 15).Build()

	_, _, err := actor.D100SkillCheck("nonexistent", roller, 0)
	if err == nil {
		t.Error("Expected error for missing skill, got nil")
	}
}
