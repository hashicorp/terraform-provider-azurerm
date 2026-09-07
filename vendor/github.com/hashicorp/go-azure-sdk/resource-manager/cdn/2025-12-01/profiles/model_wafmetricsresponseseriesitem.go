package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WafMetricsResponseSeriesItem struct {
	Data   *[]Components18OrqelSchemasWafmetricsresponsePropertiesSeriesItemsPropertiesDataItems `json:"data,omitempty"`
	Groups *[]WafMetricsResponseSeriesPropertiesItemsItem                                        `json:"groups,omitempty"`
	Metric *string                                                                               `json:"metric,omitempty"`
	Unit   *WafMetricsSeriesUnit                                                                 `json:"unit,omitempty"`
}
