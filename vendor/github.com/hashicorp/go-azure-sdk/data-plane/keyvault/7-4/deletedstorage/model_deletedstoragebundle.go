package deletedstorage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedStorageBundle struct {
	ActiveKeyName      *string                   `json:"activeKeyName,omitempty"`
	Attributes         *StorageAccountAttributes `json:"attributes,omitempty"`
	AutoRegenerateKey  *bool                     `json:"autoRegenerateKey,omitempty"`
	DeletedDate        *int64                    `json:"deletedDate,omitempty"`
	Id                 *string                   `json:"id,omitempty"`
	RecoveryId         *string                   `json:"recoveryId,omitempty"`
	RegenerationPeriod *string                   `json:"regenerationPeriod,omitempty"`
	ResourceId         *string                   `json:"resourceId,omitempty"`
	ScheduledPurgeDate *int64                    `json:"scheduledPurgeDate,omitempty"`
	Tags               *map[string]string        `json:"tags,omitempty"`
}
