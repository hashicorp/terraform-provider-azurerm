package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuickBooksLinkedServiceTypeProperties struct {
	AccessToken           SecretBase   `json:"accessToken"`
	AccessTokenSecret     SecretBase   `json:"accessTokenSecret"`
	CompanyId             *interface{} `json:"companyId,omitempty"`
	ConnectionProperties  *interface{} `json:"connectionProperties,omitempty"`
	ConsumerKey           *interface{} `json:"consumerKey,omitempty"`
	ConsumerSecret        SecretBase   `json:"consumerSecret"`
	EncryptedCredential   *string      `json:"encryptedCredential,omitempty"`
	Endpoint              *interface{} `json:"endpoint,omitempty"`
	UseEncryptedEndpoints *bool        `json:"useEncryptedEndpoints,omitempty"`
}

var _ json.Unmarshaler = &QuickBooksLinkedServiceTypeProperties{}

func (s *QuickBooksLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		CompanyId             *interface{} `json:"companyId,omitempty"`
		ConnectionProperties  *interface{} `json:"connectionProperties,omitempty"`
		ConsumerKey           *interface{} `json:"consumerKey,omitempty"`
		EncryptedCredential   *string      `json:"encryptedCredential,omitempty"`
		Endpoint              *interface{} `json:"endpoint,omitempty"`
		UseEncryptedEndpoints *bool        `json:"useEncryptedEndpoints,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.CompanyId = decoded.CompanyId
	s.ConnectionProperties = decoded.ConnectionProperties
	s.ConsumerKey = decoded.ConsumerKey
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Endpoint = decoded.Endpoint
	s.UseEncryptedEndpoints = decoded.UseEncryptedEndpoints

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling QuickBooksLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["accessToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AccessToken' for 'QuickBooksLinkedServiceTypeProperties': %+v", err)
		}
		s.AccessToken = impl
	}

	if v, ok := temp["accessTokenSecret"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AccessTokenSecret' for 'QuickBooksLinkedServiceTypeProperties': %+v", err)
		}
		s.AccessTokenSecret = impl
	}

	if v, ok := temp["consumerSecret"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ConsumerSecret' for 'QuickBooksLinkedServiceTypeProperties': %+v", err)
		}
		s.ConsumerSecret = impl
	}

	return nil
}
