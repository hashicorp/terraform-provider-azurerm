package rolemanagementpolicyassignments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyAssignmentPropertiesPolicy struct {
	Id                   *string    `json:"id,omitempty"`
	LastModifiedBy       *Principal `json:"lastModifiedBy,omitempty"`
	LastModifiedDateTime *string    `json:"lastModifiedDateTime,omitempty"`
}

func (o *PolicyAssignmentPropertiesPolicy) GetLastModifiedDateTimeAsTime() (*time.Time, error) {
	if o.LastModifiedDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *PolicyAssignmentPropertiesPolicy) SetLastModifiedDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedDateTime = &formatted
}
