package providerinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProviderSpecificProperties = PrometheusOSProviderInstanceProperties{}

type PrometheusOSProviderInstanceProperties struct {
	PrometheusURL     *string        `json:"prometheusUrl,omitempty"`
	SapSid            *string        `json:"sapSid,omitempty"`
	SslCertificateUri *string        `json:"sslCertificateUri,omitempty"`
	SslPreference     *SslPreference `json:"sslPreference,omitempty"`

	// Fields inherited from ProviderSpecificProperties

	ProviderType string `json:"providerType"`
}

func (s PrometheusOSProviderInstanceProperties) ProviderSpecificProperties() BaseProviderSpecificPropertiesImpl {
	return BaseProviderSpecificPropertiesImpl{
		ProviderType: s.ProviderType,
	}
}

var _ json.Marshaler = PrometheusOSProviderInstanceProperties{}

func (s PrometheusOSProviderInstanceProperties) MarshalJSON() ([]byte, error) {
	type wrapper PrometheusOSProviderInstanceProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PrometheusOSProviderInstanceProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PrometheusOSProviderInstanceProperties: %+v", err)
	}

	decoded["providerType"] = "PrometheusOS"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PrometheusOSProviderInstanceProperties: %+v", err)
	}

	return encoded, nil
}
