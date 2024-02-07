package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatastoreCredentials = NoneDatastoreCredentials{}

type NoneDatastoreCredentials struct {

	// Fields inherited from DatastoreCredentials
}

var _ json.Marshaler = NoneDatastoreCredentials{}

func (s NoneDatastoreCredentials) MarshalJSON() ([]byte, error) {
	type wrapper NoneDatastoreCredentials
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NoneDatastoreCredentials: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NoneDatastoreCredentials: %+v", err)
	}
	decoded["credentialsType"] = "None"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NoneDatastoreCredentials: %+v", err)
	}

	return encoded, nil
}
