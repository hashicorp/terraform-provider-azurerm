package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EndpointBaseUpdateProperties = AzureMultiCloudConnectorEndpointUpdateProperties{}

type AzureMultiCloudConnectorEndpointUpdateProperties struct {

	// Fields inherited from EndpointBaseUpdateProperties

	Description  *string      `json:"description,omitempty"`
	EndpointType EndpointType `json:"endpointType"`
}

func (s AzureMultiCloudConnectorEndpointUpdateProperties) EndpointBaseUpdateProperties() BaseEndpointBaseUpdatePropertiesImpl {
	return BaseEndpointBaseUpdatePropertiesImpl{
		Description:  s.Description,
		EndpointType: s.EndpointType,
	}
}

var _ json.Marshaler = AzureMultiCloudConnectorEndpointUpdateProperties{}

func (s AzureMultiCloudConnectorEndpointUpdateProperties) MarshalJSON() ([]byte, error) {
	type wrapper AzureMultiCloudConnectorEndpointUpdateProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureMultiCloudConnectorEndpointUpdateProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureMultiCloudConnectorEndpointUpdateProperties: %+v", err)
	}

	decoded["endpointType"] = "AzureMultiCloudConnector"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureMultiCloudConnectorEndpointUpdateProperties: %+v", err)
	}

	return encoded, nil
}
