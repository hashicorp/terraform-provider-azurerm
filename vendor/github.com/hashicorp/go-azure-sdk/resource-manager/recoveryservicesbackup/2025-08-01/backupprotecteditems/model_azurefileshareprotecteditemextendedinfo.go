package backupprotecteditems

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFileshareProtectedItemExtendedInfo struct {
	OldestRecoveryPoint   *string `json:"oldestRecoveryPoint,omitempty"`
	PolicyState           *string `json:"policyState,omitempty"`
	RecoveryPointCount    *int64  `json:"recoveryPointCount,omitempty"`
	ResourceState         *string `json:"resourceState,omitempty"`
	ResourceStateSyncTime *string `json:"resourceStateSyncTime,omitempty"`
}

func (o *AzureFileshareProtectedItemExtendedInfo) GetOldestRecoveryPointAsTime() (*time.Time, error) {
	if o.OldestRecoveryPoint == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.OldestRecoveryPoint, "2006-01-02T15:04:05Z07:00")
}

func (o *AzureFileshareProtectedItemExtendedInfo) SetOldestRecoveryPointAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.OldestRecoveryPoint = &formatted
}

func (o *AzureFileshareProtectedItemExtendedInfo) GetResourceStateSyncTimeAsTime() (*time.Time, error) {
	if o.ResourceStateSyncTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ResourceStateSyncTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AzureFileshareProtectedItemExtendedInfo) SetResourceStateSyncTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ResourceStateSyncTime = &formatted
}
