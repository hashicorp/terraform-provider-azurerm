package dataconnectors

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataConnector = CodelessApiPollingDataConnector{}

type CodelessApiPollingDataConnector struct {
	Properties *ApiPollingParameters `json:"properties,omitempty"`

	// Fields inherited from DataConnector
	Etag       *string                `json:"etag,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

var _ json.Marshaler = CodelessApiPollingDataConnector{}

func (s CodelessApiPollingDataConnector) MarshalJSON() ([]byte, error) {
	type wrapper CodelessApiPollingDataConnector
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CodelessApiPollingDataConnector: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CodelessApiPollingDataConnector: %+v", err)
	}
	decoded["kind"] = "APIPolling"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CodelessApiPollingDataConnector: %+v", err)
	}

	return encoded, nil
}
