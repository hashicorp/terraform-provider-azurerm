package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

var _ EventSubscriptionDestination = ServiceBusQueueEventSubscriptionDestination{}

type ServiceBusQueueEventSubscriptionDestination struct {
	Properties *ServiceBusQueueEventSubscriptionDestinationProperties `json:"properties,omitempty"`

	// Fields inherited from EventSubscriptionDestination
}

var _ json.Marshaler = ServiceBusQueueEventSubscriptionDestination{}

func (s ServiceBusQueueEventSubscriptionDestination) MarshalJSON() ([]byte, error) {
	type wrapper ServiceBusQueueEventSubscriptionDestination
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServiceBusQueueEventSubscriptionDestination: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServiceBusQueueEventSubscriptionDestination: %+v", err)
	}
	decoded["endpointType"] = "ServiceBusQueue"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServiceBusQueueEventSubscriptionDestination: %+v", err)
	}

	return encoded, nil
}
