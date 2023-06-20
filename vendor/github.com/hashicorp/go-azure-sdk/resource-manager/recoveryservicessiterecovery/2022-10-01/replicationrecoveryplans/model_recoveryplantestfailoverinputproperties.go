package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryPlanTestFailoverInputProperties struct {
	FailoverDirection       PossibleOperationsDirections                 `json:"failoverDirection"`
	NetworkId               *string                                      `json:"networkId,omitempty"`
	NetworkType             string                                       `json:"networkType"`
	ProviderSpecificDetails *[]RecoveryPlanProviderSpecificFailoverInput `json:"providerSpecificDetails,omitempty"`
}

var _ json.Unmarshaler = &RecoveryPlanTestFailoverInputProperties{}

func (s *RecoveryPlanTestFailoverInputProperties) UnmarshalJSON(bytes []byte) error {
	type alias RecoveryPlanTestFailoverInputProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into RecoveryPlanTestFailoverInputProperties: %+v", err)
	}

	s.FailoverDirection = decoded.FailoverDirection
	s.NetworkId = decoded.NetworkId
	s.NetworkType = decoded.NetworkType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling RecoveryPlanTestFailoverInputProperties into map[string]json.RawMessage: %+v", err)
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
				return fmt.Errorf("unmarshaling index %d field 'ProviderSpecificDetails' for 'RecoveryPlanTestFailoverInputProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.ProviderSpecificDetails = &output
	}
	return nil
}
