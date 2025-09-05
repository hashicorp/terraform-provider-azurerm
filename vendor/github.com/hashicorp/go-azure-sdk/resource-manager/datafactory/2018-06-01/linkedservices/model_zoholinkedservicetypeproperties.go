package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ZohoLinkedServiceTypeProperties struct {
	AccessToken           SecretBase   `json:"accessToken"`
	ConnectionProperties  *interface{} `json:"connectionProperties,omitempty"`
	EncryptedCredential   *string      `json:"encryptedCredential,omitempty"`
	Endpoint              *interface{} `json:"endpoint,omitempty"`
	UseEncryptedEndpoints *bool        `json:"useEncryptedEndpoints,omitempty"`
	UseHostVerification   *bool        `json:"useHostVerification,omitempty"`
	UsePeerVerification   *bool        `json:"usePeerVerification,omitempty"`
}

var _ json.Unmarshaler = &ZohoLinkedServiceTypeProperties{}

func (s *ZohoLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ConnectionProperties  *interface{} `json:"connectionProperties,omitempty"`
		EncryptedCredential   *string      `json:"encryptedCredential,omitempty"`
		Endpoint              *interface{} `json:"endpoint,omitempty"`
		UseEncryptedEndpoints *bool        `json:"useEncryptedEndpoints,omitempty"`
		UseHostVerification   *bool        `json:"useHostVerification,omitempty"`
		UsePeerVerification   *bool        `json:"usePeerVerification,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ConnectionProperties = decoded.ConnectionProperties
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Endpoint = decoded.Endpoint
	s.UseEncryptedEndpoints = decoded.UseEncryptedEndpoints
	s.UseHostVerification = decoded.UseHostVerification
	s.UsePeerVerification = decoded.UsePeerVerification

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ZohoLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["accessToken"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AccessToken' for 'ZohoLinkedServiceTypeProperties': %+v", err)
		}
		s.AccessToken = impl
	}

	return nil
}
