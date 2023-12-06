package subscription

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionContractProperties struct {
	AllowTracing     *bool             `json:"allowTracing,omitempty"`
	CreatedDate      *string           `json:"createdDate,omitempty"`
	DisplayName      *string           `json:"displayName,omitempty"`
	EndDate          *string           `json:"endDate,omitempty"`
	ExpirationDate   *string           `json:"expirationDate,omitempty"`
	NotificationDate *string           `json:"notificationDate,omitempty"`
	OwnerId          *string           `json:"ownerId,omitempty"`
	PrimaryKey       *string           `json:"primaryKey,omitempty"`
	Scope            string            `json:"scope"`
	SecondaryKey     *string           `json:"secondaryKey,omitempty"`
	StartDate        *string           `json:"startDate,omitempty"`
	State            SubscriptionState `json:"state"`
	StateComment     *string           `json:"stateComment,omitempty"`
}

func (o *SubscriptionContractProperties) GetCreatedDateAsTime() (*time.Time, error) {
	if o.CreatedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *SubscriptionContractProperties) SetCreatedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDate = &formatted
}

func (o *SubscriptionContractProperties) GetEndDateAsTime() (*time.Time, error) {
	if o.EndDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndDate, "2006-01-02T15:04:05Z07:00")
}

func (o *SubscriptionContractProperties) SetEndDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndDate = &formatted
}

func (o *SubscriptionContractProperties) GetExpirationDateAsTime() (*time.Time, error) {
	if o.ExpirationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *SubscriptionContractProperties) SetExpirationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationDate = &formatted
}

func (o *SubscriptionContractProperties) GetNotificationDateAsTime() (*time.Time, error) {
	if o.NotificationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NotificationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *SubscriptionContractProperties) SetNotificationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NotificationDate = &formatted
}

func (o *SubscriptionContractProperties) GetStartDateAsTime() (*time.Time, error) {
	if o.StartDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *SubscriptionContractProperties) SetStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartDate = &formatted
}
