package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetricsResponseSeriesItem struct {
	Data   *[]Components1Gs0LlpSchemasMetricsresponsePropertiesSeriesItemsPropertiesDataItems `json:"data,omitempty"`
	Groups *[]MetricsResponseSeriesPropertiesItemsItem                                        `json:"groups,omitempty"`
	Metric *string                                                                            `json:"metric,omitempty"`
	Unit   *MetricsSeriesUnit                                                                 `json:"unit,omitempty"`
}
