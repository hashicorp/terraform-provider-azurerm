package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SoftwareConfiguration = ServiceInitiatedSoftwareConfiguration{}

type ServiceInitiatedSoftwareConfiguration struct {
	BomURL                                string                                 `json:"bomUrl"`
	HighAvailabilitySoftwareConfiguration *HighAvailabilitySoftwareConfiguration `json:"highAvailabilitySoftwareConfiguration,omitempty"`
	SapBitsStorageAccountId               string                                 `json:"sapBitsStorageAccountId"`
	SapFqdn                               string                                 `json:"sapFqdn"`
	SoftwareVersion                       string                                 `json:"softwareVersion"`
	SshPrivateKey                         string                                 `json:"sshPrivateKey"`

	// Fields inherited from SoftwareConfiguration

	SoftwareInstallationType SAPSoftwareInstallationType `json:"softwareInstallationType"`
}

func (s ServiceInitiatedSoftwareConfiguration) SoftwareConfiguration() BaseSoftwareConfigurationImpl {
	return BaseSoftwareConfigurationImpl{
		SoftwareInstallationType: s.SoftwareInstallationType,
	}
}

var _ json.Marshaler = ServiceInitiatedSoftwareConfiguration{}

func (s ServiceInitiatedSoftwareConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper ServiceInitiatedSoftwareConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServiceInitiatedSoftwareConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServiceInitiatedSoftwareConfiguration: %+v", err)
	}

	decoded["softwareInstallationType"] = "ServiceInitiated"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServiceInitiatedSoftwareConfiguration: %+v", err)
	}

	return encoded, nil
}
