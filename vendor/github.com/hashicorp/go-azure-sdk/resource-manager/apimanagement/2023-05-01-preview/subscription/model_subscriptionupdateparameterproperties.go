package subscription

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionUpdateParameterProperties struct {
	AllowTracing   *bool              `json:"allowTracing,omitempty"`
	DisplayName    *string            `json:"displayName,omitempty"`
	ExpirationDate *string            `json:"expirationDate,omitempty"`
	OwnerId        *string            `json:"ownerId,omitempty"`
	PrimaryKey     *string            `json:"primaryKey,omitempty"`
	Scope          *string            `json:"scope,omitempty"`
	SecondaryKey   *string            `json:"secondaryKey,omitempty"`
	State          *SubscriptionState `json:"state,omitempty"`
	StateComment   *string            `json:"stateComment,omitempty"`
}

func (o *SubscriptionUpdateParameterProperties) GetExpirationDateAsTime() (*time.Time, error) {
	if o.ExpirationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *SubscriptionUpdateParameterProperties) SetExpirationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationDate = &formatted
}
