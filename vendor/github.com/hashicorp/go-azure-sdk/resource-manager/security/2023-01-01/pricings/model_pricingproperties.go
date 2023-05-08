package pricings

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PricingProperties struct {
	Deprecated             *bool        `json:"deprecated,omitempty"`
	EnablementTime         *string      `json:"enablementTime,omitempty"`
	Extensions             *[]Extension `json:"extensions,omitempty"`
	FreeTrialRemainingTime *string      `json:"freeTrialRemainingTime,omitempty"`
	PricingTier            PricingTier  `json:"pricingTier"`
	ReplacedBy             *[]string    `json:"replacedBy,omitempty"`
	SubPlan                *string      `json:"subPlan,omitempty"`
}

func (o *PricingProperties) GetEnablementTimeAsTime() (*time.Time, error) {
	if o.EnablementTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EnablementTime, "2006-01-02T15:04:05Z07:00")
}

func (o *PricingProperties) SetEnablementTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EnablementTime = &formatted
}
