package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageTemplatePropertiesValidate struct {
	ContinueDistributeOnFailure *bool                         `json:"continueDistributeOnFailure,omitempty"`
	InVMValidations             *[]ImageTemplateInVMValidator `json:"inVMValidations,omitempty"`
	SourceValidationOnly        *bool                         `json:"sourceValidationOnly,omitempty"`
}

var _ json.Unmarshaler = &ImageTemplatePropertiesValidate{}

func (s *ImageTemplatePropertiesValidate) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ContinueDistributeOnFailure *bool `json:"continueDistributeOnFailure,omitempty"`
		SourceValidationOnly        *bool `json:"sourceValidationOnly,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ContinueDistributeOnFailure = decoded.ContinueDistributeOnFailure
	s.SourceValidationOnly = decoded.SourceValidationOnly

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ImageTemplatePropertiesValidate into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["inVMValidations"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling InVMValidations into list []json.RawMessage: %+v", err)
		}

		output := make([]ImageTemplateInVMValidator, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalImageTemplateInVMValidatorImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'InVMValidations' for 'ImageTemplatePropertiesValidate': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.InVMValidations = &output
	}

	return nil
}
