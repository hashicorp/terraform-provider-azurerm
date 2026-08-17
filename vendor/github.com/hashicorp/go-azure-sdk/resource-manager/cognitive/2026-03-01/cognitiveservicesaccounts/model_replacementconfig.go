package cognitiveservicesaccounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplacementConfig struct {
	AutoUpgradeStartDate        *string `json:"autoUpgradeStartDate,omitempty"`
	TargetModelName             *string `json:"targetModelName,omitempty"`
	TargetModelVersion          *string `json:"targetModelVersion,omitempty"`
	UpgradeOnExpiryLeadTimeDays *int64  `json:"upgradeOnExpiryLeadTimeDays,omitempty"`
}

func (o *ReplacementConfig) GetAutoUpgradeStartDateAsTime() (*time.Time, error) {
	if o.AutoUpgradeStartDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.AutoUpgradeStartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ReplacementConfig) SetAutoUpgradeStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.AutoUpgradeStartDate = &formatted
}
