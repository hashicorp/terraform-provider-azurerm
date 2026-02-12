package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatasetLocation = OracleCloudStorageLocation{}

type OracleCloudStorageLocation struct {
	BucketName *interface{} `json:"bucketName,omitempty"`
	Version    *interface{} `json:"version,omitempty"`

	// Fields inherited from DatasetLocation

	FileName   *interface{} `json:"fileName,omitempty"`
	FolderPath *interface{} `json:"folderPath,omitempty"`
	Type       string       `json:"type"`
}

func (s OracleCloudStorageLocation) DatasetLocation() BaseDatasetLocationImpl {
	return BaseDatasetLocationImpl{
		FileName:   s.FileName,
		FolderPath: s.FolderPath,
		Type:       s.Type,
	}
}

var _ json.Marshaler = OracleCloudStorageLocation{}

func (s OracleCloudStorageLocation) MarshalJSON() ([]byte, error) {
	type wrapper OracleCloudStorageLocation
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling OracleCloudStorageLocation: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling OracleCloudStorageLocation: %+v", err)
	}

	decoded["type"] = "OracleCloudStorageLocation"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling OracleCloudStorageLocation: %+v", err)
	}

	return encoded, nil
}
