package views

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReportConfigDataset struct {
	Aggregation   *map[string]ReportConfigAggregation `json:"aggregation,omitempty"`
	Configuration *ReportConfigDatasetConfiguration   `json:"configuration,omitempty"`
	Filter        *ReportConfigFilter                 `json:"filter,omitempty"`
	Granularity   *ReportGranularityType              `json:"granularity,omitempty"`
	Grouping      *[]ReportConfigGrouping             `json:"grouping,omitempty"`
	Sorting       *[]ReportConfigSorting              `json:"sorting,omitempty"`
}
