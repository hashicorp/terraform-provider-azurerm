package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatasetLocation = FtpServerLocation{}

type FtpServerLocation struct {

	// Fields inherited from DatasetLocation

	FileName   *interface{} `json:"fileName,omitempty"`
	FolderPath *interface{} `json:"folderPath,omitempty"`
	Type       string       `json:"type"`
}

func (s FtpServerLocation) DatasetLocation() BaseDatasetLocationImpl {
	return BaseDatasetLocationImpl{
		FileName:   s.FileName,
		FolderPath: s.FolderPath,
		Type:       s.Type,
	}
}

var _ json.Marshaler = FtpServerLocation{}

func (s FtpServerLocation) MarshalJSON() ([]byte, error) {
	type wrapper FtpServerLocation
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FtpServerLocation: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FtpServerLocation: %+v", err)
	}

	decoded["type"] = "FtpServerLocation"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FtpServerLocation: %+v", err)
	}

	return encoded, nil
}
