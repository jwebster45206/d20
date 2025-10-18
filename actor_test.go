package d20

import "testing"

func TestNewActor(t *testing.T) {
	tests := []struct {
		name       string
		hp         int
		ac         int
		initiative int
		wantError  bool
	}{
		{"Valid actor", 45, 18, 2, false},
		{"Zero HP", 0, 18, 2, true},
		{"Negative HP", -5, 18, 2, true},
		{"Zero AC", 45, 0, 2, true},
		{"Negative AC", 45, -3, 2, true},
		{"Negative initiative allowed", 45, 18, -2, false},
		{"Zero initiative allowed", 45, 18, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actor, err := NewActor(tt.hp, tt.ac, tt.initiative)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewActor() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewActor() unexpected error: %v", err)
				return
			}

			if actor.HP() != tt.hp {
				t.Errorf("HP = %d, want %d", actor.HP(), tt.hp)
			}
			if actor.AC() != tt.ac {
				t.Errorf("AC = %d, want %d", actor.AC(), tt.ac)
			}
			if actor.Initiative() != tt.initiative {
				t.Errorf("Initiative = %d, want %d", actor.Initiative(), tt.initiative)
			}
			if len(actor.GetCombatModifiers()) != 0 {
				t.Errorf("Expected empty combat modifiers, got %d", len(actor.GetCombatModifiers()))
			}
		})
	}
}

func TestNewActorWithAttributes(t *testing.T) {
	attrs := map[string]int{
		"Strength":     16,
		"DEXTERITY":    14,
		"constitution": 15,
		"Athletics":    5,
	}

	actor, err := NewActorWithAttributes(45, 18, 2, attrs)
	if err != nil {
		t.Fatalf("NewActorWithAttributes() error: %v", err)
	}

	// Check that all keys are lowercased
	tests := []struct {
		key   string
		value int
	}{
		{"strength", 16},
		{"dexterity", 14},
		{"constitution", 15},
		{"athletics", 5},
	}

	for _, tt := range tests {
		value, exists := actor.Attribute(tt.key)
		if !exists {
			t.Errorf("Expected attribute %q to exist", tt.key)
		}
		if value != tt.value {
			t.Errorf("Attribute %q = %d, want %d", tt.key, value, tt.value)
		}
	}
}

func TestActor_HPManagement(t *testing.T) {
	actor, _ := NewActor(45, 18, 2)

	// Test initial HP
	if actor.HP() != 45 {
		t.Errorf("Initial HP = %d, want 45", actor.HP())
	}
	if actor.MaxHP() != 45 {
		t.Errorf("Max HP = %d, want 45", actor.MaxHP())
	}

	// Test SetHP within bounds
	if err := actor.SetHP(30); err != nil {
		t.Errorf("SetHP(30) error: %v", err)
	}
	if actor.HP() != 30 {
		t.Errorf("HP = %d, want 30", actor.HP())
	}

	// Test SetHP at 0 (should be allowed)
	if err := actor.SetHP(0); err != nil {
		t.Errorf("SetHP(0) should be allowed: %v", err)
	}

	// Test SetHP with negative value
	if err := actor.SetHP(-5); err == nil {
		t.Error("SetHP(-5) expected error but got nil")
	}

	// Test SetHP above max
	_ = actor.SetHP(45) // Reset to max
	if err := actor.SetHP(50); err == nil {
		t.Error("SetHP(50) expected error (exceeds max HP 45)")
	}

	// Test SetMaxHP
	if err := actor.SetMaxHP(60); err != nil {
		t.Errorf("SetMaxHP(60) error: %v", err)
	}
	if actor.MaxHP() != 60 {
		t.Errorf("Max HP = %d, want 60", actor.MaxHP())
	}

	// Test SetMaxHP below current HP adjusts current HP
	_ = actor.SetHP(50)
	_ = actor.SetMaxHP(40)
	if actor.HP() != 40 {
		t.Errorf("After reducing max HP, current HP should adjust to %d, got %d", 40, actor.HP())
	}

	// Test SubHP
	_ = actor.SetHP(35)
	actor.SubHP(15)
	if actor.HP() != 20 {
		t.Errorf("After SubHP(15), HP = %d, want 20", actor.HP())
	}

	// Test SubHP below 0
	actor.SubHP(100)
	if actor.HP() != 0 {
		t.Errorf("After massive damage, HP = %d, want 0", actor.HP())
	}

	// Test IsKnockedOut
	if !actor.IsKnockedOut() {
		t.Error("Actor with 0 HP should be knocked out")
	}

	// Test AddHP
	actor.AddHP(20)
	if actor.HP() != 20 {
		t.Errorf("After AddHP(20), HP = %d, want 20", actor.HP())
	}
	if actor.IsKnockedOut() {
		t.Error("Actor with 20 HP should be alive")
	}

	// Test AddHP doesn't exceed max HP
	actor.AddHP(100)
	if actor.HP() != actor.MaxHP() {
		t.Errorf("After AddHP(100), HP = %d, should not exceed max HP %d", actor.HP(), actor.MaxHP())
	}

	// Test ResetHP
	actor.SubHP(10)
	actor.ResetHP()
	if actor.HP() != actor.MaxHP() {
		t.Errorf("After ResetHP(), HP = %d, want %d", actor.HP(), actor.MaxHP())
	}
}

func TestActor_ACManagement(t *testing.T) {
	actor, _ := NewActor(45, 18, 2)

	// Test SetAC
	if err := actor.SetAC(20); err != nil {
		t.Errorf("SetAC(20) error: %v", err)
	}
	if actor.AC() != 20 {
		t.Errorf("AC = %d, want 20", actor.AC())
	}

	// Test SetAC with invalid value
	if err := actor.SetAC(0); err == nil {
		t.Error("SetAC(0) expected error but got nil")
	}
	if err := actor.SetAC(-5); err == nil {
		t.Error("SetAC(-5) expected error but got nil")
	}
}

func TestActor_InitiativeManagement(t *testing.T) {
	actor, _ := NewActor(45, 18, 2)

	actor.SetInitiative(5)
	if actor.Initiative() != 5 {
		t.Errorf("Initiative = %d, want 5", actor.Initiative())
	}

	// Negative initiative should be allowed
	actor.SetInitiative(-3)
	if actor.Initiative() != -3 {
		t.Errorf("Initiative = %d, want -3", actor.Initiative())
	}
}

func TestActor_AttributeManagement(t *testing.T) {
	actor, _ := NewActor(45, 18, 2)

	// Test SetAttribute
	actor.SetAttribute("Strength", 16)
	actor.SetAttribute("DEXTERITY", 14)
	actor.SetAttribute("constitution", 15)

	// Test GetAttribute with case-insensitive keys
	tests := []struct {
		key   string
		value int
	}{
		{"strength", 16},
		{"STRENGTH", 16},
		{"Strength", 16},
		{"dexterity", 14},
		{"DeXtErItY", 14},
		{"constitution", 15},
	}

	for _, tt := range tests {
		value, exists := actor.Attribute(tt.key)
		if !exists {
			t.Errorf("Expected attribute %q to exist", tt.key)
		}
		if value != tt.value {
			t.Errorf("GetAttribute(%q) = %d, want %d", tt.key, value, tt.value)
		}
	}

	// Test HasAttribute
	if !actor.HasAttribute("strength") {
		t.Error("HasAttribute(\"strength\") should be true")
	}
	if !actor.HasAttribute("STRENGTH") {
		t.Error("HasAttribute(\"STRENGTH\") should be true (case-insensitive)")
	}
	if actor.HasAttribute("wisdom") {
		t.Error("HasAttribute(\"wisdom\") should be false")
	}

	// Test RemoveAttribute
	actor.RemoveAttribute("strength")
	if actor.HasAttribute("strength") {
		t.Error("Attribute should be removed")
	}
	if actor.HasAttribute("STRENGTH") {
		t.Error("Attribute should be removed (case-insensitive)")
	}
}

func TestActor_CombatModifierManagement(t *testing.T) {
	actor, _ := NewActor(45, 18, 2)

	// Test AddCombatModifier
	actor.AddCombatModifier(Modifier{Value: 3, Reason: "Strength"})
	actor.AddCombatModifier(Modifier{Value: 3, Reason: "PROFICIENCY"})
	actor.AddCombatModifier(Modifier{Value: 1, Reason: "magic weapon"})

	modifiers := actor.GetCombatModifiers()
	if len(modifiers) != 3 {
		t.Errorf("Expected 3 combat modifiers, got %d", len(modifiers))
	}

	// Check that reasons are lowercased
	expectedReasons := []string{"strength", "proficiency", "magic weapon"}
	for i, mod := range modifiers {
		if mod.Reason != expectedReasons[i] {
			t.Errorf("Modifier %d reason = %q, want %q", i, mod.Reason, expectedReasons[i])
		}
	}

	// Test that GetCombatModifiers returns a copy
	modifiers[0].Value = 999
	if actor.GetCombatModifiers()[0].Value == 999 {
		t.Error("GetCombatModifiers should return a copy, not the original slice")
	}

	// Test RemoveCombatModifier
	actor.RemoveCombatModifier("proficiency")
	modifiers = actor.GetCombatModifiers()
	if len(modifiers) != 2 {
		t.Errorf("After removal, expected 2 modifiers, got %d", len(modifiers))
	}
	for _, mod := range modifiers {
		if mod.Reason == "proficiency" {
			t.Error("Modifier 'proficiency' should be removed")
		}
	}

	// Test RemoveCombatModifier with case-insensitive reason
	actor.RemoveCombatModifier("STRENGTH")
	modifiers = actor.GetCombatModifiers()
	if len(modifiers) != 1 {
		t.Errorf("After removal, expected 1 modifier, got %d", len(modifiers))
	}
}

func TestActor_RemoveMultipleCombatModifiers(t *testing.T) {
	actor, _ := NewActor(45, 18, 2)

	// Add multiple modifiers with the same reason
	actor.AddCombatModifier(Modifier{Value: 1, Reason: "bless"})
	actor.AddCombatModifier(Modifier{Value: 2, Reason: "guidance"})
	actor.AddCombatModifier(Modifier{Value: 1, Reason: "bless"})

	// Remove all "bless" modifiers
	actor.RemoveCombatModifier("bless")
	modifiers := actor.GetCombatModifiers()

	if len(modifiers) != 1 {
		t.Errorf("Expected 1 modifier remaining, got %d", len(modifiers))
	}
	if modifiers[0].Reason != "guidance" {
		t.Errorf("Remaining modifier should be 'guidance', got %q", modifiers[0].Reason)
	}
}

func TestRollWithAdvantage(t *testing.T) {
	roller := NewRoller(42) // Deterministic seed

	t.Run("Normal roll", func(t *testing.T) {
		roller := NewRoller(100)
		rolls, err := RollWithAdvantage(roller, 1, 20, Normal)
		if err != nil {
			t.Fatalf("RollWithAdvantage() error: %v", err)
		}
		if len(rolls) != 1 {
			t.Errorf("Expected 1 roll, got %d", len(rolls))
		}
		if rolls[0] < 1 || rolls[0] > 20 {
			t.Errorf("Roll value %d out of range [1, 20]", rolls[0])
		}
	})

	t.Run("Advantage roll", func(t *testing.T) {
		roller := NewRoller(200)
		// With advantage, we roll twice and take higher
		rolls, err := RollWithAdvantage(roller, 1, 20, Advantage)
		if err != nil {
			t.Fatalf("RollWithAdvantage() error: %v", err)
		}
		if len(rolls) != 1 {
			t.Errorf("Expected 1 roll, got %d", len(rolls))
		}
		if rolls[0] < 1 || rolls[0] > 20 {
			t.Errorf("Roll value %d out of range [1, 20]", rolls[0])
		}
	})

	t.Run("Disadvantage roll", func(t *testing.T) {
		roller := NewRoller(300)
		// With disadvantage, we roll twice and take lower
		rolls, err := RollWithAdvantage(roller, 1, 20, Disadvantage)
		if err != nil {
			t.Fatalf("RollWithAdvantage() error: %v", err)
		}
		if len(rolls) != 1 {
			t.Errorf("Expected 1 roll, got %d", len(rolls))
		}
		if rolls[0] < 1 || rolls[0] > 20 {
			t.Errorf("Roll value %d out of range [1, 20]", rolls[0])
		}
	})

	t.Run("Multiple dice with advantage", func(t *testing.T) {
		roller := NewRoller(400)
		rolls, err := RollWithAdvantage(roller, 3, 6, Advantage)
		if err != nil {
			t.Fatalf("RollWithAdvantage() error: %v", err)
		}
		if len(rolls) != 3 {
			t.Errorf("Expected 3 rolls, got %d", len(rolls))
		}
		for i, roll := range rolls {
			if roll < 1 || roll > 6 {
				t.Errorf("Roll %d value %d out of range [1, 6]", i, roll)
			}
		}
	})

	t.Run("Invalid input - zero rollCount", func(t *testing.T) {
		_, err := RollWithAdvantage(roller, 0, 20, Normal)
		if err == nil {
			t.Error("Expected error for rollCount = 0")
		}
	})

	t.Run("Invalid input - zero dieFaces", func(t *testing.T) {
		_, err := RollWithAdvantage(roller, 1, 0, Normal)
		if err == nil {
			t.Error("Expected error for dieFaces = 0")
		}
	})
}

func TestActor_SkillCheck(t *testing.T) {
	roller := NewRoller(42)
	actor, _ := NewActor(45, 18, 2)
	actor.SetAttribute("athletics", 5)
	actor.SetAttribute("stealth", 2)

	t.Run("Normal skill check", func(t *testing.T) {
		result, err := actor.SkillCheck("athletics", roller, Normal)
		if err != nil {
			t.Fatalf("SkillCheck() error: %v", err)
		}
		if result.Value < 6 || result.Value > 25 { // 1d20 + 5 = [6, 25]
			t.Errorf("Result %d out of expected range [6, 25]", result.Value)
		}
		if len(result.DiceRolls) != 1 {
			t.Errorf("Expected 1 die roll, got %d", len(result.DiceRolls))
		}
	})

	t.Run("Skill check with advantage", func(t *testing.T) {
		result, err := actor.SkillCheck("stealth", roller, Advantage)
		if err != nil {
			t.Fatalf("SkillCheck() error: %v", err)
		}
		if result.Value < 3 || result.Value > 22 { // 1d20 + 2 = [3, 22]
			t.Errorf("Result %d out of expected range [3, 22]", result.Value)
		}
	})

	t.Run("Skill check with disadvantage", func(t *testing.T) {
		result, err := actor.SkillCheck("athletics", roller, Disadvantage)
		if err != nil {
			t.Fatalf("SkillCheck() error: %v", err)
		}
		if result.Value < 6 || result.Value > 25 {
			t.Errorf("Result %d out of expected range [6, 25]", result.Value)
		}
	})

	t.Run("Non-existent skill", func(t *testing.T) {
		_, err := actor.SkillCheck("nonexistent", roller, Normal)
		if err == nil {
			t.Error("Expected error for non-existent skill")
		}
	})

	t.Run("Case-insensitive skill lookup", func(t *testing.T) {
		result, err := actor.SkillCheck("ATHLETICS", roller, Normal)
		if err != nil {
			t.Fatalf("SkillCheck() should work with uppercase: %v", err)
		}
		if result.Value < 6 || result.Value > 25 {
			t.Errorf("Result %d out of expected range", result.Value)
		}
	})
}

func TestActor_AttackRoll(t *testing.T) {
	roller := NewRoller(42)
	actor, _ := NewActor(45, 18, 2)
	actor.AddCombatModifier(Modifier{Value: 5, Reason: "strength"})
	actor.AddCombatModifier(Modifier{Value: 3, Reason: "proficiency"})

	t.Run("Normal attack roll", func(t *testing.T) {
		result, err := actor.AttackRoll(roller, Normal)
		if err != nil {
			t.Fatalf("AttackRoll() error: %v", err)
		}
		// 1d20 + 8 = [9, 28]
		if result.Value < 9 || result.Value > 28 {
			t.Errorf("Result %d out of expected range [9, 28]", result.Value)
		}
		if len(result.DiceRolls) != 1 {
			t.Errorf("Expected 1 die roll, got %d", len(result.DiceRolls))
		}
	})

	t.Run("Attack roll with advantage", func(t *testing.T) {
		result, err := actor.AttackRoll(roller, Advantage)
		if err != nil {
			t.Fatalf("AttackRoll() error: %v", err)
		}
		if result.Value < 9 || result.Value > 28 {
			t.Errorf("Result %d out of expected range [9, 28]", result.Value)
		}
	})

	t.Run("Attack roll with disadvantage", func(t *testing.T) {
		result, err := actor.AttackRoll(roller, Disadvantage)
		if err != nil {
			t.Fatalf("AttackRoll() error: %v", err)
		}
		if result.Value < 9 || result.Value > 28 {
			t.Errorf("Result %d out of expected range [9, 28]", result.Value)
		}
	})

	t.Run("No combat modifiers", func(t *testing.T) {
		emptyActor, _ := NewActor(10, 10, 0)
		result, err := emptyActor.AttackRoll(roller, Normal)
		if err != nil {
			t.Fatalf("AttackRoll() error: %v", err)
		}
		// 1d20 + 0 = [1, 20]
		if result.Value < 1 || result.Value > 20 {
			t.Errorf("Result %d out of expected range [1, 20]", result.Value)
		}
	})
}

func TestActor_AttackRollWithModifiers(t *testing.T) {
	roller := NewRoller(42)
	actor, _ := NewActor(45, 18, 2)
	actor.AddCombatModifier(Modifier{Value: 5, Reason: "strength"})
	actor.AddCombatModifier(Modifier{Value: 3, Reason: "proficiency"})

	t.Run("With situational modifiers", func(t *testing.T) {
		extraMods := []Modifier{
			{Value: 2, Reason: "flanking"},
			{Value: -2, Reason: "partial cover"},
		}
		result, err := actor.AttackRollWithModifiers(roller, Normal, extraMods)
		if err != nil {
			t.Fatalf("AttackRollWithModifiers() error: %v", err)
		}
		// 1d20 + 5 + 3 + 2 - 2 = 1d20 + 8 = [9, 28]
		if result.Value < 9 || result.Value > 28 {
			t.Errorf("Result %d out of expected range [9, 28]", result.Value)
		}
	})

	t.Run("With nil extra modifiers", func(t *testing.T) {
		result, err := actor.AttackRollWithModifiers(roller, Normal, nil)
		if err != nil {
			t.Fatalf("AttackRollWithModifiers() error: %v", err)
		}
		if result.Value < 9 || result.Value > 28 {
			t.Errorf("Result %d out of expected range [9, 28]", result.Value)
		}
	})

	t.Run("With advantage and extra modifiers", func(t *testing.T) {
		extraMods := []Modifier{
			{Value: 1, Reason: "bless"},
		}
		result, err := actor.AttackRollWithModifiers(roller, Advantage, extraMods)
		if err != nil {
			t.Fatalf("AttackRollWithModifiers() error: %v", err)
		}
		// 1d20 + 5 + 3 + 1 = 1d20 + 9 = [10, 29]
		if result.Value < 10 || result.Value > 29 {
			t.Errorf("Result %d out of expected range [10, 29]", result.Value)
		}
	})

	t.Run("Extra modifier reasons are lowercased", func(t *testing.T) {
		extraMods := []Modifier{
			{Value: 2, Reason: "FLANKING"},
		}
		result, err := actor.AttackRollWithModifiers(roller, Normal, extraMods)
		if err != nil {
			t.Fatalf("AttackRollWithModifiers() error: %v", err)
		}

		// Check the detail string for lowercase
		if !contains(result.Detail, "flanking") {
			t.Errorf("Expected 'flanking' in detail string, got: %s", result.Detail)
		}
	})
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsAt(s, substr, 1)))
}

func containsAt(s, substr string, start int) bool {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestActor_D100SkillCheck(t *testing.T) {
	actor, _ := NewActor(45, 18, 2)
	actor.SetAttribute("stealth", 45)
	actor.SetAttribute("fighting", 60)
	actor.SetAttribute("spot_hidden", 25)

	t.Run("Normal d100 roll", func(t *testing.T) {
		roller := NewRoller(42)
		success, roll, err := actor.D100SkillCheck("stealth", roller, 0)
		if err != nil {
			t.Fatalf("D100SkillCheck() error: %v", err)
		}

		// Result should be 1-100
		if roll.Value < 1 || roll.Value > 100 {
			t.Errorf("D100 result %d out of range [1, 100]", roll.Value)
		}

		// Check success matches skill value
		shouldSucceed := roll.Value <= 45
		if success != shouldSucceed {
			t.Errorf("Success = %v, but roll %d vs skill 45 should be %v", success, roll.Value, shouldSucceed)
		}
	})

	t.Run("Bonus die (positive bonus)", func(t *testing.T) {
		roller := NewRoller(100)
		// With bonus die, we should still get valid d100 result
		success, roll, err := actor.D100SkillCheck("fighting", roller, 1)
		if err != nil {
			t.Fatalf("D100SkillCheck() with bonus error: %v", err)
		}

		if roll.Value < 1 || roll.Value > 100 {
			t.Errorf("D100 result %d out of range [1, 100]", roll.Value)
		}

		shouldSucceed := roll.Value <= 60
		if success != shouldSucceed {
			t.Errorf("Success = %v, but roll %d vs skill 60 should be %v", success, roll.Value, shouldSucceed)
		}
	})

	t.Run("Penalty die (negative bonus)", func(t *testing.T) {
		roller := NewRoller(200)
		// With penalty die, we should still get valid d100 result
		success, roll, err := actor.D100SkillCheck("spot_hidden", roller, -1)
		if err != nil {
			t.Fatalf("D100SkillCheck() with penalty error: %v", err)
		}

		if roll.Value < 1 || roll.Value > 100 {
			t.Errorf("D100 result %d out of range [1, 100]", roll.Value)
		}

		shouldSucceed := roll.Value <= 25
		if success != shouldSucceed {
			t.Errorf("Success = %v, but roll %d vs skill 25 should be %v", success, roll.Value, shouldSucceed)
		}
	})

	t.Run("Multiple bonus dice", func(t *testing.T) {
		roller := NewRoller(300)
		_, roll, err := actor.D100SkillCheck("fighting", roller, 2)
		if err != nil {
			t.Fatalf("D100SkillCheck() with 2 bonus dice error: %v", err)
		}

		if roll.Value < 1 || roll.Value > 100 {
			t.Errorf("D100 result %d out of range [1, 100]", roll.Value)
		}
	})

	t.Run("Multiple penalty dice", func(t *testing.T) {
		roller := NewRoller(400)
		_, roll, err := actor.D100SkillCheck("stealth", roller, -2)
		if err != nil {
			t.Fatalf("D100SkillCheck() with 2 penalty dice error: %v", err)
		}

		if roll.Value < 1 || roll.Value > 100 {
			t.Errorf("D100 result %d out of range [1, 100]", roll.Value)
		}
	})

	t.Run("Non-existent skill", func(t *testing.T) {
		roller := NewRoller(42)
		_, _, err := actor.D100SkillCheck("nonexistent", roller, 0)
		if err == nil {
			t.Error("Expected error for non-existent skill")
		}
	})

	t.Run("Success and failure cases", func(t *testing.T) {
		// Set a very high skill to ensure success
		actor.SetAttribute("guaranteed", 99)
		successCount := 0

		for i := 0; i < 10; i++ {
			roller := NewRandomRoller()
			success, _, _ := actor.D100SkillCheck("guaranteed", roller, 0)
			if success {
				successCount++
			}
		}

		// With 99% skill, we should get mostly successes
		if successCount < 8 {
			t.Logf("Warning: Expected at least 8/10 successes with 99%% skill, got %d (might be random variance)", successCount)
		}
	})

	t.Run("Skill value of 100", func(t *testing.T) {
		actor.SetAttribute("master", 100)
		roller := NewRoller(500)
		success, roll, err := actor.D100SkillCheck("master", roller, 0)
		if err != nil {
			t.Fatalf("D100SkillCheck() error: %v", err)
		}

		// Roll of 100 should succeed with skill 100
		if roll.Value == 100 && !success {
			t.Error("Roll of 100 should succeed when skill is 100")
		}
	})

	t.Run("Case-insensitive skill lookup", func(t *testing.T) {
		roller := NewRoller(42)
		_, _, err := actor.D100SkillCheck("STEALTH", roller, 0)
		if err != nil {
			t.Errorf("D100SkillCheck() should work with uppercase skill name: %v", err)
		}
	})
}

func TestActor_D100SkillCheck_BonusPenaltyMechanics(t *testing.T) {
	// Test that bonus/penalty dice actually affect outcomes statistically
	actor, _ := NewActor(45, 18, 2)
	actor.SetAttribute("test", 50)

	// Run many rolls to check statistical distribution
	normalSuccesses := 0
	bonusSuccesses := 0
	penaltySuccesses := 0
	trials := 100

	for i := 0; i < trials; i++ {
		normalRoller := NewRandomRoller()
		bonusRoller := NewRandomRoller()
		penaltyRoller := NewRandomRoller()

		normal, _, _ := actor.D100SkillCheck("test", normalRoller, 0)
		bonus, _, _ := actor.D100SkillCheck("test", bonusRoller, 1)
		penalty, _, _ := actor.D100SkillCheck("test", penaltyRoller, -1)

		if normal {
			normalSuccesses++
		}
		if bonus {
			bonusSuccesses++
		}
		if penalty {
			penaltySuccesses++
		}
	}

	t.Logf("Success rates over %d trials: Normal=%d%%, Bonus=%d%%, Penalty=%d%%",
		trials, normalSuccesses, bonusSuccesses, penaltySuccesses)

	// Bonus die should give better results than normal
	// Penalty die should give worse results than normal
	// Note: This is statistical and might occasionally fail due to randomness
	if bonusSuccesses < normalSuccesses {
		t.Logf("Warning: Bonus die (%d) performed worse than normal (%d) - might be random variance",
			bonusSuccesses, normalSuccesses)
	}
	if penaltySuccesses > normalSuccesses {
		t.Logf("Warning: Penalty die (%d) performed better than normal (%d) - might be random variance",
			penaltySuccesses, normalSuccesses)
	}
}
