package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServicePlan struct {
	ExtendedLocation *ExtendedLocation         `json:"extendedLocation,omitempty"`
	Id               *string                   `json:"id,omitempty"`
	Kind             *string                   `json:"kind,omitempty"`
	Location         string                    `json:"location"`
	Name             *string                   `json:"name,omitempty"`
	Properties       *AppServicePlanProperties `json:"properties,omitempty"`
	Sku              *SkuDescription           `json:"sku,omitempty"`
	Tags             *map[string]string        `json:"tags,omitempty"`
	Type             *string                   `json:"type,omitempty"`
}
