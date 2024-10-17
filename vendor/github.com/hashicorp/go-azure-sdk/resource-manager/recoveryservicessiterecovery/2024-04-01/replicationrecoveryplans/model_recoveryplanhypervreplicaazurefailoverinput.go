package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RecoveryPlanProviderSpecificFailoverInput = RecoveryPlanHyperVReplicaAzureFailoverInput{}

type RecoveryPlanHyperVReplicaAzureFailoverInput struct {
	PrimaryKekCertificatePfx   *string                                `json:"primaryKekCertificatePfx,omitempty"`
	RecoveryPointType          *HyperVReplicaAzureRpRecoveryPointType `json:"recoveryPointType,omitempty"`
	SecondaryKekCertificatePfx *string                                `json:"secondaryKekCertificatePfx,omitempty"`

	// Fields inherited from RecoveryPlanProviderSpecificFailoverInput

	InstanceType string `json:"instanceType"`
}

func (s RecoveryPlanHyperVReplicaAzureFailoverInput) RecoveryPlanProviderSpecificFailoverInput() BaseRecoveryPlanProviderSpecificFailoverInputImpl {
	return BaseRecoveryPlanProviderSpecificFailoverInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = RecoveryPlanHyperVReplicaAzureFailoverInput{}

func (s RecoveryPlanHyperVReplicaAzureFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper RecoveryPlanHyperVReplicaAzureFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RecoveryPlanHyperVReplicaAzureFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanHyperVReplicaAzureFailoverInput: %+v", err)
	}

	decoded["instanceType"] = "HyperVReplicaAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RecoveryPlanHyperVReplicaAzureFailoverInput: %+v", err)
	}

	return encoded, nil
}
