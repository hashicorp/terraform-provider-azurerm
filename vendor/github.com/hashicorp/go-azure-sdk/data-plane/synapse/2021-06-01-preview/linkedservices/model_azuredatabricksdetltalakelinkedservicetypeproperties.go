package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureDatabricksDetltaLakeLinkedServiceTypeProperties struct {
	AccessToken         SecretBase   `json:"accessToken"`
	ClusterId           *interface{} `json:"clusterId,omitempty"`
	Domain              interface{}  `json:"domain"`
	EncryptedCredential *interface{} `json:"encryptedCredential,omitempty"`
}

var _ json.Unmarshaler = &AzureDatabricksDetltaLakeLinkedServiceTypeProperties{}

func (s *AzureDatabricksDetltaLakeLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ClusterId           *interface{} `json:"clusterId,omitempty"`
		Domain              interface{}  `json:"domain"`
		EncryptedCredential *interface{} `json:"encryptedCredential,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ClusterId = decoded.ClusterId
	s.Domain = decoded.Domain
	s.EncryptedCredential = decoded.EncryptedCredential

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureDatabricksDetltaLakeLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["accessToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AccessToken' for 'AzureDatabricksDetltaLakeLinkedServiceTypeProperties': %+v", err)
		}
		s.AccessToken = impl
	}

	return nil
}
