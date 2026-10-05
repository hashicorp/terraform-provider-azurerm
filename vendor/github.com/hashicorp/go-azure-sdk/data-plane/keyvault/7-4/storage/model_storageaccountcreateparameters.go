package storage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountCreateParameters struct {
	ActiveKeyName      string                    `json:"activeKeyName"`
	Attributes         *StorageAccountAttributes `json:"attributes,omitempty"`
	AutoRegenerateKey  bool                      `json:"autoRegenerateKey"`
	RegenerationPeriod *string                   `json:"regenerationPeriod,omitempty"`
	ResourceId         string                    `json:"resourceId"`
	Tags               *map[string]string        `json:"tags,omitempty"`
}
