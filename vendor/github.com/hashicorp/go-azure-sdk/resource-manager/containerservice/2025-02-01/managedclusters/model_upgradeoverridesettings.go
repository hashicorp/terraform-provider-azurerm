package managedclusters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpgradeOverrideSettings struct {
	ForceUpgrade *bool   `json:"forceUpgrade,omitempty"`
	Until        *string `json:"until,omitempty"`
}

func (o *UpgradeOverrideSettings) GetUntilAsTime() (*time.Time, error) {
	if o.Until == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Until, "2006-01-02T15:04:05Z07:00")
}

func (o *UpgradeOverrideSettings) SetUntilAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Until = &formatted
}
