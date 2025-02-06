package componentfeaturesandpricingapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationInsightsComponentBillingFeatures struct {
	CurrentBillingFeatures *[]string                                  `json:"CurrentBillingFeatures,omitempty"`
	DataVolumeCap          *ApplicationInsightsComponentDataVolumeCap `json:"DataVolumeCap,omitempty"`
}
