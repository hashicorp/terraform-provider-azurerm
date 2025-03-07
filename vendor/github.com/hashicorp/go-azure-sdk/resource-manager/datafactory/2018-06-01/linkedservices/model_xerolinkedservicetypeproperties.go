package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type XeroLinkedServiceTypeProperties struct {
	ConnectionProperties  *interface{} `json:"connectionProperties,omitempty"`
	ConsumerKey           SecretBase   `json:"consumerKey"`
	EncryptedCredential   *string      `json:"encryptedCredential,omitempty"`
	Host                  *interface{} `json:"host,omitempty"`
	PrivateKey            SecretBase   `json:"privateKey"`
	UseEncryptedEndpoints *bool        `json:"useEncryptedEndpoints,omitempty"`
	UseHostVerification   *bool        `json:"useHostVerification,omitempty"`
	UsePeerVerification   *bool        `json:"usePeerVerification,omitempty"`
}

var _ json.Unmarshaler = &XeroLinkedServiceTypeProperties{}

func (s *XeroLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ConnectionProperties  *interface{} `json:"connectionProperties,omitempty"`
		EncryptedCredential   *string      `json:"encryptedCredential,omitempty"`
		Host                  *interface{} `json:"host,omitempty"`
		UseEncryptedEndpoints *bool        `json:"useEncryptedEndpoints,omitempty"`
		UseHostVerification   *bool        `json:"useHostVerification,omitempty"`
		UsePeerVerification   *bool        `json:"usePeerVerification,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ConnectionProperties = decoded.ConnectionProperties
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Host = decoded.Host
	s.UseEncryptedEndpoints = decoded.UseEncryptedEndpoints
	s.UseHostVerification = decoded.UseHostVerification
	s.UsePeerVerification = decoded.UsePeerVerification

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling XeroLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["consumerKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ConsumerKey' for 'XeroLinkedServiceTypeProperties': %+v", err)
		}
		s.ConsumerKey = impl
	}

	if v, ok := temp["privateKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'PrivateKey' for 'XeroLinkedServiceTypeProperties': %+v", err)
		}
		s.PrivateKey = impl
	}

	return nil
}
