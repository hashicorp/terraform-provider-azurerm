package eventsources

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IoTHubEventSourceCommonProperties struct {
	ConsumerGroupName     string                    `json:"consumerGroupName"`
	CreationTime          *string                   `json:"creationTime,omitempty"`
	EventSourceResourceId string                    `json:"eventSourceResourceId"`
	IngressStartAt        *IngressStartAtProperties `json:"ingressStartAt,omitempty"`
	IotHubName            string                    `json:"iotHubName"`
	KeyName               string                    `json:"keyName"`
	LocalTimestamp        *LocalTimestamp           `json:"localTimestamp,omitempty"`
	ProvisioningState     *ProvisioningState        `json:"provisioningState,omitempty"`
	TimestampPropertyName *string                   `json:"timestampPropertyName,omitempty"`
}

func (o *IoTHubEventSourceCommonProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *IoTHubEventSourceCommonProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}
