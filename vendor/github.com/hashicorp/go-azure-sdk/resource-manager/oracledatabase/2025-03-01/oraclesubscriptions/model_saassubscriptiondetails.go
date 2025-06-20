package oraclesubscriptions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SaasSubscriptionDetails struct {
	Id                     *string `json:"id,omitempty"`
	IsAutoRenew            *bool   `json:"isAutoRenew,omitempty"`
	IsFreeTrial            *bool   `json:"isFreeTrial,omitempty"`
	OfferId                *string `json:"offerId,omitempty"`
	PlanId                 *string `json:"planId,omitempty"`
	PublisherId            *string `json:"publisherId,omitempty"`
	PurchaserEmailId       *string `json:"purchaserEmailId,omitempty"`
	PurchaserTenantId      *string `json:"purchaserTenantId,omitempty"`
	SaasSubscriptionStatus *string `json:"saasSubscriptionStatus,omitempty"`
	SubscriptionName       *string `json:"subscriptionName,omitempty"`
	TermUnit               *string `json:"termUnit,omitempty"`
	TimeCreated            *string `json:"timeCreated,omitempty"`
}

func (o *SaasSubscriptionDetails) GetTimeCreatedAsTime() (*time.Time, error) {
	if o.TimeCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *SaasSubscriptionDetails) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = &formatted
}
