package dataversion

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataVersionBase = UriFileDataVersion{}

type UriFileDataVersion struct {

	// Fields inherited from DataVersionBase
	DataUri     string             `json:"dataUri"`
	Description *string            `json:"description,omitempty"`
	IsAnonymous *bool              `json:"isAnonymous,omitempty"`
	IsArchived  *bool              `json:"isArchived,omitempty"`
	Properties  *map[string]string `json:"properties,omitempty"`
	Tags        *map[string]string `json:"tags,omitempty"`
}

var _ json.Marshaler = UriFileDataVersion{}

func (s UriFileDataVersion) MarshalJSON() ([]byte, error) {
	type wrapper UriFileDataVersion
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling UriFileDataVersion: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling UriFileDataVersion: %+v", err)
	}
	decoded["dataType"] = "uri_file"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling UriFileDataVersion: %+v", err)
	}

	return encoded, nil
}
