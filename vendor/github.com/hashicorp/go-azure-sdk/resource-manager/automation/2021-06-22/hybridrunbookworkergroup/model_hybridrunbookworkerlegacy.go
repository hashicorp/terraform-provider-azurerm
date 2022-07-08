package hybridrunbookworkergroup

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HybridRunbookWorkerLegacy struct {
	IP               *string `json:"ip,omitempty"`
	LastSeenDateTime *string `json:"lastSeenDateTime,omitempty"`
	Name             *string `json:"name,omitempty"`
	RegistrationTime *string `json:"registrationTime,omitempty"`
}

func (o *HybridRunbookWorkerLegacy) GetLastSeenDateTimeAsTime() (*time.Time, error) {
	if o.LastSeenDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSeenDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *HybridRunbookWorkerLegacy) SetLastSeenDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSeenDateTime = &formatted
}

func (o *HybridRunbookWorkerLegacy) GetRegistrationTimeAsTime() (*time.Time, error) {
	if o.RegistrationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RegistrationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *HybridRunbookWorkerLegacy) SetRegistrationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RegistrationTime = &formatted
}
