package deletedstorage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageBundle struct {
	ActiveKeyName      *string                   `json:"activeKeyName,omitempty"`
	Attributes         *StorageAccountAttributes `json:"attributes,omitempty"`
	AutoRegenerateKey  *bool                     `json:"autoRegenerateKey,omitempty"`
	Id                 *string                   `json:"id,omitempty"`
	RegenerationPeriod *string                   `json:"regenerationPeriod,omitempty"`
	ResourceId         *string                   `json:"resourceId,omitempty"`
	Tags               *map[string]string        `json:"tags,omitempty"`
}
