package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatastoreSecrets = SasDatastoreSecrets{}

type SasDatastoreSecrets struct {
	SasToken *string `json:"sasToken,omitempty"`

	// Fields inherited from DatastoreSecrets

	SecretsType SecretsType `json:"secretsType"`
}

func (s SasDatastoreSecrets) DatastoreSecrets() BaseDatastoreSecretsImpl {
	return BaseDatastoreSecretsImpl{
		SecretsType: s.SecretsType,
	}
}

var _ json.Marshaler = SasDatastoreSecrets{}

func (s SasDatastoreSecrets) MarshalJSON() ([]byte, error) {
	type wrapper SasDatastoreSecrets
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SasDatastoreSecrets: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SasDatastoreSecrets: %+v", err)
	}

	decoded["secretsType"] = "Sas"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SasDatastoreSecrets: %+v", err)
	}

	return encoded, nil
}
