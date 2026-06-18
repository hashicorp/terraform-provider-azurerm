package integrationruntime

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SsisObjectMetadata = SsisProject{}

type SsisProject struct {
	EnvironmentRefs *[]SsisEnvironmentReference `json:"environmentRefs,omitempty"`
	FolderId        *int64                      `json:"folderId,omitempty"`
	Parameters      *[]SsisParameter            `json:"parameters,omitempty"`
	Version         *int64                      `json:"version,omitempty"`

	// Fields inherited from SsisObjectMetadata

	Description *string                `json:"description,omitempty"`
	Id          *int64                 `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Type        SsisObjectMetadataType `json:"type"`
}

func (s SsisProject) SsisObjectMetadata() BaseSsisObjectMetadataImpl {
	return BaseSsisObjectMetadataImpl{
		Description: s.Description,
		Id:          s.Id,
		Name:        s.Name,
		Type:        s.Type,
	}
}

var _ json.Marshaler = SsisProject{}

func (s SsisProject) MarshalJSON() ([]byte, error) {
	type wrapper SsisProject
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SsisProject: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SsisProject: %+v", err)
	}

	decoded["type"] = "Project"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SsisProject: %+v", err)
	}

	return encoded, nil
}
