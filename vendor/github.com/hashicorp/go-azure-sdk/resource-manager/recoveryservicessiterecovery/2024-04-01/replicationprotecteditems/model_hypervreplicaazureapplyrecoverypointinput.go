package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ApplyRecoveryPointProviderSpecificInput = HyperVReplicaAzureApplyRecoveryPointInput{}

type HyperVReplicaAzureApplyRecoveryPointInput struct {
	PrimaryKekCertificatePfx   *string `json:"primaryKekCertificatePfx,omitempty"`
	SecondaryKekCertificatePfx *string `json:"secondaryKekCertificatePfx,omitempty"`

	// Fields inherited from ApplyRecoveryPointProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s HyperVReplicaAzureApplyRecoveryPointInput) ApplyRecoveryPointProviderSpecificInput() BaseApplyRecoveryPointProviderSpecificInputImpl {
	return BaseApplyRecoveryPointProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = HyperVReplicaAzureApplyRecoveryPointInput{}

func (s HyperVReplicaAzureApplyRecoveryPointInput) MarshalJSON() ([]byte, error) {
	type wrapper HyperVReplicaAzureApplyRecoveryPointInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVReplicaAzureApplyRecoveryPointInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVReplicaAzureApplyRecoveryPointInput: %+v", err)
	}

	decoded["instanceType"] = "HyperVReplicaAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVReplicaAzureApplyRecoveryPointInput: %+v", err)
	}

	return encoded, nil
}
