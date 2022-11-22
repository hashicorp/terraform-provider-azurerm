package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PlannedFailoverInputProperties struct {
	FailoverDirection       *string                                      `json:"failoverDirection,omitempty"`
	ProviderSpecificDetails PlannedFailoverProviderSpecificFailoverInput `json:"providerSpecificDetails"`
}

var _ json.Unmarshaler = &PlannedFailoverInputProperties{}

func (s *PlannedFailoverInputProperties) UnmarshalJSON(bytes []byte) error {
	type alias PlannedFailoverInputProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into PlannedFailoverInputProperties: %+v", err)
	}

	s.FailoverDirection = decoded.FailoverDirection

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling PlannedFailoverInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := unmarshalPlannedFailoverProviderSpecificFailoverInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'PlannedFailoverInputProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}
	return nil
}
