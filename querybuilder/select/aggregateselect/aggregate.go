package aggregate

import (
	"fmt"
	"strings"
)

// AggregateFunction represents different SQL aggregate functions
type AggregateFunction string

const (
	Sum   AggregateFunction = "SUM"
	Avg   AggregateFunction = "AVG"
	Count AggregateFunction = "COUNT"
	Max   AggregateFunction = "MAX"
	Min   AggregateFunction = "MIN"
)

// AggregateField represents a field with its aggregate function
type AggregateField struct {
	Function AggregateFunction
	Field    string
	Alias    string
}

// AggregateSelect implements SELECT with aggregate functions
type AggregateSelect struct {
	regularFields []string
	aggregates    []AggregateField
}

func NewAggregateSelect() *AggregateSelect {
	return &AggregateSelect{
		regularFields: make([]string, 0),
		aggregates:    make([]AggregateField, 0),
	}
}
func (as *AggregateSelect) AddRegularField(field string) *AggregateSelect {
	as.regularFields = append(as.regularFields, field)
	return as
}
func (as *AggregateSelect) AddAggregate(fn AggregateFunction, field, alias string) *AggregateSelect {
	as.aggregates = append(as.aggregates, AggregateField{
		Function: fn,
		Field:    field,
		Alias:    alias,
	})
	return as
}

func (as *AggregateSelect) Build() (string, error) {
	var fields []string

	// Add regular fields
	fields = append(fields, as.regularFields...)

	// Add aggregate fields
	for _, agg := range as.aggregates {
		if agg.Alias != "" {
			fields = append(fields, fmt.Sprintf("%s(%s) AS %s",
				agg.Function, agg.Field, agg.Alias))
		} else {
			fields = append(fields, fmt.Sprintf("%s(%s)",
				agg.Function, agg.Field))
		}
	}

	if len(fields) == 0 {
		return "", fmt.Errorf("no fields specified for select")
	}

	return fmt.Sprintf("SELECT %s", strings.Join(fields, ", ")), nil
}
