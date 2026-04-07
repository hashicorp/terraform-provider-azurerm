package storagediscoveryworkspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageDiscoveryScope struct {
	DisplayName   string                         `json:"displayName"`
	ResourceTypes []StorageDiscoveryResourceType `json:"resourceTypes"`
	TagKeysOnly   *[]string                      `json:"tagKeysOnly,omitempty"`
	Tags          *map[string]string             `json:"tags,omitempty"`
}
