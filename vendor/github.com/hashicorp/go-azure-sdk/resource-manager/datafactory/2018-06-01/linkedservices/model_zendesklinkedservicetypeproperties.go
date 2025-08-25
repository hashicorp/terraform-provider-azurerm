package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ZendeskLinkedServiceTypeProperties struct {
	ApiToken            SecretBase                `json:"apiToken"`
	AuthenticationType  ZendeskAuthenticationType `json:"authenticationType"`
	EncryptedCredential *string                   `json:"encryptedCredential,omitempty"`
	Password            SecretBase                `json:"password"`
	Url                 interface{}               `json:"url"`
	UserName            *interface{}              `json:"userName,omitempty"`
}

var _ json.Unmarshaler = &ZendeskLinkedServiceTypeProperties{}

func (s *ZendeskLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType  ZendeskAuthenticationType `json:"authenticationType"`
		EncryptedCredential *string                   `json:"encryptedCredential,omitempty"`
		Url                 interface{}               `json:"url"`
		UserName            *interface{}              `json:"userName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Url = decoded.Url
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ZendeskLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["apiToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ApiToken' for 'ZendeskLinkedServiceTypeProperties': %+v", err)
		}
		s.ApiToken = impl
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'ZendeskLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
