package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Usage struct {
	CurrentValue *int64      `json:"currentValue,omitempty"`
	Limit        *int64      `json:"limit,omitempty"`
	Name         *MetricName `json:"name,omitempty"`
	QuotaPeriod  *string     `json:"quotaPeriod,omitempty"`
	Unit         *UnitType   `json:"unit,omitempty"`
}
