package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RankingsResponseTablesPropertiesItemsItem struct {
	Metrics *[]RankingsResponseTablesPropertiesItemsMetricsItem `json:"metrics,omitempty"`
	Name    *string                                             `json:"name,omitempty"`
}
