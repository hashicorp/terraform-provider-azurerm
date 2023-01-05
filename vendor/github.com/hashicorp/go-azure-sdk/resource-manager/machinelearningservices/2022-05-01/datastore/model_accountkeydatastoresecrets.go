package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatastoreSecrets = AccountKeyDatastoreSecrets{}

type AccountKeyDatastoreSecrets struct {
	Key *string `json:"key,omitempty"`

	// Fields inherited from DatastoreSecrets
}

var _ json.Marshaler = AccountKeyDatastoreSecrets{}

func (s AccountKeyDatastoreSecrets) MarshalJSON() ([]byte, error) {
	type wrapper AccountKeyDatastoreSecrets
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AccountKeyDatastoreSecrets: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AccountKeyDatastoreSecrets: %+v", err)
	}
	decoded["secretsType"] = "AccountKey"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AccountKeyDatastoreSecrets: %+v", err)
	}

	return encoded, nil
}
