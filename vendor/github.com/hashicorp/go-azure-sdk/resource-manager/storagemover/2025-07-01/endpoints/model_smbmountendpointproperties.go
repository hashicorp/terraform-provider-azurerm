package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EndpointBaseProperties = SmbMountEndpointProperties{}

type SmbMountEndpointProperties struct {
	Credentials *AzureKeyVaultSmbCredentials `json:"credentials,omitempty"`
	Host        string                       `json:"host"`
	ShareName   string                       `json:"shareName"`

	// Fields inherited from EndpointBaseProperties

	Description       *string            `json:"description,omitempty"`
	EndpointType      EndpointType       `json:"endpointType"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}

func (s SmbMountEndpointProperties) EndpointBaseProperties() BaseEndpointBasePropertiesImpl {
	return BaseEndpointBasePropertiesImpl{
		Description:       s.Description,
		EndpointType:      s.EndpointType,
		ProvisioningState: s.ProvisioningState,
	}
}

var _ json.Marshaler = SmbMountEndpointProperties{}

func (s SmbMountEndpointProperties) MarshalJSON() ([]byte, error) {
	type wrapper SmbMountEndpointProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SmbMountEndpointProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SmbMountEndpointProperties: %+v", err)
	}

	decoded["endpointType"] = "SmbMount"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SmbMountEndpointProperties: %+v", err)
	}

	return encoded, nil
}
