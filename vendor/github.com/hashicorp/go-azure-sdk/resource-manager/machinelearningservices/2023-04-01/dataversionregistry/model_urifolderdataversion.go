package dataversionregistry

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataVersionBase = UriFolderDataVersion{}

type UriFolderDataVersion struct {

	// Fields inherited from DataVersionBase
	DataUri     string             `json:"dataUri"`
	Description *string            `json:"description,omitempty"`
	IsAnonymous *bool              `json:"isAnonymous,omitempty"`
	IsArchived  *bool              `json:"isArchived,omitempty"`
	Properties  *map[string]string `json:"properties,omitempty"`
	Tags        *map[string]string `json:"tags,omitempty"`
}

var _ json.Marshaler = UriFolderDataVersion{}

func (s UriFolderDataVersion) MarshalJSON() ([]byte, error) {
	type wrapper UriFolderDataVersion
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling UriFolderDataVersion: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling UriFolderDataVersion: %+v", err)
	}
	decoded["dataType"] = "uri_folder"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling UriFolderDataVersion: %+v", err)
	}

	return encoded, nil
}
