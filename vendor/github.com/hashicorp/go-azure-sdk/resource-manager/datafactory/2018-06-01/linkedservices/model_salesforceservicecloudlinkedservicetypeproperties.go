package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SalesforceServiceCloudLinkedServiceTypeProperties struct {
	ApiVersion          *interface{} `json:"apiVersion,omitempty"`
	EncryptedCredential *string      `json:"encryptedCredential,omitempty"`
	EnvironmentURL      *interface{} `json:"environmentUrl,omitempty"`
	ExtendedProperties  *interface{} `json:"extendedProperties,omitempty"`
	Password            SecretBase   `json:"password"`
	SecurityToken       SecretBase   `json:"securityToken"`
	Username            *interface{} `json:"username,omitempty"`
}

var _ json.Unmarshaler = &SalesforceServiceCloudLinkedServiceTypeProperties{}

func (s *SalesforceServiceCloudLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ApiVersion          *interface{} `json:"apiVersion,omitempty"`
		EncryptedCredential *string      `json:"encryptedCredential,omitempty"`
		EnvironmentURL      *interface{} `json:"environmentUrl,omitempty"`
		ExtendedProperties  *interface{} `json:"extendedProperties,omitempty"`
		Username            *interface{} `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ApiVersion = decoded.ApiVersion
	s.EncryptedCredential = decoded.EncryptedCredential
	s.EnvironmentURL = decoded.EnvironmentURL
	s.ExtendedProperties = decoded.ExtendedProperties
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SalesforceServiceCloudLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'SalesforceServiceCloudLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	if v, ok := temp["securityToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'SecurityToken' for 'SalesforceServiceCloudLinkedServiceTypeProperties': %+v", err)
		}
		s.SecurityToken = impl
	}

	return nil
}
