package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SoftwareConfiguration = SAPInstallWithoutOSConfigSoftwareConfiguration{}

type SAPInstallWithoutOSConfigSoftwareConfiguration struct {
	BomURL                                string                                 `json:"bomUrl"`
	HighAvailabilitySoftwareConfiguration *HighAvailabilitySoftwareConfiguration `json:"highAvailabilitySoftwareConfiguration,omitempty"`
	SapBitsStorageAccountId               string                                 `json:"sapBitsStorageAccountId"`
	SoftwareVersion                       string                                 `json:"softwareVersion"`

	// Fields inherited from SoftwareConfiguration

	SoftwareInstallationType SAPSoftwareInstallationType `json:"softwareInstallationType"`
}

func (s SAPInstallWithoutOSConfigSoftwareConfiguration) SoftwareConfiguration() BaseSoftwareConfigurationImpl {
	return BaseSoftwareConfigurationImpl{
		SoftwareInstallationType: s.SoftwareInstallationType,
	}
}

var _ json.Marshaler = SAPInstallWithoutOSConfigSoftwareConfiguration{}

func (s SAPInstallWithoutOSConfigSoftwareConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper SAPInstallWithoutOSConfigSoftwareConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SAPInstallWithoutOSConfigSoftwareConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SAPInstallWithoutOSConfigSoftwareConfiguration: %+v", err)
	}

	decoded["softwareInstallationType"] = "SAPInstallWithoutOSConfig"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SAPInstallWithoutOSConfigSoftwareConfiguration: %+v", err)
	}

	return encoded, nil
}
