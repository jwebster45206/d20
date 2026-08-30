package d20

import (
	"regexp"
	"testing"
)

func TestActor_NewActor(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		wantID string
	}{
		{name: "normalizes id", id: "TEST-ACTOR", wantID: "test_actor"},
		{name: "spaces become snake_case", id: "Busta the Black", wantID: "busta_the_black"},
		{name: "already simple", id: "hero", wantID: "hero"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actor := NewActor(tt.id)
			if actor.ID != tt.wantID {
				t.Errorf("ID = %q, want %q", actor.ID, tt.wantID)
			}
			if actor.Attributes == nil {
				t.Error("Attributes map is nil")
			}
			if actor.CombatModifiers == nil {
				t.Error("CombatModifiers is nil")
			}
		})
	}
}

func TestActor_HPHelpers(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*Actor)
		check func(*testing.T, *Actor)
	}{
		{
			name: "SubHP floors at zero",
			setup: func(a *Actor) {
				a.SubHP(100)
			},
			check: func(t *testing.T, a *Actor) {
				if a.HP != 0 || !a.IsKnockedOut() {
					t.Errorf("HP = %d knockedOut=%v", a.HP, a.IsKnockedOut())
				}
			},
		},
		{
			name: "AddHP caps at max",
			setup: func(a *Actor) {
				a.HP = 5
				a.AddHP(100)
			},
			check: func(t *testing.T, a *Actor) {
				if a.HP != a.MaxHP {
					t.Errorf("HP = %d, want max %d", a.HP, a.MaxHP)
				}
			},
		},
		{
			name: "ResetHP",
			setup: func(a *Actor) {
				a.HP = 3
				a.ResetHP()
			},
			check: func(t *testing.T, a *Actor) {
				if a.HP != a.MaxHP {
					t.Errorf("HP = %d, want %d", a.HP, a.MaxHP)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actor := NewActor("hero")
			actor.MaxHP = 20
			actor.HP = 20
			actor.AC = 15
			tt.setup(actor)
			tt.check(t, actor)
		})
	}
}

func TestActor_Rolls(t *testing.T) {
	tests := []struct {
		name      string
		attrs     map[string]int
		combat    []Modifier
		kind      string // skill, attack, d100, d100_mod
		skill     string
		modName   string
		modValue  int
		advantage bool
		hasError  bool
		valueMin  int
		valueMax  int
		diceCount int
		detailRE  string
		wantMods  int
	}{
		{
			name:      "skill check",
			attrs:     map[string]int{"athletics": 5},
			kind:      "skill",
			skill:     "athletics",
			valueMin:  6,
			valueMax:  25,
			diceCount: 1,
			detailRE:  `\+5 athletics`,
			wantMods:  1,
		},
		{
			name:     "skill missing",
			kind:     "skill",
			skill:    "missing",
			hasError: true,
		},
		{
			name:      "skill lookup is case-insensitive",
			attrs:     map[string]int{"stealth": 3},
			kind:      "skill",
			skill:     "Stealth",
			valueMin:  4,
			valueMax:  23,
			diceCount: 1,
			wantMods:  1,
		},
		{
			name:      "skill with advantage",
			attrs:     map[string]int{"stealth": 3},
			kind:      "skill",
			skill:     "stealth",
			advantage: true,
			valueMin:  4,
			valueMax:  23,
			diceCount: 2,
			wantMods:  1,
		},
		{
			name:      "attack with combat mods",
			combat:    []Modifier{NewModifier("strength", 3), NewModifier("proficiency", 2)},
			kind:      "attack",
			valueMin:  6,
			valueMax:  25,
			diceCount: 1,
			wantMods:  2,
		},
		{
			name:      "attack no mods",
			kind:      "attack",
			valueMin:  1,
			valueMax:  20,
			diceCount: 1,
			wantMods:  0,
		},
		{
			name:      "attack advantage",
			combat:    []Modifier{NewModifier("strength", 5)},
			kind:      "attack",
			advantage: true,
			valueMin:  6,
			valueMax:  25,
			diceCount: 2,
			wantMods:  1,
		},
		{
			name:      "d100 skill check",
			attrs:     map[string]int{"stealth": 45},
			kind:      "d100",
			skill:     "stealth",
			valueMin:  1,
			valueMax:  100,
			diceCount: 2,
			detailRE:  `Rolled 2d10\.\.\.`,
		},
		{
			name:      "d100 with modifier",
			attrs:     map[string]int{"stealth": 45},
			kind:      "d100_mod",
			skill:     "stealth",
			modName:   "penalty",
			modValue:  -5,
			valueMin:  -4,
			valueMax:  95,
			diceCount: 2,
			wantMods:  1,
			detailRE:  `-5 penalty`,
		},
		{
			name:     "d100 missing skill",
			kind:     "d100",
			skill:    "missing",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller := NewRoller(42)
			actor := NewActor("hero")
			actor.MaxHP = 20
			actor.HP = 20
			actor.AC = 15
			if tt.attrs != nil {
				actor.Attributes = tt.attrs
			}
			if tt.combat != nil {
				actor.CombatModifiers = tt.combat
			}

			var out RollOutcome
			var err error
			switch tt.kind {
			case "skill":
				var builder *DiceManager
				builder, err = actor.SkillCheck(tt.skill, roller)
				if tt.hasError {
					if err == nil {
						t.Fatal("expected error")
					}
					return
				}
				if err != nil {
					t.Fatalf("SkillCheck: %v", err)
				}
				if tt.advantage {
					builder = builder.WithAdvantage()
				}
				out, err = builder.Roll()
				if err != nil {
					t.Fatalf("Roll: %v", err)
				}
			case "attack":
				builder := actor.AttackRoll(roller)
				if tt.advantage {
					builder = builder.WithAdvantage()
				}
				out, err = builder.Roll()
				if err != nil {
					t.Fatalf("Roll: %v", err)
				}
			case "d100", "d100_mod":
				var builder *DiceManager
				builder, err = actor.D100SkillCheck(tt.skill, roller)
				if tt.hasError {
					if err == nil {
						t.Fatal("expected error")
					}
					return
				}
				if err != nil {
					t.Fatalf("D100SkillCheck: %v", err)
				}
				if tt.kind == "d100_mod" {
					builder = builder.WithModifier(tt.modName, tt.modValue)
				}
				out, err = builder.RollPercentile()
				if err != nil {
					t.Fatalf("RollPercentile: %v", err)
				}
			default:
				t.Fatalf("unknown kind %q", tt.kind)
			}

			if out.Value < tt.valueMin || out.Value > tt.valueMax {
				t.Errorf("Value %d not in [%d, %d]", out.Value, tt.valueMin, tt.valueMax)
			}
			if tt.diceCount >= 0 && len(out.DiceRolls) != tt.diceCount {
				t.Errorf("DiceRolls len = %d, want %d", len(out.DiceRolls), tt.diceCount)
			}
			if tt.wantMods >= 0 && len(out.Modifiers) != tt.wantMods {
				t.Errorf("Modifiers len = %d, want %d", len(out.Modifiers), tt.wantMods)
			}
			if tt.detailRE != "" && !regexp.MustCompile(tt.detailRE).MatchString(out.Detail()) {
				t.Errorf("Detail() = %q, want match %q", out.Detail(), tt.detailRE)
			}
		})
	}
}
