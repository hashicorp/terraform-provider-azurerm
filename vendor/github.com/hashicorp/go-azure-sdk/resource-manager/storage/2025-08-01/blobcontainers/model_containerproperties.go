package blobcontainers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerProperties struct {
	DefaultEncryptionScope         *string                         `json:"defaultEncryptionScope,omitempty"`
	Deleted                        *bool                           `json:"deleted,omitempty"`
	DeletedTime                    *string                         `json:"deletedTime,omitempty"`
	DenyEncryptionScopeOverride    *bool                           `json:"denyEncryptionScopeOverride,omitempty"`
	EnableNfsV3AllSquash           *bool                           `json:"enableNfsV3AllSquash,omitempty"`
	EnableNfsV3RootSquash          *bool                           `json:"enableNfsV3RootSquash,omitempty"`
	HasImmutabilityPolicy          *bool                           `json:"hasImmutabilityPolicy,omitempty"`
	HasLegalHold                   *bool                           `json:"hasLegalHold,omitempty"`
	ImmutabilityPolicy             *ImmutabilityPolicyProperties   `json:"immutabilityPolicy,omitempty"`
	ImmutableStorageWithVersioning *ImmutableStorageWithVersioning `json:"immutableStorageWithVersioning,omitempty"`
	LastModifiedTime               *string                         `json:"lastModifiedTime,omitempty"`
	LeaseDuration                  *LeaseDuration                  `json:"leaseDuration,omitempty"`
	LeaseState                     *LeaseState                     `json:"leaseState,omitempty"`
	LeaseStatus                    *LeaseStatus                    `json:"leaseStatus,omitempty"`
	LegalHold                      *LegalHoldProperties            `json:"legalHold,omitempty"`
	Metadata                       *map[string]string              `json:"metadata,omitempty"`
	PublicAccess                   *PublicAccess                   `json:"publicAccess,omitempty"`
	RemainingRetentionDays         *int64                          `json:"remainingRetentionDays,omitempty"`
	Version                        *string                         `json:"version,omitempty"`
}

func (o *ContainerProperties) GetDeletedTimeAsTime() (*time.Time, error) {
	if o.DeletedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DeletedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContainerProperties) SetDeletedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DeletedTime = &formatted
}

func (o *ContainerProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContainerProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
