package backupinstanceresources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthCredentials interface {
	AuthCredentials() BaseAuthCredentialsImpl
}

var _ AuthCredentials = BaseAuthCredentialsImpl{}

type BaseAuthCredentialsImpl struct {
	ObjectType string `json:"objectType"`
}

func (s BaseAuthCredentialsImpl) AuthCredentials() BaseAuthCredentialsImpl {
	return s
}

var _ AuthCredentials = RawAuthCredentialsImpl{}

// RawAuthCredentialsImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawAuthCredentialsImpl struct {
	authCredentials BaseAuthCredentialsImpl
	Type            string
	Values          map[string]interface{}
}

func (s RawAuthCredentialsImpl) AuthCredentials() BaseAuthCredentialsImpl {
	return s.authCredentials
}

func (s RawAuthCredentialsImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalAuthCredentialsImplementation(input []byte) (AuthCredentials, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AuthCredentials into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "SecretStoreBasedAuthCredentials") {
		var out SecretStoreBasedAuthCredentials
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SecretStoreBasedAuthCredentials: %+v", err)
		}
		return out, nil
	}

	var parent BaseAuthCredentialsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseAuthCredentialsImpl: %+v", err)
	}

	return RawAuthCredentialsImpl{
		authCredentials: parent,
		Type:            value,
		Values:          temp,
	}, nil

}
