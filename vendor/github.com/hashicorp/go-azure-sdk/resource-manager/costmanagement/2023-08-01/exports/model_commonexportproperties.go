package exports

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommonExportProperties struct {
	Definition          ExportDefinition           `json:"definition"`
	DeliveryInfo        ExportDeliveryInfo         `json:"deliveryInfo"`
	Format              *FormatType                `json:"format,omitempty"`
	NextRunTimeEstimate *string                    `json:"nextRunTimeEstimate,omitempty"`
	PartitionData       *bool                      `json:"partitionData,omitempty"`
	RunHistory          *ExportExecutionListResult `json:"runHistory,omitempty"`
}

func (o *CommonExportProperties) GetNextRunTimeEstimateAsTime() (*time.Time, error) {
	if o.NextRunTimeEstimate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NextRunTimeEstimate, "2006-01-02T15:04:05Z07:00")
}

func (o *CommonExportProperties) SetNextRunTimeEstimateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NextRunTimeEstimate = &formatted
}
