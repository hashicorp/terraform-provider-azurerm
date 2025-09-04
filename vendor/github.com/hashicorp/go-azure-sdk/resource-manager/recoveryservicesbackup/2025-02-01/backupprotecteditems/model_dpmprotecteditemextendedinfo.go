package backupprotecteditems

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DPMProtectedItemExtendedInfo struct {
	DiskStorageUsedInBytes       *string            `json:"diskStorageUsedInBytes,omitempty"`
	IsCollocated                 *bool              `json:"isCollocated,omitempty"`
	IsPresentOnCloud             *bool              `json:"isPresentOnCloud,omitempty"`
	LastBackupStatus             *string            `json:"lastBackupStatus,omitempty"`
	LastRefreshedAt              *string            `json:"lastRefreshedAt,omitempty"`
	OldestRecoveryPoint          *string            `json:"oldestRecoveryPoint,omitempty"`
	OnPremiseLatestRecoveryPoint *string            `json:"onPremiseLatestRecoveryPoint,omitempty"`
	OnPremiseOldestRecoveryPoint *string            `json:"onPremiseOldestRecoveryPoint,omitempty"`
	OnPremiseRecoveryPointCount  *int64             `json:"onPremiseRecoveryPointCount,omitempty"`
	ProtectableObjectLoadPath    *map[string]string `json:"protectableObjectLoadPath,omitempty"`
	Protected                    *bool              `json:"protected,omitempty"`
	ProtectionGroupName          *string            `json:"protectionGroupName,omitempty"`
	RecoveryPointCount           *int64             `json:"recoveryPointCount,omitempty"`
	TotalDiskStorageSizeInBytes  *string            `json:"totalDiskStorageSizeInBytes,omitempty"`
}

func (o *DPMProtectedItemExtendedInfo) GetLastRefreshedAtAsTime() (*time.Time, error) {
	if o.LastRefreshedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastRefreshedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *DPMProtectedItemExtendedInfo) SetLastRefreshedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastRefreshedAt = &formatted
}

func (o *DPMProtectedItemExtendedInfo) GetOldestRecoveryPointAsTime() (*time.Time, error) {
	if o.OldestRecoveryPoint == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.OldestRecoveryPoint, "2006-01-02T15:04:05Z07:00")
}

func (o *DPMProtectedItemExtendedInfo) SetOldestRecoveryPointAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.OldestRecoveryPoint = &formatted
}

func (o *DPMProtectedItemExtendedInfo) GetOnPremiseLatestRecoveryPointAsTime() (*time.Time, error) {
	if o.OnPremiseLatestRecoveryPoint == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.OnPremiseLatestRecoveryPoint, "2006-01-02T15:04:05Z07:00")
}

func (o *DPMProtectedItemExtendedInfo) SetOnPremiseLatestRecoveryPointAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.OnPremiseLatestRecoveryPoint = &formatted
}

func (o *DPMProtectedItemExtendedInfo) GetOnPremiseOldestRecoveryPointAsTime() (*time.Time, error) {
	if o.OnPremiseOldestRecoveryPoint == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.OnPremiseOldestRecoveryPoint, "2006-01-02T15:04:05Z07:00")
}

func (o *DPMProtectedItemExtendedInfo) SetOnPremiseOldestRecoveryPointAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.OnPremiseOldestRecoveryPoint = &formatted
}
