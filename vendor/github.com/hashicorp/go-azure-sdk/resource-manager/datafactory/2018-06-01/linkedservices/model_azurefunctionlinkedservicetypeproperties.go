package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFunctionLinkedServiceTypeProperties struct {
	Authentication      *interface{}         `json:"authentication,omitempty"`
	Credential          *CredentialReference `json:"credential,omitempty"`
	EncryptedCredential *string              `json:"encryptedCredential,omitempty"`
	FunctionAppURL      interface{}          `json:"functionAppUrl"`
	FunctionKey         SecretBase           `json:"functionKey"`
	ResourceId          *interface{}         `json:"resourceId,omitempty"`
}

var _ json.Unmarshaler = &AzureFunctionLinkedServiceTypeProperties{}

func (s *AzureFunctionLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Authentication      *interface{}         `json:"authentication,omitempty"`
		Credential          *CredentialReference `json:"credential,omitempty"`
		EncryptedCredential *string              `json:"encryptedCredential,omitempty"`
		FunctionAppURL      interface{}          `json:"functionAppUrl"`
		ResourceId          *interface{}         `json:"resourceId,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Authentication = decoded.Authentication
	s.Credential = decoded.Credential
	s.EncryptedCredential = decoded.EncryptedCredential
	s.FunctionAppURL = decoded.FunctionAppURL
	s.ResourceId = decoded.ResourceId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureFunctionLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["functionKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'FunctionKey' for 'AzureFunctionLinkedServiceTypeProperties': %+v", err)
		}
		s.FunctionKey = impl
	}

	return nil
}
