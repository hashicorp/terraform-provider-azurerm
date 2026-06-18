package integrationruntime

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SsisObjectMetadata = SsisFolder{}

type SsisFolder struct {

	// Fields inherited from SsisObjectMetadata

	Description *string                `json:"description,omitempty"`
	Id          *int64                 `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Type        SsisObjectMetadataType `json:"type"`
}

func (s SsisFolder) SsisObjectMetadata() BaseSsisObjectMetadataImpl {
	return BaseSsisObjectMetadataImpl{
		Description: s.Description,
		Id:          s.Id,
		Name:        s.Name,
		Type:        s.Type,
	}
}

var _ json.Marshaler = SsisFolder{}

func (s SsisFolder) MarshalJSON() ([]byte, error) {
	type wrapper SsisFolder
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SsisFolder: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SsisFolder: %+v", err)
	}

	decoded["type"] = "Folder"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SsisFolder: %+v", err)
	}

	return encoded, nil
}
