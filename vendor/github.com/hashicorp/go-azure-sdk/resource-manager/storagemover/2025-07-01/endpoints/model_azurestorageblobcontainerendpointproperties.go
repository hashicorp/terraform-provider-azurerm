package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EndpointBaseProperties = AzureStorageBlobContainerEndpointProperties{}

type AzureStorageBlobContainerEndpointProperties struct {
	BlobContainerName        string `json:"blobContainerName"`
	StorageAccountResourceId string `json:"storageAccountResourceId"`

	// Fields inherited from EndpointBaseProperties

	Description       *string            `json:"description,omitempty"`
	EndpointType      EndpointType       `json:"endpointType"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}

func (s AzureStorageBlobContainerEndpointProperties) EndpointBaseProperties() BaseEndpointBasePropertiesImpl {
	return BaseEndpointBasePropertiesImpl{
		Description:       s.Description,
		EndpointType:      s.EndpointType,
		ProvisioningState: s.ProvisioningState,
	}
}

var _ json.Marshaler = AzureStorageBlobContainerEndpointProperties{}

func (s AzureStorageBlobContainerEndpointProperties) MarshalJSON() ([]byte, error) {
	type wrapper AzureStorageBlobContainerEndpointProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureStorageBlobContainerEndpointProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureStorageBlobContainerEndpointProperties: %+v", err)
	}

	decoded["endpointType"] = "AzureStorageBlobContainer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureStorageBlobContainerEndpointProperties: %+v", err)
	}

	return encoded, nil
}
