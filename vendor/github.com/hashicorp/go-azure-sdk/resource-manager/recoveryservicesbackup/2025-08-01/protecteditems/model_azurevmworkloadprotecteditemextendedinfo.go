package protecteditems

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureVMWorkloadProtectedItemExtendedInfo struct {
	NewestRecoveryPointInArchive *string `json:"newestRecoveryPointInArchive,omitempty"`
	OldestRecoveryPoint          *string `json:"oldestRecoveryPoint,omitempty"`
	OldestRecoveryPointInArchive *string `json:"oldestRecoveryPointInArchive,omitempty"`
	OldestRecoveryPointInVault   *string `json:"oldestRecoveryPointInVault,omitempty"`
	PolicyState                  *string `json:"policyState,omitempty"`
	RecoveryModel                *string `json:"recoveryModel,omitempty"`
	RecoveryPointCount           *int64  `json:"recoveryPointCount,omitempty"`
}

func (o *AzureVMWorkloadProtectedItemExtendedInfo) GetNewestRecoveryPointInArchiveAsTime() (*time.Time, error) {
	if o.NewestRecoveryPointInArchive == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NewestRecoveryPointInArchive, "2006-01-02T15:04:05Z07:00")
}

func (o *AzureVMWorkloadProtectedItemExtendedInfo) SetNewestRecoveryPointInArchiveAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NewestRecoveryPointInArchive = &formatted
}

func (o *AzureVMWorkloadProtectedItemExtendedInfo) GetOldestRecoveryPointAsTime() (*time.Time, error) {
	if o.OldestRecoveryPoint == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.OldestRecoveryPoint, "2006-01-02T15:04:05Z07:00")
}

func (o *AzureVMWorkloadProtectedItemExtendedInfo) SetOldestRecoveryPointAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.OldestRecoveryPoint = &formatted
}

func (o *AzureVMWorkloadProtectedItemExtendedInfo) GetOldestRecoveryPointInArchiveAsTime() (*time.Time, error) {
	if o.OldestRecoveryPointInArchive == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.OldestRecoveryPointInArchive, "2006-01-02T15:04:05Z07:00")
}

func (o *AzureVMWorkloadProtectedItemExtendedInfo) SetOldestRecoveryPointInArchiveAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.OldestRecoveryPointInArchive = &formatted
}

func (o *AzureVMWorkloadProtectedItemExtendedInfo) GetOldestRecoveryPointInVaultAsTime() (*time.Time, error) {
	if o.OldestRecoveryPointInVault == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.OldestRecoveryPointInVault, "2006-01-02T15:04:05Z07:00")
}

func (o *AzureVMWorkloadProtectedItemExtendedInfo) SetOldestRecoveryPointInVaultAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.OldestRecoveryPointInVault = &formatted
}
