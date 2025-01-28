package daterange

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DateRange struct {
	startDate time.Time
	endDate   time.Time
}

func (dr DateRange) String() string {
	return fmt.Sprintf("from %s to %s",
		dr.startDate.Format("2 January 2006"),
		dr.endDate.Format("2 January 2006"))
}

func TestCalculateDateRange(t *testing.T) {
	jakartaLocation := getJakartaLocation()
	now := time.Now().In(jakartaLocation)
	currentYear := now.Year()
	currentMonth := int(now.Month())

	// Helper function to create integer pointers
	intPtr := func(i int) *int {
		p := i
		return &p
	}

	// Success test cases
	successTests := []struct {
		name        string
		config      DateConfig
		expected    DateRange
		description string
	}{
		{
			name: "Current Month",
			config: DateConfig{
				Type: CurrentMonth,
			},
			expected: DateRange{
				startDate: time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation),
				endDate:   time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation).AddDate(0, 1, 0).Add(-time.Second),
			},
			description: fmt.Sprintf("Should return the complete current month range (1-%s-%d to end of month)",
				time.Now().Month().String(), currentYear),
		},
		{
			name: "Back Month - One Month",
			config: DateConfig{
				Type: BackMonth,
				Parameters: DateParameters{
					Months: intPtr(1),
				},
			},
			expected: DateRange{
				startDate: time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation).AddDate(0, -1, 0),
				endDate:   time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation).Add(-time.Second),
			},
			description: fmt.Sprintf("Should return the complete previous month range when looking back 1 month from %s %d",
				time.Now().Month().String(), currentYear),
		},
		{
			name: "Year to Date (YTD)",
			config: DateConfig{
				Type: YTD,
			},
			expected: DateRange{
				startDate: time.Date(currentYear, 1, 1, 0, 0, 0, 0, jakartaLocation),
				endDate:   time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation).AddDate(0, 1, 0).Add(-time.Second),
			},
			description: fmt.Sprintf("Should return from January 1st to end of current month (%d)", currentYear),
		},
		{
			name: "Previous Year",
			config: DateConfig{
				Type: PreviousYear,
			},
			expected: DateRange{
				startDate: time.Date(currentYear-1, 1, 1, 0, 0, 0, 0, jakartaLocation),
				endDate:   time.Date(currentYear-1, 12, 31, 23, 59, 59, 0, jakartaLocation),
			},
			description: fmt.Sprintf("Should return the complete previous year (%d)", currentYear-1),
		},
		{
			name: "Quarter to Date (QTD)",
			config: DateConfig{
				Type: QTD,
			},
			expected: DateRange{
				startDate: time.Date(currentYear, time.Month(((currentMonth-1)/3)*3+1), 1, 0, 0, 0, 0, jakartaLocation),
				endDate:   time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation).AddDate(0, 1, 0).Add(-time.Second),
			},
			description: fmt.Sprintf("Should return from start of current quarter to end of %s %d",
				time.Now().Month().String(), currentYear),
		},
		{
			name: "Trailing Twelve Months (TTM)",
			config: DateConfig{
				Type: TTM,
			},
			expected: DateRange{
				startDate: time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation).AddDate(0, -11, 0),
				endDate:   time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation).AddDate(0, 1, 0).Add(-time.Second),
			},
			description: fmt.Sprintf("Should return last 12 months ending with current month (%s %d)",
				time.Now().Month().String(), currentYear),
		},
		{
			name: "Month over Month (MoM)",
			config: DateConfig{
				Type: MoM,
			},
			expected: DateRange{
				startDate: time.Date(currentYear, time.Month(currentMonth-1), 1, 0, 0, 0, 0, jakartaLocation),
				endDate:   time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation).AddDate(0, 1, 0).Add(-time.Second),
			},
			description: fmt.Sprintf("Should return previous month when comparing month-over-month in %s %d",
				time.Now().Month().String(), currentYear),
		},
		{
			name: "Previous Quarter",
			config: DateConfig{
				Type: PreviousQuarter,
			},
			expected: DateRange{
				startDate: func() time.Time {
					currentQuarter := (time.Month(currentMonth)-1)/3 + 1
					previousQuarter := currentQuarter - 1
					previousQuarterYear := currentYear
					if previousQuarter == 0 {
						previousQuarter = 4
						previousQuarterYear--
					}
					quarterStartMonth := time.Month((previousQuarter-1)*3 + 1)
					return time.Date(previousQuarterYear, quarterStartMonth, 1, 0, 0, 0, 0, jakartaLocation)
				}(),
				endDate: func() time.Time {
					currentQuarter := (time.Month(currentMonth)-1)/3 + 1
					previousQuarter := currentQuarter - 1
					previousQuarterYear := currentYear
					if previousQuarter == 0 {
						previousQuarter = 4
						previousQuarterYear--
					}
					quarterStartMonth := time.Month((previousQuarter-1)*3 + 1)
					return time.Date(previousQuarterYear, quarterStartMonth, 1, 0, 0, 0, 0, jakartaLocation).AddDate(0, 3, 0).Add(-time.Second)
				}(),
			},
			description: fmt.Sprintf("Should return complete previous quarter when in Q%d %d",
				(currentMonth-1)/3+1, currentYear),
		},
		{
			name: "Specific Month",
			config: DateConfig{
				Type: SpecificMonth,
				Parameters: DateParameters{
					Month: intPtr(3),
					Year:  intPtr(2024),
				},
			},
			expected: DateRange{
				startDate: time.Date(2024, 3, 1, 0, 0, 0, 0, jakartaLocation),
				endDate:   time.Date(2024, 3, 1, 0, 0, 0, 0, jakartaLocation).AddDate(0, 1, 0).Add(-time.Second),
			},
			description: "Should return complete month range for March 2024",
		},
		{
			name: "Specific Range",
			config: DateConfig{
				Type: SpecificRange,
				Parameters: DateParameters{
					Start: &MonthYear{Month: 1, Year: 2024},
					End:   &MonthYear{Month: 3, Year: 2024},
				},
			},
			expected: DateRange{
				startDate: time.Date(2024, 1, 1, 0, 0, 0, 0, jakartaLocation),
				endDate:   time.Date(2024, 3, 1, 0, 0, 0, 0, jakartaLocation).AddDate(0, 1, 0).Add(-time.Second),
			},
			description: "Should return complete range from January 2024 to March 2024",
		},
		{
			name: "Relative Range",
			config: DateConfig{
				Type: RelativeRange,
				Parameters: DateParameters{
					StartBackMonths: intPtr(6),
					EndBackMonths:   intPtr(3),
				},
			},
			expected: DateRange{
				startDate: time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation).AddDate(0, -6, 0),
				endDate:   time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation).AddDate(0, -3, 0).AddDate(0, 1, 0).Add(-time.Second),
			},
			description: "Should return range from 6 months ago to 3 months ago",
		},
	}

	// Failure test cases
	failureTests := []struct {
		name        string
		config      DateConfig
		description string
		expectedErr string
	}{
		{
			name: "Back Month - Missing Parameter",
			config: DateConfig{
				Type: BackMonth,
			},
			description: "Should fail when months parameter is not provided",
			expectedErr: "months parameter required for BACK_MONTH",
		},
		{
			name: "Specific Month - Missing Parameters",
			config: DateConfig{
				Type: SpecificMonth,
			},
			description: "Should fail when month and year parameters are not provided",
			expectedErr: "month and year parameters required for SPECIFIC_MONTH",
		},
		{
			name: "Specific Range - Missing Parameters",
			config: DateConfig{
				Type: SpecificRange,
			},
			description: "Should fail when start and end parameters are not provided",
			expectedErr: "start and end parameters required for SPECIFIC_RANGE",
		},
		{
			name: "Relative Range - Missing Parameters",
			config: DateConfig{
				Type: RelativeRange,
			},
			description: "Should fail when start_back_months and end_back_months are not provided",
			expectedErr: "start_back_months and end_back_months required for RELATIVE_RANGE",
		},
		{
			name: "Invalid Config Type",
			config: DateConfig{
				Type: "INVALID_TYPE",
			},
			description: "Should fail when an invalid date configuration type is provided",
			expectedErr: "invalid date config type",
		},
	}

	// Run success tests
	t.Run("Success Cases", func(t *testing.T) {
		for _, tt := range successTests {
			t.Run(tt.name, func(t *testing.T) {
				t.Log(tt.description)
				dateRange := NewDateRangeCondition(tt.config)
				startDate, endDate, err := dateRange.calculateDateRange()
				require.NoError(t, err)

				actualRange := DateRange{startDate, endDate}
				expectedRange := tt.expected

				assert.Equal(t, expectedRange.String(), actualRange.String(),
					"\nExpected date range %s\nBut got date range %s",
					expectedRange.String(),
					actualRange.String())

				// Validate exact time matches
				assert.True(t, expectedRange.startDate.Equal(startDate),
					"Start date should be exactly %s but got %s",
					expectedRange.startDate.Format("2006-01-02 15:04:05"),
					startDate.Format("2006-01-02 15:04:05"))
				assert.True(t, expectedRange.endDate.Equal(endDate),
					"End date should be exactly %s but got %s",
					expectedRange.endDate.Format("2006-01-02 15:04:05"),
					endDate.Format("2006-01-02 15:04:05"))
			})
		}
	})

	// Run failure tests
	t.Run("Failure Cases", func(t *testing.T) {
		for _, tt := range failureTests {
			t.Run(tt.name, func(t *testing.T) {
				t.Log(tt.description)
				dateRange := NewDateRangeCondition(tt.config)
				_, _, err := dateRange.calculateDateRange()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr,
					"Expected error containing '%s' but got '%s'",
					tt.expectedErr, err.Error())
			})
		}
	})
}
