package protecteditems

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectedItem = AzureFileshareProtectedItem{}

type AzureFileshareProtectedItem struct {
	ExtendedInfo     *AzureFileshareProtectedItemExtendedInfo `json:"extendedInfo,omitempty"`
	FriendlyName     *string                                  `json:"friendlyName,omitempty"`
	KpisHealths      *map[string]KPIResourceHealthDetails     `json:"kpisHealths,omitempty"`
	LastBackupStatus *string                                  `json:"lastBackupStatus,omitempty"`
	LastBackupTime   *string                                  `json:"lastBackupTime,omitempty"`
	ProtectionState  *ProtectionState                         `json:"protectionState,omitempty"`
	ProtectionStatus *string                                  `json:"protectionStatus,omitempty"`

	// Fields inherited from ProtectedItem

	BackupManagementType             *BackupManagementType `json:"backupManagementType,omitempty"`
	BackupSetName                    *string               `json:"backupSetName,omitempty"`
	ContainerName                    *string               `json:"containerName,omitempty"`
	CreateMode                       *CreateMode           `json:"createMode,omitempty"`
	DeferredDeleteTimeInUTC          *string               `json:"deferredDeleteTimeInUTC,omitempty"`
	DeferredDeleteTimeRemaining      *string               `json:"deferredDeleteTimeRemaining,omitempty"`
	IsArchiveEnabled                 *bool                 `json:"isArchiveEnabled,omitempty"`
	IsDeferredDeleteScheduleUpcoming *bool                 `json:"isDeferredDeleteScheduleUpcoming,omitempty"`
	IsRehydrate                      *bool                 `json:"isRehydrate,omitempty"`
	IsScheduledForDeferredDelete     *bool                 `json:"isScheduledForDeferredDelete,omitempty"`
	LastRecoveryPoint                *string               `json:"lastRecoveryPoint,omitempty"`
	PolicyId                         *string               `json:"policyId,omitempty"`
	PolicyName                       *string               `json:"policyName,omitempty"`
	ProtectedItemType                string                `json:"protectedItemType"`
	ResourceGuardOperationRequests   *[]string             `json:"resourceGuardOperationRequests,omitempty"`
	SoftDeleteRetentionPeriodInDays  *int64                `json:"softDeleteRetentionPeriodInDays,omitempty"`
	SourceResourceId                 *string               `json:"sourceResourceId,omitempty"`
	WorkloadType                     *DataSourceType       `json:"workloadType,omitempty"`
}

func (s AzureFileshareProtectedItem) ProtectedItem() BaseProtectedItemImpl {
	return BaseProtectedItemImpl{
		BackupManagementType:             s.BackupManagementType,
		BackupSetName:                    s.BackupSetName,
		ContainerName:                    s.ContainerName,
		CreateMode:                       s.CreateMode,
		DeferredDeleteTimeInUTC:          s.DeferredDeleteTimeInUTC,
		DeferredDeleteTimeRemaining:      s.DeferredDeleteTimeRemaining,
		IsArchiveEnabled:                 s.IsArchiveEnabled,
		IsDeferredDeleteScheduleUpcoming: s.IsDeferredDeleteScheduleUpcoming,
		IsRehydrate:                      s.IsRehydrate,
		IsScheduledForDeferredDelete:     s.IsScheduledForDeferredDelete,
		LastRecoveryPoint:                s.LastRecoveryPoint,
		PolicyId:                         s.PolicyId,
		PolicyName:                       s.PolicyName,
		ProtectedItemType:                s.ProtectedItemType,
		ResourceGuardOperationRequests:   s.ResourceGuardOperationRequests,
		SoftDeleteRetentionPeriodInDays:  s.SoftDeleteRetentionPeriodInDays,
		SourceResourceId:                 s.SourceResourceId,
		WorkloadType:                     s.WorkloadType,
	}
}

func (o *AzureFileshareProtectedItem) GetDeferredDeleteTimeInUTCAsTime() (*time.Time, error) {
	if o.DeferredDeleteTimeInUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DeferredDeleteTimeInUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *AzureFileshareProtectedItem) SetDeferredDeleteTimeInUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DeferredDeleteTimeInUTC = &formatted
}

func (o *AzureFileshareProtectedItem) GetLastRecoveryPointAsTime() (*time.Time, error) {
	if o.LastRecoveryPoint == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastRecoveryPoint, "2006-01-02T15:04:05Z07:00")
}

func (o *AzureFileshareProtectedItem) SetLastRecoveryPointAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastRecoveryPoint = &formatted
}

var _ json.Marshaler = AzureFileshareProtectedItem{}

func (s AzureFileshareProtectedItem) MarshalJSON() ([]byte, error) {
	type wrapper AzureFileshareProtectedItem
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureFileshareProtectedItem: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureFileshareProtectedItem: %+v", err)
	}

	decoded["protectedItemType"] = "AzureFileShareProtectedItem"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureFileshareProtectedItem: %+v", err)
	}

	return encoded, nil
}
