package servergroups

type MaintenanceWindow struct {
	CustomWindow *string `json:"customWindow,omitempty"`
	DayOfWeek    *int64  `json:"dayOfWeek,omitempty"`
	StartHour    *int64  `json:"startHour,omitempty"`
	StartMinute  *int64  `json:"startMinute,omitempty"`
}
