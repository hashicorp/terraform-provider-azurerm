package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateRecoveryPlanInputProperties struct {
	FailoverDeploymentModel *FailoverDeploymentModel             `json:"failoverDeploymentModel,omitempty"`
	Groups                  []RecoveryPlanGroup                  `json:"groups"`
	PrimaryFabricId         string                               `json:"primaryFabricId"`
	ProviderSpecificInput   *[]RecoveryPlanProviderSpecificInput `json:"providerSpecificInput,omitempty"`
	RecoveryFabricId        string                               `json:"recoveryFabricId"`
}

var _ json.Unmarshaler = &CreateRecoveryPlanInputProperties{}

func (s *CreateRecoveryPlanInputProperties) UnmarshalJSON(bytes []byte) error {
	type alias CreateRecoveryPlanInputProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into CreateRecoveryPlanInputProperties: %+v", err)
	}

	s.FailoverDeploymentModel = decoded.FailoverDeploymentModel
	s.Groups = decoded.Groups
	s.PrimaryFabricId = decoded.PrimaryFabricId
	s.RecoveryFabricId = decoded.RecoveryFabricId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CreateRecoveryPlanInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificInput"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling ProviderSpecificInput into list []json.RawMessage: %+v", err)
		}

		output := make([]RecoveryPlanProviderSpecificInput, 0)
		for i, val := range listTemp {
			impl, err := unmarshalRecoveryPlanProviderSpecificInputImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'ProviderSpecificInput' for 'CreateRecoveryPlanInputProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.ProviderSpecificInput = &output
	}
	return nil
}
