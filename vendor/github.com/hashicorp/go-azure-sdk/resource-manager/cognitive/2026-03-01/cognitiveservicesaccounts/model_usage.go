package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Usage struct {
	CurrentValue  *float64          `json:"currentValue,omitempty"`
	Limit         *float64          `json:"limit,omitempty"`
	Name          *MetricName       `json:"name,omitempty"`
	NextResetTime *string           `json:"nextResetTime,omitempty"`
	QuotaPeriod   *string           `json:"quotaPeriod,omitempty"`
	Status        *QuotaUsageStatus `json:"status,omitempty"`
	Unit          *UnitType         `json:"unit,omitempty"`
}
