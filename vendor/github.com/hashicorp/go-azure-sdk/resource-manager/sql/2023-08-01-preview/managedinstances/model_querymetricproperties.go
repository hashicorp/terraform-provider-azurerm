package managedinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryMetricProperties struct {
	Avg         *float64             `json:"avg,omitempty"`
	DisplayName *string              `json:"displayName,omitempty"`
	Max         *float64             `json:"max,omitempty"`
	Min         *float64             `json:"min,omitempty"`
	Name        *string              `json:"name,omitempty"`
	Stdev       *float64             `json:"stdev,omitempty"`
	Sum         *float64             `json:"sum,omitempty"`
	Unit        *QueryMetricUnitType `json:"unit,omitempty"`
	Value       *float64             `json:"value,omitempty"`
}
