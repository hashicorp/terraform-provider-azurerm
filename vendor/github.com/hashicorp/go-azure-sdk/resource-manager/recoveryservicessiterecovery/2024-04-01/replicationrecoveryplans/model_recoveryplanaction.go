package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryPlanAction struct {
	ActionName         string                              `json:"actionName"`
	CustomDetails      RecoveryPlanActionDetails           `json:"customDetails"`
	FailoverDirections []PossibleOperationsDirections      `json:"failoverDirections"`
	FailoverTypes      []ReplicationProtectedItemOperation `json:"failoverTypes"`
}

var _ json.Unmarshaler = &RecoveryPlanAction{}

func (s *RecoveryPlanAction) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ActionName         string                              `json:"actionName"`
		FailoverDirections []PossibleOperationsDirections      `json:"failoverDirections"`
		FailoverTypes      []ReplicationProtectedItemOperation `json:"failoverTypes"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ActionName = decoded.ActionName
	s.FailoverDirections = decoded.FailoverDirections
	s.FailoverTypes = decoded.FailoverTypes

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling RecoveryPlanAction into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["customDetails"]; ok {
		impl, err := UnmarshalRecoveryPlanActionDetailsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CustomDetails' for 'RecoveryPlanAction': %+v", err)
		}
		s.CustomDetails = impl
	}

	return nil
}
