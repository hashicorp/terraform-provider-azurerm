package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuickbaseLinkedServiceTypeProperties struct {
	EncryptedCredential *string     `json:"encryptedCredential,omitempty"`
	Url                 interface{} `json:"url"`
	UserToken           SecretBase  `json:"userToken"`
}

var _ json.Unmarshaler = &QuickbaseLinkedServiceTypeProperties{}

func (s *QuickbaseLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		EncryptedCredential *string     `json:"encryptedCredential,omitempty"`
		Url                 interface{} `json:"url"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.EncryptedCredential = decoded.EncryptedCredential
	s.Url = decoded.Url

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling QuickbaseLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["userToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'UserToken' for 'QuickbaseLinkedServiceTypeProperties': %+v", err)
		}
		s.UserToken = impl
	}

	return nil
}
