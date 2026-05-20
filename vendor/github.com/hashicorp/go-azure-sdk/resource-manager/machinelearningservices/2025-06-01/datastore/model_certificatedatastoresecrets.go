package datastore

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatastoreSecrets = CertificateDatastoreSecrets{}

type CertificateDatastoreSecrets struct {
	Certificate *string `json:"certificate,omitempty"`

	// Fields inherited from DatastoreSecrets

	SecretsType SecretsType `json:"secretsType"`
}

func (s CertificateDatastoreSecrets) DatastoreSecrets() BaseDatastoreSecretsImpl {
	return BaseDatastoreSecretsImpl{
		SecretsType: s.SecretsType,
	}
}

var _ json.Marshaler = CertificateDatastoreSecrets{}

func (s CertificateDatastoreSecrets) MarshalJSON() ([]byte, error) {
	type wrapper CertificateDatastoreSecrets
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CertificateDatastoreSecrets: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CertificateDatastoreSecrets: %+v", err)
	}

	decoded["secretsType"] = "Certificate"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CertificateDatastoreSecrets: %+v", err)
	}

	return encoded, nil
}
