package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PushInfo struct {
	DeadLetterDestinationWithResourceIdentity *DeadLetterWithResourceIdentity `json:"deadLetterDestinationWithResourceIdentity,omitempty"`
	DeliveryWithResourceIdentity              *DeliveryWithResourceIdentity   `json:"deliveryWithResourceIdentity,omitempty"`
	Destination                               EventSubscriptionDestination    `json:"destination"`
	EventTimeToLive                           *string                         `json:"eventTimeToLive,omitempty"`
	MaxDeliveryCount                          *int64                          `json:"maxDeliveryCount,omitempty"`
}

var _ json.Unmarshaler = &PushInfo{}

func (s *PushInfo) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		DeadLetterDestinationWithResourceIdentity *DeadLetterWithResourceIdentity `json:"deadLetterDestinationWithResourceIdentity,omitempty"`
		DeliveryWithResourceIdentity              *DeliveryWithResourceIdentity   `json:"deliveryWithResourceIdentity,omitempty"`
		EventTimeToLive                           *string                         `json:"eventTimeToLive,omitempty"`
		MaxDeliveryCount                          *int64                          `json:"maxDeliveryCount,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.DeadLetterDestinationWithResourceIdentity = decoded.DeadLetterDestinationWithResourceIdentity
	s.DeliveryWithResourceIdentity = decoded.DeliveryWithResourceIdentity
	s.EventTimeToLive = decoded.EventTimeToLive
	s.MaxDeliveryCount = decoded.MaxDeliveryCount

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling PushInfo into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["destination"]; ok {
		impl, err := UnmarshalEventSubscriptionDestinationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Destination' for 'PushInfo': %+v", err)
		}
		s.Destination = impl
	}

	return nil
}
