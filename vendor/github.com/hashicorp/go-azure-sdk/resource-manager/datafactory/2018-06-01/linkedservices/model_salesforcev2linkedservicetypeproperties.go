package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SalesforceV2LinkedServiceTypeProperties struct {
	ApiVersion          *interface{} `json:"apiVersion,omitempty"`
	AuthenticationType  *interface{} `json:"authenticationType,omitempty"`
	ClientId            *interface{} `json:"clientId,omitempty"`
	ClientSecret        SecretBase   `json:"clientSecret"`
	EncryptedCredential *string      `json:"encryptedCredential,omitempty"`
	EnvironmentURL      *interface{} `json:"environmentUrl,omitempty"`
}

var _ json.Unmarshaler = &SalesforceV2LinkedServiceTypeProperties{}

func (s *SalesforceV2LinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ApiVersion          *interface{} `json:"apiVersion,omitempty"`
		AuthenticationType  *interface{} `json:"authenticationType,omitempty"`
		ClientId            *interface{} `json:"clientId,omitempty"`
		EncryptedCredential *string      `json:"encryptedCredential,omitempty"`
		EnvironmentURL      *interface{} `json:"environmentUrl,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ApiVersion = decoded.ApiVersion
	s.AuthenticationType = decoded.AuthenticationType
	s.ClientId = decoded.ClientId
	s.EncryptedCredential = decoded.EncryptedCredential
	s.EnvironmentURL = decoded.EnvironmentURL

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SalesforceV2LinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["clientSecret"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ClientSecret' for 'SalesforceV2LinkedServiceTypeProperties': %+v", err)
		}
		s.ClientSecret = impl
	}

	return nil
}
