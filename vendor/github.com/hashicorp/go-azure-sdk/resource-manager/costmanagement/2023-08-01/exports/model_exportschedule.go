package exports

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportSchedule struct {
	Recurrence       *RecurrenceType         `json:"recurrence,omitempty"`
	RecurrencePeriod *ExportRecurrencePeriod `json:"recurrencePeriod,omitempty"`
	Status           *StatusType             `json:"status,omitempty"`
}
