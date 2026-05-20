package links

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretInfoBase interface {
	SecretInfoBase() BaseSecretInfoBaseImpl
}

var _ SecretInfoBase = BaseSecretInfoBaseImpl{}

type BaseSecretInfoBaseImpl struct {
	SecretType SecretType `json:"secretType"`
}

func (s BaseSecretInfoBaseImpl) SecretInfoBase() BaseSecretInfoBaseImpl {
	return s
}

var _ SecretInfoBase = RawSecretInfoBaseImpl{}

// RawSecretInfoBaseImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSecretInfoBaseImpl struct {
	secretInfoBase BaseSecretInfoBaseImpl
	Type           string
	Values         map[string]interface{}
}

func (s RawSecretInfoBaseImpl) SecretInfoBase() BaseSecretInfoBaseImpl {
	return s.secretInfoBase
}

func UnmarshalSecretInfoBaseImplementation(input []byte) (SecretInfoBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SecretInfoBase into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["secretType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "keyVaultSecretReference") {
		var out KeyVaultSecretReferenceSecretInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KeyVaultSecretReferenceSecretInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "keyVaultSecretUri") {
		var out KeyVaultSecretUriSecretInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KeyVaultSecretUriSecretInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "rawValue") {
		var out ValueSecretInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ValueSecretInfo: %+v", err)
		}
		return out, nil
	}

	var parent BaseSecretInfoBaseImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSecretInfoBaseImpl: %+v", err)
	}

	return RawSecretInfoBaseImpl{
		secretInfoBase: parent,
		Type:           value,
		Values:         temp,
	}, nil

}
