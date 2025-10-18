package d20

import (
	"testing"
)

func TestRoller_Roll(t *testing.T) {
	tests := []struct {
		name      string
		seed      int64
		expr      string
		wantTotal int
		wantRolls int // number of dice rolled
		wantErr   bool
	}{
		{
			name:      "simple 2d6",
			seed:      42,
			expr:      "2d6",
			wantTotal: 12,
			wantRolls: 2,
			wantErr:   false,
		},
		{
			name:      "1d20",
			seed:      42,
			expr:      "1d20",
			wantTotal: 6,
			wantRolls: 1,
			wantErr:   false,
		},
		{
			name:      "3d6+5",
			seed:      42,
			expr:      "3d6+5",
			wantTotal: 20,
			wantRolls: 3,
			wantErr:   false,
		},
		{
			name:      "2d8-3",
			seed:      42,
			expr:      "2d8-3",
			wantTotal: 3,
			wantRolls: 2,
			wantErr:   false,
		},
		{
			name:      "case insensitive",
			seed:      42,
			expr:      "2D6",
			wantTotal: 12,
			wantRolls: 2,
			wantErr:   false,
		},
		{
			name:      "with spaces",
			seed:      42,
			expr:      " 2d6 ",
			wantTotal: 12,
			wantRolls: 2,
			wantErr:   false,
		},
		{
			name:    "invalid format",
			seed:    42,
			expr:    "invalid",
			wantErr: true,
		},
		{
			name:    "missing dice count",
			seed:    42,
			expr:    "d6",
			wantErr: true,
		},
		{
			name:    "missing sides",
			seed:    42,
			expr:    "2d",
			wantErr: true,
		},
		{
			name:    "zero dice",
			seed:    42,
			expr:    "0d6",
			wantErr: true,
		},
		{
			name:    "zero sides",
			seed:    42,
			expr:    "2d0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRoller(tt.seed)
			result, err := r.Roll(tt.expr)

			if (err != nil) != tt.wantErr {
				t.Errorf("Roll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if result.Total != tt.wantTotal {
				t.Errorf("Roll() total = %v, want %v", result.Total, tt.wantTotal)
			}

			if len(result.Rolls) != tt.wantRolls {
				t.Errorf("Roll() rolls count = %v, want %v", len(result.Rolls), tt.wantRolls)
			}

			// Verify all rolls are within valid range
			for i, roll := range result.Rolls {
				if roll < 1 {
					t.Errorf("Roll() rolls[%d] = %v, should be >= 1", i, roll)
				}
			}
		})
	}
}

func TestRoller_Reproducibility(t *testing.T) {
	seed := int64(123)
	r1 := NewRoller(seed)
	r2 := NewRoller(seed)

	expr := "3d6+2"

	result1, err1 := r1.Roll(expr)
	if err1 != nil {
		t.Fatalf("First roll failed: %v", err1)
	}

	result2, err2 := r2.Roll(expr)
	if err2 != nil {
		t.Fatalf("Second roll failed: %v", err2)
	}

	if result1.Total != result2.Total {
		t.Errorf("Results not reproducible: %d != %d", result1.Total, result2.Total)
	}

	if len(result1.Rolls) != len(result2.Rolls) {
		t.Errorf("Roll counts don't match: %d != %d", len(result1.Rolls), len(result2.Rolls))
	}

	for i := range result1.Rolls {
		if result1.Rolls[i] != result2.Rolls[i] {
			t.Errorf("Roll[%d] not reproducible: %d != %d", i, result1.Rolls[i], result2.Rolls[i])
		}
	}
}

func TestResult_String(t *testing.T) {
	tests := []struct {
		name   string
		result Result
		want   string
	}{
		{
			name: "simple rolls",
			result: Result{
				Total: 10,
				Rolls: []int{4, 6},
			},
			want: "Rolls: [4, 6], Total: 10",
		},
		{
			name: "with positive modifier",
			result: Result{
				Total: 15,
				Rolls: []int{4, 6},
			},
			want: "Rolls: [4, 6]+5, Total: 15",
		},
		{
			name: "with negative modifier",
			result: Result{
				Total: 5,
				Rolls: []int{4, 6},
			},
			want: "Rolls: [4, 6]-5, Total: 5",
		},
		{
			name: "empty rolls",
			result: Result{
				Total: 7,
				Rolls: []int{},
			},
			want: "Total: 7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.String(); got != tt.want {
				t.Errorf("Result.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
