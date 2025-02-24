package machines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WindowsParameters struct {
	ClassificationsToInclude  *[]VMGuestPatchClassificationWindows `json:"classificationsToInclude,omitempty"`
	ExcludeKbsRequiringReboot *bool                                `json:"excludeKbsRequiringReboot,omitempty"`
	KbNumbersToExclude        *[]string                            `json:"kbNumbersToExclude,omitempty"`
	KbNumbersToInclude        *[]string                            `json:"kbNumbersToInclude,omitempty"`
	MaxPatchPublishDate       *string                              `json:"maxPatchPublishDate,omitempty"`
}

func (o *WindowsParameters) GetMaxPatchPublishDateAsTime() (*time.Time, error) {
	if o.MaxPatchPublishDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.MaxPatchPublishDate, "2006-01-02T15:04:05Z07:00")
}

func (o *WindowsParameters) SetMaxPatchPublishDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.MaxPatchPublishDate = &formatted
}
