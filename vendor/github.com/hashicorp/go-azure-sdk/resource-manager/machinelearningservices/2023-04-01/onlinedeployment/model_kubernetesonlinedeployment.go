package onlinedeployment

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OnlineDeployment = KubernetesOnlineDeployment{}

type KubernetesOnlineDeployment struct {
	ContainerResourceRequirements *ContainerResourceRequirements `json:"containerResourceRequirements,omitempty"`

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

var _ json.Marshaler = KubernetesOnlineDeployment{}

func (s KubernetesOnlineDeployment) MarshalJSON() ([]byte, error) {
	type wrapper KubernetesOnlineDeployment
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KubernetesOnlineDeployment: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KubernetesOnlineDeployment: %+v", err)
	}
	decoded["endpointComputeType"] = "Kubernetes"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KubernetesOnlineDeployment: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &KubernetesOnlineDeployment{}

func (s *KubernetesOnlineDeployment) UnmarshalJSON(bytes []byte) error {
	type alias KubernetesOnlineDeployment
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into KubernetesOnlineDeployment: %+v", err)
	}

	s.AppInsightsEnabled = decoded.AppInsightsEnabled
	s.CodeConfiguration = decoded.CodeConfiguration
	s.ContainerResourceRequirements = decoded.ContainerResourceRequirements
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
		return fmt.Errorf("unmarshaling KubernetesOnlineDeployment into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["scaleSettings"]; ok {
		impl, err := unmarshalOnlineScaleSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ScaleSettings' for 'KubernetesOnlineDeployment': %+v", err)
		}
		s.ScaleSettings = impl
	}
	return nil
}
