package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatasetLocation = AmazonS3CompatibleLocation{}

type AmazonS3CompatibleLocation struct {
	BucketName *interface{} `json:"bucketName,omitempty"`
	Version    *interface{} `json:"version,omitempty"`

	// Fields inherited from DatasetLocation

	FileName   *interface{} `json:"fileName,omitempty"`
	FolderPath *interface{} `json:"folderPath,omitempty"`
	Type       string       `json:"type"`
}

func (s AmazonS3CompatibleLocation) DatasetLocation() BaseDatasetLocationImpl {
	return BaseDatasetLocationImpl{
		FileName:   s.FileName,
		FolderPath: s.FolderPath,
		Type:       s.Type,
	}
}

var _ json.Marshaler = AmazonS3CompatibleLocation{}

func (s AmazonS3CompatibleLocation) MarshalJSON() ([]byte, error) {
	type wrapper AmazonS3CompatibleLocation
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AmazonS3CompatibleLocation: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AmazonS3CompatibleLocation: %+v", err)
	}

	decoded["type"] = "AmazonS3CompatibleLocation"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AmazonS3CompatibleLocation: %+v", err)
	}

	return encoded, nil
}
