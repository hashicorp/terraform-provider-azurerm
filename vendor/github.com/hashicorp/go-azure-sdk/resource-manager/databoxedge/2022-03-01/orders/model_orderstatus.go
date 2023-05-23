package orders

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrderStatus struct {
	AdditionalOrderDetails *map[string]string `json:"additionalOrderDetails,omitempty"`
	Comments               *string            `json:"comments,omitempty"`
	Status                 OrderState         `json:"status"`
	TrackingInformation    *TrackingInfo      `json:"trackingInformation,omitempty"`
	UpdateDateTime         *string            `json:"updateDateTime,omitempty"`
}

func (o *OrderStatus) GetUpdateDateTimeAsTime() (*time.Time, error) {
	if o.UpdateDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdateDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *OrderStatus) SetUpdateDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdateDateTime = &formatted
}
