package dateutil

import (
	"testing"
	"time"
)

func TestNewDateUtil(t *testing.T) {
	util := NewDateUtil()
	if util == nil {
		t.Error("NewDateUtil() returned nil")
	}
}

// =================== Test Parsing Methods ===================

func TestParse(t *testing.T) {
	util := NewDateUtil()

	tests := []struct {
		name        string
		dateStr     string
		formats     []string
		expectError bool
	}{
		{
			name:        "RFC3339 format",
			dateStr:     "2023-10-05T14:30:00Z",
			formats:     []string{time.RFC3339},
			expectError: false,
		},
		{
			name:        "Simple date format",
			dateStr:     "2023-10-05",
			formats:     []string{RFC3339Date},
			expectError: false,
		},
		{
			name:        "US date format",
			dateStr:     "10/05/2023",
			formats:     []string{USDate},
			expectError: false,
		},
		{
			name:        "Invalid date string",
			dateStr:     "invalid-date",
			formats:     []string{time.RFC3339},
			expectError: true,
		},
		{
			name:        "Empty date string",
			dateStr:     "",
			formats:     []string{time.RFC3339},
			expectError: true,
		},
		{
			name:        "No formats provided - use common formats",
			dateStr:     "2023-10-05T14:30:00Z",
			formats:     nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := util.Parse(tt.dateStr, tt.formats...)
			if (err != nil) != tt.expectError {
				t.Errorf("Parse(%s) error = %v, expectError %v", tt.dateStr, err, tt.expectError)
				return
			}
			if !tt.expectError && result.IsZero() {
				t.Errorf("Parse(%s) returned zero time", tt.dateStr)
			}
		})
	}
}

func TestParseUnix(t *testing.T) {
	util := NewDateUtil()

	tests := []struct {
		name        string
		timestamp   any
		expectError bool
	}{
		{
			name:        "int timestamp",
			timestamp:   1696518600,
			expectError: false,
		},
		{
			name:        "int64 timestamp",
			timestamp:   int64(1696518600),
			expectError: false,
		},
		{
			name:        "float64 timestamp",
			timestamp:   float64(1696518600),
			expectError: false,
		},
		{
			name:        "string timestamp",
			timestamp:   "1696518600",
			expectError: false,
		},
		{
			name:        "invalid string timestamp",
			timestamp:   "invalid",
			expectError: true,
		},
		{
			name:        "unsupported type",
			timestamp:   []int{1696518600},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := util.ParseUnix(tt.timestamp)
			if (err != nil) != tt.expectError {
				t.Errorf("ParseUnix(%v) error = %v, expectError %v", tt.timestamp, err, tt.expectError)
				return
			}
			if !tt.expectError && result.Unix() != int64(1696518600) {
				t.Errorf("ParseUnix(%v) = %v, want unix timestamp 1696518600", tt.timestamp, result.Unix())
			}
		})
	}
}

// =================== Test Formatting Methods ===================

func TestFormatMethods(t *testing.T) {
	util := NewDateUtil()
	testTime := time.Date(2023, 10, 5, 14, 30, 45, 123456789, time.UTC)

	t.Run("Format", func(t *testing.T) {
		result := util.Format(testTime, RFC3339Date)
		expected := "2023-10-05"
		if result != expected {
			t.Errorf("Format() = %v, want %v", result, expected)
		}
	})

	t.Run("FormatToRFC3339", func(t *testing.T) {
		result := util.FormatToRFC3339(testTime)
		expected := "2023-10-05T14:30:45Z"
		if result != expected {
			t.Errorf("FormatToRFC3339() = %v, want %v", result, expected)
		}
	})

	t.Run("FormatToUnix", func(t *testing.T) {
		result := util.FormatToUnix(testTime)
		expected := testTime.Unix()
		if result != expected {
			t.Errorf("FormatToUnix() = %v, want %v", result, expected)
		}
	})
}

// =================== Test Validation Methods ===================

func TestValidationMethods(t *testing.T) {
	util := NewDateUtil()

	t.Run("IsLeapYear", func(t *testing.T) {
		tests := []struct {
			year     int
			expected bool
		}{
			{2024, true},  // Divisible by 4
			{2023, false}, // Not divisible by 4
			{1900, false}, // Divisible by 100 but not 400
			{2000, true},  // Divisible by 400
		}

		for _, tt := range tests {
			result := util.IsLeapYear(tt.year)
			if result != tt.expected {
				t.Errorf("IsLeapYear(%d) = %v, want %v", tt.year, result, tt.expected)
			}
		}
	})

	t.Run("IsWeekday and IsWeekend", func(t *testing.T) {
		monday := time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC)   // Monday
		saturday := time.Date(2023, 10, 7, 12, 0, 0, 0, time.UTC) // Saturday

		if !util.IsWeekday(monday) {
			t.Error("IsWeekday should return true for Monday")
		}

		if util.IsWeekend(monday) {
			t.Error("IsWeekend should return false for Monday")
		}

		if util.IsWeekday(saturday) {
			t.Error("IsWeekday should return false for Saturday")
		}

		if !util.IsWeekend(saturday) {
			t.Error("IsWeekend should return true for Saturday")
		}
	})
}

// =================== Test Date Arithmetic ===================

func TestDateArithmetic(t *testing.T) {
	util := NewDateUtil()
	baseDate := time.Date(2023, 10, 5, 12, 0, 0, 0, time.UTC)

	t.Run("AddDays", func(t *testing.T) {
		result := util.AddDays(baseDate, 5)
		expected := time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("AddDays() = %v, want %v", result, expected)
		}
	})

	t.Run("AddMonths", func(t *testing.T) {
		result := util.AddMonths(baseDate, 2)
		expected := time.Date(2023, 12, 5, 12, 0, 0, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("AddMonths() = %v, want %v", result, expected)
		}
	})

	t.Run("AddYears", func(t *testing.T) {
		result := util.AddYears(baseDate, 1)
		expected := time.Date(2024, 10, 5, 12, 0, 0, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("AddYears() = %v, want %v", result, expected)
		}
	})
}

// =================== Test Date Comparison ===================

func TestDateComparison(t *testing.T) {
	util := NewDateUtil()
	date1 := time.Date(2023, 10, 5, 12, 0, 0, 0, time.UTC)
	date2 := time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC)

	t.Run("DaysBetween", func(t *testing.T) {
		result := util.DaysBetween(date1, date2)
		expected := 5
		if result != expected {
			t.Errorf("DaysBetween() = %v, want %v", result, expected)
		}

		// Test reverse order
		result = util.DaysBetween(date2, date1)
		if result != expected {
			t.Errorf("DaysBetween() reverse = %v, want %v", result, expected)
		}
	})

	t.Run("IsSameDay", func(t *testing.T) {
		sameDay := time.Date(2023, 10, 5, 18, 0, 0, 0, time.UTC)
		if !util.IsSameDay(date1, sameDay) {
			t.Error("IsSameDay should return true for same date")
		}

		if util.IsSameDay(date1, date2) {
			t.Error("IsSameDay should return false for different dates")
		}
	})

	t.Run("IsSameMonth", func(t *testing.T) {
		sameMonth := time.Date(2023, 10, 15, 12, 0, 0, 0, time.UTC)
		differentMonth := time.Date(2023, 11, 5, 12, 0, 0, 0, time.UTC)

		if !util.IsSameMonth(date1, sameMonth) {
			t.Error("IsSameMonth should return true for same month")
		}

		if util.IsSameMonth(date1, differentMonth) {
			t.Error("IsSameMonth should return false for different month")
		}
	})

	t.Run("IsAfter", func(t *testing.T) {
		// date1 is 2023-10-05, date2 is 2023-10-10
		if util.IsAfter(date1, date2) {
			t.Error("IsAfter should return false when date1 is before date2")
		}

		if !util.IsAfter(date2, date1) {
			t.Error("IsAfter should return true when date2 is after date1")
		}

		// Test with same time
		if util.IsAfter(date1, date1) {
			t.Error("IsAfter should return false when dates are equal")
		}
	})

	t.Run("IsBefore", func(t *testing.T) {
		// date1 is 2023-10-05, date2 is 2023-10-10
		if !util.IsBefore(date1, date2) {
			t.Error("IsBefore should return true when date1 is before date2")
		}

		if util.IsBefore(date2, date1) {
			t.Error("IsBefore should return false when date2 is after date1")
		}

		// Test with same time
		if util.IsBefore(date1, date1) {
			t.Error("IsBefore should return false when dates are equal")
		}
	})
}

// =================== Test Date Boundaries ===================

func TestDateBoundaries(t *testing.T) {
	util := NewDateUtil()
	testDate := time.Date(2023, 10, 15, 14, 30, 45, 123456789, time.UTC)

	t.Run("FirstDayOfMonth", func(t *testing.T) {
		result := util.FirstDayOfMonth(testDate)
		expected := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("FirstDayOfMonth() = %v, want %v", result, expected)
		}
	})

	t.Run("LastDayOfMonth", func(t *testing.T) {
		result := util.LastDayOfMonth(testDate)
		expected := time.Date(2023, 10, 31, 23, 59, 59, 999999999, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("LastDayOfMonth() = %v, want %v", result, expected)
		}
	})

	t.Run("FirstDayOfYear", func(t *testing.T) {
		result := util.FirstDayOfYear(testDate)
		expected := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("FirstDayOfYear() = %v, want %v", result, expected)
		}
	})

	t.Run("LastDayOfYear", func(t *testing.T) {
		result := util.LastDayOfYear(testDate)
		expected := time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("LastDayOfYear() = %v, want %v", result, expected)
		}
	})

	t.Run("StartOfDay", func(t *testing.T) {
		result := util.StartOfDay(testDate)
		expected := time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("StartOfDay() = %v, want %v", result, expected)
		}
	})

	t.Run("EndOfDay", func(t *testing.T) {
		result := util.EndOfDay(testDate)
		expected := time.Date(2023, 10, 15, 23, 59, 59, 999999999, time.UTC)
		if !result.Equal(expected) {
			t.Errorf("EndOfDay() = %v, want %v", result, expected)
		}
	})
}

// =================== Test Current Time Helpers ===================

func TestCurrentTimeHelpers(t *testing.T) {
	util := NewDateUtil()

	t.Run("Now", func(t *testing.T) {
		result := util.Now()
		if result.IsZero() {
			t.Error("Now() should not return zero time")
		}
	})

	t.Run("NowUTC", func(t *testing.T) {
		result := util.NowUTC()
		if result.IsZero() {
			t.Error("NowUTC() should not return zero time")
		}
		if result.Location() != time.UTC {
			t.Error("NowUTC() should return time in UTC location")
		}
	})

	t.Run("Today", func(t *testing.T) {
		result := util.Today()
		if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
			t.Error("Today() should return start of day")
		}
	})

	t.Run("Yesterday", func(t *testing.T) {
		today := util.Today()
		yesterday := util.Yesterday()
		diff := today.Sub(yesterday)
		expected := 24 * time.Hour
		if diff != expected {
			t.Errorf("Yesterday() difference from today = %v, want %v", diff, expected)
		}
	})

	t.Run("Tomorrow", func(t *testing.T) {
		today := util.Today()
		tomorrow := util.Tomorrow()
		diff := tomorrow.Sub(today)
		expected := 24 * time.Hour
		if diff != expected {
			t.Errorf("Tomorrow() difference from today = %v, want %v", diff, expected)
		}
	})

	t.Run("LastMonth", func(t *testing.T) {
		now := util.Now()
		lastMonth := util.LastMonth()

		// Check that last month is approximately 30 days ago (allowing for month variations)
		diff := now.Sub(lastMonth)
		if diff < 25*24*time.Hour || diff > 35*24*time.Hour {
			t.Errorf("LastMonth() difference from now = %v, expected between 25-35 days", diff)
		}

		// Check that it's the start of day
		if lastMonth.Hour() != 0 || lastMonth.Minute() != 0 || lastMonth.Second() != 0 {
			t.Error("LastMonth() should return start of day")
		}
	})

	t.Run("NextMonth", func(t *testing.T) {
		now := util.Now()
		nextMonth := util.NextMonth()

		// Check that next month is approximately 30 days from now (allowing for month variations)
		diff := nextMonth.Sub(now)
		if diff < 25*24*time.Hour || diff > 35*24*time.Hour {
			t.Errorf("NextMonth() difference from now = %v, expected between 25-35 days", diff)
		}

		// Check that it's the start of day
		if nextMonth.Hour() != 0 || nextMonth.Minute() != 0 || nextMonth.Second() != 0 {
			t.Error("NextMonth() should return start of day")
		}
	})
}

// =================== Test Utility Methods ===================

func TestUtilityMethods(t *testing.T) {
	util := NewDateUtil()

	t.Run("GetDaysInMonth", func(t *testing.T) {
		tests := []struct {
			year     int
			month    int
			expected int
		}{
			{2023, 10, 31}, // October
			{2023, 2, 28},  // February non-leap year
			{2024, 2, 29},  // February leap year
			{2023, 4, 30},  // April
		}

		for _, tt := range tests {
			result := util.GetDaysInMonth(tt.year, tt.month)
			if result != tt.expected {
				t.Errorf("GetDaysInMonth(%d, %d) = %v, want %v", tt.year, tt.month, result, tt.expected)
			}
		}
	})

	t.Run("IsBusinessDay", func(t *testing.T) {
		monday := time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC)
		saturday := time.Date(2023, 10, 7, 12, 0, 0, 0, time.UTC)

		if !util.IsBusinessDay(monday) {
			t.Error("IsBusinessDay should return true for Monday")
		}

		if util.IsBusinessDay(saturday) {
			t.Error("IsBusinessDay should return false for Saturday")
		}
	})

	t.Run("NextBusinessDay", func(t *testing.T) {
		friday := time.Date(2023, 10, 6, 12, 0, 0, 0, time.UTC)
		result := util.NextBusinessDay(friday)
		// Next business day after Friday should be Monday
		if result.Weekday() != time.Monday {
			t.Errorf("NextBusinessDay() weekday = %v, want %v", result.Weekday(), time.Monday)
		}
	})
}

// =================== Test GetCommonFormats ===================

func TestGetCommonFormats(t *testing.T) {
	util := NewDateUtil()
	formats := util.GetCommonFormats()

	if len(formats) == 0 {
		t.Error("GetCommonFormats() should return non-empty slice")
	}

	expectedFormats := []string{
		time.RFC3339,
		RFC3339Date,
		SimpleDateTime,
		USDate,
		EuropeanDate,
	}

	for _, expected := range expectedFormats {
		found := false
		for _, format := range formats {
			if format == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetCommonFormats() missing expected format: %s", expected)
		}
	}
}

// =================== Benchmarks ===================

func BenchmarkParse(b *testing.B) {
	util := NewDateUtil()
	dateStr := "2023-10-05T14:30:00Z"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = util.Parse(dateStr, time.RFC3339)
	}
}

func BenchmarkFormat(b *testing.B) {
	util := NewDateUtil()
	testDate := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = util.Format(testDate, RFC3339Date)
	}
}

func BenchmarkAddDays(b *testing.B) {
	util := NewDateUtil()
	testDate := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = util.AddDays(testDate, 30)
	}
}

func BenchmarkDaysBetween(b *testing.B) {
	util := NewDateUtil()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = util.DaysBetween(start, end)
	}
}
