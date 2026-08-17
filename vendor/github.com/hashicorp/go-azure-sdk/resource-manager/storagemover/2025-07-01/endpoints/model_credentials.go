package endpoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Credentials interface {
	Credentials() BaseCredentialsImpl
}

var _ Credentials = BaseCredentialsImpl{}

type BaseCredentialsImpl struct {
	Type CredentialType `json:"type"`
}

func (s BaseCredentialsImpl) Credentials() BaseCredentialsImpl {
	return s
}

var _ Credentials = RawCredentialsImpl{}

// RawCredentialsImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawCredentialsImpl struct {
	credentials BaseCredentialsImpl
	Type        string
	Values      map[string]interface{}
}

func (s RawCredentialsImpl) Credentials() BaseCredentialsImpl {
	return s.credentials
}

func (s RawCredentialsImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalCredentialsImplementation(input []byte) (Credentials, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Credentials into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureKeyVaultSmb") {
		var out AzureKeyVaultSmbCredentials
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureKeyVaultSmbCredentials: %+v", err)
		}
		return out, nil
	}

	var parent BaseCredentialsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseCredentialsImpl: %+v", err)
	}

	return RawCredentialsImpl{
		credentials: parent,
		Type:        value,
		Values:      temp,
	}, nil

}
