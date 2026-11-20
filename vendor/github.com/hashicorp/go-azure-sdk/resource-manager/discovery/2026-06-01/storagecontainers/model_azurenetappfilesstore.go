package storagecontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StorageStore = AzureNetAppFilesStore{}

type AzureNetAppFilesStore struct {
	MountProtocol  *NetAppMountProtocol `json:"mountProtocol,omitempty"`
	NetAppVolumeId string               `json:"netAppVolumeId"`

	// Fields inherited from StorageStore

	Kind StorageStoreType `json:"kind"`
}

func (s AzureNetAppFilesStore) StorageStore() BaseStorageStoreImpl {
	return BaseStorageStoreImpl{
		Kind: s.Kind,
	}
}

var _ json.Marshaler = AzureNetAppFilesStore{}

func (s AzureNetAppFilesStore) MarshalJSON() ([]byte, error) {
	type wrapper AzureNetAppFilesStore
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureNetAppFilesStore: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureNetAppFilesStore: %+v", err)
	}

	decoded["kind"] = "AzureNetAppFiles"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureNetAppFilesStore: %+v", err)
	}

	return encoded, nil
}
