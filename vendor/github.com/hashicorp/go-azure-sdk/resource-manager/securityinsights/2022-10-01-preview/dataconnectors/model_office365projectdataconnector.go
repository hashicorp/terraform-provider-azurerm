package dataconnectors

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataConnector = Office365ProjectDataConnector{}

type Office365ProjectDataConnector struct {
	Properties *Office365ProjectDataConnectorProperties `json:"properties,omitempty"`

	// Fields inherited from DataConnector

	Etag       *string                `json:"etag,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Kind       DataConnectorKind      `json:"kind"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

func (s Office365ProjectDataConnector) DataConnector() BaseDataConnectorImpl {
	return BaseDataConnectorImpl{
		Etag:       s.Etag,
		Id:         s.Id,
		Kind:       s.Kind,
		Name:       s.Name,
		SystemData: s.SystemData,
		Type:       s.Type,
	}
}

var _ json.Marshaler = Office365ProjectDataConnector{}

func (s Office365ProjectDataConnector) MarshalJSON() ([]byte, error) {
	type wrapper Office365ProjectDataConnector
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Office365ProjectDataConnector: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Office365ProjectDataConnector: %+v", err)
	}

	decoded["kind"] = "Office365Project"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Office365ProjectDataConnector: %+v", err)
	}

	return encoded, nil
}
