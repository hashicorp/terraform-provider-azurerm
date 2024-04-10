package providerinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProviderSpecificProperties = MsSqlServerProviderInstanceProperties{}

type MsSqlServerProviderInstanceProperties struct {
	DbPassword        *string        `json:"dbPassword,omitempty"`
	DbPasswordUri     *string        `json:"dbPasswordUri,omitempty"`
	DbPort            *string        `json:"dbPort,omitempty"`
	DbUsername        *string        `json:"dbUsername,omitempty"`
	Hostname          *string        `json:"hostname,omitempty"`
	SapSid            *string        `json:"sapSid,omitempty"`
	SslCertificateUri *string        `json:"sslCertificateUri,omitempty"`
	SslPreference     *SslPreference `json:"sslPreference,omitempty"`

	// Fields inherited from ProviderSpecificProperties
}

var _ json.Marshaler = MsSqlServerProviderInstanceProperties{}

func (s MsSqlServerProviderInstanceProperties) MarshalJSON() ([]byte, error) {
	type wrapper MsSqlServerProviderInstanceProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MsSqlServerProviderInstanceProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MsSqlServerProviderInstanceProperties: %+v", err)
	}
	decoded["providerType"] = "MsSqlServer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MsSqlServerProviderInstanceProperties: %+v", err)
	}

	return encoded, nil
}
