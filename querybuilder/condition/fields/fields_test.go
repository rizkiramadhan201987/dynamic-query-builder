package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldCondition_Build(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		operator      ComparisonOperator
		value         interface{}
		paramOffset   int
		expectedSQL   string
		expectedArgs  []interface{}
		expectedError string
	}{
		{
			name:         "Equals operator",
			field:        "column_name",
			operator:     Equals,
			value:        42,
			paramOffset:  1,
			expectedSQL:  "column_name = $1",
			expectedArgs: []interface{}{42},
		},
		{
			name:         "NotEquals operator",
			field:        "column_name",
			operator:     NotEquals,
			value:        "test",
			paramOffset:  2,
			expectedSQL:  "column_name != $2",
			expectedArgs: []interface{}{"test"},
		},
		{
			name:         "GreaterThan operator",
			field:        "age",
			operator:     GreaterThan,
			value:        18,
			paramOffset:  1,
			expectedSQL:  "age > $1",
			expectedArgs: []interface{}{18},
		},
		{
			name:         "LessThan operator",
			field:        "price",
			operator:     LessThan,
			value:        99.99,
			paramOffset:  3,
			expectedSQL:  "price < $3",
			expectedArgs: []interface{}{99.99},
		},
		{
			name:         "GreaterOrEqual operator",
			field:        "quantity",
			operator:     GreaterOrEqual,
			value:        10,
			paramOffset:  1,
			expectedSQL:  "quantity >= $1",
			expectedArgs: []interface{}{10},
		},
		{
			name:         "LessOrEqual operator",
			field:        "weight",
			operator:     LessOrEqual,
			value:        50.5,
			paramOffset:  4,
			expectedSQL:  "weight <= $4",
			expectedArgs: []interface{}{50.5},
		},
		{
			name:         "Like operator",
			field:        "name",
			operator:     Like,
			value:        "%John%",
			paramOffset:  1,
			expectedSQL:  "name LIKE $1",
			expectedArgs: []interface{}{"%John%"},
		},
		{
			name:         "ILike operator",
			field:        "email",
			operator:     ILike,
			value:        "%.com",
			paramOffset:  2,
			expectedSQL:  "email ILIKE $2",
			expectedArgs: []interface{}{"%.com"},
		},
		{
			name:         "In operator with valid slice",
			field:        "status",
			operator:     In,
			value:        []interface{}{"active", "pending"},
			paramOffset:  1,
			expectedSQL:  "status IN ($1,$2)",
			expectedArgs: []interface{}{"active", "pending"},
		},
		{
			name:         "NotIn operator with valid slice",
			field:        "category",
			operator:     NotIn,
			value:        []interface{}{1, 2, 3},
			paramOffset:  1,
			expectedSQL:  "category NOT IN ($1,$2,$3)",
			expectedArgs: []interface{}{1, 2, 3},
		},
		{
			name:         "IsNull operator",
			field:        "deleted_at",
			operator:     IsNull,
			value:        nil,
			paramOffset:  1,
			expectedSQL:  "deleted_at IS NULL",
			expectedArgs: nil,
		},
		{
			name:         "IsNotNull operator",
			field:        "updated_at",
			operator:     IsNotNull,
			value:        nil,
			paramOffset:  1,
			expectedSQL:  "updated_at IS NOT NULL",
			expectedArgs: nil,
		},
		{
			name:          "In operator with invalid value type",
			field:         "status",
			operator:      In,
			value:         "not a slice",
			paramOffset:   1,
			expectedError: "value for IN/NOT IN operator must be a slice",
		},
		{
			name:          "In operator with empty slice",
			field:         "status",
			operator:      In,
			value:         []interface{}{},
			paramOffset:   1,
			expectedError: "empty slice provided for IN/NOT IN operator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := NewFieldCondition(tt.field, tt.operator, tt.value)
			sql, args, err := fc.Build(tt.paramOffset)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedSQL, sql)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}

func TestNewFieldCondition(t *testing.T) {
	field := "test_field"
	operator := Equals
	value := "test_value"

	fc := NewFieldCondition(field, operator, value)

	assert.NotNil(t, fc)
	assert.Equal(t, field, fc.Field)
	assert.Equal(t, operator, fc.Operator)
	assert.Equal(t, value, fc.Value)
}
