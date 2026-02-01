package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatastoreCredentials = CertificateDatastoreCredentials{}

type CertificateDatastoreCredentials struct {
	AuthorityURL *string                     `json:"authorityUrl,omitempty"`
	ClientId     string                      `json:"clientId"`
	ResourceURL  *string                     `json:"resourceUrl,omitempty"`
	Secrets      CertificateDatastoreSecrets `json:"secrets"`
	TenantId     string                      `json:"tenantId"`
	Thumbprint   string                      `json:"thumbprint"`

	// Fields inherited from DatastoreCredentials

	CredentialsType CredentialsType `json:"credentialsType"`
}

func (s CertificateDatastoreCredentials) DatastoreCredentials() BaseDatastoreCredentialsImpl {
	return BaseDatastoreCredentialsImpl{
		CredentialsType: s.CredentialsType,
	}
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
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CertificateDatastoreCredentials: %+v", err)
	}

	decoded["credentialsType"] = "Certificate"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CertificateDatastoreCredentials: %+v", err)
	}

	return encoded, nil
}
