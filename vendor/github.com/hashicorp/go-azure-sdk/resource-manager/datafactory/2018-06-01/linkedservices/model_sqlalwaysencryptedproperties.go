package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlAlwaysEncryptedProperties struct {
	AlwaysEncryptedAkvAuthType SqlAlwaysEncryptedAkvAuthType `json:"alwaysEncryptedAkvAuthType"`
	Credential                 *CredentialReference          `json:"credential,omitempty"`
	ServicePrincipalId         *interface{}                  `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey        SecretBase                    `json:"servicePrincipalKey"`
}

var _ json.Unmarshaler = &SqlAlwaysEncryptedProperties{}

func (s *SqlAlwaysEncryptedProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AlwaysEncryptedAkvAuthType SqlAlwaysEncryptedAkvAuthType `json:"alwaysEncryptedAkvAuthType"`
		Credential                 *CredentialReference          `json:"credential,omitempty"`
		ServicePrincipalId         *interface{}                  `json:"servicePrincipalId,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AlwaysEncryptedAkvAuthType = decoded.AlwaysEncryptedAkvAuthType
	s.Credential = decoded.Credential
	s.ServicePrincipalId = decoded.ServicePrincipalId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SqlAlwaysEncryptedProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'SqlAlwaysEncryptedProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
