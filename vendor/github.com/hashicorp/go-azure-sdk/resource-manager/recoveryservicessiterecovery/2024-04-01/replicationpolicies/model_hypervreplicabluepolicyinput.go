package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificInput = HyperVReplicaBluePolicyInput{}

type HyperVReplicaBluePolicyInput struct {
	ReplicationFrequencyInSeconds *int64 `json:"replicationFrequencyInSeconds,omitempty"`

	// Fields inherited from PolicyProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s HyperVReplicaBluePolicyInput) PolicyProviderSpecificInput() BasePolicyProviderSpecificInputImpl {
	return BasePolicyProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = HyperVReplicaBluePolicyInput{}

func (s HyperVReplicaBluePolicyInput) MarshalJSON() ([]byte, error) {
	type wrapper HyperVReplicaBluePolicyInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVReplicaBluePolicyInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVReplicaBluePolicyInput: %+v", err)
	}

	decoded["instanceType"] = "HyperVReplica2012R2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVReplicaBluePolicyInput: %+v", err)
	}

	return encoded, nil
}
