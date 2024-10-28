package providerinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderInstanceProperties struct {
	Errors            *Error                            `json:"errors,omitempty"`
	ProviderSettings  ProviderSpecificProperties        `json:"providerSettings"`
	ProvisioningState *WorkloadMonitorProvisioningState `json:"provisioningState,omitempty"`
}

var _ json.Unmarshaler = &ProviderInstanceProperties{}

func (s *ProviderInstanceProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Errors            *Error                            `json:"errors,omitempty"`
		ProvisioningState *WorkloadMonitorProvisioningState `json:"provisioningState,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Errors = decoded.Errors
	s.ProvisioningState = decoded.ProvisioningState

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ProviderInstanceProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSettings"]; ok {
		impl, err := UnmarshalProviderSpecificPropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSettings' for 'ProviderInstanceProperties': %+v", err)
		}
		s.ProviderSettings = impl
	}

	return nil
}
