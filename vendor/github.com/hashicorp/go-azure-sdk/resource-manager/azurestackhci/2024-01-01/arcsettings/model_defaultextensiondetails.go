package arcsettings

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DefaultExtensionDetails struct {
	Category    *string `json:"category,omitempty"`
	ConsentTime *string `json:"consentTime,omitempty"`
}

func (o *DefaultExtensionDetails) GetConsentTimeAsTime() (*time.Time, error) {
	if o.ConsentTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ConsentTime, "2006-01-02T15:04:05Z07:00")
}

func (o *DefaultExtensionDetails) SetConsentTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ConsentTime = &formatted
}
