package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RecoveryPlanProviderSpecificFailoverInput = RecoveryPlanInMageRcmFailbackFailoverInput{}

type RecoveryPlanInMageRcmFailbackFailoverInput struct {
	RecoveryPointType   InMageRcmFailbackRecoveryPointType `json:"recoveryPointType"`
	UseMultiVMSyncPoint *string                            `json:"useMultiVmSyncPoint,omitempty"`

	// Fields inherited from RecoveryPlanProviderSpecificFailoverInput
}

var _ json.Marshaler = RecoveryPlanInMageRcmFailbackFailoverInput{}

func (s RecoveryPlanInMageRcmFailbackFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper RecoveryPlanInMageRcmFailbackFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RecoveryPlanInMageRcmFailbackFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanInMageRcmFailbackFailoverInput: %+v", err)
	}
	decoded["instanceType"] = "InMageRcmFailback"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RecoveryPlanInMageRcmFailbackFailoverInput: %+v", err)
	}

	return encoded, nil
}
