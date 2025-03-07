package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EloquaLinkedServiceTypeProperties struct {
	EncryptedCredential   *string     `json:"encryptedCredential,omitempty"`
	Endpoint              interface{} `json:"endpoint"`
	Password              SecretBase  `json:"password"`
	UseEncryptedEndpoints *bool       `json:"useEncryptedEndpoints,omitempty"`
	UseHostVerification   *bool       `json:"useHostVerification,omitempty"`
	UsePeerVerification   *bool       `json:"usePeerVerification,omitempty"`
	Username              interface{} `json:"username"`
}

var _ json.Unmarshaler = &EloquaLinkedServiceTypeProperties{}

func (s *EloquaLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		EncryptedCredential   *string     `json:"encryptedCredential,omitempty"`
		Endpoint              interface{} `json:"endpoint"`
		UseEncryptedEndpoints *bool       `json:"useEncryptedEndpoints,omitempty"`
		UseHostVerification   *bool       `json:"useHostVerification,omitempty"`
		UsePeerVerification   *bool       `json:"usePeerVerification,omitempty"`
		Username              interface{} `json:"username"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.EncryptedCredential = decoded.EncryptedCredential
	s.Endpoint = decoded.Endpoint
	s.UseEncryptedEndpoints = decoded.UseEncryptedEndpoints
	s.UseHostVerification = decoded.UseHostVerification
	s.UsePeerVerification = decoded.UsePeerVerification
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling EloquaLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'EloquaLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
