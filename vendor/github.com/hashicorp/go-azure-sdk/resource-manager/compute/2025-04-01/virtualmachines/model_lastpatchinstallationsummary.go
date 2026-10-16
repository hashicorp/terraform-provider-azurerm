package virtualmachines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LastPatchInstallationSummary struct {
	Error                     *ApiError             `json:"error,omitempty"`
	ExcludedPatchCount        *int64                `json:"excludedPatchCount,omitempty"`
	FailedPatchCount          *int64                `json:"failedPatchCount,omitempty"`
	InstallationActivityId    *string               `json:"installationActivityId,omitempty"`
	InstalledPatchCount       *int64                `json:"installedPatchCount,omitempty"`
	LastModifiedTime          *string               `json:"lastModifiedTime,omitempty"`
	MaintenanceWindowExceeded *bool                 `json:"maintenanceWindowExceeded,omitempty"`
	NotSelectedPatchCount     *int64                `json:"notSelectedPatchCount,omitempty"`
	PendingPatchCount         *int64                `json:"pendingPatchCount,omitempty"`
	StartTime                 *string               `json:"startTime,omitempty"`
	Status                    *PatchOperationStatus `json:"status,omitempty"`
}

func (o *LastPatchInstallationSummary) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *LastPatchInstallationSummary) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}

func (o *LastPatchInstallationSummary) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *LastPatchInstallationSummary) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
