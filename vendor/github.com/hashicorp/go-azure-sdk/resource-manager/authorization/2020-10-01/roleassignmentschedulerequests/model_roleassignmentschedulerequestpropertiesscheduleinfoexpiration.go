package roleassignmentschedulerequests

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleAssignmentScheduleRequestPropertiesScheduleInfoExpiration struct {
	Duration    *string `json:"duration,omitempty"`
	EndDateTime *string `json:"endDateTime,omitempty"`
	Type        *Type   `json:"type,omitempty"`
}

func (o *RoleAssignmentScheduleRequestPropertiesScheduleInfoExpiration) GetEndDateTimeAsTime() (*time.Time, error) {
	if o.EndDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RoleAssignmentScheduleRequestPropertiesScheduleInfoExpiration) SetEndDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndDateTime = &formatted
}
