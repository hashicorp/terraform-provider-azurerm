package providerinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProviderSpecificProperties = SapNetWeaverProviderInstanceProperties{}

type SapNetWeaverProviderInstanceProperties struct {
	SapClientId        *string        `json:"sapClientId,omitempty"`
	SapHostFileEntries *[]string      `json:"sapHostFileEntries,omitempty"`
	SapHostname        *string        `json:"sapHostname,omitempty"`
	SapInstanceNr      *string        `json:"sapInstanceNr,omitempty"`
	SapPassword        *string        `json:"sapPassword,omitempty"`
	SapPasswordUri     *string        `json:"sapPasswordUri,omitempty"`
	SapPortNumber      *string        `json:"sapPortNumber,omitempty"`
	SapSid             *string        `json:"sapSid,omitempty"`
	SapUsername        *string        `json:"sapUsername,omitempty"`
	SslCertificateUri  *string        `json:"sslCertificateUri,omitempty"`
	SslPreference      *SslPreference `json:"sslPreference,omitempty"`

	// Fields inherited from ProviderSpecificProperties

	ProviderType string `json:"providerType"`
}

func (s SapNetWeaverProviderInstanceProperties) ProviderSpecificProperties() BaseProviderSpecificPropertiesImpl {
	return BaseProviderSpecificPropertiesImpl{
		ProviderType: s.ProviderType,
	}
}

var _ json.Marshaler = SapNetWeaverProviderInstanceProperties{}

func (s SapNetWeaverProviderInstanceProperties) MarshalJSON() ([]byte, error) {
	type wrapper SapNetWeaverProviderInstanceProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SapNetWeaverProviderInstanceProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SapNetWeaverProviderInstanceProperties: %+v", err)
	}

	decoded["providerType"] = "SapNetWeaver"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SapNetWeaverProviderInstanceProperties: %+v", err)
	}

	return encoded, nil
}
