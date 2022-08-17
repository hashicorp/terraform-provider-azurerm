package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthCredentials interface {
}

func unmarshalAuthCredentialsImplementation(input []byte) (AuthCredentials, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AuthCredentials into map[string]interface: %+v", err)
	}

	value, ok := temp["objectType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "SecretStoreBasedAuthCredentials") {
		var out SecretStoreBasedAuthCredentials
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SecretStoreBasedAuthCredentials: %+v", err)
		}
		return out, nil
	}

	type RawAuthCredentialsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawAuthCredentialsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
