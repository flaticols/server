package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSemVer2Compliance validates complete compliance with SemVer 2.0.0 specification
// as defined at https://semver.org
func TestSemVer2Compliance(t *testing.T) {
	t.Run("Official Examples", func(t *testing.T) {
		// All examples from https://semver.org should parse correctly
		examples := []string{
			"1.0.0",
			"1.0.0-alpha",
			"1.0.0-alpha.1",
			"1.0.0-alpha.beta",
			"1.0.0-beta",
			"1.0.0-beta.2",
			"1.0.0-beta.11",
			"1.0.0-rc.1",
			"1.0.0-0.3.7",
			"1.0.0-x.7.z.92",
			"1.0.0+20130313144700",
			"1.0.0-beta+exp.sha.5114f85",
			"1.0.0+21AF26D3-117B344092BD",
		}

		for _, ex := range examples {
			t.Run(ex, func(t *testing.T) {
				v, err := Parse(ex)
				require.NoError(t, err, "Failed to parse official example: %s", ex)
				require.NotNil(t, v)
			})
		}
	})

	t.Run("Spec Example Precedence", func(t *testing.T) {
		// From https://semver.org/#spec-item-11
		// "1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta < 1.0.0-beta <
		// 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0"
		versions := []string{
			"1.0.0-alpha",
			"1.0.0-alpha.1",
			"1.0.0-alpha.beta",
			"1.0.0-beta",
			"1.0.0-beta.2",
			"1.0.0-beta.11",
			"1.0.0-rc.1",
			"1.0.0",
		}

		for i := 0; i < len(versions)-1; i++ {
			v1, err1 := Parse(versions[i])
			v2, err2 := Parse(versions[i+1])

			require.NoError(t, err1)
			require.NoError(t, err2)
			require.True(t, v1.LessThan(v2),
				"Expected %s < %s per spec example", versions[i], versions[i+1])
		}
	})

	t.Run("Parsing Rules", func(t *testing.T) {
		tests := []struct {
			name    string
			input   string
			valid   bool
			reason  string
		}{
			// Valid versions
			{"basic", "1.2.3", true, "basic version"},
			{"zero", "0.0.0", true, "zero version"},
			{"large numbers", "999999999.999999999.999999999", true, "large valid numbers"},
			{"prerelease single", "1.0.0-alpha", true, "single prerelease identifier"},
			{"prerelease multiple", "1.0.0-alpha.beta.gamma", true, "multiple prerelease identifiers"},
			{"prerelease numeric", "1.0.0-0", true, "single zero prerelease"},
			{"prerelease mixed", "1.0.0-alpha.1.beta.2", true, "mixed prerelease identifiers"},
			{"metadata single", "1.0.0+build", true, "single metadata identifier"},
			{"metadata multiple", "1.0.0+build.1.2.3", true, "multiple metadata identifiers"},
			{"both", "1.0.0-alpha+build", true, "both prerelease and metadata"},

			// Invalid versions per spec
			{"missing components 1", "1", false, "missing minor and patch"},
			{"missing components 2", "1.2", false, "missing patch"},
			{"too many components", "1.2.3.4", false, "too many components"},
			{"leading zero major", "01.2.3", false, "leading zero in major"},
			{"leading zero minor", "1.02.3", false, "leading zero in minor"},
			{"leading zero patch", "1.2.03", false, "leading zero in patch"},
			{"empty prerelease", "1.2.3-", false, "empty prerelease"},
			{"empty metadata", "1.2.3+", false, "empty metadata"},
			{"numeric prerelease leading zero", "1.2.3-01", false, "leading zero in numeric prerelease"},
			{"empty identifier in prerelease", "1.2.3-alpha..beta", false, "empty identifier"},
			{"underscore", "1.2.3-alpha_1", false, "underscore not allowed"},
			{"invalid char", "1.2.3-alpha!1", false, "invalid character"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := Parse(tt.input)
				if tt.valid {
					require.NoError(t, err, "Expected valid: %s (%s)", tt.input, tt.reason)
				} else {
					require.Error(t, err, "Expected invalid: %s (%s)", tt.input, tt.reason)
				}
			})
		}
	})

	t.Run("Identifier Character Set", func(t *testing.T) {
		// Spec: "Identifiers MUST comprise only ASCII alphanumerics and hyphens [0-9A-Za-z-]"
		
		// Valid characters
		valid := []string{
			"1.0.0-abc",
			"1.0.0-ABC",
			"1.0.0-123",
			"1.0.0-a-b-c",
			"1.0.0-aB1-Cd2",
			"1.0.0+abc-123-XYZ",
		}
		for _, v := range valid {
			_, err := Parse(v)
			require.NoError(t, err, "Should accept valid characters: %s", v)
		}

		// Invalid characters
		invalid := []string{
			"1.0.0-alpha_beta",  // underscore
			"1.0.0-alpha.beta!", // exclamation
			"1.0.0-alpha@beta",  // at sign
			"1.0.0+build#123",   // hash
			"1.0.0-alpha$beta",  // dollar sign
		}
		for _, v := range invalid {
			_, err := Parse(v)
			require.Error(t, err, "Should reject invalid characters: %s", v)
		}
	})

	t.Run("Precedence Rules", func(t *testing.T) {
		tests := []struct {
			v1   string
			v2   string
			cmp  int // -1: v1 < v2, 0: equal, 1: v1 > v2
			desc string
		}{
			// Major, minor, patch precedence
			{"1.0.0", "2.0.0", -1, "major version precedence"},
			{"2.0.0", "2.1.0", -1, "minor version precedence"},
			{"2.1.0", "2.1.1", -1, "patch version precedence"},

			// Prerelease precedence
			{"1.0.0-alpha", "1.0.0", -1, "prerelease < release"},
			{"1.0.0", "1.0.0-alpha", 1, "release > prerelease"},

			// Numeric identifier comparison
			{"1.0.0-1", "1.0.0-2", -1, "numeric identifiers compared as integers"},
			{"1.0.0-2", "1.0.0-10", -1, "numeric: 2 < 10"},
			{"1.0.0-10", "1.0.0-100", -1, "numeric: 10 < 100"},

			// Alphanumeric identifier comparison
			{"1.0.0-alpha", "1.0.0-beta", -1, "alphanumeric lexical comparison"},
			{"1.0.0-alpha", "1.0.0-alpha", 0, "identical alphanumeric"},

			// Numeric < alphanumeric
			{"1.0.0-1", "1.0.0-alpha", -1, "numeric < alphanumeric"},
			{"1.0.0-99", "1.0.0-alpha", -1, "numeric always less than alphanumeric"},

			// Larger set has higher precedence
			{"1.0.0-alpha", "1.0.0-alpha.1", -1, "larger set > smaller set"},
			{"1.0.0-alpha.beta", "1.0.0-alpha.beta.gamma", -1, "3 identifiers > 2 identifiers"},

			// Build metadata ignored
			{"1.0.0+build1", "1.0.0+build2", 0, "metadata ignored"},
			{"1.0.0-alpha+build1", "1.0.0-alpha+build2", 0, "metadata ignored with prerelease"},
		}

		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				v1, err1 := Parse(tt.v1)
				v2, err2 := Parse(tt.v2)

				require.NoError(t, err1)
				require.NoError(t, err2)

				result := Compare(v1, v2)
				require.Equal(t, tt.cmp, result,
					"Precedence rule failed: %s vs %s (%s)", tt.v1, tt.v2, tt.desc)
			})
		}
	})

	t.Run("Leading Zeros in Identifiers", func(t *testing.T) {
		// Numeric identifiers MUST NOT include leading zeros
		_, err := Parse("1.0.0-01")
		require.Error(t, err, "Numeric identifier with leading zero should fail")

		// But alphanumeric identifiers can have leading zeros
		_, err = Parse("1.0.0-01alpha")
		require.NoError(t, err, "Alphanumeric identifier can have leading zero pattern")

		// Metadata can have any pattern (no restriction on leading zeros)
		_, err = Parse("1.0.0+001")
		require.NoError(t, err, "Metadata can have leading zeros")
	})

	t.Run("Empty Identifiers", func(t *testing.T) {
		// Empty identifiers not allowed
		invalid := []string{
			"1.0.0-",           // empty prerelease
			"1.0.0+",           // empty metadata
			"1.0.0-alpha..beta", // empty identifier in middle
			"1.0.0-.alpha",     // empty first identifier
			"1.0.0-alpha.",     // empty last identifier
			"1.0.0+alpha..beta", // empty identifier in metadata
		}

		for _, v := range invalid {
			_, err := Parse(v)
			require.Error(t, err, "Empty identifier should be rejected: %s", v)
		}
	})

	t.Run("Build Metadata Ignored in Comparison", func(t *testing.T) {
		// Build metadata SHOULD be ignored when determining version precedence
		pairs := []struct {
			v1 string
			v2 string
		}{
			{"1.0.0+a", "1.0.0+b"},
			{"1.0.0+build1", "1.0.0+build2"},
			{"1.0.0+20130313144700", "1.0.0+20140313144700"},
			{"1.0.0-alpha+001", "1.0.0-alpha+002"},
		}

		for _, p := range pairs {
			v1, _ := Parse(p.v1)
			v2, _ := Parse(p.v2)

			require.Equal(t, 0, Compare(v1, v2),
				"Build metadata should be ignored: %s vs %s", p.v1, p.v2)
		}
	})
}
