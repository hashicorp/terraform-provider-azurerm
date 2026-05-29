package componentfeaturesandpricingapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationInsightsComponentQuotaStatus struct {
	AppId             *string `json:"AppId,omitempty"`
	ExpirationTime    *string `json:"ExpirationTime,omitempty"`
	ShouldBeThrottled *bool   `json:"ShouldBeThrottled,omitempty"`
}
