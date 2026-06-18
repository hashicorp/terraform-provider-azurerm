package integrationruntime

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SsisObjectMetadata = SsisEnvironment{}

type SsisEnvironment struct {
	FolderId  *int64          `json:"folderId,omitempty"`
	Variables *[]SsisVariable `json:"variables,omitempty"`

	// Fields inherited from SsisObjectMetadata

	Description *string                `json:"description,omitempty"`
	Id          *int64                 `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Type        SsisObjectMetadataType `json:"type"`
}

func (s SsisEnvironment) SsisObjectMetadata() BaseSsisObjectMetadataImpl {
	return BaseSsisObjectMetadataImpl{
		Description: s.Description,
		Id:          s.Id,
		Name:        s.Name,
		Type:        s.Type,
	}
}

var _ json.Marshaler = SsisEnvironment{}

func (s SsisEnvironment) MarshalJSON() ([]byte, error) {
	type wrapper SsisEnvironment
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SsisEnvironment: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SsisEnvironment: %+v", err)
	}

	decoded["type"] = "Environment"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SsisEnvironment: %+v", err)
	}

	return encoded, nil
}
