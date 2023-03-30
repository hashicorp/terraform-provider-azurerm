package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SwitchProviderInputProperties struct {
	ProviderSpecificDetails SwitchProviderProviderSpecificInput `json:"providerSpecificDetails"`
	TargetInstanceType      *string                             `json:"targetInstanceType,omitempty"`
}

var _ json.Unmarshaler = &SwitchProviderInputProperties{}

func (s *SwitchProviderInputProperties) UnmarshalJSON(bytes []byte) error {
	type alias SwitchProviderInputProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into SwitchProviderInputProperties: %+v", err)
	}

	s.TargetInstanceType = decoded.TargetInstanceType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SwitchProviderInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := unmarshalSwitchProviderProviderSpecificInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'SwitchProviderInputProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}
	return nil
}
