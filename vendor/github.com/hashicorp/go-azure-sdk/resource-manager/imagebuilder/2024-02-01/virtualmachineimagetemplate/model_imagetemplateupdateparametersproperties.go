package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageTemplateUpdateParametersProperties struct {
	Distribute *[]ImageTemplateDistributor `json:"distribute,omitempty"`
	VMProfile  *ImageTemplateVMProfile     `json:"vmProfile,omitempty"`
}

var _ json.Unmarshaler = &ImageTemplateUpdateParametersProperties{}

func (s *ImageTemplateUpdateParametersProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		VMProfile *ImageTemplateVMProfile `json:"vmProfile,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.VMProfile = decoded.VMProfile

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ImageTemplateUpdateParametersProperties into map[string]json.RawMessage: %+v", err)
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
				return fmt.Errorf("unmarshaling index %d field 'Distribute' for 'ImageTemplateUpdateParametersProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Distribute = &output
	}

	return nil
}
