package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WafRankingsResponseDataItem struct {
	GroupValues *[]string                                                                              `json:"groupValues,omitempty"`
	Metrics     *[]ComponentsKpo1PjSchemasWafrankingsresponsePropertiesDataItemsPropertiesMetricsItems `json:"metrics,omitempty"`
}
