package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageTemplateProperties struct {
	AutoRun                   *ImageTemplateAutoRun                 `json:"autoRun,omitempty"`
	BuildTimeoutInMinutes     *int64                                `json:"buildTimeoutInMinutes,omitempty"`
	Customize                 *[]ImageTemplateCustomizer            `json:"customize,omitempty"`
	Distribute                []ImageTemplateDistributor            `json:"distribute"`
	ErrorHandling             *ImageTemplatePropertiesErrorHandling `json:"errorHandling,omitempty"`
	ExactStagingResourceGroup *string                               `json:"exactStagingResourceGroup,omitempty"`
	LastRunStatus             *ImageTemplateLastRunStatus           `json:"lastRunStatus,omitempty"`
	ManagedResourceTags       *map[string]string                    `json:"managedResourceTags,omitempty"`
	Optimize                  *ImageTemplatePropertiesOptimize      `json:"optimize,omitempty"`
	ProvisioningError         *ProvisioningError                    `json:"provisioningError,omitempty"`
	ProvisioningState         *ProvisioningState                    `json:"provisioningState,omitempty"`
	Source                    ImageTemplateSource                   `json:"source"`
	StagingResourceGroup      *string                               `json:"stagingResourceGroup,omitempty"`
	VMProfile                 *ImageTemplateVMProfile               `json:"vmProfile,omitempty"`
	Validate                  *ImageTemplatePropertiesValidate      `json:"validate,omitempty"`
}

var _ json.Unmarshaler = &ImageTemplateProperties{}

func (s *ImageTemplateProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AutoRun                   *ImageTemplateAutoRun                 `json:"autoRun,omitempty"`
		BuildTimeoutInMinutes     *int64                                `json:"buildTimeoutInMinutes,omitempty"`
		ErrorHandling             *ImageTemplatePropertiesErrorHandling `json:"errorHandling,omitempty"`
		ExactStagingResourceGroup *string                               `json:"exactStagingResourceGroup,omitempty"`
		LastRunStatus             *ImageTemplateLastRunStatus           `json:"lastRunStatus,omitempty"`
		ManagedResourceTags       *map[string]string                    `json:"managedResourceTags,omitempty"`
		Optimize                  *ImageTemplatePropertiesOptimize      `json:"optimize,omitempty"`
		ProvisioningError         *ProvisioningError                    `json:"provisioningError,omitempty"`
		ProvisioningState         *ProvisioningState                    `json:"provisioningState,omitempty"`
		StagingResourceGroup      *string                               `json:"stagingResourceGroup,omitempty"`
		VMProfile                 *ImageTemplateVMProfile               `json:"vmProfile,omitempty"`
		Validate                  *ImageTemplatePropertiesValidate      `json:"validate,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AutoRun = decoded.AutoRun
	s.BuildTimeoutInMinutes = decoded.BuildTimeoutInMinutes
	s.ErrorHandling = decoded.ErrorHandling
	s.ExactStagingResourceGroup = decoded.ExactStagingResourceGroup
	s.LastRunStatus = decoded.LastRunStatus
	s.ManagedResourceTags = decoded.ManagedResourceTags
	s.Optimize = decoded.Optimize
	s.ProvisioningError = decoded.ProvisioningError
	s.ProvisioningState = decoded.ProvisioningState
	s.StagingResourceGroup = decoded.StagingResourceGroup
	s.VMProfile = decoded.VMProfile
	s.Validate = decoded.Validate

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ImageTemplateProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["customize"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Customize into list []json.RawMessage: %+v", err)
		}

		output := make([]ImageTemplateCustomizer, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalImageTemplateCustomizerImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Customize' for 'ImageTemplateProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Customize = &output
	}

	if v, ok := temp["distribute"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Distribute into list []json.RawMessage: %+v", err)
		}

		output := make([]ImageTemplateDistributor, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalImageTemplateDistributorImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Distribute' for 'ImageTemplateProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Distribute = output
	}

	if v, ok := temp["source"]; ok {
		impl, err := UnmarshalImageTemplateSourceImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Source' for 'ImageTemplateProperties': %+v", err)
		}
		s.Source = impl
	}

	return nil
}
