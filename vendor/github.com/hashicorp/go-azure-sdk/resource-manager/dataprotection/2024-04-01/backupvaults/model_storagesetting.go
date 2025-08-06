package backupvaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageSetting struct {
	DatastoreType *StorageSettingStoreTypes `json:"datastoreType,omitempty"`
	Type          *StorageSettingTypes      `json:"type,omitempty"`
}
