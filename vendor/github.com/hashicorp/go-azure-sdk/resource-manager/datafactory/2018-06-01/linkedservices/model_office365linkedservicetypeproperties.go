package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Office365LinkedServiceTypeProperties struct {
	EncryptedCredential      *string     `json:"encryptedCredential,omitempty"`
	Office365TenantId        interface{} `json:"office365TenantId"`
	ServicePrincipalId       interface{} `json:"servicePrincipalId"`
	ServicePrincipalKey      SecretBase  `json:"servicePrincipalKey"`
	ServicePrincipalTenantId interface{} `json:"servicePrincipalTenantId"`
}

var _ json.Unmarshaler = &Office365LinkedServiceTypeProperties{}

func (s *Office365LinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		EncryptedCredential      *string     `json:"encryptedCredential,omitempty"`
		Office365TenantId        interface{} `json:"office365TenantId"`
		ServicePrincipalId       interface{} `json:"servicePrincipalId"`
		ServicePrincipalTenantId interface{} `json:"servicePrincipalTenantId"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.EncryptedCredential = decoded.EncryptedCredential
	s.Office365TenantId = decoded.Office365TenantId
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.ServicePrincipalTenantId = decoded.ServicePrincipalTenantId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling Office365LinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'Office365LinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
