package servicelinker

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureResourcePropertiesBase interface {
	AzureResourcePropertiesBase() BaseAzureResourcePropertiesBaseImpl
}

var _ AzureResourcePropertiesBase = BaseAzureResourcePropertiesBaseImpl{}

type BaseAzureResourcePropertiesBaseImpl struct {
	Type AzureResourceType `json:"type"`
}

func (s BaseAzureResourcePropertiesBaseImpl) AzureResourcePropertiesBase() BaseAzureResourcePropertiesBaseImpl {
	return s
}

var _ AzureResourcePropertiesBase = RawAzureResourcePropertiesBaseImpl{}

// RawAzureResourcePropertiesBaseImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawAzureResourcePropertiesBaseImpl struct {
	azureResourcePropertiesBase BaseAzureResourcePropertiesBaseImpl
	Type                        string
	Values                      map[string]interface{}
}

func (s RawAzureResourcePropertiesBaseImpl) AzureResourcePropertiesBase() BaseAzureResourcePropertiesBaseImpl {
	return s.azureResourcePropertiesBase
}

func UnmarshalAzureResourcePropertiesBaseImplementation(input []byte) (AzureResourcePropertiesBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureResourcePropertiesBase into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "KeyVault") {
		var out AzureKeyVaultProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureKeyVaultProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseAzureResourcePropertiesBaseImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseAzureResourcePropertiesBaseImpl: %+v", err)
	}

	return RawAzureResourcePropertiesBaseImpl{
		azureResourcePropertiesBase: parent,
		Type:                        value,
		Values:                      temp,
	}, nil

}
