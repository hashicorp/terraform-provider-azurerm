package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SoftwareConfiguration = ExternalInstallationSoftwareConfiguration{}

type ExternalInstallationSoftwareConfiguration struct {
	CentralServerVMId *string `json:"centralServerVmId,omitempty"`

	// Fields inherited from SoftwareConfiguration

	SoftwareInstallationType SAPSoftwareInstallationType `json:"softwareInstallationType"`
}

func (s ExternalInstallationSoftwareConfiguration) SoftwareConfiguration() BaseSoftwareConfigurationImpl {
	return BaseSoftwareConfigurationImpl{
		SoftwareInstallationType: s.SoftwareInstallationType,
	}
}

var _ json.Marshaler = ExternalInstallationSoftwareConfiguration{}

func (s ExternalInstallationSoftwareConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper ExternalInstallationSoftwareConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ExternalInstallationSoftwareConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ExternalInstallationSoftwareConfiguration: %+v", err)
	}

	decoded["softwareInstallationType"] = "External"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ExternalInstallationSoftwareConfiguration: %+v", err)
	}

	return encoded, nil
}
