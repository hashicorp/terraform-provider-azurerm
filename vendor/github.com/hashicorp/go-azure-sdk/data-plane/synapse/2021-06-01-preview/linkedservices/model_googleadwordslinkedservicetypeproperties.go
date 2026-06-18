package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GoogleAdWordsLinkedServiceTypeProperties struct {
	AuthenticationType   *GoogleAdWordsAuthenticationType `json:"authenticationType,omitempty"`
	ClientCustomerID     *interface{}                     `json:"clientCustomerID,omitempty"`
	ClientId             *interface{}                     `json:"clientId,omitempty"`
	ClientSecret         SecretBase                       `json:"clientSecret"`
	ConnectionProperties *interface{}                     `json:"connectionProperties,omitempty"`
	DeveloperToken       SecretBase                       `json:"developerToken"`
	Email                *interface{}                     `json:"email,omitempty"`
	EncryptedCredential  *interface{}                     `json:"encryptedCredential,omitempty"`
	KeyFilePath          *interface{}                     `json:"keyFilePath,omitempty"`
	RefreshToken         SecretBase                       `json:"refreshToken"`
	TrustedCertPath      *interface{}                     `json:"trustedCertPath,omitempty"`
	UseSystemTrustStore  *interface{}                     `json:"useSystemTrustStore,omitempty"`
}

var _ json.Unmarshaler = &GoogleAdWordsLinkedServiceTypeProperties{}

func (s *GoogleAdWordsLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType   *GoogleAdWordsAuthenticationType `json:"authenticationType,omitempty"`
		ClientCustomerID     *interface{}                     `json:"clientCustomerID,omitempty"`
		ClientId             *interface{}                     `json:"clientId,omitempty"`
		ConnectionProperties *interface{}                     `json:"connectionProperties,omitempty"`
		Email                *interface{}                     `json:"email,omitempty"`
		EncryptedCredential  *interface{}                     `json:"encryptedCredential,omitempty"`
		KeyFilePath          *interface{}                     `json:"keyFilePath,omitempty"`
		TrustedCertPath      *interface{}                     `json:"trustedCertPath,omitempty"`
		UseSystemTrustStore  *interface{}                     `json:"useSystemTrustStore,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.ClientCustomerID = decoded.ClientCustomerID
	s.ClientId = decoded.ClientId
	s.ConnectionProperties = decoded.ConnectionProperties
	s.Email = decoded.Email
	s.EncryptedCredential = decoded.EncryptedCredential
	s.KeyFilePath = decoded.KeyFilePath
	s.TrustedCertPath = decoded.TrustedCertPath
	s.UseSystemTrustStore = decoded.UseSystemTrustStore

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling GoogleAdWordsLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["clientSecret"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ClientSecret' for 'GoogleAdWordsLinkedServiceTypeProperties': %+v", err)
		}
		s.ClientSecret = impl
	}

	if v, ok := temp["developerToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'DeveloperToken' for 'GoogleAdWordsLinkedServiceTypeProperties': %+v", err)
		}
		s.DeveloperToken = impl
	}

	if v, ok := temp["refreshToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RefreshToken' for 'GoogleAdWordsLinkedServiceTypeProperties': %+v", err)
		}
		s.RefreshToken = impl
	}

	return nil
}
