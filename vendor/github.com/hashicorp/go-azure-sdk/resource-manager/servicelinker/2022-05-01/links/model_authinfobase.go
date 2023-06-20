package links

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthInfoBase interface {
}

func unmarshalAuthInfoBaseImplementation(input []byte) (AuthInfoBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AuthInfoBase into map[string]interface: %+v", err)
	}

	value, ok := temp["authType"].(string)
	if !ok {
		return nil, nil
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

	if strings.EqualFold(value, "userAssignedIdentity") {
		var out UserAssignedIdentityAuthInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UserAssignedIdentityAuthInfo: %+v", err)
		}
		return out, nil
	}

	type RawAuthInfoBaseImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawAuthInfoBaseImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
