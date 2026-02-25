package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EndpointBaseUpdateProperties = AzureStorageSmbFileShareEndpointUpdateProperties{}

type AzureStorageSmbFileShareEndpointUpdateProperties struct {

	// Fields inherited from EndpointBaseUpdateProperties

	Description  *string      `json:"description,omitempty"`
	EndpointType EndpointType `json:"endpointType"`
}

func (s AzureStorageSmbFileShareEndpointUpdateProperties) EndpointBaseUpdateProperties() BaseEndpointBaseUpdatePropertiesImpl {
	return BaseEndpointBaseUpdatePropertiesImpl{
		Description:  s.Description,
		EndpointType: s.EndpointType,
	}
}

var _ json.Marshaler = AzureStorageSmbFileShareEndpointUpdateProperties{}

func (s AzureStorageSmbFileShareEndpointUpdateProperties) MarshalJSON() ([]byte, error) {
	type wrapper AzureStorageSmbFileShareEndpointUpdateProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureStorageSmbFileShareEndpointUpdateProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureStorageSmbFileShareEndpointUpdateProperties: %+v", err)
	}

	decoded["endpointType"] = "AzureStorageSmbFileShare"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureStorageSmbFileShareEndpointUpdateProperties: %+v", err)
	}

	return encoded, nil
}
