package databases

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SuspensionDetails struct {
	SuspensionStartDate *string `json:"suspensionStartDate,omitempty"`
}

func (o *SuspensionDetails) GetSuspensionStartDateAsTime() (*time.Time, error) {
	if o.SuspensionStartDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SuspensionStartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *SuspensionDetails) SetSuspensionStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SuspensionStartDate = &formatted
}
