package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrestoLinkedServiceTypeProperties struct {
	AllowHostNameCNMismatch           *bool                    `json:"allowHostNameCNMismatch,omitempty"`
	AllowSelfSignedServerCert         *bool                    `json:"allowSelfSignedServerCert,omitempty"`
	AuthenticationType                PrestoAuthenticationType `json:"authenticationType"`
	Catalog                           interface{}              `json:"catalog"`
	EnableServerCertificateValidation *bool                    `json:"enableServerCertificateValidation,omitempty"`
	EnableSsl                         *bool                    `json:"enableSsl,omitempty"`
	EncryptedCredential               *string                  `json:"encryptedCredential,omitempty"`
	Host                              interface{}              `json:"host"`
	Password                          SecretBase               `json:"password"`
	Port                              *int64                   `json:"port,omitempty"`
	ServerVersion                     *interface{}             `json:"serverVersion,omitempty"`
	TimeZoneID                        *interface{}             `json:"timeZoneID,omitempty"`
	TrustedCertPath                   *interface{}             `json:"trustedCertPath,omitempty"`
	UseSystemTrustStore               *bool                    `json:"useSystemTrustStore,omitempty"`
	Username                          *interface{}             `json:"username,omitempty"`
}

var _ json.Unmarshaler = &PrestoLinkedServiceTypeProperties{}

func (s *PrestoLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AllowHostNameCNMismatch           *bool                    `json:"allowHostNameCNMismatch,omitempty"`
		AllowSelfSignedServerCert         *bool                    `json:"allowSelfSignedServerCert,omitempty"`
		AuthenticationType                PrestoAuthenticationType `json:"authenticationType"`
		Catalog                           interface{}              `json:"catalog"`
		EnableServerCertificateValidation *bool                    `json:"enableServerCertificateValidation,omitempty"`
		EnableSsl                         *bool                    `json:"enableSsl,omitempty"`
		EncryptedCredential               *string                  `json:"encryptedCredential,omitempty"`
		Host                              interface{}              `json:"host"`
		Port                              *int64                   `json:"port,omitempty"`
		ServerVersion                     *interface{}             `json:"serverVersion,omitempty"`
		TimeZoneID                        *interface{}             `json:"timeZoneID,omitempty"`
		TrustedCertPath                   *interface{}             `json:"trustedCertPath,omitempty"`
		UseSystemTrustStore               *bool                    `json:"useSystemTrustStore,omitempty"`
		Username                          *interface{}             `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AllowHostNameCNMismatch = decoded.AllowHostNameCNMismatch
	s.AllowSelfSignedServerCert = decoded.AllowSelfSignedServerCert
	s.AuthenticationType = decoded.AuthenticationType
	s.Catalog = decoded.Catalog
	s.EnableServerCertificateValidation = decoded.EnableServerCertificateValidation
	s.EnableSsl = decoded.EnableSsl
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Host = decoded.Host
	s.Port = decoded.Port
	s.ServerVersion = decoded.ServerVersion
	s.TimeZoneID = decoded.TimeZoneID
	s.TrustedCertPath = decoded.TrustedCertPath
	s.UseSystemTrustStore = decoded.UseSystemTrustStore
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling PrestoLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'PrestoLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
