package simpleselect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSimpleSelect(t *testing.T) {
	tests := []struct {
		name          string
		fields        []string
		expectedField []string
		description   string
	}{
		{
			name:          "No fields provided",
			fields:        []string{},
			expectedField: []string{"*"},
			description:   "Should default to SELECT * when no fields are provided",
		},
		{
			name:          "Single field",
			fields:        []string{"name"},
			expectedField: []string{"name"},
			description:   "Should create SELECT with a single field",
		},
		{
			name:          "Multiple fields",
			fields:        []string{"id", "name", "email"},
			expectedField: []string{"id", "name", "email"},
			description:   "Should create SELECT with multiple fields",
		},
		{
			name:          "With table qualified fields",
			fields:        []string{"users.id", "users.name", "profiles.avatar"},
			expectedField: []string{"users.id", "users.name", "profiles.avatar"},
			description:   "Should handle table qualified field names",
		},
		{
			name:          "With aliased fields",
			fields:        []string{"id AS user_id", "name AS full_name"},
			expectedField: []string{"id AS user_id", "name AS full_name"},
			description:   "Should handle aliased field names",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.description)
			select_ := NewSimpleSelect(tt.fields...)
			assert.Equal(t, tt.expectedField, select_.fields,
				"Fields should match expected values")
		})
	}
}

func TestSimpleSelectBuild(t *testing.T) {
	tests := []struct {
		name        string
		fields      []string
		expected    string
		expectError bool
		description string
	}{
		{
			name:        "Default select all",
			fields:      []string{},
			expected:    "SELECT *",
			expectError: false,
			description: "Should build SELECT * query when no fields specified",
		},
		{
			name:        "Single field",
			fields:      []string{"name"},
			expected:    "SELECT name",
			expectError: false,
			description: "Should build query with single field",
		},
		{
			name:        "Multiple fields",
			fields:      []string{"id", "name", "email"},
			expected:    "SELECT id, name, email",
			expectError: false,
			description: "Should build query with multiple fields",
		},
		{
			name:        "Table qualified fields",
			fields:      []string{"users.id", "users.name", "profiles.avatar"},
			expected:    "SELECT users.id, users.name, profiles.avatar",
			expectError: false,
			description: "Should build query with table qualified fields",
		},
		{
			name:        "Aliased fields",
			fields:      []string{"id AS user_id", "name AS full_name"},
			expected:    "SELECT id AS user_id, name AS full_name",
			expectError: false,
			description: "Should build query with aliased fields",
		},
		{
			name:        "With functions",
			fields:      []string{"COUNT(*)", "MAX(score)", "MIN(created_at)"},
			expected:    "SELECT COUNT(*), MAX(score), MIN(created_at)",
			expectError: false,
			description: "Should build query with SQL functions",
		},
		{
			name:        "Mixed field types",
			fields:      []string{"users.id AS user_id", "COUNT(*) as total", "MAX(score) as high_score"},
			expected:    "SELECT users.id AS user_id, COUNT(*) as total, MAX(score) as high_score",
			expectError: false,
			description: "Should build query with mixed field types including qualified names, aliases, and functions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.description)
			select_ := NewSimpleSelect(tt.fields...)
			result, err := select_.Build()

			if tt.expectError {
				assert.Error(t, err, "Should return an error")
			} else {
				assert.NoError(t, err, "Should not return an error")
				assert.Equal(t, tt.expected, result,
					"Built query should match expected SQL")
			}
		})
	}
}

func TestSimpleSelectEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		fields      []string
		expected    string
		expectError bool
		description string
	}{
		{
			name:        "Empty fields array",
			fields:      []string{},
			expected:    "SELECT *",
			expectError: false,
			description: "Should handle empty fields array gracefully",
		},
		{
			name:        "Fields with special characters",
			fields:      []string{"`special.field`", "\"quoted.field\"", "[bracketed.field]"},
			expected:    "SELECT `special.field`, \"quoted.field\", [bracketed.field]",
			expectError: false,
			description: "Should handle fields with special characters and different quote styles",
		},
		{
			name:        "Fields with whitespace",
			fields:      []string{"  id  ", " name ", "email  "},
			expected:    "SELECT   id  ,  name , email  ",
			expectError: false,
			description: "Should preserve whitespace in field names",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.description)
			select_ := NewSimpleSelect(tt.fields...)
			result, err := select_.Build()

			if tt.expectError {
				assert.Error(t, err, "Should return an error")
			} else {
				assert.NoError(t, err, "Should not return an error")
				assert.Equal(t, tt.expected, result,
					"Built query should match expected SQL")
			}
		})
	}
}
