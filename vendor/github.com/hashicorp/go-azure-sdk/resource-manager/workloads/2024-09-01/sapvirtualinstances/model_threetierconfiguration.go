package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ InfrastructureConfiguration = ThreeTierConfiguration{}

type ThreeTierConfiguration struct {
	ApplicationServer      ApplicationServerConfiguration `json:"applicationServer"`
	CentralServer          CentralServerConfiguration     `json:"centralServer"`
	CustomResourceNames    ThreeTierCustomResourceNames   `json:"customResourceNames"`
	DatabaseServer         DatabaseConfiguration          `json:"databaseServer"`
	HighAvailabilityConfig *HighAvailabilityConfiguration `json:"highAvailabilityConfig,omitempty"`
	NetworkConfiguration   *NetworkConfiguration          `json:"networkConfiguration,omitempty"`
	StorageConfiguration   *StorageConfiguration          `json:"storageConfiguration,omitempty"`

	// Fields inherited from InfrastructureConfiguration

	AppResourceGroup string            `json:"appResourceGroup"`
	DeploymentType   SAPDeploymentType `json:"deploymentType"`
}

func (s ThreeTierConfiguration) InfrastructureConfiguration() BaseInfrastructureConfigurationImpl {
	return BaseInfrastructureConfigurationImpl{
		AppResourceGroup: s.AppResourceGroup,
		DeploymentType:   s.DeploymentType,
	}
}

var _ json.Marshaler = ThreeTierConfiguration{}

func (s ThreeTierConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper ThreeTierConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ThreeTierConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ThreeTierConfiguration: %+v", err)
	}

	decoded["deploymentType"] = "ThreeTier"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ThreeTierConfiguration: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &ThreeTierConfiguration{}

func (s *ThreeTierConfiguration) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ApplicationServer      ApplicationServerConfiguration `json:"applicationServer"`
		CentralServer          CentralServerConfiguration     `json:"centralServer"`
		DatabaseServer         DatabaseConfiguration          `json:"databaseServer"`
		HighAvailabilityConfig *HighAvailabilityConfiguration `json:"highAvailabilityConfig,omitempty"`
		NetworkConfiguration   *NetworkConfiguration          `json:"networkConfiguration,omitempty"`
		StorageConfiguration   *StorageConfiguration          `json:"storageConfiguration,omitempty"`
		AppResourceGroup       string                         `json:"appResourceGroup"`
		DeploymentType         SAPDeploymentType              `json:"deploymentType"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ApplicationServer = decoded.ApplicationServer
	s.CentralServer = decoded.CentralServer
	s.DatabaseServer = decoded.DatabaseServer
	s.HighAvailabilityConfig = decoded.HighAvailabilityConfig
	s.NetworkConfiguration = decoded.NetworkConfiguration
	s.StorageConfiguration = decoded.StorageConfiguration
	s.AppResourceGroup = decoded.AppResourceGroup
	s.DeploymentType = decoded.DeploymentType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ThreeTierConfiguration into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["customResourceNames"]; ok {
		impl, err := UnmarshalThreeTierCustomResourceNamesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CustomResourceNames' for 'ThreeTierConfiguration': %+v", err)
		}
		s.CustomResourceNames = impl
	}

	return nil
}
