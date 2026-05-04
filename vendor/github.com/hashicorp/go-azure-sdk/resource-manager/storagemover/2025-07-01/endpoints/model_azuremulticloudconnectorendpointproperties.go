package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EndpointBaseProperties = AzureMultiCloudConnectorEndpointProperties{}

type AzureMultiCloudConnectorEndpointProperties struct {
	AwsS3BucketId         string `json:"awsS3BucketId"`
	MultiCloudConnectorId string `json:"multiCloudConnectorId"`

	// Fields inherited from EndpointBaseProperties

	Description       *string            `json:"description,omitempty"`
	EndpointType      EndpointType       `json:"endpointType"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}

func (s AzureMultiCloudConnectorEndpointProperties) EndpointBaseProperties() BaseEndpointBasePropertiesImpl {
	return BaseEndpointBasePropertiesImpl{
		Description:       s.Description,
		EndpointType:      s.EndpointType,
		ProvisioningState: s.ProvisioningState,
	}
}

var _ json.Marshaler = AzureMultiCloudConnectorEndpointProperties{}

func (s AzureMultiCloudConnectorEndpointProperties) MarshalJSON() ([]byte, error) {
	type wrapper AzureMultiCloudConnectorEndpointProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureMultiCloudConnectorEndpointProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureMultiCloudConnectorEndpointProperties: %+v", err)
	}

	decoded["endpointType"] = "AzureMultiCloudConnector"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureMultiCloudConnectorEndpointProperties: %+v", err)
	}

	return encoded, nil
}
