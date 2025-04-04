package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatasetLocation = LakeHouseLocation{}

type LakeHouseLocation struct {

	// Fields inherited from DatasetLocation

	FileName   *string `json:"fileName,omitempty"`
	FolderPath *string `json:"folderPath,omitempty"`
	Type       string  `json:"type"`
}

func (s LakeHouseLocation) DatasetLocation() BaseDatasetLocationImpl {
	return BaseDatasetLocationImpl{
		FileName:   s.FileName,
		FolderPath: s.FolderPath,
		Type:       s.Type,
	}
}

var _ json.Marshaler = LakeHouseLocation{}

func (s LakeHouseLocation) MarshalJSON() ([]byte, error) {
	type wrapper LakeHouseLocation
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LakeHouseLocation: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LakeHouseLocation: %+v", err)
	}

	decoded["type"] = "LakeHouseLocation"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LakeHouseLocation: %+v", err)
	}

	return encoded, nil
}
