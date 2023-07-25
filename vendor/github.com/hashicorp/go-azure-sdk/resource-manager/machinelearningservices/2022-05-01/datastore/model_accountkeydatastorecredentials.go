package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatastoreCredentials = AccountKeyDatastoreCredentials{}

type AccountKeyDatastoreCredentials struct {
	Secrets DatastoreSecrets `json:"secrets"`

	// Fields inherited from DatastoreCredentials
}

var _ json.Marshaler = AccountKeyDatastoreCredentials{}

func (s AccountKeyDatastoreCredentials) MarshalJSON() ([]byte, error) {
	type wrapper AccountKeyDatastoreCredentials
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AccountKeyDatastoreCredentials: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AccountKeyDatastoreCredentials: %+v", err)
	}
	decoded["credentialsType"] = "AccountKey"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AccountKeyDatastoreCredentials: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &AccountKeyDatastoreCredentials{}

func (s *AccountKeyDatastoreCredentials) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AccountKeyDatastoreCredentials into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["secrets"]; ok {
		impl, err := unmarshalDatastoreSecretsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Secrets' for 'AccountKeyDatastoreCredentials': %+v", err)
		}
		s.Secrets = impl
	}
	return nil
}
