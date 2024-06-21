package workspaces

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OutboundRule = PrivateEndpointOutboundRule{}

type PrivateEndpointOutboundRule struct {
	Destination *PrivateEndpointDestination `json:"destination,omitempty"`

	// Fields inherited from OutboundRule
	Category *RuleCategory `json:"category,omitempty"`
	Status   *RuleStatus   `json:"status,omitempty"`
}

var _ json.Marshaler = PrivateEndpointOutboundRule{}

func (s PrivateEndpointOutboundRule) MarshalJSON() ([]byte, error) {
	type wrapper PrivateEndpointOutboundRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PrivateEndpointOutboundRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PrivateEndpointOutboundRule: %+v", err)
	}
	decoded["type"] = "PrivateEndpoint"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PrivateEndpointOutboundRule: %+v", err)
	}

	return encoded, nil
}
