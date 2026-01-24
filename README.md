# Server

[![SemVer 2.0.0](https://img.shields.io/badge/SemVer-2.0.0-blue)](https://semver.org)
[![Go Reference](https://pkg.go.dev/badge/github.com/flaticols/server.svg)](https://pkg.go.dev/github.com/flaticols/server)

A Go package for parsing and comparing semantic versions with **full [SemVer 2.0.0](https://semver.org) compliance**.

## SemVer 2.0.0 Compliance

This package is **fully compliant** with the [Semantic Versioning 2.0.0 specification](https://semver.org). See [COMPLIANCE.md](COMPLIANCE.md) for detailed compliance information.

**Key compliance features:**
- All parsing rules correctly implemented
- All precedence rules correctly implemented  
- All validation rules correctly implemented
- All official examples from semver.org work correctly

**Optional 'v' prefix support:** While the SemVer spec doesn't include a 'v' prefix, this library optionally accepts it (e.g., `v1.2.3`) for compatibility with common tooling conventions (Git tags, npm, Go modules). Use `ParseStrict()` to enforce strict spec compliance.

## Features

- **Strict SemVer 2.0.0 Parsing**: Parse version strings with major.minor.patch format
- **Prerelease & Metadata Support**: Handle prerelease identifiers and build metadata
- **Version Comparison**: Compare versions with proper precedence rules
- **Constraint Matching**: Support for common constraint operators (`=`, `!=`, `>`, `<`, `>=`, `<=`, `~`, `^`)
- **Flexible Prefix Handling**: Parse versions with or without the `v` prefix (or use `ParseStrict()` to reject it)
- **Zero-allocation Operations**: Efficient string parsing and comparison

## Installation

```bash
go get github.com/flaticols/server
```

## Usage

### Parsing Versions

```go
package main

import (
    "fmt"
    "github.com/flaticols/server"
)

func main() {
    // Parse a version string (accepts optional 'v' prefix)
    v, err := server.Parse("1.2.3-beta.1+build.123")
    if err != nil {
        panic(err)
    }

    fmt.Println(v.String())   // "1.2.3-beta.1+build.123"
    fmt.Println(v.Stringv())  // "v1.2.3-beta.1+build.123"

    // Parse with 'v' prefix (also works)
    v2, _ := server.Parse("v1.2.3")
    fmt.Println(v2.String())  // "1.2.3"

    // Strict parsing (rejects 'v' prefix for SemVer 2.0.0 compliance)
    v3, err := server.ParseStrict("1.2.3")  // OK
    v4, err := server.ParseStrict("v1.2.3") // Error: 'v' prefix not allowed

    // Access version components
    fmt.Println(v.Major, v.Minor, v.Patch)  // 1 2 3
    fmt.Println(v.Prerelease)               // ["beta", "1"]
    fmt.Println(v.Metadata)                 // ["build", "123"]
}
```

### Comparing Versions

```go
v1, _ := server.Parse("1.0.0")
v2, _ := server.Parse("2.0.0")

// Direct comparison
fmt.Println(v1.LessThan(v2))           // true
fmt.Println(v1.Equal(v2))              // false
fmt.Println(v2.GreaterThan(v1))        // true

// Compare function returns -1, 0, or 1
result := server.Compare(v1, v2)       // -1
```

### Version Constraints

```go
// Parse a constraint
constraint, err := server.ParseConstraint(">= 1.0.0")
if err != nil {
    panic(err)
}

v, _ := server.Parse("1.5.0")
fmt.Println(constraint.Check(v))  // true

// Tilde operator: allows patch-level changes
c1, _ := server.ParseConstraint("~1.2.3")
c1.Check(server.Parse("1.2.5"))  // true
c1.Check(server.Parse("1.3.0"))  // false

// Caret operator: allows non-breaking changes
c2, _ := server.ParseConstraint("^1.2.3")
c2.Check(server.Parse("1.9.0"))  // true
c2.Check(server.Parse("2.0.0"))  // false

// Multiple constraints
cs, _ := server.ParseConstraintSet(">= 1.0.0, < 2.0.0")
cs.Check(server.Parse("1.5.0"))  // true
```

### Incrementing Versions

```go
v, _ := server.Parse("1.2.3-beta+build")

major := v.IncrementMajor()  // 2.0.0
minor := v.IncrementMinor()  // 1.3.0
patch := v.IncrementPatch()  // 1.2.4
```

### Working with Prerelease and Metadata

```go
v, _ := server.Parse("1.0.0")

// Set prerelease
v = server.SetPrerelease(v, []string{"alpha", "1"})
fmt.Println(v.String())  // "1.0.0-alpha.1"

// Set metadata
v = server.SetMetadata(v, []string{"commit", "abc123"})
fmt.Println(v.String())  // "1.0.0-alpha.1+commit.abc123"

// Get as map (assumes key-value pairs)
prereleaseMap := v.GetPrereleaseMap()
metadataMap := v.GetMetadataMap()
```

## Constraint Operators

| Operator | Description | Example | Matches |
|----------|-------------|---------|---------|
| `=`      | Equal to | `= 1.0.0` | 1.0.0 only |
| `!=`     | Not equal to | `!= 1.0.0` | Any except 1.0.0 |
| `>`      | Greater than | `> 1.0.0` | 1.0.1, 1.1.0, 2.0.0, etc. |
| `<`      | Less than | `< 2.0.0` | 1.x.x |
| `>=`     | Greater than or equal | `>= 1.0.0` | 1.0.0 and above |
| `<=`     | Less than or equal | `<= 2.0.0` | 2.0.0 and below |
| `~`      | Tilde (patch changes) | `~1.2.3` | 1.2.x (1.2.3 ≤ version < 1.3.0) |
| `^`      | Caret (compatible) | `^1.2.3` | 1.x.x (1.2.3 ≤ version < 2.0.0) |

## Error Handling

The package provides specific error types for common parsing failures:

```go
v, err := server.Parse("01.2.3")
if err == server.ErrLeadingZeros {
    fmt.Println("Version numbers cannot have leading zeros")
}
```

Available error types:
- `ErrEmptyVersion` - Empty version string
- `ErrMalformedCore` - Invalid major.minor.patch format
- `ErrLeadingZeros` - Leading zeros in version numbers
- `ErrNonDigitComponent` - Non-numeric characters in version numbers
- `ErrEmptyIdentifier` - Empty prerelease or metadata identifier
- `ErrLeadingZeroesIdentifier` - Leading zeros in numeric identifiers
- `ErrInvalidIdentifierChars` - Invalid characters in identifiers
- `ErrInvalidVPrefix` - 'v' prefix used with ParseStrict

## Validation

```go
v, valid := server.IsValid("1.2.3")
if valid {
    fmt.Println("Version is valid:", v.String())
}
```

## SemVer 2.0.0 Compliance

This package implements all requirements from the [Semantic Versioning 2.0.0 specification](https://semver.org):

- **Parsing**: All version formats defined in the spec are correctly parsed  
- **Precedence**: Version comparison follows all precedence rules exactly  
- **Validation**: All invalid versions are correctly rejected  
- **Official Examples**: All examples from semver.org work correctly

### Strict Mode

For applications requiring strict SemVer 2.0.0 compliance (no 'v' prefix), use `ParseStrict`:

```go
// Strict mode rejects 'v' prefix
v, err := server.ParseStrict("1.2.3")   // ✓ OK
v, err := server.ParseStrict("v1.2.3")  // ✗ Error: ErrInvalidVPrefix
```

### 'v' Prefix Support

The SemVer 2.0.0 specification does not include a 'v' prefix. However, many tools use it (Git tags, npm, Go modules). This library:

- `Parse()`: Accepts both `1.2.3` and `v1.2.3` for convenience
- `ParseStrict()`: Only accepts `1.2.3` (strict spec compliance)
- `String()`: Always outputs without 'v' (e.g., `1.2.3`)
- `Stringv()`: Always outputs with 'v' (e.g., `v1.2.3`)

See [COMPLIANCE.md](COMPLIANCE.md) for detailed compliance documentation.

## License

See LICENSE file for details
