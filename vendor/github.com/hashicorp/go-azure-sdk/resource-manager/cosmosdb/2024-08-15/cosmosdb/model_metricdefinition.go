package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetricDefinition struct {
	MetricAvailabilities   *[]MetricAvailability   `json:"metricAvailabilities,omitempty"`
	Name                   *MetricName             `json:"name,omitempty"`
	PrimaryAggregationType *PrimaryAggregationType `json:"primaryAggregationType,omitempty"`
	ResourceUri            *string                 `json:"resourceUri,omitempty"`
	Unit                   *UnitType               `json:"unit,omitempty"`
}
