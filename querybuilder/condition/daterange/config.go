package daterange

// DateConfigType defines the various types of date ranges used in financial reporting and analysis
type DateConfigType string

const (
	// CurrentMonth represents the date range from the first to the last day of the current month.
	// For example, if today is April 15, 2024, the range would be April 1, 2024 to April 30, 2024.
	CurrentMonth DateConfigType = "CURRENT_MONTH"

	// BackMonth provides a date range for a specific month in the past, relative to the current month.
	// Requires a 'Months' parameter specifying how many months to look back.
	// Example: With Months=3 in April 2024, returns January 1-31, 2024.
	BackMonth DateConfigType = "BACK_MONTH"

	// RelativeRange allows flexible date ranges based on months back from the current month.
	// Requires 'StartBackMonths' and 'EndBackMonths' parameters to define the range.
	// Example: StartBackMonths=6, EndBackMonths=3 in April 2024 returns October 2023 to January 2024.
	RelativeRange DateConfigType = "RELATIVE_RANGE"

	// YTD (Year-to-Date) represents the period from January 1st of the current year to the current date.
	// Commonly used to track annual performance and compare against previous years' YTD figures.
	// Example: In April 2024, returns January 1, 2024 to April 30, 2024.
	YTD DateConfigType = "YTD"

	// PreviousYear represents the complete previous calendar year.
	// Used for year-over-year comparisons and annual reporting.
	// Example: In 2024, returns January 1, 2023 to December 31, 2023.
	PreviousYear DateConfigType = "PREVIOUS_YEAR"

	// SpecificMonth allows querying data for a particular month and year.
	// Requires 'Month' and 'Year' parameters to specify the exact month.
	// Example: Month=3, Year=2024 returns March 1-31, 2024.
	SpecificMonth DateConfigType = "SPECIFIC_MONTH"

	// SpecificRange enables custom date ranges with specific start and end months.
	// Requires 'Start' and 'End' parameters containing Month and Year values.
	// Example: Start={Month:1, Year:2024}, End={Month:3, Year:2024} returns January 1 to March 31, 2024.
	SpecificRange DateConfigType = "SPECIFIC_RANGE"

	// QTD (Quarter-to-Date) represents the period from the start of the current quarter to the current date.
	// Used for quarterly performance tracking and interim reporting.
	// Example: In May 2024 (Q2), returns April 1, 2024 to May 31, 2024.
	QTD DateConfigType = "QTD"

	// TTM (Trailing Twelve Months) provides a rolling 12-month window ending with the current month.
	// Also known as Last Twelve Months (LTM), useful for analyzing annual performance on a rolling basis.
	// Example: In April 2024, returns May 1, 2023 to April 30, 2024.
	TTM DateConfigType = "TTM"

	// MoM (Month-over-Month) represents the previous month's complete date range.
	// Used for comparing sequential monthly performance and identifying short-term trends.
	// Example: In April 2024, returns March 1-31, 2024.
	MoM DateConfigType = "MOM"

	// PreviousQuarter represents the complete previous quarter's date range.
	// Used for quarter-over-quarter comparisons and quarterly trend analysis.
	// Example: In Q2 2024, returns January 1, 2024 to March 31, 2024 (Q1).
	PreviousQuarter DateConfigType = "PREVIOUS_QUARTER"
)

// DateParameters defines the configuration options for various date range calculations
// used in financial reporting and analysis. Each field is optional and its usage
// depends on the specific DateConfigType being employed.
type DateParameters struct {
	// Months specifies the number of months to look back from the current month.
	// Used with DateConfigType.BackMonth to retrieve data for a specific past month.
	// Example: A value of 3 in April 2024 would target January 2024.
	Months *int `json:"months,omitempty"`

	// StartBackMonths defines the starting point for a relative date range,
	// expressed as the number of months back from the current month.
	// Used with DateConfigType.RelativeRange in conjunction with EndBackMonths.
	// Example: A value of 6 in April 2024 would start the range from October 2023.
	StartBackMonths *int `json:"start_back_months,omitempty"`

	// EndBackMonths defines the ending point for a relative date range,
	// expressed as the number of months back from the current month.
	// Used with DateConfigType.RelativeRange in conjunction with StartBackMonths.
	// Example: A value of 3 in April 2024 would end the range at January 2024.
	EndBackMonths *int `json:"end_back_months,omitempty"`

	// Month specifies a particular month (1-12) for date range calculations.
	// Used with DateConfigType.SpecificMonth along with Year to target a specific month.
	// Example: A value of 3 represents March.
	Month *int `json:"month,omitempty"`

	// Year specifies the target year for date range calculations.
	// Used with DateConfigType.SpecificMonth along with Month to target a specific month.
	// Example: A value of 2024 represents the year 2024.
	Year *int `json:"year,omitempty"`

	// Start defines the beginning of a custom date range using a MonthYear structure.
	// Used with DateConfigType.SpecificRange to set a precise starting point.
	// The MonthYear structure contains both month and year values.
	Start *MonthYear `json:"start,omitempty"`

	// End defines the conclusion of a custom date range using a MonthYear structure.
	// Used with DateConfigType.SpecificRange to set a precise ending point.
	// The MonthYear structure contains both month and year values.
	End *MonthYear `json:"end,omitempty"`
}

// All fields in DateParameters are pointers to allow for optional values in JSON
// representation and to distinguish between zero values and unset fields. This is
// particularly important when parsing JSON requests where some fields may be
// intentionally omitted based on the selected DateConfigType.

type MonthYear struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type DateConfig struct {
	Type       DateConfigType `json:"type"`
	Parameters DateParameters `json:"parameters"`
}
