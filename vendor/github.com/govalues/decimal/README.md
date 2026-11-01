# decimal

[![githubb]][github]
[![codecovb]][codecov]
[![goreportb]][goreport]
[![godocb]][godoc]
[![licenseb]][license]
[![versionb]][version]
[![awesomeb]][awesome]

Package decimal implements correctly rounded decimal floating-point numbers for Go.
This package is designed specifically for use in transactional financial systems.

## Key Features

- **BSON, JSON, XML, SQL** - Implements the necessary interfaces for direct compatibility
  with the [mongo-driver/bson], [encoding/json], [encoding/xml], and [database/sql] packages.
- **No Heap Allocations** - Optimized to avoid heap allocations,
  preventing garbage collector impact during arithmetic operations.
- **Correct Rounding** - For all methods, the result is the one that would
  be obtained if the true mathematical value were rounded to 19 digits of
  precision using the [half-to-even] rounding (a.k.a. "banker's rounding").
- **No Panics** - All methods are panic-free, returning errors instead of crashing
  your application in cases such as overflow or division by zero.
- **Immutability** - Once set, a decimal remains constant,
  ensuring safe concurrent access across goroutines.
- **Simple String Representation** - Decimals are represented in a straightforward
  format avoiding the complexities of scientific or engineering notations.
- **Rigorous Testing** - All methods are cross-validated against
  the [cockroachdb/apd] and [shopspring/decimal] packages through extensive [fuzz testing].

## Getting Started

### Installation

To add the decimal package to your Go workspace:

```bash
go get github.com/govalues/decimal
```

### Basic Usage

Create decimal values using one of the constructors.
After creating a decimal, you can perform various operations as shown below:

```go
package main

import (
    "fmt"
    "github.com/govalues/decimal"
)

func main() {
    // Constructors
    d, _ := decimal.New(8, 0)               // d = 8
    e, _ := decimal.Parse("12.5")           // e = 12.5
    f, _ := decimal.NewFromFloat64(2.567)   // f = 2.567
    g, _ := decimal.NewFromInt64(7, 896, 3) // g = 7.896

    // Arithmetic operations
    fmt.Println(d.Add(e))              // 8 + 12.5
    fmt.Println(d.Sub(e))              // 8 - 12.5
    fmt.Println(d.SubAbs(e))           // abs(8 - 12.5)

    fmt.Println(d.Mul(e))              // 8 * 12.5
    fmt.Println(d.AddMul(e, f))        // 8 + 12.5 * 2.567
    fmt.Println(d.SubMul(e, f))        // 8 - 12.5 * 2.567
    fmt.Println(d.PowInt(2))           // 8²

    fmt.Println(d.Quo(e))              // 8 / 12.5
    fmt.Println(d.AddQuo(e, f))        // 8 + 12.5 / 2.567
    fmt.Println(d.SubQuo(e, f))        // 8 - 12.5 / 2.567
    fmt.Println(d.QuoRem(e))           // 8 div 12.5, 8 mod 12.5
    fmt.Println(d.Inv())               // 1 / 8

    fmt.Println(decimal.Sum(d, e, f))  // 8 + 12.5 + 2.567
    fmt.Println(decimal.Mean(d, e, f)) // (8 + 12.5 + 2.567) / 3
    fmt.Println(decimal.Prod(d, e, f)) // 8 * 12.5 * 2.567

    // Transcendental functions
    fmt.Println(e.Sqrt())              // √12.5
    fmt.Println(e.Exp())               // exp(12.5)
    fmt.Println(e.Expm1())             // exp(12.5) - 1
    fmt.Println(e.Log())               // ln(12.5)
    fmt.Println(e.Log1p())             // ln(12.5 + 1)
    fmt.Println(e.Log2())              // log₂(12.5)
    fmt.Println(e.Log10())             // log₁₀(12.5)
    fmt.Println(e.Pow(d))              // 12.5⁸

    // Rounding to 2 decimal places
    fmt.Println(g.Round(2))            // 7.90
    fmt.Println(g.Ceil(2))             // 7.90
    fmt.Println(g.Floor(2))            // 7.89
    fmt.Println(g.Trunc(2))            // 7.89

    // Conversions
    fmt.Println(f.Int64(9))            // 2 567000000
    fmt.Println(f.Float64())           // 2.567
    fmt.Println(f.String())            // 2.567

    // Formatting
    fmt.Printf("%.2f", f)              // 2.57
    fmt.Printf("%.2k", f)              // 256.70%
}
```

## Documentation

For detailed documentation and additional examples, visit the package
[documentation](https://pkg.go.dev/github.com/govalues/decimal#section-documentation).
For examples related to financial calculations, see the `money` package
[documentation](https://pkg.go.dev/github.com/govalues/money#section-documentation).

## Comparison

Comparison with other popular packages:

| Feature              | govalues  | [cockroachdb/apd] v3.2.1 | [shopspring/decimal] v1.4.0 |
| -------------------- | --------- | ------------------------ | --------------------------- |
| Correctly Rounded    | Yes       | No                       | No                          |
| Speed                | High      | Medium                   | Low[^reason]                |
| Heap Allocations     | No        | Medium                   | High                        |
| Precision            | 19 digits | Arbitrary                | Arbitrary                   |
| Panic Free           | Yes       | Yes                      | No[^divzero]                |
| Mutability           | Immutable | Mutable[^reason]         | Immutable                   |
| Mathematical Context | Implicit  | Explicit                 | Implicit                    |

[^reason]: decimal package was created simply because [shopspring/decimal] was
too slow and [cockroachdb/apd] was mutable.

[^divzero]: [shopspring/decimal] panics on division by zero.

### Benchmarks

```text
goos: linux
goarch: amd64
pkg: github.com/govalues/decimal-tests
cpu: AMD Ryzen 7 3700C  with Radeon Vega Mobile Gfx 
```

| Test Case | Expression            | govalues | [cockroachdb/apd] v3.2.1 | [shopspring/decimal] v1.4.0 | govalues vs cockroachdb | govalues vs shopspring |
| --------- | --------------------- | -------: | -----------------------: | --------------------------: | ----------------------: | ---------------------: |
| Add       | 5 + 6                 |   16.06n |                   74.88n |                     140.90n |                +366.22% |               +777.33% |
| Mul       | 2 * 3                 |   16.93n |                   62.20n |                     146.00n |                +267.40% |               +762.37% |
| Quo       | 2 / 4 (exact)         |   59.52n |                  176.95n |                     657.40n |                +197.30% |              +1004.50% |
| Quo       | 2 / 3 (inexact)       |  391.60n |                  976.80n |                    2962.50n |                +149.39% |               +656.42% |
| PowInt    | 1.1^60                |  950.90n |                 3302.50n |                    4599.50n |                +247.32% |               +383.73% |
| PowInt    | 1.01^600              |    3.45µ |                   10.67µ |                      18.67µ |                +209.04% |               +440.89% |
| PowInt    | 1.001^6000            |    5.94µ |                   20.50µ |                     722.22µ |                +244.88% |             +12052.44% |
| Sqrt      | √2                    |    3.40µ |                    4.96µ |                    2101.86µ |                 +46.00% |             +61755.71% |
| Exp       | exp(0.5)              |    8.35µ |                   39.28µ |                      20.06µ |                +370.58% |               +140.32% |
| Log       | ln(0.5)               |   54.89µ |                  129.01µ |                     151.55µ |                +135.03% |               +176.10% |
| Parse     | 1                     |   16.52n |                   76.30n |                     136.55n |                +362.00% |               +726.82% |
| Parse     | 123.456               |   47.37n |                  176.90n |                     242.60n |                +273.44% |               +412.14% |
| Parse     | 123456789.1234567890  |   85.49n |                  224.15n |                     497.95n |                +162.19% |               +482.47% |
| String    | 1                     |    5.11n |                   19.57n |                     198.25n |                +283.21% |              +3783.07% |
| String    | 123.456               |   35.78n |                   77.12n |                     228.85n |                +115.52% |               +539.51% |
| String    | 123456789.1234567890  |   70.72n |                  239.10n |                     337.25n |                +238.12% |               +376.91% |
| Telco     | (see [specification]) |  137.00n |                  969.40n |                    3981.00n |                +607.33% |              +2804.78% |

The benchmark results shown in the table are provided for informational purposes only and may vary depending on your specific use case.

[codecov]: https://codecov.io/gh/govalues/decimal
[codecovb]: https://img.shields.io/codecov/c/github/govalues/decimal/main?color=brightcolor
[goreport]: https://goreportcard.com/report/github.com/govalues/decimal
[goreportb]: https://goreportcard.com/badge/github.com/govalues/decimal
[github]: https://github.com/govalues/decimal/actions/workflows/go.yml
[githubb]: https://img.shields.io/github/actions/workflow/status/govalues/decimal/go.yml
[godoc]: https://pkg.go.dev/github.com/govalues/decimal#section-documentation
[godocb]: https://img.shields.io/badge/go.dev-reference-blue
[version]: https://go.dev/dl
[versionb]: https://img.shields.io/github/go-mod/go-version/govalues/decimal?label=go
[license]: https://en.wikipedia.org/wiki/MIT_License
[licenseb]: https://img.shields.io/github/license/govalues/decimal?color=blue
[awesome]: https://github.com/avelino/awesome-go#financial
[awesomeb]: https://awesome.re/mentioned-badge.svg
[cockroachdb/apd]: https://pkg.go.dev/github.com/cockroachdb/apd
[shopspring/decimal]: https://pkg.go.dev/github.com/shopspring/decimal
[mongo-driver/bson]: https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/bson#ValueUnmarshaler
[encoding/json]: https://pkg.go.dev/encoding/json#Unmarshaler
[encoding/xml]: https://pkg.go.dev/encoding#TextUnmarshaler
[database/sql]: https://pkg.go.dev/database/sql#Scanner
[specification]: https://speleotrove.com/decimal/telcoSpec.html
[fuzz testing]: https://github.com/govalues/decimal-tests
[half-to-even]: https://en.wikipedia.org/wiki/Rounding#Rounding_half_to_even
