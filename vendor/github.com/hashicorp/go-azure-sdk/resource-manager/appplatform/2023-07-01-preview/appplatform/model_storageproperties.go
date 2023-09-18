package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageProperties interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawStoragePropertiesImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalStoragePropertiesImplementation(input []byte) (StorageProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling StorageProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["storageType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "StorageAccount") {
		var out StorageAccount
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StorageAccount: %+v", err)
		}
		return out, nil
	}

	out := RawStoragePropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
