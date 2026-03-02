package views

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReportConfigDefinition struct {
	DataSet                   *ReportConfigDataset    `json:"dataSet,omitempty"`
	IncludeMonetaryCommitment *bool                   `json:"includeMonetaryCommitment,omitempty"`
	TimePeriod                *ReportConfigTimePeriod `json:"timePeriod,omitempty"`
	Timeframe                 ReportTimeframeType     `json:"timeframe"`
	Type                      ReportType              `json:"type"`
}
