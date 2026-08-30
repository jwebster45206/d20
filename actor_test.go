package d20

import (
	"regexp"
	"testing"
)

func TestActor_Build(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		hp       int
		ac       int
		setHP    bool
		setAC    bool
		attrs    map[string]int
		hasError bool
		wantID   string
		wantHP   int
		wantAC   int
		wantInit int
	}{
		{
			name:     "normalizes id and sets stats",
			id:       "TEST-ACTOR",
			hp:       20,
			ac:       15,
			setHP:    true,
			setAC:    true,
			wantID:   "test_actor",
			wantHP:   20,
			wantAC:   15,
			wantInit: 0,
		},
		{
			name:   "spaces become snake_case",
			id:     "Busta the Black",
			hp:     10,
			ac:     10,
			setHP:  true,
			setAC:  true,
			wantID: "busta_the_black",
			wantHP: 10,
			wantAC: 10,
		},
		{
			name:   "with attributes",
			id:     "hero",
			hp:     20,
			ac:     15,
			setHP:  true,
			setAC:  true,
			attrs:  map[string]int{"strength": 16, "Stealth": 5},
			wantID: "hero",
			wantHP: 20,
			wantAC: 15,
		},
		{
			name:     "missing hp",
			id:       "test",
			ac:       10,
			setAC:    true,
			hasError: true,
		},
		{
			name:     "missing ac",
			id:       "test",
			hp:       10,
			setHP:    true,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewActor(tt.id)
			if tt.setHP {
				b = b.WithHP(tt.hp)
			}
			if tt.setAC {
				b = b.WithAC(tt.ac)
			}
			if tt.attrs != nil {
				b = b.WithAttributes(tt.attrs)
			}
			actor, err := b.Build()
			if tt.hasError {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actor.ID() != tt.wantID {
				t.Errorf("ID = %q, want %q", actor.ID(), tt.wantID)
			}
			if actor.HP() != tt.wantHP || actor.MaxHP() != tt.wantHP {
				t.Errorf("HP/MaxHP = %d/%d, want %d", actor.HP(), actor.MaxHP(), tt.wantHP)
			}
			if actor.AC() != tt.wantAC {
				t.Errorf("AC = %d, want %d", actor.AC(), tt.wantAC)
			}
			if actor.Initiative() != tt.wantInit {
				t.Errorf("Initiative = %d, want %d", actor.Initiative(), tt.wantInit)
			}
			for k, v := range tt.attrs {
				got, ok := actor.Attribute(k)
				if !ok || got != v {
					t.Errorf("Attribute(%q) = %d,%v, want %d,true", k, got, ok, v)
				}
			}
		})
	}
}

func TestActor_State(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Actor)
		hasError bool
		check    func(*testing.T, *Actor)
	}{
		{
			name: "SetHP valid",
			setup: func(a *Actor) {
				_ = a.SetHP(10)
			},
			check: func(t *testing.T, a *Actor) {
				if a.HP() != 10 {
					t.Errorf("HP = %d, want 10", a.HP())
				}
			},
		},
		{
			name:     "SetHP negative errors",
			setup:    func(a *Actor) {},
			hasError: true,
			check: func(t *testing.T, a *Actor) {
				if err := a.SetHP(-1); err == nil {
					t.Fatal("expected error")
				}
			},
		},
		{
			name:     "SetHP above max errors",
			setup:    func(a *Actor) {},
			hasError: true,
			check: func(t *testing.T, a *Actor) {
				if err := a.SetHP(100); err == nil {
					t.Fatal("expected error")
				}
			},
		},
		{
			name: "SetMaxHP clamps current",
			setup: func(a *Actor) {
				_ = a.SetHP(20)
				_ = a.SetMaxHP(15)
			},
			check: func(t *testing.T, a *Actor) {
				if a.MaxHP() != 15 || a.HP() != 15 {
					t.Errorf("HP/Max = %d/%d, want 15/15", a.HP(), a.MaxHP())
				}
			},
		},
		{
			name:     "SetMaxHP zero errors",
			hasError: true,
			check: func(t *testing.T, a *Actor) {
				if err := a.SetMaxHP(0); err == nil {
					t.Fatal("expected error")
				}
			},
		},
		{
			name: "SubHP floors at zero",
			setup: func(a *Actor) {
				a.SubHP(100)
			},
			check: func(t *testing.T, a *Actor) {
				if a.HP() != 0 || !a.IsKnockedOut() {
					t.Errorf("HP = %d knockedOut=%v", a.HP(), a.IsKnockedOut())
				}
			},
		},
		{
			name: "AddHP caps at max",
			setup: func(a *Actor) {
				_ = a.SetHP(5)
				a.AddHP(100)
			},
			check: func(t *testing.T, a *Actor) {
				if a.HP() != a.MaxHP() {
					t.Errorf("HP = %d, want max %d", a.HP(), a.MaxHP())
				}
			},
		},
		{
			name: "ResetHP",
			setup: func(a *Actor) {
				_ = a.SetHP(3)
				a.ResetHP()
			},
			check: func(t *testing.T, a *Actor) {
				if a.HP() != a.MaxHP() {
					t.Errorf("HP = %d, want %d", a.HP(), a.MaxHP())
				}
			},
		},
		{
			name: "SetAC",
			setup: func(a *Actor) {
				_ = a.SetAC(18)
			},
			check: func(t *testing.T, a *Actor) {
				if a.AC() != 18 {
					t.Errorf("AC = %d, want 18", a.AC())
				}
			},
		},
		{
			name:     "SetAC negative errors",
			hasError: true,
			check: func(t *testing.T, a *Actor) {
				if err := a.SetAC(-1); err == nil {
					t.Fatal("expected error")
				}
			},
		},
		{
			name: "SetInitiative",
			setup: func(a *Actor) {
				a.SetInitiative(7)
			},
			check: func(t *testing.T, a *Actor) {
				if a.Initiative() != 7 {
					t.Errorf("Initiative = %d, want 7", a.Initiative())
				}
			},
		},
		{
			name: "attribute set get remove",
			setup: func(a *Actor) {
				a.SetAttribute("stealth", 45)
				a.IncrementAttribute("stealth", 5)
				a.DecrementAttribute("stealth", 2)
			},
			check: func(t *testing.T, a *Actor) {
				v, ok := a.Attribute("stealth")
				if !ok || v != 48 {
					t.Errorf("stealth = %d,%v, want 48,true", v, ok)
				}
				a.RemoveAttribute("stealth")
				if a.HasAttribute("stealth") {
					t.Error("stealth still present after remove")
				}
			},
		},
		{
			name: "combat modifiers add remove",
			setup: func(a *Actor) {
				a.AddCombatModifier("strength", 3)
				a.AddCombatModifier("proficiency", 2)
				a.RemoveCombatModifier("strength")
			},
			check: func(t *testing.T, a *Actor) {
				mods := a.GetCombatModifiers()
				if len(mods) != 1 || mods[0].Reason != "proficiency" || mods[0].Value != 2 {
					t.Errorf("mods = %+v, want proficiency +2 only", mods)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actor, err := NewActor("hero").WithHP(20).WithAC(15).Build()
			if err != nil {
				t.Fatalf("Build: %v", err)
			}
			if tt.setup != nil {
				tt.setup(actor)
			}
			tt.check(t, actor)
		})
	}
}

func TestActor_Rolls(t *testing.T) {
	tests := []struct {
		name      string
		attrs     map[string]int
		combat    map[string]int
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
			combat:    map[string]int{"strength": 3, "proficiency": 2},
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
			combat:    map[string]int{"strength": 5},
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
			b := NewActor("hero").WithHP(20).WithAC(15)
			if tt.attrs != nil {
				b = b.WithAttributes(tt.attrs)
			}
			if tt.combat != nil {
				b = b.WithCombatModifiers(tt.combat)
			}
			actor, err := b.Build()
			if err != nil {
				t.Fatalf("Build: %v", err)
			}

			var out RollOutcome
			switch tt.kind {
			case "skill":
				builder, err := actor.SkillCheck(tt.skill, roller)
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
				builder, err := actor.D100SkillCheck(tt.skill, roller)
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
