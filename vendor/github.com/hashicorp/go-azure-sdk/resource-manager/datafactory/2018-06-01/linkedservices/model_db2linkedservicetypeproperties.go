package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Db2LinkedServiceTypeProperties struct {
	AuthenticationType    *Db2AuthenticationType `json:"authenticationType,omitempty"`
	CertificateCommonName *interface{}           `json:"certificateCommonName,omitempty"`
	ConnectionString      *interface{}           `json:"connectionString,omitempty"`
	Database              *interface{}           `json:"database,omitempty"`
	EncryptedCredential   *string                `json:"encryptedCredential,omitempty"`
	PackageCollection     *interface{}           `json:"packageCollection,omitempty"`
	Password              SecretBase             `json:"password"`
	Server                *interface{}           `json:"server,omitempty"`
	Username              *interface{}           `json:"username,omitempty"`
}

var _ json.Unmarshaler = &Db2LinkedServiceTypeProperties{}

func (s *Db2LinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthenticationType    *Db2AuthenticationType `json:"authenticationType,omitempty"`
		CertificateCommonName *interface{}           `json:"certificateCommonName,omitempty"`
		ConnectionString      *interface{}           `json:"connectionString,omitempty"`
		Database              *interface{}           `json:"database,omitempty"`
		EncryptedCredential   *string                `json:"encryptedCredential,omitempty"`
		PackageCollection     *interface{}           `json:"packageCollection,omitempty"`
		Server                *interface{}           `json:"server,omitempty"`
		Username              *interface{}           `json:"username,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthenticationType = decoded.AuthenticationType
	s.CertificateCommonName = decoded.CertificateCommonName
	s.ConnectionString = decoded.ConnectionString
	s.Database = decoded.Database
	s.EncryptedCredential = decoded.EncryptedCredential
	s.PackageCollection = decoded.PackageCollection
	s.Server = decoded.Server
	s.Username = decoded.Username

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling Db2LinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["password"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Password' for 'Db2LinkedServiceTypeProperties': %+v", err)
		}
		s.Password = impl
	}

	return nil
}
