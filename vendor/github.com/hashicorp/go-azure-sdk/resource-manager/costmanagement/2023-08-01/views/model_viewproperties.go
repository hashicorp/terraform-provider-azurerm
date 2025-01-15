package views

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ViewProperties struct {
	Accumulated *AccumulatedType        `json:"accumulated,omitempty"`
	Chart       *ChartType              `json:"chart,omitempty"`
	CreatedOn   *string                 `json:"createdOn,omitempty"`
	Currency    *string                 `json:"currency,omitempty"`
	DateRange   *string                 `json:"dateRange,omitempty"`
	DisplayName *string                 `json:"displayName,omitempty"`
	Kpis        *[]KpiProperties        `json:"kpis,omitempty"`
	Metric      *MetricType             `json:"metric,omitempty"`
	ModifiedOn  *string                 `json:"modifiedOn,omitempty"`
	Pivots      *[]PivotProperties      `json:"pivots,omitempty"`
	Query       *ReportConfigDefinition `json:"query,omitempty"`
	Scope       *string                 `json:"scope,omitempty"`
}

func (o *ViewProperties) GetCreatedOnAsTime() (*time.Time, error) {
	if o.CreatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *ViewProperties) SetCreatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedOn = &formatted
}

func (o *ViewProperties) GetModifiedOnAsTime() (*time.Time, error) {
	if o.ModifiedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ModifiedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *ViewProperties) SetModifiedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ModifiedOn = &formatted
}
