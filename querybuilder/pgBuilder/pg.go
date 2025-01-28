package pgbuilder

import (
	"dynamic-sqlbuilder/querybuilder"
	"dynamic-sqlbuilder/querybuilder/condition/wheregroups"
	"dynamic-sqlbuilder/querybuilder/from/simplefrom"
	aggregate "dynamic-sqlbuilder/querybuilder/select/aggregateselect"
	"dynamic-sqlbuilder/querybuilder/select/simpleselect"
	"fmt"
	"strings"
)

// PostgresQueryBuilder implementation with aggregate support
type PostgresQueryBuilder struct {
	query           *querybuilder.Query
	aggregateSelect *aggregate.AggregateSelect
	whereGroups     *wheregroups.WhereGroups // Keep track of where groups
	currentGroup    *querybuilder.WhereGroup
}

func NewPostgresQueryBuilder() *PostgresQueryBuilder {
	initialGroup := querybuilder.NewWhereGroup(querybuilder.AND)
	whereGroups := wheregroups.NewWhereGroups()

	return &PostgresQueryBuilder{
		query: &querybuilder.Query{
			SelectClause: simpleselect.NewSimpleSelect("*"),
			WhereClause:  whereGroups,
			Args:         make([]interface{}, 0),
		},
		whereGroups:  whereGroups,
		currentGroup: initialGroup,
	}
}

func (b *PostgresQueryBuilder) Where(condition querybuilder.QueryCondition) querybuilder.QueryBuilder {
	// fmt.Printf("Adding condition to current group: %+v\n", condition)
	// If this is the first condition overall
	if len(b.whereGroups.Groups) == 0 {
		b.currentGroup.Add(condition)
		b.whereGroups.Add(*b.currentGroup)
	} else {
		// Add to the current group (which was either created by And() or is existing)
		b.currentGroup.Add(condition)

		// If this group isn't in whereGroups yet, add it
		if b.currentGroup.IsNew {
			b.whereGroups.Add(*b.currentGroup)
			b.currentGroup.IsNew = false
		}
	}
	return b
}
func (b *PostgresQueryBuilder) WhereGroup(operator querybuilder.LogicalOperator, buildGroup func(*querybuilder.WhereGroup)) querybuilder.QueryBuilder {
	group := querybuilder.NewWhereGroup(operator)
	buildGroup(group)
	b.whereGroups.Add(*group)
	b.currentGroup = group
	return b
}

func (b *PostgresQueryBuilder) And() querybuilder.QueryBuilder {
	// Create a new group but mark it as new
	b.currentGroup = querybuilder.NewWhereGroup(querybuilder.AND)
	b.currentGroup.IsNew = true
	b.whereGroups.Add(*b.currentGroup)
	return b
}
func (b *PostgresQueryBuilder) Or() querybuilder.QueryBuilder {
	b.currentGroup = querybuilder.NewWhereGroup(querybuilder.OR)
	b.whereGroups.Add(*b.currentGroup)
	return b
}

func (b *PostgresQueryBuilder) From(table string) querybuilder.QueryBuilder {
	b.query.FromClause = simplefrom.NewSimpleFrom(table)
	return b
}
func (b *PostgresQueryBuilder) Select(fields ...string) querybuilder.QueryBuilder {
	b.query.SelectClause = simpleselect.NewSimpleSelect(fields...)
	return b
}

// Now returns QueryBuilder to satisfy interface
func (b *PostgresQueryBuilder) SelectAggregate() querybuilder.QueryBuilder {
	aggSelect := aggregate.NewAggregateSelect()
	b.query.SelectClause = aggSelect
	b.aggregateSelect = aggSelect // Store reference
	return b
}

// Add methods to access aggregate functions
func (b *PostgresQueryBuilder) AddRegularField(field string) querybuilder.QueryBuilder {
	if b.aggregateSelect != nil {
		b.aggregateSelect.AddRegularField(field)
	}
	return b
}

func (b *PostgresQueryBuilder) AddAggregate(fn aggregate.AggregateFunction, field, alias string) querybuilder.QueryBuilder {
	if b.aggregateSelect != nil {
		b.aggregateSelect.AddAggregate(fn, field, alias)
	}
	return b
}
func (b *PostgresQueryBuilder) Build() (string, []interface{}, error) {
	var queryParts []string
	var args []interface{}

	// fmt.Printf("Building query with %d where groups\n", len(b.whereGroups.Groups))

	// Build SELECT clause
	selectSQL, err := b.query.SelectClause.Build()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build SELECT clause: %w", err)
	}
	queryParts = append(queryParts, selectSQL)

	// Build FROM clause
	if b.query.FromClause != nil {
		queryParts = append(queryParts, b.query.FromClause.Build())
	}
	// Build WHERE clause
	if b.query.WhereClause != nil {
		// fmt.Printf("Building WHERE clause...\n")
		whereSQL, whereArgs, err := b.query.WhereClause.Build(len(args) + 1)
		if err != nil {
			return "", nil, fmt.Errorf("failed to build WHERE clause: %w", err)
		}
		if whereSQL != "" {
			// fmt.Printf("WHERE clause generated: %s\n", whereSQL)
			queryParts = append(queryParts, whereSQL)
			args = append(args, whereArgs...)
		} else {
			fmt.Printf("No WHERE clause generated\n")
		}
	}
	return strings.Join(queryParts, " "), args, nil
}
