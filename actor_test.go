package d20

import (
	"errors"
	"maps"
	"regexp"
	"testing"
)

func TestNormalizeID(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"Ironpants", "ironpants"},
		{"Busta the Black", "busta_the_black"},
		{"Fighter-1", "fighter_1"},
		{"TEST-ACTOR", "test_actor"},
		{"hero", "hero"},
		{"Goblin-#3", "goblin_3"},
		{"  __Foo__  ", "foo"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := normalizeID(tt.in); got != tt.want {
				t.Errorf("normalizeID(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestActor_Normalize(t *testing.T) {
	tests := []struct {
		name      string
		in        Actor
		wantID    string
		wantAttrs map[string]int
		wantMods  map[string]int
		wantHP    int
		wantMaxHP int
		checkHP   bool
		errIs     error
	}{
		{
			name:   "id",
			in:     Actor{ID: "Foo"},
			wantID: "foo",
		},
		{
			name:      "attribute keys",
			in:        Actor{Attributes: map[string]int{"str": 6, "Building climbing": 12}},
			wantAttrs: map[string]int{"str": 6, "building_climbing": 12},
		},
		{
			name:     "modifier keys",
			in:       Actor{Modifiers: map[string]int{"STR": 3, "strike-bonus": 1}},
			wantMods: map[string]int{"str": 3, "strike_bonus": 1},
		},
		{
			name: "same key in both maps is not a duplicate",
			in: Actor{
				Attributes: map[string]int{"Strength": 16},
				Modifiers:  map[string]int{"strength": 3},
			},
			wantAttrs: map[string]int{"strength": 16},
			wantMods:  map[string]int{"strength": 3},
		},
		{
			name:      "duplicate attribute keys",
			in:        Actor{Attributes: map[string]int{"Foo Bar": 1, "foo-bar": 2}},
			errIs:     ErrDuplicateKey,
			wantAttrs: map[string]int{"Foo Bar": 1, "foo-bar": 2},
		},
		{
			name:     "duplicate modifier keys",
			in:       Actor{Modifiers: map[string]int{"A B": 1, "a_b": 2}},
			errIs:    ErrDuplicateKey,
			wantMods: map[string]int{"A B": 1, "a_b": 2},
		},
		{
			name:      "hp greater than maxhp",
			in:        Actor{HP: 45, MaxHP: 10},
			wantHP:    45,
			wantMaxHP: 45,
			checkHP:   true,
		},
		{
			name:      "hp sets maxhp when maxhp is zero",
			in:        Actor{HP: 45},
			wantHP:    45,
			wantMaxHP: 45,
			checkHP:   true,
		},
		{
			name:      "hp less than maxhp unchanged",
			in:        Actor{HP: 10, MaxHP: 20},
			wantHP:    10,
			wantMaxHP: 20,
			checkHP:   true,
		},
		{
			name:      "equal hp and maxhp unchanged",
			in:        Actor{HP: 30, MaxHP: 30},
			wantHP:    30,
			wantMaxHP: 30,
			checkHP:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.in
			err := a.Normalize()
			if tt.errIs != nil {
				if !errors.Is(err, tt.errIs) {
					t.Fatalf("err = %v, want %v", err, tt.errIs)
				}
				if tt.wantAttrs != nil && !maps.Equal(a.Attributes, tt.wantAttrs) {
					t.Errorf("Attributes mutated on error: %v", a.Attributes)
				}
				if tt.wantMods != nil && !maps.Equal(a.Modifiers, tt.wantMods) {
					t.Errorf("Modifiers mutated on error: %v", a.Modifiers)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if tt.wantID != "" && a.ID != tt.wantID {
				t.Errorf("ID = %q, want %q", a.ID, tt.wantID)
			}
			if tt.wantAttrs != nil && !maps.Equal(a.Attributes, tt.wantAttrs) {
				t.Errorf("Attributes = %v, want %v", a.Attributes, tt.wantAttrs)
			}
			if tt.wantMods != nil && !maps.Equal(a.Modifiers, tt.wantMods) {
				t.Errorf("Modifiers = %v, want %v", a.Modifiers, tt.wantMods)
			}
			if tt.checkHP && (a.HP != tt.wantHP || a.MaxHP != tt.wantMaxHP) {
				t.Errorf("HP/MaxHP = %d/%d, want %d/%d", a.HP, a.MaxHP, tt.wantHP, tt.wantMaxHP)
			}
		})
	}
}

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
		{
			name:      "lookup normalizes keys",
			mods:      map[string]int{"building_climbing": 2},
			keys:      []string{"Building climbing"},
			valueMin:  3,
			valueMax:  22,
			diceCount: 1,
			detailRE:  `\+2 building_climbing`,
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
