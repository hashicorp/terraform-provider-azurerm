package appplatform

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TriggeredBuildResult struct {
	Id                   *string                                `json:"id,omitempty"`
	Image                *string                                `json:"image,omitempty"`
	LastTransitionReason *string                                `json:"lastTransitionReason,omitempty"`
	LastTransitionStatus *string                                `json:"lastTransitionStatus,omitempty"`
	LastTransitionTime   *string                                `json:"lastTransitionTime,omitempty"`
	ProvisioningState    *TriggeredBuildResultProvisioningState `json:"provisioningState,omitempty"`
}

func (o *TriggeredBuildResult) GetLastTransitionTimeAsTime() (*time.Time, error) {
	if o.LastTransitionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastTransitionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *TriggeredBuildResult) SetLastTransitionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastTransitionTime = &formatted
}
