package endpoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointBaseProperties interface {
}

func unmarshalEndpointBasePropertiesImplementation(input []byte) (EndpointBaseProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EndpointBaseProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["endpointType"].(string)
	if !ok {
		return nil, nil
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

	type RawEndpointBasePropertiesImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawEndpointBasePropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
