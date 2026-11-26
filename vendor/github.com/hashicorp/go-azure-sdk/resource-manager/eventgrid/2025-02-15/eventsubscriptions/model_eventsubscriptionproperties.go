package eventsubscriptions

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSubscriptionProperties struct {
	DeadLetterDestination          DeadLetterDestination               `json:"deadLetterDestination"`
	DeadLetterWithResourceIdentity *DeadLetterWithResourceIdentity     `json:"deadLetterWithResourceIdentity,omitempty"`
	DeliveryWithResourceIdentity   *DeliveryWithResourceIdentity       `json:"deliveryWithResourceIdentity,omitempty"`
	Destination                    EventSubscriptionDestination        `json:"destination"`
	EventDeliverySchema            *EventDeliverySchema                `json:"eventDeliverySchema,omitempty"`
	ExpirationTimeUtc              *string                             `json:"expirationTimeUtc,omitempty"`
	Filter                         *EventSubscriptionFilter            `json:"filter,omitempty"`
	Labels                         *[]string                           `json:"labels,omitempty"`
	ProvisioningState              *EventSubscriptionProvisioningState `json:"provisioningState,omitempty"`
	RetryPolicy                    *RetryPolicy                        `json:"retryPolicy,omitempty"`
	Topic                          *string                             `json:"topic,omitempty"`
}

func (o *EventSubscriptionProperties) GetExpirationTimeUtcAsTime() (*time.Time, error) {
	if o.ExpirationTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *EventSubscriptionProperties) SetExpirationTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationTimeUtc = &formatted
}

var _ json.Unmarshaler = &EventSubscriptionProperties{}

func (s *EventSubscriptionProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		DeadLetterWithResourceIdentity *DeadLetterWithResourceIdentity     `json:"deadLetterWithResourceIdentity,omitempty"`
		DeliveryWithResourceIdentity   *DeliveryWithResourceIdentity       `json:"deliveryWithResourceIdentity,omitempty"`
		EventDeliverySchema            *EventDeliverySchema                `json:"eventDeliverySchema,omitempty"`
		ExpirationTimeUtc              *string                             `json:"expirationTimeUtc,omitempty"`
		Filter                         *EventSubscriptionFilter            `json:"filter,omitempty"`
		Labels                         *[]string                           `json:"labels,omitempty"`
		ProvisioningState              *EventSubscriptionProvisioningState `json:"provisioningState,omitempty"`
		RetryPolicy                    *RetryPolicy                        `json:"retryPolicy,omitempty"`
		Topic                          *string                             `json:"topic,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.DeadLetterWithResourceIdentity = decoded.DeadLetterWithResourceIdentity
	s.DeliveryWithResourceIdentity = decoded.DeliveryWithResourceIdentity
	s.EventDeliverySchema = decoded.EventDeliverySchema
	s.ExpirationTimeUtc = decoded.ExpirationTimeUtc
	s.Filter = decoded.Filter
	s.Labels = decoded.Labels
	s.ProvisioningState = decoded.ProvisioningState
	s.RetryPolicy = decoded.RetryPolicy
	s.Topic = decoded.Topic

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling EventSubscriptionProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["deadLetterDestination"]; ok {
		impl, err := UnmarshalDeadLetterDestinationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'DeadLetterDestination' for 'EventSubscriptionProperties': %+v", err)
		}
		s.DeadLetterDestination = impl
	}

	if v, ok := temp["destination"]; ok {
		impl, err := UnmarshalEventSubscriptionDestinationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Destination' for 'EventSubscriptionProperties': %+v", err)
		}
		s.Destination = impl
	}

	return nil
}
