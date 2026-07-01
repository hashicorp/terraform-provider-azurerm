package servicelinker

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthInfoBase interface {
	AuthInfoBase() BaseAuthInfoBaseImpl
}

var _ AuthInfoBase = BaseAuthInfoBaseImpl{}

type BaseAuthInfoBaseImpl struct {
	AuthMode *AuthMode `json:"authMode,omitempty"`
	AuthType AuthType  `json:"authType"`
}

func (s BaseAuthInfoBaseImpl) AuthInfoBase() BaseAuthInfoBaseImpl {
	return s
}

var _ AuthInfoBase = RawAuthInfoBaseImpl{}

// RawAuthInfoBaseImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawAuthInfoBaseImpl struct {
	authInfoBase BaseAuthInfoBaseImpl
	Type         string
	Values       map[string]interface{}
}

func (s RawAuthInfoBaseImpl) AuthInfoBase() BaseAuthInfoBaseImpl {
	return s.authInfoBase
}

func (s RawAuthInfoBaseImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalAuthInfoBaseImplementation(input []byte) (AuthInfoBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AuthInfoBase into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["authType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "accessKey") {
		var out AccessKeyInfoBase
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AccessKeyInfoBase: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "easyAuthMicrosoftEntraID") {
		var out EasyAuthMicrosoftEntraIDAuthInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EasyAuthMicrosoftEntraIDAuthInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "secret") {
		var out SecretAuthInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SecretAuthInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "servicePrincipalCertificate") {
		var out ServicePrincipalCertificateAuthInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServicePrincipalCertificateAuthInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "servicePrincipalSecret") {
		var out ServicePrincipalSecretAuthInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServicePrincipalSecretAuthInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "systemAssignedIdentity") {
		var out SystemAssignedIdentityAuthInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SystemAssignedIdentityAuthInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "userAccount") {
		var out UserAccountAuthInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UserAccountAuthInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "userAssignedIdentity") {
		var out UserAssignedIdentityAuthInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UserAssignedIdentityAuthInfo: %+v", err)
		}
		return out, nil
	}

	var parent BaseAuthInfoBaseImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseAuthInfoBaseImpl: %+v", err)
	}

	return RawAuthInfoBaseImpl{
		authInfoBase: parent,
		Type:         value,
		Values:       temp,
	}, nil

}
