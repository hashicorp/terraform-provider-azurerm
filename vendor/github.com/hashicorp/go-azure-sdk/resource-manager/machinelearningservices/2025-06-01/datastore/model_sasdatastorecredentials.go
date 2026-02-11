package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatastoreCredentials = SasDatastoreCredentials{}

type SasDatastoreCredentials struct {
	Secrets SasDatastoreSecrets `json:"secrets"`

	// Fields inherited from DatastoreCredentials

	CredentialsType CredentialsType `json:"credentialsType"`
}

func (s SasDatastoreCredentials) DatastoreCredentials() BaseDatastoreCredentialsImpl {
	return BaseDatastoreCredentialsImpl{
		CredentialsType: s.CredentialsType,
	}
}

var _ json.Marshaler = SasDatastoreCredentials{}

func (s SasDatastoreCredentials) MarshalJSON() ([]byte, error) {
	type wrapper SasDatastoreCredentials
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SasDatastoreCredentials: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SasDatastoreCredentials: %+v", err)
	}

	decoded["credentialsType"] = "Sas"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SasDatastoreCredentials: %+v", err)
	}

	return encoded, nil
}
