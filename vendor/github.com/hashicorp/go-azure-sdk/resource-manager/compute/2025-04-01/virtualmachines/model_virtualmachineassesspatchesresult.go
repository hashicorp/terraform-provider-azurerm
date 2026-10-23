package virtualmachines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineAssessPatchesResult struct {
	AssessmentActivityId          *string                                  `json:"assessmentActivityId,omitempty"`
	AvailablePatches              *[]VirtualMachineSoftwarePatchProperties `json:"availablePatches,omitempty"`
	CriticalAndSecurityPatchCount *int64                                   `json:"criticalAndSecurityPatchCount,omitempty"`
	Error                         *ApiError                                `json:"error,omitempty"`
	OtherPatchCount               *int64                                   `json:"otherPatchCount,omitempty"`
	RebootPending                 *bool                                    `json:"rebootPending,omitempty"`
	StartDateTime                 *string                                  `json:"startDateTime,omitempty"`
	Status                        *PatchOperationStatus                    `json:"status,omitempty"`
}

func (o *VirtualMachineAssessPatchesResult) GetStartDateTimeAsTime() (*time.Time, error) {
	if o.StartDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *VirtualMachineAssessPatchesResult) SetStartDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartDateTime = &formatted
}
