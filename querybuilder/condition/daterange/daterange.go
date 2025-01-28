package daterange

import (
	"dynamic-sqlbuilder/querybuilder"
	"errors"
	"fmt"
	"time"
)

const jakartaTimezone = "Asia/Jakarta"

func getJakartaLocation() *time.Location {
	loc, _ := time.LoadLocation(jakartaTimezone)
	return loc
}

// DateRangeCondition specifically for date filtering using your existing date d.DateConfig
type DateRangeCondition struct {
	DateConfig DateConfig
}

// Verify interface implementation at compile time
var _ querybuilder.QueryCondition = (*DateRangeCondition)(nil)

func NewDateRangeCondition(config DateConfig) *DateRangeCondition {
	return &DateRangeCondition{
		DateConfig: config,
	}
}
func (d *DateRangeCondition) calculateDateRange() (startDate, endDate time.Time, err error) {
	jakartaLocation := getJakartaLocation()
	now := time.Now().In(jakartaLocation)
	currentYear := now.Year()
	currentMonth := int(now.Month())

	switch d.DateConfig.Type {
	case CurrentMonth:
		startDate = time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation)
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)

	case BackMonth:
		if d.DateConfig.Parameters.Months == nil {
			return time.Time{}, time.Time{}, errors.New("months parameter required for BACK_MONTH")
		}
		startDate = time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation)
		startDate = startDate.AddDate(0, -(*d.DateConfig.Parameters.Months), 0)
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)

	case RelativeRange:
		if d.DateConfig.Parameters.StartBackMonths == nil || d.DateConfig.Parameters.EndBackMonths == nil {
			return time.Time{}, time.Time{}, errors.New("start_back_months and end_back_months required for RELATIVE_RANGE")
		}
		startDate = time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation)
		startDate = startDate.AddDate(0, -(*d.DateConfig.Parameters.StartBackMonths), 0)
		endDate = time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation)
		endDate = endDate.AddDate(0, -(*d.DateConfig.Parameters.EndBackMonths), 0).AddDate(0, 1, 0).Add(-time.Second)

	case YTD:
		startDate = time.Date(currentYear, 1, 1, 0, 0, 0, 0, jakartaLocation)
		endDate = time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation)
		endDate = endDate.AddDate(0, 1, 0).Add(-time.Second)

	case PreviousYear:
		startDate = time.Date(currentYear-1, 1, 1, 0, 0, 0, 0, jakartaLocation)
		endDate = time.Date(currentYear-1, 12, 31, 23, 59, 59, 0, jakartaLocation)

	case SpecificMonth:
		if d.DateConfig.Parameters.Month == nil || d.DateConfig.Parameters.Year == nil {
			return time.Time{}, time.Time{}, errors.New("month and year parameters required for SPECIFIC_MONTH")
		}
		startDate = time.Date(*d.DateConfig.Parameters.Year, time.Month(*d.DateConfig.Parameters.Month), 1, 0, 0, 0, 0, jakartaLocation)
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)

	case SpecificRange:
		if d.DateConfig.Parameters.Start == nil || d.DateConfig.Parameters.End == nil {
			return time.Time{}, time.Time{}, errors.New("start and end parameters required for SPECIFIC_RANGE")
		}
		startDate = time.Date(d.DateConfig.Parameters.Start.Year, time.Month(d.DateConfig.Parameters.Start.Month), 1, 0, 0, 0, 0, jakartaLocation)
		endDate = time.Date(d.DateConfig.Parameters.End.Year, time.Month(d.DateConfig.Parameters.End.Month), 1, 0, 0, 0, 0, jakartaLocation)
		endDate = endDate.AddDate(0, 1, 0).Add(-time.Second)
	case QTD:
		currentQuarter := (time.Month(currentMonth)-1)/3 + 1
		quarterStartMonth := time.Month((currentQuarter-1)*3 + 1)
		startDate = time.Date(currentYear, quarterStartMonth, 1, 0, 0, 0, 0, jakartaLocation)
		endDate = time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation)
		endDate = endDate.AddDate(0, 1, 0).Add(-time.Second)
	case TTM:
		startDate = time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation)
		startDate = startDate.AddDate(0, -11, 0)
		endDate = time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation)
		endDate = endDate.AddDate(0, 1, 0).Add(-time.Second)
	case MoM:
		// Current month
		endDate = time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, jakartaLocation)
		endDate = endDate.AddDate(0, 1, 0).Add(-time.Second)
		// Previous month
		startDate = time.Date(currentYear, time.Month(currentMonth-1), 1, 0, 0, 0, 0, jakartaLocation)
	case PreviousQuarter:
		currentQuarter := (time.Month(currentMonth)-1)/3 + 1
		previousQuarter := currentQuarter - 1
		previousQuarterYear := currentYear

		if previousQuarter == 0 {
			previousQuarter = 4
			previousQuarterYear--
		}

		quarterStartMonth := time.Month((previousQuarter-1)*3 + 1)
		startDate = time.Date(previousQuarterYear, quarterStartMonth, 1, 0, 0, 0, 0, jakartaLocation)
		endDate = startDate.AddDate(0, 3, 0).Add(-time.Second)
	default:
		return time.Time{}, time.Time{}, errors.New("invalid date config type")
	}

	return startDate, endDate, nil
}
func (d *DateRangeCondition) Build(paramOffset int) (string, []interface{}, error) {
	fmt.Printf("DateRangeCondition Build called with paramOffset: %d\n", paramOffset)
	startDate, endDate, err := d.calculateDateRange()
	if err != nil {
		return "", nil, fmt.Errorf("failed to calculate date range: %w", err)
	}

	startYear, endYear := startDate.Year(), endDate.Year()
	startMonth, endMonth := int(startDate.Month()), int(endDate.Month())
	params := make([]interface{}, 0)

	if startYear == endYear {
		params = append(params, startYear, startMonth, endMonth)
		return fmt.Sprintf(
			"period_year = $%d AND period_month BETWEEN $%d AND $%d",
			paramOffset,
			paramOffset+1,
			paramOffset+2,
		), params, nil
	}
	params = append(params, startYear, startMonth, endYear, endMonth)
	return fmt.Sprintf(
		"(period_year = $%d AND period_month >= $%d) OR (period_year = $%d AND period_month <= $%d)",
		paramOffset,
		paramOffset+1,
		paramOffset+2,
		paramOffset+3,
	), params, nil
}
