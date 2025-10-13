package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResponsysLinkedServiceTypeProperties struct {
	ClientId              interface{} `json:"clientId"`
	ClientSecret          SecretBase  `json:"clientSecret"`
	EncryptedCredential   *string     `json:"encryptedCredential,omitempty"`
	Endpoint              interface{} `json:"endpoint"`
	UseEncryptedEndpoints *bool       `json:"useEncryptedEndpoints,omitempty"`
	UseHostVerification   *bool       `json:"useHostVerification,omitempty"`
	UsePeerVerification   *bool       `json:"usePeerVerification,omitempty"`
}

var _ json.Unmarshaler = &ResponsysLinkedServiceTypeProperties{}

func (s *ResponsysLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ClientId              interface{} `json:"clientId"`
		EncryptedCredential   *string     `json:"encryptedCredential,omitempty"`
		Endpoint              interface{} `json:"endpoint"`
		UseEncryptedEndpoints *bool       `json:"useEncryptedEndpoints,omitempty"`
		UseHostVerification   *bool       `json:"useHostVerification,omitempty"`
		UsePeerVerification   *bool       `json:"usePeerVerification,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ClientId = decoded.ClientId
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Endpoint = decoded.Endpoint
	s.UseEncryptedEndpoints = decoded.UseEncryptedEndpoints
	s.UseHostVerification = decoded.UseHostVerification
	s.UsePeerVerification = decoded.UsePeerVerification

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ResponsysLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["clientSecret"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ClientSecret' for 'ResponsysLinkedServiceTypeProperties': %+v", err)
		}
		s.ClientSecret = impl
	}

	return nil
}
