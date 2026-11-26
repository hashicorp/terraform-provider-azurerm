package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EndpointBaseProperties = AzureStorageNfsFileShareEndpointProperties{}

type AzureStorageNfsFileShareEndpointProperties struct {
	FileShareName            string `json:"fileShareName"`
	StorageAccountResourceId string `json:"storageAccountResourceId"`

	// Fields inherited from EndpointBaseProperties

	Description       *string            `json:"description,omitempty"`
	EndpointType      EndpointType       `json:"endpointType"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}

func (s AzureStorageNfsFileShareEndpointProperties) EndpointBaseProperties() BaseEndpointBasePropertiesImpl {
	return BaseEndpointBasePropertiesImpl{
		Description:       s.Description,
		EndpointType:      s.EndpointType,
		ProvisioningState: s.ProvisioningState,
	}
}

var _ json.Marshaler = AzureStorageNfsFileShareEndpointProperties{}

func (s AzureStorageNfsFileShareEndpointProperties) MarshalJSON() ([]byte, error) {
	type wrapper AzureStorageNfsFileShareEndpointProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureStorageNfsFileShareEndpointProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureStorageNfsFileShareEndpointProperties: %+v", err)
	}

	decoded["endpointType"] = "AzureStorageNfsFileShare"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureStorageNfsFileShareEndpointProperties: %+v", err)
	}

	return encoded, nil
}
