package storagecontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StorageStore = AzureStorageBlobStore{}

type AzureStorageBlobStore struct {
	MountProtocol    *BlobStorageMountProtocol `json:"mountProtocol,omitempty"`
	StorageAccountId string                    `json:"storageAccountId"`

	// Fields inherited from StorageStore

	Kind StorageStoreType `json:"kind"`
}

func (s AzureStorageBlobStore) StorageStore() BaseStorageStoreImpl {
	return BaseStorageStoreImpl{
		Kind: s.Kind,
	}
}

var _ json.Marshaler = AzureStorageBlobStore{}

func (s AzureStorageBlobStore) MarshalJSON() ([]byte, error) {
	type wrapper AzureStorageBlobStore
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureStorageBlobStore: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureStorageBlobStore: %+v", err)
	}

	decoded["kind"] = "AzureStorageBlob"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureStorageBlobStore: %+v", err)
	}

	return encoded, nil
}
