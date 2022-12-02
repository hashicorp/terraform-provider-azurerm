package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryPlanProviderSpecificDetails interface {
}

func unmarshalRecoveryPlanProviderSpecificDetailsImplementation(input []byte) (RecoveryPlanProviderSpecificDetails, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanProviderSpecificDetails into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "A2A") {
		var out RecoveryPlanA2ADetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RecoveryPlanA2ADetails: %+v", err)
		}
		return out, nil
	}

	type RawRecoveryPlanProviderSpecificDetailsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawRecoveryPlanProviderSpecificDetailsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
