package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EndpointBaseUpdateProperties = NfsMountEndpointUpdateProperties{}

type NfsMountEndpointUpdateProperties struct {

	// Fields inherited from EndpointBaseUpdateProperties

	Description  *string      `json:"description,omitempty"`
	EndpointType EndpointType `json:"endpointType"`
}

func (s NfsMountEndpointUpdateProperties) EndpointBaseUpdateProperties() BaseEndpointBaseUpdatePropertiesImpl {
	return BaseEndpointBaseUpdatePropertiesImpl{
		Description:  s.Description,
		EndpointType: s.EndpointType,
	}
}

var _ json.Marshaler = NfsMountEndpointUpdateProperties{}

func (s NfsMountEndpointUpdateProperties) MarshalJSON() ([]byte, error) {
	type wrapper NfsMountEndpointUpdateProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NfsMountEndpointUpdateProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NfsMountEndpointUpdateProperties: %+v", err)
	}

	decoded["endpointType"] = "NfsMount"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NfsMountEndpointUpdateProperties: %+v", err)
	}

	return encoded, nil
}
