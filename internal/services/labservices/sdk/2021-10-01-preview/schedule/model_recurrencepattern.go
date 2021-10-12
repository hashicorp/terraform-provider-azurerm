package schedule

type RecurrencePattern struct {
	ExpirationDate string              `json:"expirationDate"`
	Frequency      RecurrenceFrequency `json:"frequency"`
	Interval       *int64              `json:"interval,omitempty"`
	WeekDays       *[]WeekDay          `json:"weekDays,omitempty"`
}
