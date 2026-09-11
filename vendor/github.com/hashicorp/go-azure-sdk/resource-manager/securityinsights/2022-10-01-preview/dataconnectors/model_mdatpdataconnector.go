package dataconnectors

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataConnector = MDATPDataConnector{}

type MDATPDataConnector struct {
	Properties *MDATPDataConnectorProperties `json:"properties,omitempty"`

	// Fields inherited from DataConnector

	Etag       *string                `json:"etag,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Kind       DataConnectorKind      `json:"kind"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

func (s MDATPDataConnector) DataConnector() BaseDataConnectorImpl {
	return BaseDataConnectorImpl{
		Etag:       s.Etag,
		Id:         s.Id,
		Kind:       s.Kind,
		Name:       s.Name,
		SystemData: s.SystemData,
		Type:       s.Type,
	}
}

var _ json.Marshaler = MDATPDataConnector{}

func (s MDATPDataConnector) MarshalJSON() ([]byte, error) {
	type wrapper MDATPDataConnector
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MDATPDataConnector: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MDATPDataConnector: %+v", err)
	}

	decoded["kind"] = "MicrosoftDefenderAdvancedThreatProtection"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MDATPDataConnector: %+v", err)
	}

	return encoded, nil
}
