# d20

[![CI](https://github.com/jwebster45206/d20/actions/workflows/ci.yml/badge.svg)](https://github.com/jwebster45206/d20/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jwebster45206/d20.svg)](https://pkg.go.dev/github.com/jwebster45206/d20)

Go library for dice rolling and common D20 mechanics.

## Features

- Simple and intuitive API for dice rolling
- Seedable random number generator for reproducible results
- Support for standard dice notation (e.g., "2d6", "1d20+5")
- Pure Go implementation with no external dependencies

## Installation

```bash
go get github.com/jwebster45206/d20
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/jwebster45206/d20"
)

func main() {
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
}
```

## Supported Dice Notation

- `NdS` - Roll N dice with S sides (e.g., "2d6", "3d8")
- `NdS+M` - Roll N dice with S sides and add modifier M (e.g., "1d20+5")
- `NdS-M` - Roll N dice with S sides and subtract modifier M (e.g., "3d6-2")

## API Documentation

### `NewRoller(seed int64) *Roller`

Creates a new Roller with the given seed. Use the same seed to get reproducible results.

### `(*Roller) Roll(expr string) (Result, error)`

Evaluates a dice expression and returns the result. Returns an error if the expression is invalid.

### `type Result`

```go
type Result struct {
    Total int   // The final result including modifiers
    Rolls []int // Individual dice rolls
}
```

The `Result` type has a `String()` method that formats the result in a human-readable way.

## License

MIT License - see [LICENSE](LICENSE) file for details.
