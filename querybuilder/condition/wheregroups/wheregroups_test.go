package wheregroups

import (
	"dynamic-sqlbuilder/querybuilder"
	"dynamic-sqlbuilder/querybuilder/condition/fields"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhereGroups(t *testing.T) {
	tests := []struct {
		name          string
		buildGroups   func() *WhereGroups
		paramOffset   int
		expectedSQL   string
		expectedArgs  []interface{}
		expectedError error
	}{
		{
			name: "Empty Where Groups",
			buildGroups: func() *WhereGroups {
				return NewWhereGroups()
			},
			paramOffset:   1,
			expectedSQL:   "",
			expectedArgs:  nil,
			expectedError: nil,
		},
		{
			name: "Single Group with Single Condition",
			buildGroups: func() *WhereGroups {
				wg := NewWhereGroups()
				group := querybuilder.NewWhereGroup(querybuilder.AND)
				condition := fields.NewFieldCondition("column1", fields.Equals, "value1")
				group.Add(condition)
				wg.Add(*group)
				return wg
			},
			paramOffset:   1,
			expectedSQL:   "WHERE column1 = $1",
			expectedArgs:  []interface{}{"value1"},
			expectedError: nil,
		},
		{
			name: "Multiple Groups with AND",
			buildGroups: func() *WhereGroups {
				wg := NewWhereGroups()

				// First group
				group1 := querybuilder.NewWhereGroup(querybuilder.AND)
				group1.Add(fields.NewFieldCondition("column1", fields.Equals, "value1"))
				wg.Add(*group1)

				// Second group
				group2 := querybuilder.NewWhereGroup(querybuilder.AND)
				group2.Add(fields.NewFieldCondition("column2", fields.GreaterThan, 10))
				wg.Add(*group2)

				return wg
			},
			paramOffset:   1,
			expectedSQL:   "WHERE column1 = $1 AND column2 > $2",
			expectedArgs:  []interface{}{"value1", 10},
			expectedError: nil,
		},
		{
			name: "Multiple Groups with OR",
			buildGroups: func() *WhereGroups {
				wg := NewWhereGroups()

				group1 := querybuilder.NewWhereGroup(querybuilder.OR)
				group1.Add(fields.NewFieldCondition("column1", fields.Equals, "value1"))
				wg.Add(*group1)

				group2 := querybuilder.NewWhereGroup(querybuilder.OR)
				group2.Add(fields.NewFieldCondition("column2", fields.LessThan, 20))
				wg.Add(*group2)

				return wg
			},
			paramOffset:   1,
			expectedSQL:   "WHERE column1 = $1 OR column2 < $2",
			expectedArgs:  []interface{}{"value1", 20},
			expectedError: nil,
		},
		{
			name: "Mixed AND/OR Groups",
			buildGroups: func() *WhereGroups {
				wg := NewWhereGroups()

				group1 := querybuilder.NewWhereGroup(querybuilder.AND)
				group1.Add(fields.NewFieldCondition("column1", fields.Equals, "value1"))
				wg.Add(*group1)

				group2 := querybuilder.NewWhereGroup(querybuilder.OR)
				group2.Add(fields.NewFieldCondition("column2", fields.GreaterThan, 10))
				wg.Add(*group2)

				group3 := querybuilder.NewWhereGroup(querybuilder.AND)
				group3.Add(fields.NewFieldCondition("column3", fields.LessThan, 20))
				wg.Add(*group3)

				return wg
			},
			paramOffset:   1,
			expectedSQL:   "WHERE column1 = $1 OR column2 > $2 AND column3 < $3",
			expectedArgs:  []interface{}{"value1", 10, 20},
			expectedError: nil,
		},
		{
			name: "Group with Multiple Conditions",
			buildGroups: func() *WhereGroups {
				wg := NewWhereGroups()
				group := querybuilder.NewWhereGroup(querybuilder.AND)

				group.Add(fields.NewFieldCondition("column1", fields.Equals, "value1"))
				group.Add(fields.NewFieldCondition("column2", fields.GreaterThan, 10))
				group.Add(fields.NewFieldCondition("column3", fields.Like, "%pattern%"))

				wg.Add(*group)
				return wg
			},
			paramOffset:   1,
			expectedSQL:   "WHERE (column1 = $1 AND column2 > $2 AND column3 LIKE $3)",
			expectedArgs:  []interface{}{"value1", 10, "%pattern%"},
			expectedError: nil,
		},
		{
			name: "Custom Parameter Offset",
			buildGroups: func() *WhereGroups {
				wg := NewWhereGroups()
				group := querybuilder.NewWhereGroup(querybuilder.AND)
				group.Add(fields.NewFieldCondition("column1", fields.Equals, "value1"))
				wg.Add(*group)
				return wg
			},
			paramOffset:   5,
			expectedSQL:   "WHERE column1 = $5",
			expectedArgs:  []interface{}{"value1"},
			expectedError: nil,
		},
		{
			name: "Complex Nested Conditions",
			buildGroups: func() *WhereGroups {
				wg := NewWhereGroups()

				// First group (id = 1 AND status = 'active')
				group1 := querybuilder.NewWhereGroup(querybuilder.AND)
				group1.Add(fields.NewFieldCondition("id", fields.Equals, 1))
				group1.Add(fields.NewFieldCondition("status", fields.Equals, "active"))
				wg.Add(*group1)

				// Second group (price > 100 OR category = 'premium')
				group2 := querybuilder.NewWhereGroup(querybuilder.OR)
				group2.Add(fields.NewFieldCondition("price", fields.GreaterThan, 100))
				group2.Add(fields.NewFieldCondition("category", fields.Equals, "premium"))
				wg.Add(*group2)

				return wg
			},
			paramOffset:   1,
			expectedSQL:   "WHERE (id = $1 AND status = $2) OR (price > $3 OR category = $4)",
			expectedArgs:  []interface{}{1, "active", 100, "premium"},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := tt.buildGroups()
			sql, args, err := wg.Build(tt.paramOffset)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedSQL, sql)
				assert.Equal(t, tt.expectedArgs, args)
			}
		})
	}
}
