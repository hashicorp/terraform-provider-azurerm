package backupprotecteditems

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectedItem = GenericProtectedItem{}

type GenericProtectedItem struct {
	FabricName         *string            `json:"fabricName,omitempty"`
	FriendlyName       *string            `json:"friendlyName,omitempty"`
	PolicyState        *string            `json:"policyState,omitempty"`
	ProtectedItemId    *int64             `json:"protectedItemId,omitempty"`
	ProtectionState    *ProtectionState   `json:"protectionState,omitempty"`
	SourceAssociations *map[string]string `json:"sourceAssociations,omitempty"`

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

func (s GenericProtectedItem) ProtectedItem() BaseProtectedItemImpl {
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

func (o *GenericProtectedItem) GetDeferredDeleteTimeInUTCAsTime() (*time.Time, error) {
	if o.DeferredDeleteTimeInUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DeferredDeleteTimeInUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *GenericProtectedItem) SetDeferredDeleteTimeInUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DeferredDeleteTimeInUTC = &formatted
}

func (o *GenericProtectedItem) GetLastRecoveryPointAsTime() (*time.Time, error) {
	if o.LastRecoveryPoint == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastRecoveryPoint, "2006-01-02T15:04:05Z07:00")
}

func (o *GenericProtectedItem) SetLastRecoveryPointAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastRecoveryPoint = &formatted
}

var _ json.Marshaler = GenericProtectedItem{}

func (s GenericProtectedItem) MarshalJSON() ([]byte, error) {
	type wrapper GenericProtectedItem
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling GenericProtectedItem: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling GenericProtectedItem: %+v", err)
	}

	decoded["protectedItemType"] = "GenericProtectedItem"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling GenericProtectedItem: %+v", err)
	}

	return encoded, nil
}
