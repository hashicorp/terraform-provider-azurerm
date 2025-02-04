package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrestoLinkedServiceTypeProperties struct {
	AllowHostNameCNMismatch   *bool                    `json:"allowHostNameCNMismatch,omitempty"`
	AllowSelfSignedServerCert *bool                    `json:"allowSelfSignedServerCert,omitempty"`
	AuthenticationType        PrestoAuthenticationType `json:"authenticationType"`
	Catalog                   string                   `json:"catalog"`
	EnableSsl                 *bool                    `json:"enableSsl,omitempty"`
	EncryptedCredential       *string                  `json:"encryptedCredential,omitempty"`
	Host                      string                   `json:"host"`
	Password                  SecretBase               `json:"password"`
	Port                      *int64                   `json:"port,omitempty"`
	ServerVersion             string                   `json:"serverVersion"`
	TimeZoneID                *string                  `json:"timeZoneID,omitempty"`
	TrustedCertPath           *string                  `json:"trustedCertPath,omitempty"`
	UseSystemTrustStore       *bool                    `json:"useSystemTrustStore,omitempty"`
	Username                  *string                  `json:"username,omitempty"`
}

var _ json.Unmarshaler = &PrestoLinkedServiceTypeProperties{}

func (s *PrestoLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AllowHostNameCNMismatch   *bool                    `json:"allowHostNameCNMismatch,omitempty"`
		AllowSelfSignedServerCert *bool                    `json:"allowSelfSignedServerCert,omitempty"`
		AuthenticationType        PrestoAuthenticationType `json:"authenticationType"`
		Catalog                   string                   `json:"catalog"`
		EnableSsl                 *bool                    `json:"enableSsl,omitempty"`
		EncryptedCredential       *string                  `json:"encryptedCredential,omitempty"`
		Host                      string                   `json:"host"`
		Port                      *int64                   `json:"port,omitempty"`
		ServerVersion             string                   `json:"serverVersion"`
		TimeZoneID                *string                  `json:"timeZoneID,omitempty"`
		TrustedCertPath           *string                  `json:"trustedCertPath,omitempty"`
		UseSystemTrustStore       *bool                    `json:"useSystemTrustStore,omitempty"`
		Username                  *string                  `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AllowHostNameCNMismatch = decoded.AllowHostNameCNMismatch
	s.AllowSelfSignedServerCert = decoded.AllowSelfSignedServerCert
	s.AuthenticationType = decoded.AuthenticationType
	s.Catalog = decoded.Catalog
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
