package vaults

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpgradeDetails struct {
	EndTimeUtc         *string            `json:"endTimeUtc,omitempty"`
	LastUpdatedTimeUtc *string            `json:"lastUpdatedTimeUtc,omitempty"`
	Message            *string            `json:"message,omitempty"`
	OperationId        *string            `json:"operationId,omitempty"`
	PreviousResourceId *string            `json:"previousResourceId,omitempty"`
	StartTimeUtc       *string            `json:"startTimeUtc,omitempty"`
	Status             *VaultUpgradeState `json:"status,omitempty"`
	TriggerType        *TriggerType       `json:"triggerType,omitempty"`
	UpgradedResourceId *string            `json:"upgradedResourceId,omitempty"`
}

func (o *UpgradeDetails) GetEndTimeUtcAsTime() (*time.Time, error) {
	if o.EndTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *UpgradeDetails) SetEndTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTimeUtc = &formatted
}

func (o *UpgradeDetails) GetLastUpdatedTimeUtcAsTime() (*time.Time, error) {
	if o.LastUpdatedTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *UpgradeDetails) SetLastUpdatedTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTimeUtc = &formatted
}

func (o *UpgradeDetails) GetStartTimeUtcAsTime() (*time.Time, error) {
	if o.StartTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *UpgradeDetails) SetStartTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTimeUtc = &formatted
}
