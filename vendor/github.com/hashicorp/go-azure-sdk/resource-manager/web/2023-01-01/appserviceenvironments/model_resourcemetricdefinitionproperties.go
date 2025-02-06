package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceMetricDefinitionProperties struct {
	MetricAvailabilities   *[]ResourceMetricAvailability `json:"metricAvailabilities,omitempty"`
	PrimaryAggregationType *string                       `json:"primaryAggregationType,omitempty"`
	Properties             *map[string]string            `json:"properties,omitempty"`
	ResourceUri            *string                       `json:"resourceUri,omitempty"`
	Unit                   *string                       `json:"unit,omitempty"`
}
