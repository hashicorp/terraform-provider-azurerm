package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryPlanUnplannedFailoverInputProperties struct {
	FailoverDirection       PossibleOperationsDirections                 `json:"failoverDirection"`
	ProviderSpecificDetails *[]RecoveryPlanProviderSpecificFailoverInput `json:"providerSpecificDetails,omitempty"`
	SourceSiteOperations    SourceSiteOperations                         `json:"sourceSiteOperations"`
}

var _ json.Unmarshaler = &RecoveryPlanUnplannedFailoverInputProperties{}

func (s *RecoveryPlanUnplannedFailoverInputProperties) UnmarshalJSON(bytes []byte) error {
	type alias RecoveryPlanUnplannedFailoverInputProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into RecoveryPlanUnplannedFailoverInputProperties: %+v", err)
	}

	s.FailoverDirection = decoded.FailoverDirection
	s.SourceSiteOperations = decoded.SourceSiteOperations

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling RecoveryPlanUnplannedFailoverInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling ProviderSpecificDetails into list []json.RawMessage: %+v", err)
		}

		output := make([]RecoveryPlanProviderSpecificFailoverInput, 0)
		for i, val := range listTemp {
			impl, err := unmarshalRecoveryPlanProviderSpecificFailoverInputImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'ProviderSpecificDetails' for 'RecoveryPlanUnplannedFailoverInputProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.ProviderSpecificDetails = &output
	}
	return nil
}
