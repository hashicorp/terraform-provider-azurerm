package protecteditems

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MabFileFolderProtectedItemExtendedInfo struct {
	LastRefreshedAt     *string `json:"lastRefreshedAt,omitempty"`
	OldestRecoveryPoint *string `json:"oldestRecoveryPoint,omitempty"`
	RecoveryPointCount  *int64  `json:"recoveryPointCount,omitempty"`
}

func (o *MabFileFolderProtectedItemExtendedInfo) GetLastRefreshedAtAsTime() (*time.Time, error) {
	if o.LastRefreshedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastRefreshedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *MabFileFolderProtectedItemExtendedInfo) SetLastRefreshedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastRefreshedAt = &formatted
}

func (o *MabFileFolderProtectedItemExtendedInfo) GetOldestRecoveryPointAsTime() (*time.Time, error) {
	if o.OldestRecoveryPoint == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.OldestRecoveryPoint, "2006-01-02T15:04:05Z07:00")
}

func (o *MabFileFolderProtectedItemExtendedInfo) SetOldestRecoveryPointAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.OldestRecoveryPoint = &formatted
}
