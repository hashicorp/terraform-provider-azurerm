package secrets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SecretParameters = AfdSecretMtlsCertificateChain{}

type AfdSecretMtlsCertificateChain struct {
	ExpirationDate *string           `json:"expirationDate,omitempty"`
	SecretSource   ResourceReference `json:"secretSource"`
	SecretVersion  string            `json:"secretVersion"`

	// Fields inherited from SecretParameters

	Type SecretType `json:"type"`
}

func (s AfdSecretMtlsCertificateChain) SecretParameters() BaseSecretParametersImpl {
	return BaseSecretParametersImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AfdSecretMtlsCertificateChain{}

func (s AfdSecretMtlsCertificateChain) MarshalJSON() ([]byte, error) {
	type wrapper AfdSecretMtlsCertificateChain
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AfdSecretMtlsCertificateChain: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AfdSecretMtlsCertificateChain: %+v", err)
	}

	decoded["type"] = "MtlsCertificateChain"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AfdSecretMtlsCertificateChain: %+v", err)
	}

	return encoded, nil
}
