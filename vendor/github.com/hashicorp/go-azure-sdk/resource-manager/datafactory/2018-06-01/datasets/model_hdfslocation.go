package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatasetLocation = HdfsLocation{}

type HdfsLocation struct {

	// Fields inherited from DatasetLocation

	FileName   *interface{} `json:"fileName,omitempty"`
	FolderPath *interface{} `json:"folderPath,omitempty"`
	Type       string       `json:"type"`
}

func (s HdfsLocation) DatasetLocation() BaseDatasetLocationImpl {
	return BaseDatasetLocationImpl{
		FileName:   s.FileName,
		FolderPath: s.FolderPath,
		Type:       s.Type,
	}
}

var _ json.Marshaler = HdfsLocation{}

func (s HdfsLocation) MarshalJSON() ([]byte, error) {
	type wrapper HdfsLocation
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HdfsLocation: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HdfsLocation: %+v", err)
	}

	decoded["type"] = "HdfsLocation"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HdfsLocation: %+v", err)
	}

	return encoded, nil
}
