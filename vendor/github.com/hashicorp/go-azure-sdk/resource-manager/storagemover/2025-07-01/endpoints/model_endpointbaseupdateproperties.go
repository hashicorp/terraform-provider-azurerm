package endpoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointBaseUpdateProperties interface {
	EndpointBaseUpdateProperties() BaseEndpointBaseUpdatePropertiesImpl
}

var _ EndpointBaseUpdateProperties = BaseEndpointBaseUpdatePropertiesImpl{}

type BaseEndpointBaseUpdatePropertiesImpl struct {
	Description  *string      `json:"description,omitempty"`
	EndpointType EndpointType `json:"endpointType"`
}

func (s BaseEndpointBaseUpdatePropertiesImpl) EndpointBaseUpdateProperties() BaseEndpointBaseUpdatePropertiesImpl {
	return s
}

var _ EndpointBaseUpdateProperties = RawEndpointBaseUpdatePropertiesImpl{}

// RawEndpointBaseUpdatePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEndpointBaseUpdatePropertiesImpl struct {
	endpointBaseUpdateProperties BaseEndpointBaseUpdatePropertiesImpl
	Type                         string
	Values                       map[string]interface{}
}

func (s RawEndpointBaseUpdatePropertiesImpl) EndpointBaseUpdateProperties() BaseEndpointBaseUpdatePropertiesImpl {
	return s.endpointBaseUpdateProperties
}

func UnmarshalEndpointBaseUpdatePropertiesImplementation(input []byte) (EndpointBaseUpdateProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EndpointBaseUpdateProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["endpointType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureMultiCloudConnector") {
		var out AzureMultiCloudConnectorEndpointUpdateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureMultiCloudConnectorEndpointUpdateProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureStorageBlobContainer") {
		var out AzureStorageBlobContainerEndpointUpdateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureStorageBlobContainerEndpointUpdateProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureStorageNfsFileShare") {
		var out AzureStorageNfsFileShareEndpointUpdateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureStorageNfsFileShareEndpointUpdateProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureStorageSmbFileShare") {
		var out AzureStorageSmbFileShareEndpointUpdateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureStorageSmbFileShareEndpointUpdateProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NfsMount") {
		var out NfsMountEndpointUpdateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NfsMountEndpointUpdateProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SmbMount") {
		var out SmbMountEndpointUpdateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SmbMountEndpointUpdateProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseEndpointBaseUpdatePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseEndpointBaseUpdatePropertiesImpl: %+v", err)
	}

	return RawEndpointBaseUpdatePropertiesImpl{
		endpointBaseUpdateProperties: parent,
		Type:                         value,
		Values:                       temp,
	}, nil

}
