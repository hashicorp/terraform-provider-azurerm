package dataconnectors

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataConnector = OfficeIRMDataConnector{}

type OfficeIRMDataConnector struct {
	Properties *OfficeIRMDataConnectorProperties `json:"properties,omitempty"`

	// Fields inherited from DataConnector
	Etag       *string                `json:"etag,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

var _ json.Marshaler = OfficeIRMDataConnector{}

func (s OfficeIRMDataConnector) MarshalJSON() ([]byte, error) {
	type wrapper OfficeIRMDataConnector
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling OfficeIRMDataConnector: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling OfficeIRMDataConnector: %+v", err)
	}
	decoded["kind"] = "OfficeIRM"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling OfficeIRMDataConnector: %+v", err)
	}

	return encoded, nil
}
