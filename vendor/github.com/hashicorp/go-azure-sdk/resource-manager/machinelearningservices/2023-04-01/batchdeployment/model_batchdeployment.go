package batchdeployment

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchDeployment struct {
	CodeConfiguration         *CodeConfiguration           `json:"codeConfiguration,omitempty"`
	Compute                   *string                      `json:"compute,omitempty"`
	Description               *string                      `json:"description,omitempty"`
	EnvironmentId             *string                      `json:"environmentId,omitempty"`
	EnvironmentVariables      *map[string]string           `json:"environmentVariables,omitempty"`
	ErrorThreshold            *int64                       `json:"errorThreshold,omitempty"`
	LoggingLevel              *BatchLoggingLevel           `json:"loggingLevel,omitempty"`
	MaxConcurrencyPerInstance *int64                       `json:"maxConcurrencyPerInstance,omitempty"`
	MiniBatchSize             *int64                       `json:"miniBatchSize,omitempty"`
	Model                     AssetReferenceBase           `json:"model"`
	OutputAction              *BatchOutputAction           `json:"outputAction,omitempty"`
	OutputFileName            *string                      `json:"outputFileName,omitempty"`
	Properties                *map[string]string           `json:"properties,omitempty"`
	ProvisioningState         *DeploymentProvisioningState `json:"provisioningState,omitempty"`
	Resources                 *ResourceConfiguration       `json:"resources,omitempty"`
	RetrySettings             *BatchRetrySettings          `json:"retrySettings,omitempty"`
}

var _ json.Unmarshaler = &BatchDeployment{}

func (s *BatchDeployment) UnmarshalJSON(bytes []byte) error {
	type alias BatchDeployment
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into BatchDeployment: %+v", err)
	}

	s.CodeConfiguration = decoded.CodeConfiguration
	s.Compute = decoded.Compute
	s.Description = decoded.Description
	s.EnvironmentId = decoded.EnvironmentId
	s.EnvironmentVariables = decoded.EnvironmentVariables
	s.ErrorThreshold = decoded.ErrorThreshold
	s.LoggingLevel = decoded.LoggingLevel
	s.MaxConcurrencyPerInstance = decoded.MaxConcurrencyPerInstance
	s.MiniBatchSize = decoded.MiniBatchSize
	s.OutputAction = decoded.OutputAction
	s.OutputFileName = decoded.OutputFileName
	s.Properties = decoded.Properties
	s.ProvisioningState = decoded.ProvisioningState
	s.Resources = decoded.Resources
	s.RetrySettings = decoded.RetrySettings

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BatchDeployment into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["model"]; ok {
		impl, err := unmarshalAssetReferenceBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Model' for 'BatchDeployment': %+v", err)
		}
		s.Model = impl
	}
	return nil
}
