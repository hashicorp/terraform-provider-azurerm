package secrets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SecretParameters = ManagedCertificateParameters{}

type ManagedCertificateParameters struct {
	ExpirationDate *string `json:"expirationDate,omitempty"`
	Subject        *string `json:"subject,omitempty"`

	// Fields inherited from SecretParameters

	Type SecretType `json:"type"`
}

func (s ManagedCertificateParameters) SecretParameters() BaseSecretParametersImpl {
	return BaseSecretParametersImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ManagedCertificateParameters{}

func (s ManagedCertificateParameters) MarshalJSON() ([]byte, error) {
	type wrapper ManagedCertificateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ManagedCertificateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ManagedCertificateParameters: %+v", err)
	}

	decoded["type"] = "ManagedCertificate"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ManagedCertificateParameters: %+v", err)
	}

	return encoded, nil
}
