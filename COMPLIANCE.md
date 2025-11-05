# SemVer 2.0.0 Compliance

This document details how this package complies with the [Semantic Versioning 2.0.0](https://semver.org) specification.

## Compliance Status

✅ **FULLY COMPLIANT** with SemVer 2.0.0 specification

## Specification Requirements

### 1. Version Format (§2)

✅ **Compliant**: Version format MUST be X.Y.Z where X, Y, and Z are non-negative integers.

```go
v, _ := server.Parse("1.2.3")  // Valid
v, _ := server.Parse("0.0.0")  // Valid (minimum version)
```

### 2. Leading Zeros (§2)

✅ **Compliant**: Version numbers MUST NOT contain leading zeroes.

```go
_, err := server.Parse("01.2.3")   // Error: leading zeros not allowed
_, err := server.Parse("1.02.3")   // Error: leading zeros not allowed
_, err := server.Parse("1.2.03")   // Error: leading zeros not allowed
```

### 3. Prerelease Versions (§9)

✅ **Compliant**: Prerelease version MAY be denoted by appending a hyphen and a series of dot-separated identifiers.

```go
v, _ := server.Parse("1.0.0-alpha")
v, _ := server.Parse("1.0.0-alpha.1")
v, _ := server.Parse("1.0.0-0.3.7")
v, _ := server.Parse("1.0.0-x.7.z.92")
```

✅ **Compliant**: Identifiers MUST comprise only ASCII alphanumerics and hyphens [0-9A-Za-z-].

```go
v, _ := server.Parse("1.0.0-alpha-beta")  // Valid
_, err := server.Parse("1.0.0-alpha_1")   // Error: underscore not allowed
```

✅ **Compliant**: Identifiers MUST NOT be empty.

```go
_, err := server.Parse("1.0.0-")           // Error: empty identifier
_, err := server.Parse("1.0.0-alpha..beta") // Error: empty identifier between dots
```

✅ **Compliant**: Numeric identifiers MUST NOT include leading zeroes.

```go
v, _ := server.Parse("1.0.0-0")     // Valid (single zero)
_, err := server.Parse("1.0.0-01")  // Error: leading zero in numeric identifier
v, _ := server.Parse("1.0.0-01alpha") // Valid (alphanumeric, not pure numeric)
```

### 4. Build Metadata (§10)

✅ **Compliant**: Build metadata MAY be denoted by appending a plus sign and a series of dot-separated identifiers.

```go
v, _ := server.Parse("1.0.0+build")
v, _ := server.Parse("1.0.0+20130313144700")
v, _ := server.Parse("1.0.0-beta+exp.sha.5114f85")
```

✅ **Compliant**: Identifiers MUST comprise only ASCII alphanumerics and hyphens [0-9A-Za-z-].

✅ **Compliant**: Build metadata SHOULD be ignored when determining version precedence.

```go
v1, _ := server.Parse("1.0.0+build1")
v2, _ := server.Parse("1.0.0+build2")
server.Compare(v1, v2) // Returns 0 (equal)
```

### 5. Precedence Rules (§11)

✅ **Compliant**: All precedence rules are correctly implemented.

#### 5.1 Major, Minor, Patch Precedence

```go
// 1.0.0 < 2.0.0 < 2.1.0 < 2.1.1
v1, _ := server.Parse("1.0.0")
v2, _ := server.Parse("2.0.0")
v1.LessThan(v2) // true
```

#### 5.2 Prerelease Precedence

✅ **Compliant**: Prerelease versions have lower precedence than normal versions.

```go
// 1.0.0-alpha < 1.0.0
v1, _ := server.Parse("1.0.0-alpha")
v2, _ := server.Parse("1.0.0")
v1.LessThan(v2) // true
```

#### 5.3 Prerelease Comparison

✅ **Compliant**: Official example from semver.org:

```go
// 1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta < 1.0.0-beta < 
// 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0
```

All comparisons work correctly:
- Numeric identifiers are compared as integers
- Alphanumeric identifiers are compared lexically in ASCII sort order
- Numeric identifiers always have lower precedence than non-numeric identifiers
- A larger set of prerelease fields has higher precedence than a smaller set

## Deviations from Specification

### 'v' Prefix (Non-Standard Extension)

⚠️ **Intentional Deviation**: This library accepts an optional 'v' prefix on version strings.

**Rationale**: While the SemVer 2.0.0 specification only defines versions in the format `X.Y.Z[-prerelease][+metadata]`, many tools (Git tags, npm, Go modules) commonly use a `v` prefix (e.g., `v1.2.3`).

**Behavior**:
```go
v1, _ := server.Parse("1.2.3")   // Parses successfully
v2, _ := server.Parse("v1.2.3")  // Also parses successfully

v1.Equal(v2)          // true - treated as identical
v1.String()           // "1.2.3" - canonical form without 'v'
v2.Stringv()          // "v1.2.3" - with 'v' prefix
```

**Strict Mode**: For strict SemVer 2.0.0 compliance, use `ParseStrict()`:
```go
v, err := server.ParseStrict("v1.2.3")  // Error: 'v' prefix not allowed in strict mode
v, _   := server.ParseStrict("1.2.3")   // OK
```

## Testing

The package includes comprehensive tests that validate:
- All official examples from semver.org
- All parsing rules from the specification
- All precedence rules from the specification
- All validation rules from the specification
- Edge cases and corner cases

Run compliance tests:
```bash
go test -v -run TestPrereleaseComparisons
go test -v -run TestBuildMetadataHandling
go test -v -run TestInvalidVersionCharacters
```

## References

- [Semantic Versioning 2.0.0 Specification](https://semver.org)
- [SemVer FAQ](https://semver.org/#faq)

## Validation

Last validated: 2025-11-05

All requirements from SemVer 2.0.0 specification have been verified and are correctly implemented.
