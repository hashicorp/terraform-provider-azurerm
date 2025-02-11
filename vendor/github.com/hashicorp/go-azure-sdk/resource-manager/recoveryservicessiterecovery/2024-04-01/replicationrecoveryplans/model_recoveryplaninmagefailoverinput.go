package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RecoveryPlanProviderSpecificFailoverInput = RecoveryPlanInMageFailoverInput{}

type RecoveryPlanInMageFailoverInput struct {
	RecoveryPointType RpInMageRecoveryPointType `json:"recoveryPointType"`

	// Fields inherited from RecoveryPlanProviderSpecificFailoverInput

	InstanceType string `json:"instanceType"`
}

func (s RecoveryPlanInMageFailoverInput) RecoveryPlanProviderSpecificFailoverInput() BaseRecoveryPlanProviderSpecificFailoverInputImpl {
	return BaseRecoveryPlanProviderSpecificFailoverInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = RecoveryPlanInMageFailoverInput{}

func (s RecoveryPlanInMageFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper RecoveryPlanInMageFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RecoveryPlanInMageFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanInMageFailoverInput: %+v", err)
	}

	decoded["instanceType"] = "InMage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RecoveryPlanInMageFailoverInput: %+v", err)
	}

	return encoded, nil
}
