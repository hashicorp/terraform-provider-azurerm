package exports

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportSuspensionContext struct {
	SuspensionCode   *string `json:"suspensionCode,omitempty"`
	SuspensionReason *string `json:"suspensionReason,omitempty"`
	SuspensionTime   *string `json:"suspensionTime,omitempty"`
}

func (o *ExportSuspensionContext) GetSuspensionTimeAsTime() (*time.Time, error) {
	if o.SuspensionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SuspensionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ExportSuspensionContext) SetSuspensionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SuspensionTime = &formatted
}
