package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatasetLocation = GoogleCloudStorageLocation{}

type GoogleCloudStorageLocation struct {
	BucketName *interface{} `json:"bucketName,omitempty"`
	Version    *interface{} `json:"version,omitempty"`

	// Fields inherited from DatasetLocation

	FileName   *interface{} `json:"fileName,omitempty"`
	FolderPath *interface{} `json:"folderPath,omitempty"`
	Type       string       `json:"type"`
}

func (s GoogleCloudStorageLocation) DatasetLocation() BaseDatasetLocationImpl {
	return BaseDatasetLocationImpl{
		FileName:   s.FileName,
		FolderPath: s.FolderPath,
		Type:       s.Type,
	}
}

var _ json.Marshaler = GoogleCloudStorageLocation{}

func (s GoogleCloudStorageLocation) MarshalJSON() ([]byte, error) {
	type wrapper GoogleCloudStorageLocation
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling GoogleCloudStorageLocation: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling GoogleCloudStorageLocation: %+v", err)
	}

	decoded["type"] = "GoogleCloudStorageLocation"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling GoogleCloudStorageLocation: %+v", err)
	}

	return encoded, nil
}
