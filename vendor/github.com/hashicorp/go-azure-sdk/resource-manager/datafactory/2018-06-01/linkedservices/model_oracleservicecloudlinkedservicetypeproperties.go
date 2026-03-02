package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OracleServiceCloudLinkedServiceTypeProperties struct {
	EncryptedCredential   *string     `json:"encryptedCredential,omitempty"`
	Host                  interface{} `json:"host"`
	Password              SecretBase  `json:"password"`
	UseEncryptedEndpoints *bool       `json:"useEncryptedEndpoints,omitempty"`
	UseHostVerification   *bool       `json:"useHostVerification,omitempty"`
	UsePeerVerification   *bool       `json:"usePeerVerification,omitempty"`
	Username              interface{} `json:"username"`
}

var _ json.Unmarshaler = &OracleServiceCloudLinkedServiceTypeProperties{}

func (s *OracleServiceCloudLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		EncryptedCredential   *string     `json:"encryptedCredential,omitempty"`
		Host                  interface{} `json:"host"`
		UseEncryptedEndpoints *bool       `json:"useEncryptedEndpoints,omitempty"`
		UseHostVerification   *bool       `json:"useHostVerification,omitempty"`
		UsePeerVerification   *bool       `json:"usePeerVerification,omitempty"`
		Username              interface{} `json:"username"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.EncryptedCredential = decoded.EncryptedCredential
	s.Host = decoded.Host
	s.UseEncryptedEndpoints = decoded.UseEncryptedEndpoints
	s.UseHostVerification = decoded.UseHostVerification
	s.UsePeerVerification = decoded.UsePeerVerification
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling OracleServiceCloudLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'OracleServiceCloudLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
