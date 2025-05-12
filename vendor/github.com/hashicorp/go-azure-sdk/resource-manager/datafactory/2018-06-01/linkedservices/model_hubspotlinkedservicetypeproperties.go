package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HubspotLinkedServiceTypeProperties struct {
	AccessToken           SecretBase  `json:"accessToken"`
	ClientId              interface{} `json:"clientId"`
	ClientSecret          SecretBase  `json:"clientSecret"`
	EncryptedCredential   *string     `json:"encryptedCredential,omitempty"`
	RefreshToken          SecretBase  `json:"refreshToken"`
	UseEncryptedEndpoints *bool       `json:"useEncryptedEndpoints,omitempty"`
	UseHostVerification   *bool       `json:"useHostVerification,omitempty"`
	UsePeerVerification   *bool       `json:"usePeerVerification,omitempty"`
}

var _ json.Unmarshaler = &HubspotLinkedServiceTypeProperties{}

func (s *HubspotLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ClientId              interface{} `json:"clientId"`
		EncryptedCredential   *string     `json:"encryptedCredential,omitempty"`
		UseEncryptedEndpoints *bool       `json:"useEncryptedEndpoints,omitempty"`
		UseHostVerification   *bool       `json:"useHostVerification,omitempty"`
		UsePeerVerification   *bool       `json:"usePeerVerification,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ClientId = decoded.ClientId
	s.EncryptedCredential = decoded.EncryptedCredential
	s.UseEncryptedEndpoints = decoded.UseEncryptedEndpoints
	s.UseHostVerification = decoded.UseHostVerification
	s.UsePeerVerification = decoded.UsePeerVerification

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling HubspotLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["accessToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AccessToken' for 'HubspotLinkedServiceTypeProperties': %+v", err)
		}
		s.AccessToken = impl
	}

	if v, ok := temp["clientSecret"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ClientSecret' for 'HubspotLinkedServiceTypeProperties': %+v", err)
		}
		s.ClientSecret = impl
	}

	if v, ok := temp["refreshToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RefreshToken' for 'HubspotLinkedServiceTypeProperties': %+v", err)
		}
		s.RefreshToken = impl
	}

	return nil
}
