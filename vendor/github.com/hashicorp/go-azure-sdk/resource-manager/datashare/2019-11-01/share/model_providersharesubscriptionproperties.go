package share

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderShareSubscriptionProperties struct {
	ConsumerEmail             *string                  `json:"consumerEmail,omitempty"`
	ConsumerName              *string                  `json:"consumerName,omitempty"`
	ConsumerTenantName        *string                  `json:"consumerTenantName,omitempty"`
	CreatedAt                 *string                  `json:"createdAt,omitempty"`
	ProviderEmail             *string                  `json:"providerEmail,omitempty"`
	ProviderName              *string                  `json:"providerName,omitempty"`
	ShareSubscriptionObjectId *string                  `json:"shareSubscriptionObjectId,omitempty"`
	ShareSubscriptionStatus   *ShareSubscriptionStatus `json:"shareSubscriptionStatus,omitempty"`
	SharedAt                  *string                  `json:"sharedAt,omitempty"`
}

func (o *ProviderShareSubscriptionProperties) GetCreatedAtAsTime() (*time.Time, error) {
	if o.CreatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *ProviderShareSubscriptionProperties) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o *ProviderShareSubscriptionProperties) GetSharedAtAsTime() (*time.Time, error) {
	if o.SharedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SharedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *ProviderShareSubscriptionProperties) SetSharedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SharedAt = &formatted
}
