package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

type DeliveryWithResourceIdentity struct {
	Destination EventSubscriptionDestination `json:"destination"`
	Identity    *EventSubscriptionIdentity   `json:"identity,omitempty"`
}

var _ json.Unmarshaler = &DeliveryWithResourceIdentity{}

func (s *DeliveryWithResourceIdentity) UnmarshalJSON(bytes []byte) error {
	type alias DeliveryWithResourceIdentity
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into DeliveryWithResourceIdentity: %+v", err)
	}

	s.Identity = decoded.Identity

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DeliveryWithResourceIdentity into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["destination"]; ok {
		impl, err := unmarshalEventSubscriptionDestinationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Destination' for 'DeliveryWithResourceIdentity': %+v", err)
		}
		s.Destination = impl
	}
	return nil
}
