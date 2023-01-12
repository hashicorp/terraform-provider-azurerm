package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RecoveryPlanProviderSpecificFailoverInput = RecoveryPlanInMageAzureV2FailoverInput{}

type RecoveryPlanInMageAzureV2FailoverInput struct {
	RecoveryPointType   InMageV2RpRecoveryPointType `json:"recoveryPointType"`
	UseMultiVMSyncPoint *string                     `json:"useMultiVmSyncPoint,omitempty"`

	// Fields inherited from RecoveryPlanProviderSpecificFailoverInput
}

var _ json.Marshaler = RecoveryPlanInMageAzureV2FailoverInput{}

func (s RecoveryPlanInMageAzureV2FailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper RecoveryPlanInMageAzureV2FailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RecoveryPlanInMageAzureV2FailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanInMageAzureV2FailoverInput: %+v", err)
	}
	decoded["instanceType"] = "InMageAzureV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RecoveryPlanInMageAzureV2FailoverInput: %+v", err)
	}

	return encoded, nil
}
