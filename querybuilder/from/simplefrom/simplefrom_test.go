package simplefrom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleFrom(t *testing.T) {
	tests := []struct {
		name           string
		tableName      string
		expectedOutput string
		description    string
	}{
		{
			name:           "Basic Table Name",
			tableName:      "users",
			expectedOutput: "FROM users",
			description:    "Should handle basic table name correctly",
		},
		{
			name:           "Table Name with Schema",
			tableName:      "public.users",
			expectedOutput: "FROM public.users",
			description:    "Should handle schema-qualified table names",
		},
		{
			name:           "Table Name with Special Characters",
			tableName:      "user_data",
			expectedOutput: "FROM user_data",
			description:    "Should handle table names with underscores",
		},
		{
			name:           "Table Name with Mixed Case",
			tableName:      "UserData",
			expectedOutput: "FROM UserData",
			description:    "Should preserve case in table names",
		},
		{
			name:           "Complex Table Name",
			tableName:      "myschema.user_data_2023",
			expectedOutput: "FROM myschema.user_data_2023",
			description:    "Should handle complex table names with schema, underscores and numbers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test setup
			t.Log(tt.description)
			simpleFrom := NewSimpleFrom(tt.tableName)

			// Execute test
			result := simpleFrom.Build()

			// Assert results
			assert.Equal(t, tt.expectedOutput, result,
				"Expected '%s' but got '%s'",
				tt.expectedOutput, result)
		})
	}
}

func TestNewSimpleFrom(t *testing.T) {
	t.Run("Constructor Test", func(t *testing.T) {
		// Test setup
		tableName := "test_table"

		// Execute test
		simpleFrom := NewSimpleFrom(tableName)

		// Assert results
		assert.NotNil(t, simpleFrom, "NewSimpleFrom should not return nil")
		assert.Equal(t, tableName, simpleFrom.Table,
			"Table name not correctly set in constructor")
	})
}

func TestSimpleFromImplementation(t *testing.T) {
	t.Run("Multiple Builds", func(t *testing.T) {
		// Test setup
		simpleFrom := NewSimpleFrom("test_table")

		// Execute test multiple times
		result1 := simpleFrom.Build()
		result2 := simpleFrom.Build()

		// Assert results
		assert.Equal(t, result1, result2,
			"Multiple Build() calls should return consistent results")
	})
}
