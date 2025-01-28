package querybuilder

import (
	"fmt"
	"strings"
)

// LogicalOperator represents SQL logical operators
type LogicalOperator string

const (
	AND LogicalOperator = "AND"
	OR  LogicalOperator = "OR"
)

// WhereGroup represents a group of conditions with a logical operator
type WhereGroup struct {
	Conditions []QueryCondition
	Operator   LogicalOperator
	IsNew      bool // Track if this group has been added to whereGroups
}

// NewWhereGroup creates a new where group with specified operator
func NewWhereGroup(operator LogicalOperator) *WhereGroup {
	return &WhereGroup{
		Conditions: make([]QueryCondition, 0),
		Operator:   operator,
		IsNew:      true,
	}
}

// Add adds a condition to the where group
func (wg *WhereGroup) Add(condition QueryCondition) *WhereGroup {
	wg.Conditions = append(wg.Conditions, condition)
	return wg
}

// Build implements QueryCondition interface
func (wg *WhereGroup) Build(paramOffset int) (string, []interface{}, error) {
	if len(wg.Conditions) == 0 {
		return "", nil, nil
	}

	var clauses []string
	var args []interface{}

	for _, cond := range wg.Conditions {
		clause, condArgs, err := cond.Build(paramOffset + len(args))
		if err != nil {
			return "", nil, err
		}
		if clause != "" {
			clauses = append(clauses, clause)
			args = append(args, condArgs...)
		}
	}

	if len(clauses) == 0 {
		return "", nil, nil
	}

	// If there's only one condition, don't wrap in parentheses
	if len(clauses) == 1 {
		return clauses[0], args, nil
	}

	// Multiple conditions are wrapped in parentheses
	return fmt.Sprintf("(%s)", strings.Join(clauses, fmt.Sprintf(" %s ", wg.Operator))), args, nil
}
