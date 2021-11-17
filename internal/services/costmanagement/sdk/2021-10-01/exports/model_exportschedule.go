package exports

type ExportSchedule struct {
	Recurrence       *RecurrenceType         `json:"recurrence,omitempty"`
	RecurrencePeriod *ExportRecurrencePeriod `json:"recurrencePeriod,omitempty"`
	Status           *StatusType             `json:"status,omitempty"`
}
