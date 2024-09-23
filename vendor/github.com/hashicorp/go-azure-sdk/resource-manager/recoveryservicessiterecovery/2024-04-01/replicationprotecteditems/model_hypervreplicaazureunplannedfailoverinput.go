package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UnplannedFailoverProviderSpecificInput = HyperVReplicaAzureUnplannedFailoverInput{}

type HyperVReplicaAzureUnplannedFailoverInput struct {
	PrimaryKekCertificatePfx   *string `json:"primaryKekCertificatePfx,omitempty"`
	RecoveryPointId            *string `json:"recoveryPointId,omitempty"`
	SecondaryKekCertificatePfx *string `json:"secondaryKekCertificatePfx,omitempty"`

	// Fields inherited from UnplannedFailoverProviderSpecificInput
}

var _ json.Marshaler = HyperVReplicaAzureUnplannedFailoverInput{}

func (s HyperVReplicaAzureUnplannedFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper HyperVReplicaAzureUnplannedFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVReplicaAzureUnplannedFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVReplicaAzureUnplannedFailoverInput: %+v", err)
	}
	decoded["instanceType"] = "HyperVReplicaAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVReplicaAzureUnplannedFailoverInput: %+v", err)
	}

	return encoded, nil
}
