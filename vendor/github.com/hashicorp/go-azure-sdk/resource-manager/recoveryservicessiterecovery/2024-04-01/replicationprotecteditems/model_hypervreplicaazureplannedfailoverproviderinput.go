package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PlannedFailoverProviderSpecificFailoverInput = HyperVReplicaAzurePlannedFailoverProviderInput{}

type HyperVReplicaAzurePlannedFailoverProviderInput struct {
	OsUpgradeVersion           *string `json:"osUpgradeVersion,omitempty"`
	PrimaryKekCertificatePfx   *string `json:"primaryKekCertificatePfx,omitempty"`
	RecoveryPointId            *string `json:"recoveryPointId,omitempty"`
	SecondaryKekCertificatePfx *string `json:"secondaryKekCertificatePfx,omitempty"`

	// Fields inherited from PlannedFailoverProviderSpecificFailoverInput

	InstanceType string `json:"instanceType"`
}

func (s HyperVReplicaAzurePlannedFailoverProviderInput) PlannedFailoverProviderSpecificFailoverInput() BasePlannedFailoverProviderSpecificFailoverInputImpl {
	return BasePlannedFailoverProviderSpecificFailoverInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = HyperVReplicaAzurePlannedFailoverProviderInput{}

func (s HyperVReplicaAzurePlannedFailoverProviderInput) MarshalJSON() ([]byte, error) {
	type wrapper HyperVReplicaAzurePlannedFailoverProviderInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVReplicaAzurePlannedFailoverProviderInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVReplicaAzurePlannedFailoverProviderInput: %+v", err)
	}

	decoded["instanceType"] = "HyperVReplicaAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVReplicaAzurePlannedFailoverProviderInput: %+v", err)
	}

	return encoded, nil
}
