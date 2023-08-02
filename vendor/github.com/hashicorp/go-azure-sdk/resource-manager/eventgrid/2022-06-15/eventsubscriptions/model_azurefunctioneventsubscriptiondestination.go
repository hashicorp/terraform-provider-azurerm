package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSubscriptionDestination = AzureFunctionEventSubscriptionDestination{}

type AzureFunctionEventSubscriptionDestination struct {
	Properties *AzureFunctionEventSubscriptionDestinationProperties `json:"properties,omitempty"`

	// Fields inherited from EventSubscriptionDestination
}

var _ json.Marshaler = AzureFunctionEventSubscriptionDestination{}

func (s AzureFunctionEventSubscriptionDestination) MarshalJSON() ([]byte, error) {
	type wrapper AzureFunctionEventSubscriptionDestination
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureFunctionEventSubscriptionDestination: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureFunctionEventSubscriptionDestination: %+v", err)
	}
	decoded["endpointType"] = "AzureFunction"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureFunctionEventSubscriptionDestination: %+v", err)
	}

	return encoded, nil
}
