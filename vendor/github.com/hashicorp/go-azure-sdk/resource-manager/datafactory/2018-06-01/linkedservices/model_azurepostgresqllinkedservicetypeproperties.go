package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzurePostgreSqlLinkedServiceTypeProperties struct {
	AzureCloudType                       *interface{}                  `json:"azureCloudType,omitempty"`
	CommandTimeout                       *int64                        `json:"commandTimeout,omitempty"`
	ConnectionString                     *interface{}                  `json:"connectionString,omitempty"`
	Credential                           *CredentialReference          `json:"credential,omitempty"`
	Database                             *interface{}                  `json:"database,omitempty"`
	Encoding                             *interface{}                  `json:"encoding,omitempty"`
	EncryptedCredential                  *string                       `json:"encryptedCredential,omitempty"`
	Password                             *AzureKeyVaultSecretReference `json:"password,omitempty"`
	Port                                 *int64                        `json:"port,omitempty"`
	ReadBufferSize                       *int64                        `json:"readBufferSize,omitempty"`
	Server                               *interface{}                  `json:"server,omitempty"`
	ServicePrincipalCredentialType       *interface{}                  `json:"servicePrincipalCredentialType,omitempty"`
	ServicePrincipalEmbeddedCert         SecretBase                    `json:"servicePrincipalEmbeddedCert"`
	ServicePrincipalEmbeddedCertPassword SecretBase                    `json:"servicePrincipalEmbeddedCertPassword"`
	ServicePrincipalId                   *interface{}                  `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey                  SecretBase                    `json:"servicePrincipalKey"`
	SslMode                              *int64                        `json:"sslMode,omitempty"`
	Tenant                               *interface{}                  `json:"tenant,omitempty"`
	Timeout                              *int64                        `json:"timeout,omitempty"`
	Timezone                             *interface{}                  `json:"timezone,omitempty"`
	TrustServerCertificate               *bool                         `json:"trustServerCertificate,omitempty"`
	Username                             *interface{}                  `json:"username,omitempty"`
}

var _ json.Unmarshaler = &AzurePostgreSqlLinkedServiceTypeProperties{}

func (s *AzurePostgreSqlLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AzureCloudType                 *interface{}                  `json:"azureCloudType,omitempty"`
		CommandTimeout                 *int64                        `json:"commandTimeout,omitempty"`
		ConnectionString               *interface{}                  `json:"connectionString,omitempty"`
		Credential                     *CredentialReference          `json:"credential,omitempty"`
		Database                       *interface{}                  `json:"database,omitempty"`
		Encoding                       *interface{}                  `json:"encoding,omitempty"`
		EncryptedCredential            *string                       `json:"encryptedCredential,omitempty"`
		Password                       *AzureKeyVaultSecretReference `json:"password,omitempty"`
		Port                           *int64                        `json:"port,omitempty"`
		ReadBufferSize                 *int64                        `json:"readBufferSize,omitempty"`
		Server                         *interface{}                  `json:"server,omitempty"`
		ServicePrincipalCredentialType *interface{}                  `json:"servicePrincipalCredentialType,omitempty"`
		ServicePrincipalId             *interface{}                  `json:"servicePrincipalId,omitempty"`
		SslMode                        *int64                        `json:"sslMode,omitempty"`
		Tenant                         *interface{}                  `json:"tenant,omitempty"`
		Timeout                        *int64                        `json:"timeout,omitempty"`
		Timezone                       *interface{}                  `json:"timezone,omitempty"`
		TrustServerCertificate         *bool                         `json:"trustServerCertificate,omitempty"`
		Username                       *interface{}                  `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AzureCloudType = decoded.AzureCloudType
	s.CommandTimeout = decoded.CommandTimeout
	s.ConnectionString = decoded.ConnectionString
	s.Credential = decoded.Credential
	s.Database = decoded.Database
	s.Encoding = decoded.Encoding
	s.EncryptedCredential = decoded.EncryptedCredential
	s.Password = decoded.Password
	s.Port = decoded.Port
	s.ReadBufferSize = decoded.ReadBufferSize
	s.Server = decoded.Server
	s.ServicePrincipalCredentialType = decoded.ServicePrincipalCredentialType
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.SslMode = decoded.SslMode
	s.Tenant = decoded.Tenant
	s.Timeout = decoded.Timeout
	s.Timezone = decoded.Timezone
	s.TrustServerCertificate = decoded.TrustServerCertificate
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzurePostgreSqlLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["servicePrincipalEmbeddedCert"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalEmbeddedCert' for 'AzurePostgreSqlLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalEmbeddedCert = impl
	}

	if v, ok := temp["servicePrincipalEmbeddedCertPassword"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalEmbeddedCertPassword' for 'AzurePostgreSqlLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalEmbeddedCertPassword = impl
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'AzurePostgreSqlLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
