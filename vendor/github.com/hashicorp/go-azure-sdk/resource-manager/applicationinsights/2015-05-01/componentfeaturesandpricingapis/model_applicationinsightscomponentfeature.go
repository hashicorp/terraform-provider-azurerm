package componentfeaturesandpricingapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationInsightsComponentFeature struct {
	Capabilities           *[]ApplicationInsightsComponentFeatureCapability `json:"Capabilities,omitempty"`
	FeatureName            *string                                          `json:"FeatureName,omitempty"`
	IsHidden               *bool                                            `json:"IsHidden,omitempty"`
	IsMainFeature          *bool                                            `json:"IsMainFeature,omitempty"`
	MeterId                *string                                          `json:"MeterId,omitempty"`
	MeterRateFrequency     *string                                          `json:"MeterRateFrequency,omitempty"`
	ResouceId              *string                                          `json:"ResouceId,omitempty"`
	SupportedAddonFeatures *string                                          `json:"SupportedAddonFeatures,omitempty"`
	Title                  *string                                          `json:"Title,omitempty"`
}
