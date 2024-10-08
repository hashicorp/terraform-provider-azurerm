package endpoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointBaseProperties interface {
	EndpointBaseProperties() BaseEndpointBasePropertiesImpl
}

var _ EndpointBaseProperties = BaseEndpointBasePropertiesImpl{}

type BaseEndpointBasePropertiesImpl struct {
	Description       *string            `json:"description,omitempty"`
	EndpointType      EndpointType       `json:"endpointType"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}

func (s BaseEndpointBasePropertiesImpl) EndpointBaseProperties() BaseEndpointBasePropertiesImpl {
	return s
}

var _ EndpointBaseProperties = RawEndpointBasePropertiesImpl{}

// RawEndpointBasePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEndpointBasePropertiesImpl struct {
	endpointBaseProperties BaseEndpointBasePropertiesImpl
	Type                   string
	Values                 map[string]interface{}
}

func (s RawEndpointBasePropertiesImpl) EndpointBaseProperties() BaseEndpointBasePropertiesImpl {
	return s.endpointBaseProperties
}

func UnmarshalEndpointBasePropertiesImplementation(input []byte) (EndpointBaseProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EndpointBaseProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["endpointType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureStorageBlobContainer") {
		var out AzureStorageBlobContainerEndpointProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureStorageBlobContainerEndpointProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NfsMount") {
		var out NfsMountEndpointProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NfsMountEndpointProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseEndpointBasePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseEndpointBasePropertiesImpl: %+v", err)
	}

	return RawEndpointBasePropertiesImpl{
		endpointBaseProperties: parent,
		Type:                   value,
		Values:                 temp,
	}, nil

}
