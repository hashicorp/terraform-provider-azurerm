package grafanaresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MarketplaceTrialQuota struct {
	AvailablePromotion *AvailablePromotion `json:"availablePromotion,omitempty"`
	GrafanaResourceId  *string             `json:"grafanaResourceId,omitempty"`
	TrialEndAt         *string             `json:"trialEndAt,omitempty"`
	TrialStartAt       *string             `json:"trialStartAt,omitempty"`
}

func (o *MarketplaceTrialQuota) GetTrialEndAtAsTime() (*time.Time, error) {
	if o.TrialEndAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TrialEndAt, "2006-01-02T15:04:05Z07:00")
}

func (o *MarketplaceTrialQuota) SetTrialEndAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TrialEndAt = &formatted
}

func (o *MarketplaceTrialQuota) GetTrialStartAtAsTime() (*time.Time, error) {
	if o.TrialStartAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TrialStartAt, "2006-01-02T15:04:05Z07:00")
}

func (o *MarketplaceTrialQuota) SetTrialStartAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TrialStartAt = &formatted
}
