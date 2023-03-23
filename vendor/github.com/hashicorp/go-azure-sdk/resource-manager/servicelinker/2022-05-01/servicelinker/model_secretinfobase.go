package servicelinker

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretInfoBase interface {
}

func unmarshalSecretInfoBaseImplementation(input []byte) (SecretInfoBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SecretInfoBase into map[string]interface: %+v", err)
	}

	value, ok := temp["secretType"].(string)
	if !ok {
		return nil, nil
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

	type RawSecretInfoBaseImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawSecretInfoBaseImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
