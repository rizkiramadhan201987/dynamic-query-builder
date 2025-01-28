package querybuilder

import (
	aggregate "dynamic-sqlbuilder/querybuilder/select/aggregateselect"
)

type Query struct {
	SelectClause SelectClause
	FromClause   FromClause
	WhereClause  WhereClause
	Args         []interface{}
}

// QueryBuilder interface with enhanced aggregate support
type QueryBuilder interface {
	// Select operations
	Select(fields ...string) QueryBuilder
	SelectAggregate() QueryBuilder
	AddRegularField(field string) QueryBuilder
	AddAggregate(fn aggregate.AggregateFunction, field, alias string) QueryBuilder

	// FROM operation
	From(table string) QueryBuilder

	// Where operations
	Where(condition QueryCondition) QueryBuilder
	WhereGroup(operator LogicalOperator,
		buildGroup func(*WhereGroup)) QueryBuilder
	Or() QueryBuilder  // Starts a new OR group
	And() QueryBuilder // Starts a new AND group

	Build() (string, []interface{}, error)
}

// QueryCondition represents a condition in the query
type QueryCondition interface {
	Build(paramOffset int) (string, []interface{}, error)
}

// SelectClause interface defines how select clauses should be built
type SelectClause interface {
	Build() (string, error)
}

// FromClause defines the interface for building FROM part of query
type FromClause interface {
	Build() string
}

// WhereClause interface
type WhereClause interface {
	Build(paramOffset int) (string, []interface{}, error)
}
