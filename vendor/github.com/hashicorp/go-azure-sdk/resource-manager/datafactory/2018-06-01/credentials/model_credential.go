package credentials

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Credential interface {
	Credential() BaseCredentialImpl
}

var _ Credential = BaseCredentialImpl{}

type BaseCredentialImpl struct {
	Annotations *[]interface{} `json:"annotations,omitempty"`
	Description *string        `json:"description,omitempty"`
	Type        string         `json:"type"`
}

func (s BaseCredentialImpl) Credential() BaseCredentialImpl {
	return s
}

var _ Credential = RawCredentialImpl{}

// RawCredentialImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawCredentialImpl struct {
	credential BaseCredentialImpl
	Type       string
	Values     map[string]interface{}
}

func (s RawCredentialImpl) Credential() BaseCredentialImpl {
	return s.credential
}

func UnmarshalCredentialImplementation(input []byte) (Credential, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Credential into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "ManagedIdentity") {
		var out ManagedIdentityCredential
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ManagedIdentityCredential: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServicePrincipal") {
		var out ServicePrincipalCredential
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServicePrincipalCredential: %+v", err)
		}
		return out, nil
	}

	var parent BaseCredentialImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseCredentialImpl: %+v", err)
	}

	return RawCredentialImpl{
		credential: parent,
		Type:       value,
		Values:     temp,
	}, nil

}
