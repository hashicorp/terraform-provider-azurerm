package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatastoreSecrets = ServicePrincipalDatastoreSecrets{}

type ServicePrincipalDatastoreSecrets struct {
	ClientSecret *string `json:"clientSecret,omitempty"`

	// Fields inherited from DatastoreSecrets

	SecretsType SecretsType `json:"secretsType"`
}

func (s ServicePrincipalDatastoreSecrets) DatastoreSecrets() BaseDatastoreSecretsImpl {
	return BaseDatastoreSecretsImpl{
		SecretsType: s.SecretsType,
	}
}

var _ json.Marshaler = ServicePrincipalDatastoreSecrets{}

func (s ServicePrincipalDatastoreSecrets) MarshalJSON() ([]byte, error) {
	type wrapper ServicePrincipalDatastoreSecrets
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServicePrincipalDatastoreSecrets: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServicePrincipalDatastoreSecrets: %+v", err)
	}

	decoded["secretsType"] = "ServicePrincipal"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServicePrincipalDatastoreSecrets: %+v", err)
	}

	return encoded, nil
}
