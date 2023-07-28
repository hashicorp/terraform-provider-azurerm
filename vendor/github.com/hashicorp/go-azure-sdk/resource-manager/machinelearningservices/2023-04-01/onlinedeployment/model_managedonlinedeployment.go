package onlinedeployment

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OnlineDeployment = ManagedOnlineDeployment{}

type ManagedOnlineDeployment struct {

	// Fields inherited from OnlineDeployment
	AppInsightsEnabled        *bool                          `json:"appInsightsEnabled,omitempty"`
	CodeConfiguration         *CodeConfiguration             `json:"codeConfiguration,omitempty"`
	Description               *string                        `json:"description,omitempty"`
	EgressPublicNetworkAccess *EgressPublicNetworkAccessType `json:"egressPublicNetworkAccess,omitempty"`
	EnvironmentId             *string                        `json:"environmentId,omitempty"`
	EnvironmentVariables      *map[string]string             `json:"environmentVariables,omitempty"`
	InstanceType              *string                        `json:"instanceType,omitempty"`
	LivenessProbe             *ProbeSettings                 `json:"livenessProbe,omitempty"`
	Model                     *string                        `json:"model,omitempty"`
	ModelMountPath            *string                        `json:"modelMountPath,omitempty"`
	Properties                *map[string]string             `json:"properties,omitempty"`
	ProvisioningState         *DeploymentProvisioningState   `json:"provisioningState,omitempty"`
	ReadinessProbe            *ProbeSettings                 `json:"readinessProbe,omitempty"`
	RequestSettings           *OnlineRequestSettings         `json:"requestSettings,omitempty"`
	ScaleSettings             OnlineScaleSettings            `json:"scaleSettings"`
}

var _ json.Marshaler = ManagedOnlineDeployment{}

func (s ManagedOnlineDeployment) MarshalJSON() ([]byte, error) {
	type wrapper ManagedOnlineDeployment
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ManagedOnlineDeployment: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ManagedOnlineDeployment: %+v", err)
	}
	decoded["endpointComputeType"] = "Managed"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ManagedOnlineDeployment: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &ManagedOnlineDeployment{}

func (s *ManagedOnlineDeployment) UnmarshalJSON(bytes []byte) error {
	type alias ManagedOnlineDeployment
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ManagedOnlineDeployment: %+v", err)
	}

	s.AppInsightsEnabled = decoded.AppInsightsEnabled
	s.CodeConfiguration = decoded.CodeConfiguration
	s.Description = decoded.Description
	s.EgressPublicNetworkAccess = decoded.EgressPublicNetworkAccess
	s.EnvironmentId = decoded.EnvironmentId
	s.EnvironmentVariables = decoded.EnvironmentVariables
	s.InstanceType = decoded.InstanceType
	s.LivenessProbe = decoded.LivenessProbe
	s.Model = decoded.Model
	s.ModelMountPath = decoded.ModelMountPath
	s.Properties = decoded.Properties
	s.ProvisioningState = decoded.ProvisioningState
	s.ReadinessProbe = decoded.ReadinessProbe
	s.RequestSettings = decoded.RequestSettings

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ManagedOnlineDeployment into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["scaleSettings"]; ok {
		impl, err := unmarshalOnlineScaleSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ScaleSettings' for 'ManagedOnlineDeployment': %+v", err)
		}
		s.ScaleSettings = impl
	}
	return nil
}
