package querybuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockCondition implements QueryCondition for testing
type MockCondition struct {
	sql  string
	args []interface{}
	err  error
}

func (m MockCondition) Build(paramOffset int) (string, []interface{}, error) {
	if m.err != nil {
		return "", nil, m.err
	}

	// Adjust parameter placeholders based on offset
	if len(m.args) > 0 {
		finalSQL := m.sql
		for i := range m.args {
			finalSQL = replaceNthOccurrence(finalSQL, "$1", "$"+string(rune('0'+paramOffset+i)), 1)
		}
		return finalSQL, m.args, nil
	}
	return m.sql, m.args, nil
}

// Helper function to replace nth occurrence of a substring
func replaceNthOccurrence(s, old, new string, n int) string {
	count := 0
	for i := 0; i < len(s); i++ {
		if len(s[i:]) < len(old) {
			break
		}
		if s[i:i+len(old)] == old {
			count++
			if count == n {
				return s[:i] + new + s[i+len(old):]
			}
		}
	}
	return s
}

func TestWhereGroup_Build(t *testing.T) {
	tests := []struct {
		name          string
		operator      LogicalOperator
		conditions    []QueryCondition
		paramOffset   int
		expectedSQL   string
		expectedArgs  []interface{}
		expectedError bool
		errorContains string
	}{
		{
			name:          "Empty group",
			operator:      AND,
			conditions:    []QueryCondition{},
			paramOffset:   1,
			expectedSQL:   "",
			expectedArgs:  nil,
			expectedError: false,
		},
		{
			name:     "Single condition",
			operator: AND,
			conditions: []QueryCondition{
				MockCondition{
					sql:  "column = $1",
					args: []interface{}{"value"},
				},
			},
			paramOffset:   1,
			expectedSQL:   "column = $1",
			expectedArgs:  []interface{}{"value"},
			expectedError: false,
		},
		{
			name:     "Two AND conditions",
			operator: AND,
			conditions: []QueryCondition{
				MockCondition{
					sql:  "column1 = $1",
					args: []interface{}{"value1"},
				},
				MockCondition{
					sql:  "column2 = $1",
					args: []interface{}{"value2"},
				},
			},
			paramOffset:   1,
			expectedSQL:   "(column1 = $1 AND column2 = $2)",
			expectedArgs:  []interface{}{"value1", "value2"},
			expectedError: false,
		},
		{
			name:     "Two OR conditions",
			operator: OR,
			conditions: []QueryCondition{
				MockCondition{
					sql:  "column1 = $1",
					args: []interface{}{"value1"},
				},
				MockCondition{
					sql:  "column2 = $1",
					args: []interface{}{"value2"},
				},
			},
			paramOffset:   1,
			expectedSQL:   "(column1 = $1 OR column2 = $2)",
			expectedArgs:  []interface{}{"value1", "value2"},
			expectedError: false,
		},
		{
			name:     "Mixed conditions with offset",
			operator: AND,
			conditions: []QueryCondition{
				MockCondition{
					sql:  "column1 = $1",
					args: []interface{}{"value1"},
				},
				MockCondition{
					sql:  "column2 > $1",
					args: []interface{}{42},
				},
			},
			paramOffset:   3,
			expectedSQL:   "(column1 = $3 AND column2 > $4)",
			expectedArgs:  []interface{}{"value1", 42},
			expectedError: false,
		},
		{
			name:     "Condition with error",
			operator: AND,
			conditions: []QueryCondition{
				MockCondition{
					err: assert.AnError,
				},
			},
			paramOffset:   1,
			expectedError: true,
			errorContains: assert.AnError.Error(),
		},
		{
			name:     "Multiple conditions with no args",
			operator: AND,
			conditions: []QueryCondition{
				MockCondition{
					sql: "column1 IS NULL",
				},
				MockCondition{
					sql: "column2 IS NOT NULL",
				},
			},
			paramOffset:   1,
			expectedSQL:   "(column1 IS NULL AND column2 IS NOT NULL)",
			expectedArgs:  nil,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			whereGroup := NewWhereGroup(tt.operator)
			for _, condition := range tt.conditions {
				whereGroup.Add(condition)
			}

			sql, args, err := whereGroup.Build(tt.paramOffset)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedSQL, sql)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}

func TestNewWhereGroup(t *testing.T) {
	tests := []struct {
		name          string
		operator      LogicalOperator
		expectedGroup *WhereGroup
	}{
		{
			name:     "New AND group",
			operator: AND,
			expectedGroup: &WhereGroup{
				Conditions: make([]QueryCondition, 0),
				Operator:   AND,
				IsNew:      true,
			},
		},
		{
			name:     "New OR group",
			operator: OR,
			expectedGroup: &WhereGroup{
				Conditions: make([]QueryCondition, 0),
				Operator:   OR,
				IsNew:      true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := NewWhereGroup(tt.operator)
			assert.Equal(t, tt.expectedGroup.Operator, group.Operator)
			assert.Equal(t, tt.expectedGroup.IsNew, group.IsNew)
			assert.Empty(t, group.Conditions)
		})
	}
}
