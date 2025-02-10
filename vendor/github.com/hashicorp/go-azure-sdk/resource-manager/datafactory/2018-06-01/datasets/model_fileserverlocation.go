package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatasetLocation = FileServerLocation{}

type FileServerLocation struct {

	// Fields inherited from DatasetLocation

	FileName   *string `json:"fileName,omitempty"`
	FolderPath *string `json:"folderPath,omitempty"`
	Type       string  `json:"type"`
}

func (s FileServerLocation) DatasetLocation() BaseDatasetLocationImpl {
	return BaseDatasetLocationImpl{
		FileName:   s.FileName,
		FolderPath: s.FolderPath,
		Type:       s.Type,
	}
}

var _ json.Marshaler = FileServerLocation{}

func (s FileServerLocation) MarshalJSON() ([]byte, error) {
	type wrapper FileServerLocation
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FileServerLocation: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FileServerLocation: %+v", err)
	}

	decoded["type"] = "FileServerLocation"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FileServerLocation: %+v", err)
	}

	return encoded, nil
}
