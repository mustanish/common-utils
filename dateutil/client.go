package dateutil

import (
	"fmt"
	"strconv"
	"time"
)

// Common date formats
const (
	RFC3339Date     = "2006-01-02"
	RFC3339DateTime = "2006-01-02T15:04:05Z07:00"
	ISO8601DateTime = "2006-01-02T15:04:05Z"
	USDate          = "01/02/2006"
	EuropeanDate    = "02/01/2006"
	SimpleDateTime  = "2006-01-02 15:04:05"
)

// DateClient defines the interface for most commonly used date utility operations
type DateClient interface {
	// Parsing methods
	Parse(dateStr string, formats ...string) (time.Time, error)
	ParseUnix(timestamp any) (time.Time, error)

	// Formatting methods
	Format(date time.Time, format string) string
	FormatToRFC3339(date time.Time) string
	FormatToUnix(date time.Time) int64

	// Validation methods
	IsLeapYear(year int) bool
	IsWeekday(date time.Time) bool
	IsWeekend(date time.Time) bool

	// Date arithmetic - most common operations
	AddDays(date time.Time, days int) time.Time
	AddMonths(date time.Time, months int) time.Time
	AddYears(date time.Time, years int) time.Time

	// Date comparison
	DaysBetween(start, end time.Time) int
	IsSameDay(date1, date2 time.Time) bool
	IsSameMonth(date1, date2 time.Time) bool
	IsAfter(date1, date2 time.Time) bool
	IsBefore(date1, date2 time.Time) bool

	// Most commonly used date boundaries
	FirstDayOfMonth(date time.Time) time.Time
	LastDayOfMonth(date time.Time) time.Time
	FirstDayOfYear(date time.Time) time.Time
	LastDayOfYear(date time.Time) time.Time
	StartOfDay(date time.Time) time.Time
	EndOfDay(date time.Time) time.Time

	// Current time helpers - most essential
	Now() time.Time
	NowUTC() time.Time
	Today() time.Time
	Yesterday() time.Time
	Tomorrow() time.Time
	LastMonth() time.Time
	NextMonth() time.Time

	// Most common utility methods
	GetDaysInMonth(year, month int) int
	IsBusinessDay(date time.Time) bool
	NextBusinessDay(date time.Time) time.Time

	// Essential formats
	GetCommonFormats() []string
}

// DateUtil provides comprehensive date utility operations
type DateUtil struct{}

// NewDateUtil creates a new instance of DateUtil
func NewDateUtil() DateClient {
	return &DateUtil{}
}

// Parse attempts to parse a date string using the provided formats or common formats
func (d *DateUtil) Parse(dateStr string, formats ...string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}

	parseFormats := formats
	if len(parseFormats) == 0 {
		parseFormats = d.GetCommonFormats()
	}

	var lastErr error
	for _, format := range parseFormats {
		if parsedTime, err := time.Parse(format, dateStr); err == nil {
			return parsedTime, nil
		} else {
			lastErr = err
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date '%s': %v", dateStr, lastErr)
}

// ParseUnix parses a Unix timestamp (int, int64, float64, or string)
func (d *DateUtil) ParseUnix(timestamp any) (time.Time, error) {
	switch v := timestamp.(type) {
	case int:
		return time.Unix(int64(v), 0), nil
	case int64:
		return time.Unix(v, 0), nil
	case float64:
		return time.Unix(int64(v), 0), nil
	case string:
		ts, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid unix timestamp string: %v", err)
		}
		return time.Unix(ts, 0), nil
	default:
		return time.Time{}, fmt.Errorf("unsupported timestamp type: %T", timestamp)
	}
}

// Format formats a time using the specified format
func (d *DateUtil) Format(date time.Time, format string) string {
	return date.Format(format)
}

// FormatToRFC3339 formats a time to RFC3339 format
func (d *DateUtil) FormatToRFC3339(date time.Time) string {
	return date.Format(time.RFC3339)
}

// FormatToUnix formats a time to Unix timestamp
func (d *DateUtil) FormatToUnix(date time.Time) int64 {
	return date.Unix()
}

// IsLeapYear checks if the given year is a leap year
func (d *DateUtil) IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// IsWeekday checks if the given date is a weekday (Monday-Friday)
func (d *DateUtil) IsWeekday(date time.Time) bool {
	weekday := date.Weekday()
	return weekday >= time.Monday && weekday <= time.Friday
}

// IsWeekend checks if the given date is a weekend (Saturday-Sunday)
func (d *DateUtil) IsWeekend(date time.Time) bool {
	return !d.IsWeekday(date)
}

// AddDays adds the specified number of days to a date
func (d *DateUtil) AddDays(date time.Time, days int) time.Time {
	return date.AddDate(0, 0, days)
}

// AddMonths adds the specified number of months to a date
func (d *DateUtil) AddMonths(date time.Time, months int) time.Time {
	return date.AddDate(0, months, 0)
}

// AddYears adds the specified number of years to a date
func (d *DateUtil) AddYears(date time.Time, years int) time.Time {
	return date.AddDate(years, 0, 0)
}

// DaysBetween calculates the number of days between two dates
func (d *DateUtil) DaysBetween(start, end time.Time) int {
	if start.After(end) {
		start, end = end, start
	}
	return int(end.Sub(start).Hours() / 24)
}

// IsSameDay checks if two dates are on the same day
func (d *DateUtil) IsSameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsSameMonth checks if two dates are in the same month and year
func (d *DateUtil) IsSameMonth(date1, date2 time.Time) bool {
	return date1.Year() == date2.Year() && date1.Month() == date2.Month()
}

// IsAfter checks if date1 is after date2
func (d *DateUtil) IsAfter(date1, date2 time.Time) bool {
	return date1.After(date2)
}

// IsBefore checks if date1 is before date2
func (d *DateUtil) IsBefore(date1, date2 time.Time) bool {
	return date1.Before(date2)
}

// StartOfDay returns the start of the day (00:00:00)
func (d *DateUtil) StartOfDay(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, date.Location())
}

// EndOfDay returns the end of the day (23:59:59.999999999)
func (d *DateUtil) EndOfDay(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, date.Location())
}

// FirstDayOfMonth returns the first day of the month (1st day 00:00:00)
func (d *DateUtil) FirstDayOfMonth(date time.Time) time.Time {
	year, month, _ := date.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, date.Location())
}

// LastDayOfMonth returns the last day of the month (last day 23:59:59.999999999)
func (d *DateUtil) LastDayOfMonth(date time.Time) time.Time {
	year, month, _ := date.Date()
	lastDay := d.GetDaysInMonth(year, int(month))
	return time.Date(year, month, lastDay, 23, 59, 59, 999999999, date.Location())
}

// FirstDayOfYear returns the first day of the year (January 1st 00:00:00)
func (d *DateUtil) FirstDayOfYear(date time.Time) time.Time {
	return time.Date(date.Year(), time.January, 1, 0, 0, 0, 0, date.Location())
}

// LastDayOfYear returns the last day of the year (December 31st 23:59:59.999999999)
func (d *DateUtil) LastDayOfYear(date time.Time) time.Time {
	return time.Date(date.Year(), time.December, 31, 23, 59, 59, 999999999, date.Location())
}

// GetDaysInMonth returns the number of days in the specified month and year
func (d *DateUtil) GetDaysInMonth(year, month int) int {
	// Create a date for the first day of the next month, then subtract a day
	nextMonth := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
	lastDay := nextMonth.AddDate(0, 0, -1)
	return lastDay.Day()
}

// IsBusinessDay checks if the given date is a business day (Monday-Friday)
func (d *DateUtil) IsBusinessDay(date time.Time) bool {
	// Basic implementation - only checks for weekdays
	// In a real implementation, you might want to include holiday checking
	return d.IsWeekday(date)
}

// NextBusinessDay returns the next business day
func (d *DateUtil) NextBusinessDay(date time.Time) time.Time {
	next := d.AddDays(date, 1)
	for !d.IsBusinessDay(next) {
		next = d.AddDays(next, 1)
	}
	return next
}

// Now returns the current time
func (d *DateUtil) Now() time.Time {
	return time.Now()
}

// NowUTC returns the current time in UTC
func (d *DateUtil) NowUTC() time.Time {
	return time.Now().UTC()
}

// Today returns today's date at 00:00:00
func (d *DateUtil) Today() time.Time {
	return d.StartOfDay(time.Now())
}

// Yesterday returns yesterday's date at 00:00:00
func (d *DateUtil) Yesterday() time.Time {
	return d.StartOfDay(d.AddDays(time.Now(), -1))
}

// Tomorrow returns tomorrow's date at 00:00:00
func (d *DateUtil) Tomorrow() time.Time {
	return d.StartOfDay(d.AddDays(time.Now(), 1))
}

// LastMonth returns the same day last month at 00:00:00
func (d *DateUtil) LastMonth() time.Time {
	return d.StartOfDay(d.AddMonths(time.Now(), -1))
}

// NextMonth returns the same day next month at 00:00:00
func (d *DateUtil) NextMonth() time.Time {
	return d.StartOfDay(d.AddMonths(time.Now(), 1))
}

// GetCommonFormats returns a list of commonly used date formats
func (d *DateUtil) GetCommonFormats() []string {
	return []string{
		time.RFC3339,
		RFC3339Date,
		SimpleDateTime,
		USDate,
		EuropeanDate,
	}
}
