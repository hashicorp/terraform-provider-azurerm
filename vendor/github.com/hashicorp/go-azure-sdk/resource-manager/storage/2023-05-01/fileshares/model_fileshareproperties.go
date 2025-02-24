package fileshares

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileShareProperties struct {
	AccessTier             *ShareAccessTier    `json:"accessTier,omitempty"`
	AccessTierChangeTime   *string             `json:"accessTierChangeTime,omitempty"`
	AccessTierStatus       *string             `json:"accessTierStatus,omitempty"`
	Deleted                *bool               `json:"deleted,omitempty"`
	DeletedTime            *string             `json:"deletedTime,omitempty"`
	EnabledProtocols       *EnabledProtocols   `json:"enabledProtocols,omitempty"`
	LastModifiedTime       *string             `json:"lastModifiedTime,omitempty"`
	LeaseDuration          *LeaseDuration      `json:"leaseDuration,omitempty"`
	LeaseState             *LeaseState         `json:"leaseState,omitempty"`
	LeaseStatus            *LeaseStatus        `json:"leaseStatus,omitempty"`
	Metadata               *map[string]string  `json:"metadata,omitempty"`
	RemainingRetentionDays *int64              `json:"remainingRetentionDays,omitempty"`
	RootSquash             *RootSquashType     `json:"rootSquash,omitempty"`
	ShareQuota             *int64              `json:"shareQuota,omitempty"`
	ShareUsageBytes        *int64              `json:"shareUsageBytes,omitempty"`
	SignedIdentifiers      *[]SignedIdentifier `json:"signedIdentifiers,omitempty"`
	SnapshotTime           *string             `json:"snapshotTime,omitempty"`
	Version                *string             `json:"version,omitempty"`
}

func (o *FileShareProperties) GetAccessTierChangeTimeAsTime() (*time.Time, error) {
	if o.AccessTierChangeTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.AccessTierChangeTime, "2006-01-02T15:04:05Z07:00")
}

func (o *FileShareProperties) SetAccessTierChangeTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.AccessTierChangeTime = &formatted
}

func (o *FileShareProperties) GetDeletedTimeAsTime() (*time.Time, error) {
	if o.DeletedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DeletedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *FileShareProperties) SetDeletedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DeletedTime = &formatted
}

func (o *FileShareProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *FileShareProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}

func (o *FileShareProperties) GetSnapshotTimeAsTime() (*time.Time, error) {
	if o.SnapshotTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SnapshotTime, "2006-01-02T15:04:05Z07:00")
}

func (o *FileShareProperties) SetSnapshotTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SnapshotTime = &formatted
}
