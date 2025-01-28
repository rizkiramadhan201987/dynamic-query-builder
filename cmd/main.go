package main

import (
	"dynamic-sqlbuilder/querybuilder"
	"dynamic-sqlbuilder/querybuilder/condition/daterange"
	"dynamic-sqlbuilder/querybuilder/condition/fields"
	pgbuilder "dynamic-sqlbuilder/querybuilder/pgBuilder"
	aggregate "dynamic-sqlbuilder/querybuilder/select/aggregateselect"
	"fmt"
)

func buildBalanceSheetQuery() {
	builder := pgbuilder.NewPostgresQueryBuilder()

	// Create date range condition for YTD
	ytdConfig := daterange.DateConfig{
		Type: daterange.YTD,
	}
	dateCondition := daterange.NewDateRangeCondition(ytdConfig)

	// Example of active accounts for balance sheet
	activeAccounts := []interface{}{"1001", "1002", "1003", "2001", "2002"}

	query, args, err := builder.SelectAggregate().
		AddRegularField("coa.coadescription").
		AddAggregate(aggregate.Sum, "amount", "total_amount").
		From("transactions").
		Where(dateCondition).
		And().
		Where(fields.NewFieldCondition("account_code", fields.In, activeAccounts)).
		And().
		Where(fields.NewFieldCondition("is_active", fields.Equals, true)).
		Build()

	if err != nil {
		fmt.Printf("Error building balance sheet query: %v\n", err)
		return
	}

	fmt.Printf("Balance Sheet Query: %s\n", query)
	fmt.Printf("Arguments: %v\n", args)
}

func buildProfitLossQuery() {
	builder := pgbuilder.NewPostgresQueryBuilder()

	// Create date range for last 3 months
	monthsBack := 3
	backMonthConfig := daterange.DateConfig{
		Type: daterange.BackMonth,
		Parameters: daterange.DateParameters{
			Months: &monthsBack,
		},
	}
	dateCondition := daterange.NewDateRangeCondition(backMonthConfig)

	// Revenue and expense account types
	accountTypes := []interface{}{"REVENUE", "EXPENSE"}

	query, args, err := builder.SelectAggregate().
		AddRegularField("account_type").
		AddRegularField("department").
		AddAggregate(aggregate.Sum, "amount", "total_amount").
		From("financial_transactions").
		Where(dateCondition).
		And().
		Where(fields.NewFieldCondition("account_type", fields.In, accountTypes)).
		Build()

	if err != nil {
		fmt.Printf("Error building profit & loss query: %v\n", err)
		return
	}

	fmt.Printf("Profit & Loss Query: %s\n", query)
	fmt.Printf("Arguments: %v\n", args)
}

func buildCashFlowQuery() {
	builder := pgbuilder.NewPostgresQueryBuilder()

	// Create specific date range
	startMonth := &daterange.MonthYear{Month: 1, Year: 2024}
	endMonth := &daterange.MonthYear{Month: 12, Year: 2024}

	rangeConfig := daterange.DateConfig{
		Type: daterange.SpecificRange,
		Parameters: daterange.DateParameters{
			Start: startMonth,
			End:   endMonth,
		},
	}
	dateCondition := daterange.NewDateRangeCondition(rangeConfig)

	query, args, err := builder.SelectAggregate().
		AddRegularField("transaction_type").
		AddAggregate(aggregate.Sum, "inflow", "total_inflow").
		AddAggregate(aggregate.Sum, "outflow", "total_outflow").
		From("cash_transactions").
		Where(dateCondition).
		WhereGroup(querybuilder.OR, func(group *querybuilder.WhereGroup) {
			group.Add(fields.NewFieldCondition("transaction_type", fields.Equals, "OPERATING"))
			group.Add(fields.NewFieldCondition("transaction_type", fields.Equals, "INVESTING"))
			group.Add(fields.NewFieldCondition("transaction_type", fields.Equals, "FINANCING"))
		}).
		Build()

	if err != nil {
		fmt.Printf("Error building cash flow query: %v\n", err)
		return
	}

	fmt.Printf("Cash Flow Query: %s\n", query)
	fmt.Printf("Arguments: %v\n", args)
}

func buildCustomReportQuery() {
	builder := pgbuilder.NewPostgresQueryBuilder()

	// Create TTM (Trailing Twelve Months) date range
	ttmConfig := daterange.DateConfig{
		Type: daterange.TTM,
	}
	dateCondition := daterange.NewDateRangeCondition(ttmConfig)

	// Complex query with multiple conditions and groupings
	query, args, err := builder.SelectAggregate().
		AddRegularField("department").
		AddRegularField("cost_center").
		AddAggregate(aggregate.Sum, "revenue", "total_revenue").
		AddAggregate(aggregate.Sum, "expense", "total_expense").
		AddAggregate(aggregate.Avg, "profit_margin", "avg_margin").
		From("financial_metrics").
		Where(dateCondition).
		WhereGroup(querybuilder.AND, func(group *querybuilder.WhereGroup) {
			group.Add(fields.NewFieldCondition("is_active", fields.Equals, true))
			group.Add(fields.NewFieldCondition("profit_margin", fields.GreaterThan, 0.15))
		}).
		Build()

	if err != nil {
		fmt.Printf("Error building custom report query: %v\n", err)
		return
	}

	fmt.Printf("Custom Report Query: %s\n", query)
	fmt.Printf("Arguments: %v\n", args)
}

func main() {
	fmt.Println("=== Building Various Financial Queries ===")

	fmt.Println("1. Balance Sheet Query:")
	buildBalanceSheetQuery()
	fmt.Println()

	fmt.Println("2. Profit & Loss Query:")
	buildProfitLossQuery()
	fmt.Println()

	fmt.Println("3. Cash Flow Query:")
	buildCashFlowQuery()
	fmt.Println()

	fmt.Println("4. Custom Report Query:")
	buildCustomReportQuery()
}
