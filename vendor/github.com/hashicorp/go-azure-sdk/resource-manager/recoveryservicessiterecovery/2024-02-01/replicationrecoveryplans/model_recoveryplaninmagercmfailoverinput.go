package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RecoveryPlanProviderSpecificFailoverInput = RecoveryPlanInMageRcmFailoverInput{}

type RecoveryPlanInMageRcmFailoverInput struct {
	RecoveryPointType   RecoveryPlanPointType `json:"recoveryPointType"`
	UseMultiVMSyncPoint *string               `json:"useMultiVmSyncPoint,omitempty"`

	// Fields inherited from RecoveryPlanProviderSpecificFailoverInput
}

var _ json.Marshaler = RecoveryPlanInMageRcmFailoverInput{}

func (s RecoveryPlanInMageRcmFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper RecoveryPlanInMageRcmFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RecoveryPlanInMageRcmFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanInMageRcmFailoverInput: %+v", err)
	}
	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RecoveryPlanInMageRcmFailoverInput: %+v", err)
	}

	return encoded, nil
}
