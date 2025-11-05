package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseStrict(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expectedErr error
	}{
		// Valid versions without 'v' prefix
		{
			name:        "basic version",
			input:       "1.2.3",
			expectError: false,
		},
		{
			name:        "version with prerelease",
			input:       "1.2.3-alpha",
			expectError: false,
		},
		{
			name:        "version with metadata",
			input:       "1.2.3+build",
			expectError: false,
		},
		{
			name:        "version with both prerelease and metadata",
			input:       "1.2.3-alpha+build",
			expectError: false,
		},
		{
			name:        "zero version",
			input:       "0.0.0",
			expectError: false,
		},
		{
			name:        "complex prerelease",
			input:       "1.0.0-alpha.beta.1",
			expectError: false,
		},

		// Invalid: versions with 'v' prefix
		{
			name:        "reject v prefix",
			input:       "v1.2.3",
			expectError: true,
			expectedErr: ErrInvalidVPrefix,
		},
		{
			name:        "reject v prefix with prerelease",
			input:       "v1.2.3-alpha",
			expectError: true,
			expectedErr: ErrInvalidVPrefix,
		},
		{
			name:        "reject v prefix with metadata",
			input:       "v1.2.3+build",
			expectError: true,
			expectedErr: ErrInvalidVPrefix,
		},

		// Other invalid versions
		{
			name:        "empty string",
			input:       "",
			expectError: true,
			expectedErr: ErrEmptyVersion,
		},
		{
			name:        "leading zeros",
			input:       "01.2.3",
			expectError: true,
			expectedErr: ErrLeadingZeros,
		},
		{
			name:        "malformed core",
			input:       "1.2",
			expectError: true,
			expectedErr: ErrMalformedCore,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := ParseStrict(tt.input)

			if tt.expectError {
				require.Error(t, err, "Expected error for input: %s", tt.input)
				if tt.expectedErr != nil {
					require.ErrorIs(t, err, tt.expectedErr,
						"Expected error %v, got %v", tt.expectedErr, err)
				}
			} else {
				require.NoError(t, err, "Unexpected error for input: %s", tt.input)
				require.NotNil(t, v)
			}
		})
	}
}

func TestParseVsParsStrict(t *testing.T) {
	// Test that Parse accepts 'v' prefix but ParseStrict rejects it
	t.Run("Parse accepts v prefix", func(t *testing.T) {
		v, err := Parse("v1.2.3")
		require.NoError(t, err)
		require.Equal(t, 1, v.Major)
		require.Equal(t, 2, v.Minor)
		require.Equal(t, 3, v.Patch)
	})

	t.Run("ParseStrict rejects v prefix", func(t *testing.T) {
		_, err := ParseStrict("v1.2.3")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrInvalidVPrefix)
	})

	t.Run("Both accept standard format", func(t *testing.T) {
		v1, err1 := Parse("1.2.3")
		require.NoError(t, err1)

		v2, err2 := ParseStrict("1.2.3")
		require.NoError(t, err2)

		require.Equal(t, v1, v2)
	})
}

func TestParseStrictCompliance(t *testing.T) {
	// All official examples from semver.org should work with ParseStrict
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
			v, err := ParseStrict(ex)
			require.NoError(t, err, "Failed to parse official example: %s", ex)
			require.NotNil(t, v)
		})
	}
}
