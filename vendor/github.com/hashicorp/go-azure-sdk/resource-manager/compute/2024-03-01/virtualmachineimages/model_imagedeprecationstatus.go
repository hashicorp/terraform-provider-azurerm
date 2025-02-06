package virtualmachineimages

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageDeprecationStatus struct {
	AlternativeOption        *AlternativeOption `json:"alternativeOption,omitempty"`
	ImageState               *ImageState        `json:"imageState,omitempty"`
	ScheduledDeprecationTime *string            `json:"scheduledDeprecationTime,omitempty"`
}

func (o *ImageDeprecationStatus) GetScheduledDeprecationTimeAsTime() (*time.Time, error) {
	if o.ScheduledDeprecationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ScheduledDeprecationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ImageDeprecationStatus) SetScheduledDeprecationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ScheduledDeprecationTime = &formatted
}
