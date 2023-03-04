package dataconnectors

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataConnector = OfficeDataConnector{}

type OfficeDataConnector struct {
	Properties *OfficeDataConnectorProperties `json:"properties,omitempty"`

	// Fields inherited from DataConnector
	Etag       *string                `json:"etag,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

var _ json.Marshaler = OfficeDataConnector{}

func (s OfficeDataConnector) MarshalJSON() ([]byte, error) {
	type wrapper OfficeDataConnector
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling OfficeDataConnector: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling OfficeDataConnector: %+v", err)
	}
	decoded["kind"] = "Office365"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling OfficeDataConnector: %+v", err)
	}

	return encoded, nil
}
