package eventsources

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHubEventSourceCommonProperties struct {
	ConsumerGroupName     string                    `json:"consumerGroupName"`
	CreationTime          *string                   `json:"creationTime,omitempty"`
	EventHubName          string                    `json:"eventHubName"`
	EventSourceResourceId string                    `json:"eventSourceResourceId"`
	IngressStartAt        *IngressStartAtProperties `json:"ingressStartAt,omitempty"`
	KeyName               string                    `json:"keyName"`
	LocalTimestamp        *LocalTimestamp           `json:"localTimestamp,omitempty"`
	ProvisioningState     *ProvisioningState        `json:"provisioningState,omitempty"`
	ServiceBusNamespace   string                    `json:"serviceBusNamespace"`
	TimestampPropertyName *string                   `json:"timestampPropertyName,omitempty"`
}

func (o *EventHubEventSourceCommonProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *EventHubEventSourceCommonProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}
