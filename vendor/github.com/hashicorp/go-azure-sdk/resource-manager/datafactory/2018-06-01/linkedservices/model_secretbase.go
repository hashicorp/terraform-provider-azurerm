package linkedservices

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretBase interface {
	SecretBase() BaseSecretBaseImpl
}

var _ SecretBase = BaseSecretBaseImpl{}

type BaseSecretBaseImpl struct {
	Type string `json:"type"`
}

func (s BaseSecretBaseImpl) SecretBase() BaseSecretBaseImpl {
	return s
}

var _ SecretBase = RawSecretBaseImpl{}

// RawSecretBaseImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSecretBaseImpl struct {
	secretBase BaseSecretBaseImpl
	Type       string
	Values     map[string]interface{}
}

func (s RawSecretBaseImpl) SecretBase() BaseSecretBaseImpl {
	return s.secretBase
}

func UnmarshalSecretBaseImplementation(input []byte) (SecretBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SecretBase into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureKeyVaultSecret") {
		var out AzureKeyVaultSecretReference
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureKeyVaultSecretReference: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SecureString") {
		var out SecureString
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SecureString: %+v", err)
		}
		return out, nil
	}

	var parent BaseSecretBaseImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSecretBaseImpl: %+v", err)
	}

	return RawSecretBaseImpl{
		secretBase: parent,
		Type:       value,
		Values:     temp,
	}, nil

}
