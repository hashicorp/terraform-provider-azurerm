package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSubscriptionDestination = HybridConnectionEventSubscriptionDestination{}

type HybridConnectionEventSubscriptionDestination struct {
	Properties *HybridConnectionEventSubscriptionDestinationProperties `json:"properties,omitempty"`

	// Fields inherited from EventSubscriptionDestination
}

var _ json.Marshaler = HybridConnectionEventSubscriptionDestination{}

func (s HybridConnectionEventSubscriptionDestination) MarshalJSON() ([]byte, error) {
	type wrapper HybridConnectionEventSubscriptionDestination
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HybridConnectionEventSubscriptionDestination: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HybridConnectionEventSubscriptionDestination: %+v", err)
	}
	decoded["endpointType"] = "HybridConnection"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HybridConnectionEventSubscriptionDestination: %+v", err)
	}

	return encoded, nil
}
