package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlServerLinkedServiceTypeProperties struct {
	AlwaysEncryptedSettings  *SqlAlwaysEncryptedProperties `json:"alwaysEncryptedSettings,omitempty"`
	ApplicationIntent        *interface{}                  `json:"applicationIntent,omitempty"`
	AuthenticationType       *SqlServerAuthenticationType  `json:"authenticationType,omitempty"`
	CommandTimeout           *interface{}                  `json:"commandTimeout,omitempty"`
	ConnectRetryCount        *interface{}                  `json:"connectRetryCount,omitempty"`
	ConnectRetryInterval     *interface{}                  `json:"connectRetryInterval,omitempty"`
	ConnectTimeout           *interface{}                  `json:"connectTimeout,omitempty"`
	ConnectionString         *interface{}                  `json:"connectionString,omitempty"`
	Credential               *CredentialReference          `json:"credential,omitempty"`
	Database                 *interface{}                  `json:"database,omitempty"`
	Encrypt                  *interface{}                  `json:"encrypt,omitempty"`
	EncryptedCredential      *string                       `json:"encryptedCredential,omitempty"`
	FailoverPartner          *interface{}                  `json:"failoverPartner,omitempty"`
	HostNameInCertificate    *interface{}                  `json:"hostNameInCertificate,omitempty"`
	IntegratedSecurity       *interface{}                  `json:"integratedSecurity,omitempty"`
	LoadBalanceTimeout       *interface{}                  `json:"loadBalanceTimeout,omitempty"`
	MaxPoolSize              *interface{}                  `json:"maxPoolSize,omitempty"`
	MinPoolSize              *interface{}                  `json:"minPoolSize,omitempty"`
	MultiSubnetFailover      *interface{}                  `json:"multiSubnetFailover,omitempty"`
	MultipleActiveResultSets *interface{}                  `json:"multipleActiveResultSets,omitempty"`
	PacketSize               *interface{}                  `json:"packetSize,omitempty"`
	Password                 SecretBase                    `json:"password"`
	Pooling                  *interface{}                  `json:"pooling,omitempty"`
	Server                   *interface{}                  `json:"server,omitempty"`
	TrustServerCertificate   *interface{}                  `json:"trustServerCertificate,omitempty"`
	UserName                 *interface{}                  `json:"userName,omitempty"`
}

var _ json.Unmarshaler = &SqlServerLinkedServiceTypeProperties{}

func (s *SqlServerLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AlwaysEncryptedSettings  *SqlAlwaysEncryptedProperties `json:"alwaysEncryptedSettings,omitempty"`
		ApplicationIntent        *interface{}                  `json:"applicationIntent,omitempty"`
		AuthenticationType       *SqlServerAuthenticationType  `json:"authenticationType,omitempty"`
		CommandTimeout           *interface{}                  `json:"commandTimeout,omitempty"`
		ConnectRetryCount        *interface{}                  `json:"connectRetryCount,omitempty"`
		ConnectRetryInterval     *interface{}                  `json:"connectRetryInterval,omitempty"`
		ConnectTimeout           *interface{}                  `json:"connectTimeout,omitempty"`
		ConnectionString         *interface{}                  `json:"connectionString,omitempty"`
		Credential               *CredentialReference          `json:"credential,omitempty"`
		Database                 *interface{}                  `json:"database,omitempty"`
		Encrypt                  *interface{}                  `json:"encrypt,omitempty"`
		EncryptedCredential      *string                       `json:"encryptedCredential,omitempty"`
		FailoverPartner          *interface{}                  `json:"failoverPartner,omitempty"`
		HostNameInCertificate    *interface{}                  `json:"hostNameInCertificate,omitempty"`
		IntegratedSecurity       *interface{}                  `json:"integratedSecurity,omitempty"`
		LoadBalanceTimeout       *interface{}                  `json:"loadBalanceTimeout,omitempty"`
		MaxPoolSize              *interface{}                  `json:"maxPoolSize,omitempty"`
		MinPoolSize              *interface{}                  `json:"minPoolSize,omitempty"`
		MultiSubnetFailover      *interface{}                  `json:"multiSubnetFailover,omitempty"`
		MultipleActiveResultSets *interface{}                  `json:"multipleActiveResultSets,omitempty"`
		PacketSize               *interface{}                  `json:"packetSize,omitempty"`
		Pooling                  *interface{}                  `json:"pooling,omitempty"`
		Server                   *interface{}                  `json:"server,omitempty"`
		TrustServerCertificate   *interface{}                  `json:"trustServerCertificate,omitempty"`
		UserName                 *interface{}                  `json:"userName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AlwaysEncryptedSettings = decoded.AlwaysEncryptedSettings
	s.ApplicationIntent = decoded.ApplicationIntent
	s.AuthenticationType = decoded.AuthenticationType
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
	s.Pooling = decoded.Pooling
	s.Server = decoded.Server
	s.TrustServerCertificate = decoded.TrustServerCertificate
	s.UserName = decoded.UserName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SqlServerLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'SqlServerLinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
