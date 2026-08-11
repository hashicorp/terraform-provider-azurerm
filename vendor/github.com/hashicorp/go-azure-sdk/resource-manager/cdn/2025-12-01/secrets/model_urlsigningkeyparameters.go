package secrets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SecretParameters = URLSigningKeyParameters{}

type URLSigningKeyParameters struct {
	KeyId         string            `json:"keyId"`
	SecretSource  ResourceReference `json:"secretSource"`
	SecretVersion string            `json:"secretVersion"`

	// Fields inherited from SecretParameters

	Type SecretType `json:"type"`
}

func (s URLSigningKeyParameters) SecretParameters() BaseSecretParametersImpl {
	return BaseSecretParametersImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = URLSigningKeyParameters{}

func (s URLSigningKeyParameters) MarshalJSON() ([]byte, error) {
	type wrapper URLSigningKeyParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling URLSigningKeyParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling URLSigningKeyParameters: %+v", err)
	}

	decoded["type"] = "UrlSigningKey"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling URLSigningKeyParameters: %+v", err)
	}

	return encoded, nil
}
