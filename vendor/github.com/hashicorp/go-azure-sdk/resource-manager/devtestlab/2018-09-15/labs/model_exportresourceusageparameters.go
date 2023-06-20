package labs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportResourceUsageParameters struct {
	BlobStorageAbsoluteSasUri *string `json:"blobStorageAbsoluteSasUri,omitempty"`
	UsageStartDate            *string `json:"usageStartDate,omitempty"`
}

func (o *ExportResourceUsageParameters) GetUsageStartDateAsTime() (*time.Time, error) {
	if o.UsageStartDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UsageStartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ExportResourceUsageParameters) SetUsageStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UsageStartDate = &formatted
}
