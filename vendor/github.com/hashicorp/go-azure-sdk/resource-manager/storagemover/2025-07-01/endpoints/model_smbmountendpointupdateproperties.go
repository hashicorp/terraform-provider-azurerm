package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EndpointBaseUpdateProperties = SmbMountEndpointUpdateProperties{}

type SmbMountEndpointUpdateProperties struct {
	Credentials *AzureKeyVaultSmbCredentials `json:"credentials,omitempty"`

	// Fields inherited from EndpointBaseUpdateProperties

	Description  *string      `json:"description,omitempty"`
	EndpointType EndpointType `json:"endpointType"`
}

func (s SmbMountEndpointUpdateProperties) EndpointBaseUpdateProperties() BaseEndpointBaseUpdatePropertiesImpl {
	return BaseEndpointBaseUpdatePropertiesImpl{
		Description:  s.Description,
		EndpointType: s.EndpointType,
	}
}

var _ json.Marshaler = SmbMountEndpointUpdateProperties{}

func (s SmbMountEndpointUpdateProperties) MarshalJSON() ([]byte, error) {
	type wrapper SmbMountEndpointUpdateProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SmbMountEndpointUpdateProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SmbMountEndpointUpdateProperties: %+v", err)
	}

	decoded["endpointType"] = "SmbMount"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SmbMountEndpointUpdateProperties: %+v", err)
	}

	return encoded, nil
}
