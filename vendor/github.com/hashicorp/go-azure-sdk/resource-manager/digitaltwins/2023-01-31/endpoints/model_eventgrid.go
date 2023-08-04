package endpoints

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DigitalTwinsEndpointResourceProperties = EventGrid{}

type EventGrid struct {
	AccessKey1    string  `json:"accessKey1"`
	AccessKey2    *string `json:"accessKey2,omitempty"`
	TopicEndpoint string  `json:"TopicEndpoint"`

	// Fields inherited from DigitalTwinsEndpointResourceProperties
	AuthenticationType *AuthenticationType        `json:"authenticationType,omitempty"`
	CreatedTime        *string                    `json:"createdTime,omitempty"`
	DeadLetterSecret   *string                    `json:"deadLetterSecret,omitempty"`
	DeadLetterUri      *string                    `json:"deadLetterUri,omitempty"`
	Identity           *ManagedIdentityReference  `json:"identity,omitempty"`
	ProvisioningState  *EndpointProvisioningState `json:"provisioningState,omitempty"`
}

func (o *EventGrid) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *EventGrid) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}

var _ json.Marshaler = EventGrid{}

func (s EventGrid) MarshalJSON() ([]byte, error) {
	type wrapper EventGrid
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EventGrid: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EventGrid: %+v", err)
	}
	decoded["endpointType"] = "EventGrid"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EventGrid: %+v", err)
	}

	return encoded, nil
}
