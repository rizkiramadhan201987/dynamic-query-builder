package aggregate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregateSelect_Build(t *testing.T) {
	tests := []struct {
		name          string
		regularFields []string
		aggregates    []struct {
			function AggregateFunction
			field    string
			alias    string
		}
		expectedSQL   string
		expectedError bool
		description   string
	}{
		{
			name: "Basic Aggregate",
			aggregates: []struct {
				function AggregateFunction
				field    string
				alias    string
			}{
				{Sum, "amount", "total_amount"},
			},
			expectedSQL: "SELECT SUM(amount) AS total_amount",
			description: "Should generate SQL with a single aggregate function",
		},
		{
			name: "Multiple Aggregates",
			aggregates: []struct {
				function AggregateFunction
				field    string
				alias    string
			}{
				{Sum, "amount", "total_amount"},
				{Avg, "price", "avg_price"},
				{Count, "id", "count"},
			},
			expectedSQL: "SELECT SUM(amount) AS total_amount, AVG(price) AS avg_price, COUNT(id) AS count",
			description: "Should generate SQL with multiple aggregate functions",
		},
		{
			name:          "Regular Fields Only",
			regularFields: []string{"id", "name", "date"},
			expectedSQL:   "SELECT id, name, date",
			description:   "Should generate SQL with only regular fields",
		},
		{
			name:          "Mixed Regular and Aggregate Fields",
			regularFields: []string{"date", "category"},
			aggregates: []struct {
				function AggregateFunction
				field    string
				alias    string
			}{
				{Sum, "amount", "total"},
				{Count, "*", "count"},
			},
			expectedSQL: "SELECT date, category, SUM(amount) AS total, COUNT(*) AS count",
			description: "Should generate SQL with both regular and aggregate fields",
		},
		{
			name: "Aggregate Without Alias",
			aggregates: []struct {
				function AggregateFunction
				field    string
				alias    string
			}{
				{Min, "price", ""},
				{Max, "price", ""},
			},
			expectedSQL: "SELECT MIN(price), MAX(price)",
			description: "Should generate SQL with aggregate functions without aliases",
		},
		{
			name:          "No Fields",
			expectedError: true,
			description:   "Should return error when no fields are specified",
		},
		{
			name: "All Aggregate Functions",
			aggregates: []struct {
				function AggregateFunction
				field    string
				alias    string
			}{
				{Sum, "amount", "sum"},
				{Avg, "amount", "avg"},
				{Count, "id", "count"},
				{Max, "amount", "max"},
				{Min, "amount", "min"},
			},
			expectedSQL: "SELECT SUM(amount) AS sum, AVG(amount) AS avg, COUNT(id) AS count, MAX(amount) AS max, MIN(amount) AS min",
			description: "Should handle all supported aggregate functions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.description)

			// Create new AggregateSelect
			as := NewAggregateSelect()

			// Add regular fields
			for _, field := range tt.regularFields {
				as.AddRegularField(field)
			}

			// Add aggregate fields
			for _, agg := range tt.aggregates {
				as.AddAggregate(agg.function, agg.field, agg.alias)
			}

			// Build the SQL
			sql, err := as.Build()

			if tt.expectedError {
				assert.Error(t, err, "Expected an error but got none")
				return
			}

			require.NoError(t, err, "Unexpected error building SQL")
			assert.Equal(t, tt.expectedSQL, sql, "SQL does not match expected output")
		})
	}
}

func TestAggregateSelect_Chaining(t *testing.T) {
	t.Run("Method Chaining", func(t *testing.T) {
		as := NewAggregateSelect().
			AddRegularField("date").
			AddAggregate(Sum, "amount", "total").
			AddAggregate(Count, "*", "count").
			AddRegularField("category")

		sql, err := as.Build()
		require.NoError(t, err)
		assert.Equal(t, "SELECT date, category, SUM(amount) AS total, COUNT(*) AS count", sql)
	})
}

func TestNewAggregateSelect(t *testing.T) {
	t.Run("New Instance", func(t *testing.T) {
		as := NewAggregateSelect()
		assert.NotNil(t, as)
		assert.Empty(t, as.regularFields)
		assert.Empty(t, as.aggregates)
	})
}
