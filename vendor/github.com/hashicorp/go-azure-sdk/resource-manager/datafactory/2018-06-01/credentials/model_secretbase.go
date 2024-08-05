package credentials

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretBase interface {
}

// RawSecretBaseImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSecretBaseImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalSecretBaseImplementation(input []byte) (SecretBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SecretBase into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureKeyVaultSecret") {
		var out AzureKeyVaultSecretReference
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureKeyVaultSecretReference: %+v", err)
		}
		return out, nil
	}

	out := RawSecretBaseImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
