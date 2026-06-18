package integrationruntime

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SsisObjectMetadata = SsisPackage{}

type SsisPackage struct {
	FolderId       *int64           `json:"folderId,omitempty"`
	Parameters     *[]SsisParameter `json:"parameters,omitempty"`
	ProjectId      *int64           `json:"projectId,omitempty"`
	ProjectVersion *int64           `json:"projectVersion,omitempty"`

	// Fields inherited from SsisObjectMetadata

	Description *string                `json:"description,omitempty"`
	Id          *int64                 `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Type        SsisObjectMetadataType `json:"type"`
}

func (s SsisPackage) SsisObjectMetadata() BaseSsisObjectMetadataImpl {
	return BaseSsisObjectMetadataImpl{
		Description: s.Description,
		Id:          s.Id,
		Name:        s.Name,
		Type:        s.Type,
	}
}

var _ json.Marshaler = SsisPackage{}

func (s SsisPackage) MarshalJSON() ([]byte, error) {
	type wrapper SsisPackage
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SsisPackage: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SsisPackage: %+v", err)
	}

	decoded["type"] = "Package"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SsisPackage: %+v", err)
	}

	return encoded, nil
}
