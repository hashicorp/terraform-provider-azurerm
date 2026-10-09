package virtualmachines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineInstallPatchesResult struct {
	Error                     *ApiError                  `json:"error,omitempty"`
	ExcludedPatchCount        *int64                     `json:"excludedPatchCount,omitempty"`
	FailedPatchCount          *int64                     `json:"failedPatchCount,omitempty"`
	InstallationActivityId    *string                    `json:"installationActivityId,omitempty"`
	InstalledPatchCount       *int64                     `json:"installedPatchCount,omitempty"`
	MaintenanceWindowExceeded *bool                      `json:"maintenanceWindowExceeded,omitempty"`
	NotSelectedPatchCount     *int64                     `json:"notSelectedPatchCount,omitempty"`
	Patches                   *[]PatchInstallationDetail `json:"patches,omitempty"`
	PendingPatchCount         *int64                     `json:"pendingPatchCount,omitempty"`
	RebootStatus              *VMGuestPatchRebootStatus  `json:"rebootStatus,omitempty"`
	StartDateTime             *string                    `json:"startDateTime,omitempty"`
	Status                    *PatchOperationStatus      `json:"status,omitempty"`
}

func (o *VirtualMachineInstallPatchesResult) GetStartDateTimeAsTime() (*time.Time, error) {
	if o.StartDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *VirtualMachineInstallPatchesResult) SetStartDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartDateTime = &formatted
}
