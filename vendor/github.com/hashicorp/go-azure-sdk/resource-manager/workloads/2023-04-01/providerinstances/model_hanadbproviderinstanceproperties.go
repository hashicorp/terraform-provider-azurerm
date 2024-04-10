package providerinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProviderSpecificProperties = HanaDbProviderInstanceProperties{}

type HanaDbProviderInstanceProperties struct {
	DbName                   *string        `json:"dbName,omitempty"`
	DbPassword               *string        `json:"dbPassword,omitempty"`
	DbPasswordUri            *string        `json:"dbPasswordUri,omitempty"`
	DbUsername               *string        `json:"dbUsername,omitempty"`
	Hostname                 *string        `json:"hostname,omitempty"`
	InstanceNumber           *string        `json:"instanceNumber,omitempty"`
	SapSid                   *string        `json:"sapSid,omitempty"`
	SqlPort                  *string        `json:"sqlPort,omitempty"`
	SslCertificateUri        *string        `json:"sslCertificateUri,omitempty"`
	SslHostNameInCertificate *string        `json:"sslHostNameInCertificate,omitempty"`
	SslPreference            *SslPreference `json:"sslPreference,omitempty"`

	// Fields inherited from ProviderSpecificProperties
}

var _ json.Marshaler = HanaDbProviderInstanceProperties{}

func (s HanaDbProviderInstanceProperties) MarshalJSON() ([]byte, error) {
	type wrapper HanaDbProviderInstanceProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HanaDbProviderInstanceProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HanaDbProviderInstanceProperties: %+v", err)
	}
	decoded["providerType"] = "SapHana"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HanaDbProviderInstanceProperties: %+v", err)
	}

	return encoded, nil
}
