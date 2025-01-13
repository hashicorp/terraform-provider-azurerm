package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TestFailoverProviderSpecificInput = HyperVReplicaAzureTestFailoverInput{}

type HyperVReplicaAzureTestFailoverInput struct {
	OsUpgradeVersion           *string `json:"osUpgradeVersion,omitempty"`
	PrimaryKekCertificatePfx   *string `json:"primaryKekCertificatePfx,omitempty"`
	RecoveryPointId            *string `json:"recoveryPointId,omitempty"`
	SecondaryKekCertificatePfx *string `json:"secondaryKekCertificatePfx,omitempty"`

	// Fields inherited from TestFailoverProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s HyperVReplicaAzureTestFailoverInput) TestFailoverProviderSpecificInput() BaseTestFailoverProviderSpecificInputImpl {
	return BaseTestFailoverProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = HyperVReplicaAzureTestFailoverInput{}

func (s HyperVReplicaAzureTestFailoverInput) MarshalJSON() ([]byte, error) {
	type wrapper HyperVReplicaAzureTestFailoverInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVReplicaAzureTestFailoverInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVReplicaAzureTestFailoverInput: %+v", err)
	}

	decoded["instanceType"] = "HyperVReplicaAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVReplicaAzureTestFailoverInput: %+v", err)
	}

	return encoded, nil
}
