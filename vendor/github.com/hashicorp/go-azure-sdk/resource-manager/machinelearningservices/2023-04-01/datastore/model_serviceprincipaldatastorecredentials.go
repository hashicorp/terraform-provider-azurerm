package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatastoreCredentials = ServicePrincipalDatastoreCredentials{}

type ServicePrincipalDatastoreCredentials struct {
	AuthorityUrl *string          `json:"authorityUrl,omitempty"`
	ClientId     string           `json:"clientId"`
	ResourceUrl  *string          `json:"resourceUrl,omitempty"`
	Secrets      DatastoreSecrets `json:"secrets"`
	TenantId     string           `json:"tenantId"`

	// Fields inherited from DatastoreCredentials
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
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServicePrincipalDatastoreCredentials: %+v", err)
	}
	decoded["credentialsType"] = "ServicePrincipal"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServicePrincipalDatastoreCredentials: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &ServicePrincipalDatastoreCredentials{}

func (s *ServicePrincipalDatastoreCredentials) UnmarshalJSON(bytes []byte) error {
	type alias ServicePrincipalDatastoreCredentials
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ServicePrincipalDatastoreCredentials: %+v", err)
	}

	s.AuthorityUrl = decoded.AuthorityUrl
	s.ClientId = decoded.ClientId
	s.ResourceUrl = decoded.ResourceUrl
	s.TenantId = decoded.TenantId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ServicePrincipalDatastoreCredentials into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["secrets"]; ok {
		impl, err := unmarshalDatastoreSecretsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Secrets' for 'ServicePrincipalDatastoreCredentials': %+v", err)
		}
		s.Secrets = impl
	}
	return nil
}
