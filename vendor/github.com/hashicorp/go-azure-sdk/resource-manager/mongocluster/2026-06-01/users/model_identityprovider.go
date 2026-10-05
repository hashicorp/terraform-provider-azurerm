package users

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IdentityProvider interface {
	IdentityProvider() BaseIdentityProviderImpl
}

var _ IdentityProvider = BaseIdentityProviderImpl{}

type BaseIdentityProviderImpl struct {
	Type IdentityProviderType `json:"type"`
}

func (s BaseIdentityProviderImpl) IdentityProvider() BaseIdentityProviderImpl {
	return s
}

var _ IdentityProvider = RawIdentityProviderImpl{}

// RawIdentityProviderImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawIdentityProviderImpl struct {
	identityProvider BaseIdentityProviderImpl
	Type             string
	Values           map[string]interface{}
}

func (s RawIdentityProviderImpl) IdentityProvider() BaseIdentityProviderImpl {
	return s.identityProvider
}

func (s RawIdentityProviderImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalIdentityProviderImplementation(input []byte) (IdentityProvider, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling IdentityProvider into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "MicrosoftEntraID") {
		var out EntraIdentityProvider
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EntraIdentityProvider: %+v", err)
		}
		return out, nil
	}

	var parent BaseIdentityProviderImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseIdentityProviderImpl: %+v", err)
	}

	return RawIdentityProviderImpl{
		identityProvider: parent,
		Type:             value,
		Values:           temp,
	}, nil

}
