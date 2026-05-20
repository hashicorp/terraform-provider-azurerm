package eventsubscriptions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionUpdateParametersProperties struct {
	DeliveryConfiguration *DeliveryConfiguration `json:"deliveryConfiguration,omitempty"`
	EventDeliverySchema   *DeliverySchema        `json:"eventDeliverySchema,omitempty"`
	ExpirationTimeUtc     *string                `json:"expirationTimeUtc,omitempty"`
	FiltersConfiguration  *FiltersConfiguration  `json:"filtersConfiguration,omitempty"`
}

func (o *SubscriptionUpdateParametersProperties) GetExpirationTimeUtcAsTime() (*time.Time, error) {
	if o.ExpirationTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *SubscriptionUpdateParametersProperties) SetExpirationTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationTimeUtc = &formatted
}
