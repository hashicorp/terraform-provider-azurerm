package virtualmachines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailablePatchSummary struct {
	AssessmentActivityId          *string               `json:"assessmentActivityId,omitempty"`
	CriticalAndSecurityPatchCount *int64                `json:"criticalAndSecurityPatchCount,omitempty"`
	Error                         *ApiError             `json:"error,omitempty"`
	LastModifiedTime              *string               `json:"lastModifiedTime,omitempty"`
	OtherPatchCount               *int64                `json:"otherPatchCount,omitempty"`
	RebootPending                 *bool                 `json:"rebootPending,omitempty"`
	StartTime                     *string               `json:"startTime,omitempty"`
	Status                        *PatchOperationStatus `json:"status,omitempty"`
}

func (o *AvailablePatchSummary) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AvailablePatchSummary) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}

func (o *AvailablePatchSummary) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AvailablePatchSummary) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
