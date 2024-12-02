package providerinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProviderSpecificProperties = PrometheusHaClusterProviderInstanceProperties{}

type PrometheusHaClusterProviderInstanceProperties struct {
	ClusterName       *string        `json:"clusterName,omitempty"`
	Hostname          *string        `json:"hostname,omitempty"`
	PrometheusURL     *string        `json:"prometheusUrl,omitempty"`
	Sid               *string        `json:"sid,omitempty"`
	SslCertificateUri *string        `json:"sslCertificateUri,omitempty"`
	SslPreference     *SslPreference `json:"sslPreference,omitempty"`

	// Fields inherited from ProviderSpecificProperties

	ProviderType string `json:"providerType"`
}

func (s PrometheusHaClusterProviderInstanceProperties) ProviderSpecificProperties() BaseProviderSpecificPropertiesImpl {
	return BaseProviderSpecificPropertiesImpl{
		ProviderType: s.ProviderType,
	}
}

var _ json.Marshaler = PrometheusHaClusterProviderInstanceProperties{}

func (s PrometheusHaClusterProviderInstanceProperties) MarshalJSON() ([]byte, error) {
	type wrapper PrometheusHaClusterProviderInstanceProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PrometheusHaClusterProviderInstanceProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PrometheusHaClusterProviderInstanceProperties: %+v", err)
	}

	decoded["providerType"] = "PrometheusHaCluster"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PrometheusHaClusterProviderInstanceProperties: %+v", err)
	}

	return encoded, nil
}
