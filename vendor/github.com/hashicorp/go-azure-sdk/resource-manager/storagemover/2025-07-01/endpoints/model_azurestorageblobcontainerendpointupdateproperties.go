package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EndpointBaseUpdateProperties = AzureStorageBlobContainerEndpointUpdateProperties{}

type AzureStorageBlobContainerEndpointUpdateProperties struct {

	// Fields inherited from EndpointBaseUpdateProperties

	Description  *string      `json:"description,omitempty"`
	EndpointType EndpointType `json:"endpointType"`
}

func (s AzureStorageBlobContainerEndpointUpdateProperties) EndpointBaseUpdateProperties() BaseEndpointBaseUpdatePropertiesImpl {
	return BaseEndpointBaseUpdatePropertiesImpl{
		Description:  s.Description,
		EndpointType: s.EndpointType,
	}
}

var _ json.Marshaler = AzureStorageBlobContainerEndpointUpdateProperties{}

func (s AzureStorageBlobContainerEndpointUpdateProperties) MarshalJSON() ([]byte, error) {
	type wrapper AzureStorageBlobContainerEndpointUpdateProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureStorageBlobContainerEndpointUpdateProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureStorageBlobContainerEndpointUpdateProperties: %+v", err)
	}

	decoded["endpointType"] = "AzureStorageBlobContainer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureStorageBlobContainerEndpointUpdateProperties: %+v", err)
	}

	return encoded, nil
}
