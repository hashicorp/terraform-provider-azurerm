package machines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MachineAssessPatchesResult struct {
	AssessmentActivityId                *string                              `json:"assessmentActivityId,omitempty"`
	AvailablePatchCountByClassification *AvailablePatchCountByClassification `json:"availablePatchCountByClassification,omitempty"`
	ErrorDetails                        *ErrorDetail                         `json:"errorDetails,omitempty"`
	LastModifiedDateTime                *string                              `json:"lastModifiedDateTime,omitempty"`
	OsType                              *OsType                              `json:"osType,omitempty"`
	PatchServiceUsed                    *PatchServiceUsed                    `json:"patchServiceUsed,omitempty"`
	RebootPending                       *bool                                `json:"rebootPending,omitempty"`
	StartDateTime                       *string                              `json:"startDateTime,omitempty"`
	StartedBy                           *PatchOperationStartedBy             `json:"startedBy,omitempty"`
	Status                              *PatchOperationStatus                `json:"status,omitempty"`
}

func (o *MachineAssessPatchesResult) GetLastModifiedDateTimeAsTime() (*time.Time, error) {
	if o.LastModifiedDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MachineAssessPatchesResult) SetLastModifiedDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedDateTime = &formatted
}

func (o *MachineAssessPatchesResult) GetStartDateTimeAsTime() (*time.Time, error) {
	if o.StartDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MachineAssessPatchesResult) SetStartDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartDateTime = &formatted
}
