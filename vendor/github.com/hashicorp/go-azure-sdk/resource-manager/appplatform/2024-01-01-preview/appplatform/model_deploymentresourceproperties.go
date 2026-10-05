package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentResourceProperties struct {
	Active             *bool                                `json:"active,omitempty"`
	DeploymentSettings *DeploymentSettings                  `json:"deploymentSettings,omitempty"`
	Instances          *[]DeploymentInstance                `json:"instances,omitempty"`
	ProvisioningState  *DeploymentResourceProvisioningState `json:"provisioningState,omitempty"`
	Source             UserSourceInfo                       `json:"source"`
	Status             *DeploymentResourceStatus            `json:"status,omitempty"`
}

var _ json.Unmarshaler = &DeploymentResourceProperties{}

func (s *DeploymentResourceProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Active             *bool                                `json:"active,omitempty"`
		DeploymentSettings *DeploymentSettings                  `json:"deploymentSettings,omitempty"`
		Instances          *[]DeploymentInstance                `json:"instances,omitempty"`
		ProvisioningState  *DeploymentResourceProvisioningState `json:"provisioningState,omitempty"`
		Status             *DeploymentResourceStatus            `json:"status,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Active = decoded.Active
	s.DeploymentSettings = decoded.DeploymentSettings
	s.Instances = decoded.Instances
	s.ProvisioningState = decoded.ProvisioningState
	s.Status = decoded.Status

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DeploymentResourceProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["source"]; ok {
		impl, err := UnmarshalUserSourceInfoImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Source' for 'DeploymentResourceProperties': %+v", err)
		}
		s.Source = impl
	}

	return nil
}
