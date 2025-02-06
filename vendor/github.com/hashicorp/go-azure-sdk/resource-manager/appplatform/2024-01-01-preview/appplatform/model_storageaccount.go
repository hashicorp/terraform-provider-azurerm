package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StorageProperties = StorageAccount{}

type StorageAccount struct {
	AccountKey  string `json:"accountKey"`
	AccountName string `json:"accountName"`

	// Fields inherited from StorageProperties

	StorageType StorageType `json:"storageType"`
}

func (s StorageAccount) StorageProperties() BaseStoragePropertiesImpl {
	return BaseStoragePropertiesImpl{
		StorageType: s.StorageType,
	}
}

var _ json.Marshaler = StorageAccount{}

func (s StorageAccount) MarshalJSON() ([]byte, error) {
	type wrapper StorageAccount
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StorageAccount: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StorageAccount: %+v", err)
	}

	decoded["storageType"] = "StorageAccount"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StorageAccount: %+v", err)
	}

	return encoded, nil
}
