package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SftpServerLinkedServiceTypeProperties struct {
	AuthenticationType    *SftpAuthenticationType `json:"authenticationType,omitempty"`
	EncryptedCredential   *string                 `json:"encryptedCredential,omitempty"`
	Host                  interface{}             `json:"host"`
	HostKeyFingerprint    *interface{}            `json:"hostKeyFingerprint,omitempty"`
	PassPhrase            SecretBase              `json:"passPhrase"`
	Password              SecretBase              `json:"password"`
	Port                  *int64                  `json:"port,omitempty"`
	PrivateKeyContent     SecretBase              `json:"privateKeyContent"`
	PrivateKeyPath        *interface{}            `json:"privateKeyPath,omitempty"`
	SkipHostKeyValidation *bool                   `json:"skipHostKeyValidation,omitempty"`
	UserName              *interface{}            `json:"userName,omitempty"`
}

var _ json.Unmarshaler = &SftpServerLinkedServiceTypeProperties{}

func (s *SftpServerLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType    *SftpAuthenticationType `json:"authenticationType,omitempty"`
		EncryptedCredential   *string                 `json:"encryptedCredential,omitempty"`
		Host                  interface{}             `json:"host"`
		HostKeyFingerprint    *interface{}            `json:"hostKeyFingerprint,omitempty"`
		Port                  *int64                  `json:"port,omitempty"`
		PrivateKeyPath        *interface{}            `json:"privateKeyPath,omitempty"`
		SkipHostKeyValidation *bool                   `json:"skipHostKeyValidation,omitempty"`
		UserName              *interface{}            `json:"userName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Host = decoded.Host
	s.HostKeyFingerprint = decoded.HostKeyFingerprint
	s.Port = decoded.Port
	s.PrivateKeyPath = decoded.PrivateKeyPath
	s.SkipHostKeyValidation = decoded.SkipHostKeyValidation
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SftpServerLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["passPhrase"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'PassPhrase' for 'SftpServerLinkedServiceTypeProperties': %+v", err)
		}
		s.PassPhrase = impl
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'SftpServerLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	if v, ok := temp["privateKeyContent"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'PrivateKeyContent' for 'SftpServerLinkedServiceTypeProperties': %+v", err)
		}
		s.PrivateKeyContent = impl
	}

	return nil
}
