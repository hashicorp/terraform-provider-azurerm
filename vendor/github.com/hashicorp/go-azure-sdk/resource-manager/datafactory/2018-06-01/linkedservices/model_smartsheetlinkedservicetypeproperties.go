package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SmartsheetLinkedServiceTypeProperties struct {
	ApiToken            SecretBase `json:"apiToken"`
	EncryptedCredential *string    `json:"encryptedCredential,omitempty"`
}

var _ json.Unmarshaler = &SmartsheetLinkedServiceTypeProperties{}

func (s *SmartsheetLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		EncryptedCredential *string `json:"encryptedCredential,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.EncryptedCredential = decoded.EncryptedCredential

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SmartsheetLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["apiToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ApiToken' for 'SmartsheetLinkedServiceTypeProperties': %+v", err)
		}
		s.ApiToken = impl
	}

	return nil
}
