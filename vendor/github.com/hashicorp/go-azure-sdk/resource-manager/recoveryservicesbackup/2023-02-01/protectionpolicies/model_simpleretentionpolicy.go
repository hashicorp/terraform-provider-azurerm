package protectionpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RetentionPolicy = SimpleRetentionPolicy{}

type SimpleRetentionPolicy struct {
	RetentionDuration *RetentionDuration `json:"retentionDuration,omitempty"`

	// Fields inherited from RetentionPolicy
}

var _ json.Marshaler = SimpleRetentionPolicy{}

func (s SimpleRetentionPolicy) MarshalJSON() ([]byte, error) {
	type wrapper SimpleRetentionPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SimpleRetentionPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SimpleRetentionPolicy: %+v", err)
	}
	decoded["retentionPolicyType"] = "SimpleRetentionPolicy"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SimpleRetentionPolicy: %+v", err)
	}

	return encoded, nil
}
