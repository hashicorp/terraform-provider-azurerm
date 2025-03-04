package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSubscriptionDestination = MonitorAlertEventSubscriptionDestination{}

type MonitorAlertEventSubscriptionDestination struct {
	Properties *MonitorAlertEventSubscriptionDestinationProperties `json:"properties,omitempty"`

	// Fields inherited from EventSubscriptionDestination

	EndpointType EndpointType `json:"endpointType"`
}

func (s MonitorAlertEventSubscriptionDestination) EventSubscriptionDestination() BaseEventSubscriptionDestinationImpl {
	return BaseEventSubscriptionDestinationImpl{
		EndpointType: s.EndpointType,
	}
}

var _ json.Marshaler = MonitorAlertEventSubscriptionDestination{}

func (s MonitorAlertEventSubscriptionDestination) MarshalJSON() ([]byte, error) {
	type wrapper MonitorAlertEventSubscriptionDestination
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MonitorAlertEventSubscriptionDestination: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MonitorAlertEventSubscriptionDestination: %+v", err)
	}

	decoded["endpointType"] = "MonitorAlert"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MonitorAlertEventSubscriptionDestination: %+v", err)
	}

	return encoded, nil
}
