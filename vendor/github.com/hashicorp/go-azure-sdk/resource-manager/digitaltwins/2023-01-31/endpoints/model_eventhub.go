package endpoints

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DigitalTwinsEndpointResourceProperties = EventHub{}

type EventHub struct {
	ConnectionStringPrimaryKey   *string `json:"connectionStringPrimaryKey,omitempty"`
	ConnectionStringSecondaryKey *string `json:"connectionStringSecondaryKey,omitempty"`
	EndpointUri                  *string `json:"endpointUri,omitempty"`
	EntityPath                   *string `json:"entityPath,omitempty"`

	// Fields inherited from DigitalTwinsEndpointResourceProperties

	AuthenticationType *AuthenticationType        `json:"authenticationType,omitempty"`
	CreatedTime        *string                    `json:"createdTime,omitempty"`
	DeadLetterSecret   *string                    `json:"deadLetterSecret,omitempty"`
	DeadLetterUri      *string                    `json:"deadLetterUri,omitempty"`
	EndpointType       EndpointType               `json:"endpointType"`
	Identity           *ManagedIdentityReference  `json:"identity,omitempty"`
	ProvisioningState  *EndpointProvisioningState `json:"provisioningState,omitempty"`
}

func (s EventHub) DigitalTwinsEndpointResourceProperties() BaseDigitalTwinsEndpointResourcePropertiesImpl {
	return BaseDigitalTwinsEndpointResourcePropertiesImpl{
		AuthenticationType: s.AuthenticationType,
		CreatedTime:        s.CreatedTime,
		DeadLetterSecret:   s.DeadLetterSecret,
		DeadLetterUri:      s.DeadLetterUri,
		EndpointType:       s.EndpointType,
		Identity:           s.Identity,
		ProvisioningState:  s.ProvisioningState,
	}
}

func (o *EventHub) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *EventHub) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}

var _ json.Marshaler = EventHub{}

func (s EventHub) MarshalJSON() ([]byte, error) {
	type wrapper EventHub
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EventHub: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EventHub: %+v", err)
	}

	decoded["endpointType"] = "EventHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EventHub: %+v", err)
	}

	return encoded, nil
}
