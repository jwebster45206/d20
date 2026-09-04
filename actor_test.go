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
			if actor.Modifiers == nil {
				t.Error("Modifiers map is nil")
			}
		})
	}
}

func TestActor_D20Dice(t *testing.T) {
	tests := []struct {
		name      string
		mods      map[string]int
		keys      []string
		advantage bool
		valueMin  int
		valueMax  int
		diceCount int
		detailRE  string
		wantMods  int
	}{
		{
			name:      "named mods",
			mods:      map[string]int{"strength": 3, "striking": 2},
			keys:      []string{"strength", "striking"},
			valueMin:  6,
			valueMax:  25,
			diceCount: 1,
			wantMods:  2,
		},
		{
			name:      "no keys",
			valueMin:  1,
			valueMax:  20,
			diceCount: 1,
			wantMods:  0,
		},
		{
			name:      "skips missing and unselected keys",
			mods:      map[string]int{"strength": 5, "damage": 10},
			keys:      []string{"strength", "striking"},
			valueMin:  6,
			valueMax:  25,
			diceCount: 1,
			wantMods:  1,
		},
		{
			name:      "advantage",
			mods:      map[string]int{"strength": 5},
			keys:      []string{"strength"},
			advantage: true,
			valueMin:  6,
			valueMax:  25,
			diceCount: 2,
			wantMods:  1,
		},
		{
			name:      "lookup is case-insensitive",
			mods:      map[string]int{"strength": 3},
			keys:      []string{"Strength"},
			valueMin:  4,
			valueMax:  23,
			diceCount: 1,
			detailRE:  `\+3 strength`,
			wantMods:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller := NewRoller(42)
			actor := NewActor("hero")
			if tt.mods != nil {
				actor.Modifiers = tt.mods
			}

			d, err := actor.D20Dice(tt.keys...)
			if err != nil {
				t.Fatalf("D20Dice: %v", err)
			}
			if tt.advantage {
				d = d.WithAdvantage()
			}
			out, err := roller.Roll(d)
			if err != nil {
				t.Fatalf("Roll: %v", err)
			}

			if out.Value < tt.valueMin || out.Value > tt.valueMax {
				t.Errorf("Value %d not in [%d, %d]", out.Value, tt.valueMin, tt.valueMax)
			}
			if len(out.DiceRolls) != tt.diceCount {
				t.Errorf("DiceRolls len = %d, want %d", len(out.DiceRolls), tt.diceCount)
			}
			if len(out.Modifiers) != tt.wantMods {
				t.Errorf("Modifiers len = %d, want %d", len(out.Modifiers), tt.wantMods)
			}
			if tt.detailRE != "" && !regexp.MustCompile(tt.detailRE).MatchString(out.Detail()) {
				t.Errorf("Detail() = %q, want match %q", out.Detail(), tt.detailRE)
			}
		})
	}
}
