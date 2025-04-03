package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureSqlMILinkedServiceTypeProperties struct {
	AlwaysEncryptedSettings        *SqlAlwaysEncryptedProperties `json:"alwaysEncryptedSettings,omitempty"`
	ApplicationIntent              *interface{}                  `json:"applicationIntent,omitempty"`
	AuthenticationType             *AzureSqlMIAuthenticationType `json:"authenticationType,omitempty"`
	AzureCloudType                 *interface{}                  `json:"azureCloudType,omitempty"`
	CommandTimeout                 *int64                        `json:"commandTimeout,omitempty"`
	ConnectRetryCount              *int64                        `json:"connectRetryCount,omitempty"`
	ConnectRetryInterval           *int64                        `json:"connectRetryInterval,omitempty"`
	ConnectTimeout                 *int64                        `json:"connectTimeout,omitempty"`
	ConnectionString               *interface{}                  `json:"connectionString,omitempty"`
	Credential                     *CredentialReference          `json:"credential,omitempty"`
	Database                       *interface{}                  `json:"database,omitempty"`
	Encrypt                        *interface{}                  `json:"encrypt,omitempty"`
	EncryptedCredential            *string                       `json:"encryptedCredential,omitempty"`
	FailoverPartner                *interface{}                  `json:"failoverPartner,omitempty"`
	HostNameInCertificate          *interface{}                  `json:"hostNameInCertificate,omitempty"`
	IntegratedSecurity             *bool                         `json:"integratedSecurity,omitempty"`
	LoadBalanceTimeout             *int64                        `json:"loadBalanceTimeout,omitempty"`
	MaxPoolSize                    *int64                        `json:"maxPoolSize,omitempty"`
	MinPoolSize                    *int64                        `json:"minPoolSize,omitempty"`
	MultiSubnetFailover            *bool                         `json:"multiSubnetFailover,omitempty"`
	MultipleActiveResultSets       *bool                         `json:"multipleActiveResultSets,omitempty"`
	PacketSize                     *int64                        `json:"packetSize,omitempty"`
	Password                       *AzureKeyVaultSecretReference `json:"password,omitempty"`
	Pooling                        *bool                         `json:"pooling,omitempty"`
	Server                         *interface{}                  `json:"server,omitempty"`
	ServicePrincipalCredential     SecretBase                    `json:"servicePrincipalCredential"`
	ServicePrincipalCredentialType *interface{}                  `json:"servicePrincipalCredentialType,omitempty"`
	ServicePrincipalId             *interface{}                  `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey            SecretBase                    `json:"servicePrincipalKey"`
	Tenant                         *interface{}                  `json:"tenant,omitempty"`
	TrustServerCertificate         *bool                         `json:"trustServerCertificate,omitempty"`
	UserName                       *interface{}                  `json:"userName,omitempty"`
}

var _ json.Unmarshaler = &AzureSqlMILinkedServiceTypeProperties{}

func (s *AzureSqlMILinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AlwaysEncryptedSettings        *SqlAlwaysEncryptedProperties `json:"alwaysEncryptedSettings,omitempty"`
		ApplicationIntent              *interface{}                  `json:"applicationIntent,omitempty"`
		AuthenticationType             *AzureSqlMIAuthenticationType `json:"authenticationType,omitempty"`
		AzureCloudType                 *interface{}                  `json:"azureCloudType,omitempty"`
		CommandTimeout                 *int64                        `json:"commandTimeout,omitempty"`
		ConnectRetryCount              *int64                        `json:"connectRetryCount,omitempty"`
		ConnectRetryInterval           *int64                        `json:"connectRetryInterval,omitempty"`
		ConnectTimeout                 *int64                        `json:"connectTimeout,omitempty"`
		ConnectionString               *interface{}                  `json:"connectionString,omitempty"`
		Credential                     *CredentialReference          `json:"credential,omitempty"`
		Database                       *interface{}                  `json:"database,omitempty"`
		Encrypt                        *interface{}                  `json:"encrypt,omitempty"`
		EncryptedCredential            *string                       `json:"encryptedCredential,omitempty"`
		FailoverPartner                *interface{}                  `json:"failoverPartner,omitempty"`
		HostNameInCertificate          *interface{}                  `json:"hostNameInCertificate,omitempty"`
		IntegratedSecurity             *bool                         `json:"integratedSecurity,omitempty"`
		LoadBalanceTimeout             *int64                        `json:"loadBalanceTimeout,omitempty"`
		MaxPoolSize                    *int64                        `json:"maxPoolSize,omitempty"`
		MinPoolSize                    *int64                        `json:"minPoolSize,omitempty"`
		MultiSubnetFailover            *bool                         `json:"multiSubnetFailover,omitempty"`
		MultipleActiveResultSets       *bool                         `json:"multipleActiveResultSets,omitempty"`
		PacketSize                     *int64                        `json:"packetSize,omitempty"`
		Password                       *AzureKeyVaultSecretReference `json:"password,omitempty"`
		Pooling                        *bool                         `json:"pooling,omitempty"`
		Server                         *interface{}                  `json:"server,omitempty"`
		ServicePrincipalCredentialType *interface{}                  `json:"servicePrincipalCredentialType,omitempty"`
		ServicePrincipalId             *interface{}                  `json:"servicePrincipalId,omitempty"`
		Tenant                         *interface{}                  `json:"tenant,omitempty"`
		TrustServerCertificate         *bool                         `json:"trustServerCertificate,omitempty"`
		UserName                       *interface{}                  `json:"userName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AlwaysEncryptedSettings = decoded.AlwaysEncryptedSettings
	s.ApplicationIntent = decoded.ApplicationIntent
	s.AuthenticationType = decoded.AuthenticationType
	s.AzureCloudType = decoded.AzureCloudType
	s.CommandTimeout = decoded.CommandTimeout
	s.ConnectRetryCount = decoded.ConnectRetryCount
	s.ConnectRetryInterval = decoded.ConnectRetryInterval
	s.ConnectTimeout = decoded.ConnectTimeout
	s.ConnectionString = decoded.ConnectionString
	s.Credential = decoded.Credential
	s.Database = decoded.Database
	s.Encrypt = decoded.Encrypt
	s.EncryptedCredential = decoded.EncryptedCredential
	s.FailoverPartner = decoded.FailoverPartner
	s.HostNameInCertificate = decoded.HostNameInCertificate
	s.IntegratedSecurity = decoded.IntegratedSecurity
	s.LoadBalanceTimeout = decoded.LoadBalanceTimeout
	s.MaxPoolSize = decoded.MaxPoolSize
	s.MinPoolSize = decoded.MinPoolSize
	s.MultiSubnetFailover = decoded.MultiSubnetFailover
	s.MultipleActiveResultSets = decoded.MultipleActiveResultSets
	s.PacketSize = decoded.PacketSize
	s.Password = decoded.Password
	s.Pooling = decoded.Pooling
	s.Server = decoded.Server
	s.ServicePrincipalCredentialType = decoded.ServicePrincipalCredentialType
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.Tenant = decoded.Tenant
	s.TrustServerCertificate = decoded.TrustServerCertificate
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureSqlMILinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["servicePrincipalCredential"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalCredential' for 'AzureSqlMILinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalCredential = impl
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'AzureSqlMILinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
