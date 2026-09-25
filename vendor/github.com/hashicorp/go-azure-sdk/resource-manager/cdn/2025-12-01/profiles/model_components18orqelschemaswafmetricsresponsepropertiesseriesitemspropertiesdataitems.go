package profiles

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Components18OrqelSchemasWafmetricsresponsePropertiesSeriesItemsPropertiesDataItems struct {
	DateTime *string  `json:"dateTime,omitempty"`
	Value    *float64 `json:"value,omitempty"`
}

func (o *Components18OrqelSchemasWafmetricsresponsePropertiesSeriesItemsPropertiesDataItems) GetDateTimeAsTime() (*time.Time, error) {
	if o.DateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *Components18OrqelSchemasWafmetricsresponsePropertiesSeriesItemsPropertiesDataItems) SetDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DateTime = &formatted
}
