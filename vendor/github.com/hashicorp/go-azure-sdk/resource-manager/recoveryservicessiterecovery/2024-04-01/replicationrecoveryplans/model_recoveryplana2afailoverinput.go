package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RecoveryPlanProviderSpecificFailoverInput = RecoveryPlanA2AFailoverInput{}

type RecoveryPlanA2AFailoverInput struct {
	CloudServiceCreationOption *string                 `json:"cloudServiceCreationOption,omitempty"`
	MultiVMSyncPointOption     *MultiVMSyncPointOption `json:"multiVmSyncPointOption,omitempty"`
	RecoveryPointType          A2ARpRecoveryPointType  `json:"recoveryPointType"`

	// Fields inherited from RecoveryPlanProviderSpecificFailoverInput

	InstanceType string `json:"instanceType"`
}

func (s RecoveryPlanA2AFailoverInput) RecoveryPlanProviderSpecificFailoverInput() BaseRecoveryPlanProviderSpecificFailoverInputImpl {
	return BaseRecoveryPlanProviderSpecificFailoverInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = RecoveryPlanA2AFailoverInput{}

func (s RecoveryPlanA2AFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper RecoveryPlanA2AFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RecoveryPlanA2AFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanA2AFailoverInput: %+v", err)
	}

	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RecoveryPlanA2AFailoverInput: %+v", err)
	}

	return encoded, nil
}
