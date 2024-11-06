package providerinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProviderSpecificProperties = DB2ProviderInstanceProperties{}

type DB2ProviderInstanceProperties struct {
	DbName            *string        `json:"dbName,omitempty"`
	DbPassword        *string        `json:"dbPassword,omitempty"`
	DbPasswordUri     *string        `json:"dbPasswordUri,omitempty"`
	DbPort            *string        `json:"dbPort,omitempty"`
	DbUsername        *string        `json:"dbUsername,omitempty"`
	Hostname          *string        `json:"hostname,omitempty"`
	SapSid            *string        `json:"sapSid,omitempty"`
	SslCertificateUri *string        `json:"sslCertificateUri,omitempty"`
	SslPreference     *SslPreference `json:"sslPreference,omitempty"`

	// Fields inherited from ProviderSpecificProperties

	ProviderType string `json:"providerType"`
}

func (s DB2ProviderInstanceProperties) ProviderSpecificProperties() BaseProviderSpecificPropertiesImpl {
	return BaseProviderSpecificPropertiesImpl{
		ProviderType: s.ProviderType,
	}
}

var _ json.Marshaler = DB2ProviderInstanceProperties{}

func (s DB2ProviderInstanceProperties) MarshalJSON() ([]byte, error) {
	type wrapper DB2ProviderInstanceProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DB2ProviderInstanceProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DB2ProviderInstanceProperties: %+v", err)
	}

	decoded["providerType"] = "Db2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DB2ProviderInstanceProperties: %+v", err)
	}

	return encoded, nil
}
