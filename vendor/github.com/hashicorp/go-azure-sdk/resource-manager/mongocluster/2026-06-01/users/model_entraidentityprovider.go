package users

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ IdentityProvider = EntraIdentityProvider{}

type EntraIdentityProvider struct {
	Properties EntraIdentityProviderProperties `json:"properties"`

	// Fields inherited from IdentityProvider

	Type IdentityProviderType `json:"type"`
}

func (s EntraIdentityProvider) IdentityProvider() BaseIdentityProviderImpl {
	return BaseIdentityProviderImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = EntraIdentityProvider{}

func (s EntraIdentityProvider) MarshalJSON() ([]byte, error) {
	type wrapper EntraIdentityProvider
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EntraIdentityProvider: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EntraIdentityProvider: %+v", err)
	}

	decoded["type"] = "MicrosoftEntraID"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EntraIdentityProvider: %+v", err)
	}

	return encoded, nil
}
