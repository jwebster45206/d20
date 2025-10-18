package d20_test

import (
	"fmt"

	"github.com/jwebster45206/d20"
)

func Example() {
	// Create a new roller with a seed for reproducible results
	roller := d20.NewRoller(42)

	// Roll 2d6 (two six-sided dice)
	result, err := roller.Roll("2d6")
	if err != nil {
		panic(err)
	}
	fmt.Printf("2d6: %s\n", result)

	// Roll 1d20+5 (one twenty-sided die plus 5)
	result, err = roller.Roll("1d20+5")
	if err != nil {
		panic(err)
	}
	fmt.Printf("1d20+5: %s\n", result)

	// Output:
	// 2d6: Rolls: [6, 6], Total: 12
	// 1d20+5: Rolls: [9]+5, Total: 14
}

func ExampleRoller_Roll() {
	roller := d20.NewRoller(123)

	result, err := roller.Roll("3d6")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Total: %d, Individual rolls: %v\n", result.Total, result.Rolls)
	// Output:
	// Total: 12, Individual rolls: [6 4 2]
}

func ExampleNewRoller() {
	// Create a roller with a specific seed for reproducible results
	roller := d20.NewRoller(42)

	// Same seed produces same results
	result1, _ := roller.Roll("1d20")
	fmt.Printf("First roll: %d\n", result1.Total)

	// Create another roller with the same seed
	roller2 := d20.NewRoller(42)
	result2, _ := roller2.Roll("1d20")
	fmt.Printf("Second roll: %d\n", result2.Total)

	// Output:
	// First roll: 6
	// Second roll: 6
}
