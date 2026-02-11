package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatastoreCredentials = ServicePrincipalDatastoreCredentials{}

type ServicePrincipalDatastoreCredentials struct {
	AuthorityURL *string                          `json:"authorityUrl,omitempty"`
	ClientId     string                           `json:"clientId"`
	ResourceURL  *string                          `json:"resourceUrl,omitempty"`
	Secrets      ServicePrincipalDatastoreSecrets `json:"secrets"`
	TenantId     string                           `json:"tenantId"`

	// Fields inherited from DatastoreCredentials

	CredentialsType CredentialsType `json:"credentialsType"`
}

func (s ServicePrincipalDatastoreCredentials) DatastoreCredentials() BaseDatastoreCredentialsImpl {
	return BaseDatastoreCredentialsImpl{
		CredentialsType: s.CredentialsType,
	}
}

var _ json.Marshaler = ServicePrincipalDatastoreCredentials{}

func (s ServicePrincipalDatastoreCredentials) MarshalJSON() ([]byte, error) {
	type wrapper ServicePrincipalDatastoreCredentials
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServicePrincipalDatastoreCredentials: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServicePrincipalDatastoreCredentials: %+v", err)
	}

	decoded["credentialsType"] = "ServicePrincipal"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServicePrincipalDatastoreCredentials: %+v", err)
	}

	return encoded, nil
}
