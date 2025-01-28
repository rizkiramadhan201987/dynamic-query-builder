package fields

import (
	"fmt"
	"strings"
)

// ComparisonOperator represents SQL comparison operators
type ComparisonOperator string

const (
	Equals         ComparisonOperator = "="
	NotEquals      ComparisonOperator = "!="
	GreaterThan    ComparisonOperator = ">"
	LessThan       ComparisonOperator = "<"
	GreaterOrEqual ComparisonOperator = ">="
	LessOrEqual    ComparisonOperator = "<="
	Like           ComparisonOperator = "LIKE"
	ILike          ComparisonOperator = "ILIKE"
	In             ComparisonOperator = "IN"
	NotIn          ComparisonOperator = "NOT IN"
	IsNull         ComparisonOperator = "IS NULL"
	IsNotNull      ComparisonOperator = "IS NOT NULL"
)

// FieldCondition represents a single field condition
type FieldCondition struct {
	Field    string             // Column name
	Operator ComparisonOperator // Comparison operator
	Value    interface{}        // Value to compare against
}

// NewFieldCondition creates a new field condition
func NewFieldCondition(field string, operator ComparisonOperator, value interface{}) *FieldCondition {
	return &FieldCondition{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

// Build implements the QueryCondition interface
func (fc *FieldCondition) Build(paramOffset int) (string, []interface{}, error) {
	// Handle special cases for NULL comparisons
	switch fc.Operator {
	case IsNull:
		return fmt.Sprintf("%s IS NULL", fc.Field), nil, nil
	case IsNotNull:
		return fmt.Sprintf("%s IS NOT NULL", fc.Field), nil, nil
	}

	// Handle IN and NOT IN operators
	if fc.Operator == In || fc.Operator == NotIn {
		values, ok := fc.Value.([]interface{})
		if !ok {
			return "", nil, fmt.Errorf("value for IN/NOT IN operator must be a slice")
		}
		if len(values) == 0 {
			return "", nil, fmt.Errorf("empty slice provided for IN/NOT IN operator")
		}

		// Build the parameter placeholders
		placeholders := make([]string, len(values))
		for i := range values {
			placeholders[i] = fmt.Sprintf("$%d", paramOffset+i)
		}
		return fmt.Sprintf("%s %s (%s)", fc.Field, fc.Operator,
			strings.Join(placeholders, ",")), values, nil
	}

	// Regular comparison
	return fmt.Sprintf("%s %s $%d", fc.Field, fc.Operator, paramOffset),
		[]interface{}{fc.Value}, nil
}
