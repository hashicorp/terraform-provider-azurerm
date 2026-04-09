package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EndpointBaseProperties = NfsMountEndpointProperties{}

type NfsMountEndpointProperties struct {
	Export     string      `json:"export"`
	Host       string      `json:"host"`
	NfsVersion *NfsVersion `json:"nfsVersion,omitempty"`

	// Fields inherited from EndpointBaseProperties

	Description       *string            `json:"description,omitempty"`
	EndpointType      EndpointType       `json:"endpointType"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}

func (s NfsMountEndpointProperties) EndpointBaseProperties() BaseEndpointBasePropertiesImpl {
	return BaseEndpointBasePropertiesImpl{
		Description:       s.Description,
		EndpointType:      s.EndpointType,
		ProvisioningState: s.ProvisioningState,
	}
}

var _ json.Marshaler = NfsMountEndpointProperties{}

func (s NfsMountEndpointProperties) MarshalJSON() ([]byte, error) {
	type wrapper NfsMountEndpointProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NfsMountEndpointProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NfsMountEndpointProperties: %+v", err)
	}

	decoded["endpointType"] = "NfsMount"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NfsMountEndpointProperties: %+v", err)
	}

	return encoded, nil
}
