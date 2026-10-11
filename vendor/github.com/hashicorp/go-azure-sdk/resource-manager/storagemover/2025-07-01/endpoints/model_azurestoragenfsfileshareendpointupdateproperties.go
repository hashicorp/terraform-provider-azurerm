package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EndpointBaseUpdateProperties = AzureStorageNfsFileShareEndpointUpdateProperties{}

type AzureStorageNfsFileShareEndpointUpdateProperties struct {

	// Fields inherited from EndpointBaseUpdateProperties

	Description  *string      `json:"description,omitempty"`
	EndpointType EndpointType `json:"endpointType"`
}

func (s AzureStorageNfsFileShareEndpointUpdateProperties) EndpointBaseUpdateProperties() BaseEndpointBaseUpdatePropertiesImpl {
	return BaseEndpointBaseUpdatePropertiesImpl{
		Description:  s.Description,
		EndpointType: s.EndpointType,
	}
}

var _ json.Marshaler = AzureStorageNfsFileShareEndpointUpdateProperties{}

func (s AzureStorageNfsFileShareEndpointUpdateProperties) MarshalJSON() ([]byte, error) {
	type wrapper AzureStorageNfsFileShareEndpointUpdateProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureStorageNfsFileShareEndpointUpdateProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureStorageNfsFileShareEndpointUpdateProperties: %+v", err)
	}

	decoded["endpointType"] = "AzureStorageNfsFileShare"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureStorageNfsFileShareEndpointUpdateProperties: %+v", err)
	}

	return encoded, nil
}
