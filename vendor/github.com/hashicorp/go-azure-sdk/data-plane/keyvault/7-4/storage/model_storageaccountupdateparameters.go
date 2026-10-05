package storage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountUpdateParameters struct {
	ActiveKeyName      *string                   `json:"activeKeyName,omitempty"`
	Attributes         *StorageAccountAttributes `json:"attributes,omitempty"`
	AutoRegenerateKey  *bool                     `json:"autoRegenerateKey,omitempty"`
	RegenerationPeriod *string                   `json:"regenerationPeriod,omitempty"`
	Tags               *map[string]string        `json:"tags,omitempty"`
}
