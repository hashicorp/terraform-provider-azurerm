package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatastoreCredentials = CertificateDatastoreCredentials{}

type CertificateDatastoreCredentials struct {
	AuthorityUrl *string          `json:"authorityUrl,omitempty"`
	ClientId     string           `json:"clientId"`
	ResourceUrl  *string          `json:"resourceUrl,omitempty"`
	Secrets      DatastoreSecrets `json:"secrets"`
	TenantId     string           `json:"tenantId"`
	Thumbprint   string           `json:"thumbprint"`

	// Fields inherited from DatastoreCredentials
}

var _ json.Marshaler = CertificateDatastoreCredentials{}

func (s CertificateDatastoreCredentials) MarshalJSON() ([]byte, error) {
	type wrapper CertificateDatastoreCredentials
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CertificateDatastoreCredentials: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CertificateDatastoreCredentials: %+v", err)
	}
	decoded["credentialsType"] = "Certificate"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CertificateDatastoreCredentials: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &CertificateDatastoreCredentials{}

func (s *CertificateDatastoreCredentials) UnmarshalJSON(bytes []byte) error {
	type alias CertificateDatastoreCredentials
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into CertificateDatastoreCredentials: %+v", err)
	}

	s.AuthorityUrl = decoded.AuthorityUrl
	s.ClientId = decoded.ClientId
	s.ResourceUrl = decoded.ResourceUrl
	s.TenantId = decoded.TenantId
	s.Thumbprint = decoded.Thumbprint

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CertificateDatastoreCredentials into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["secrets"]; ok {
		impl, err := unmarshalDatastoreSecretsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Secrets' for 'CertificateDatastoreCredentials': %+v", err)
		}
		s.Secrets = impl
	}
	return nil
}
