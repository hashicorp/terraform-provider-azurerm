package licenseprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProductProfileUpdateProperties struct {
	ProductFeatures    *[]ProductFeatureUpdate                 `json:"productFeatures,omitempty"`
	ProductType        *LicenseProfileProductType              `json:"productType,omitempty"`
	SubscriptionStatus *LicenseProfileSubscriptionStatusUpdate `json:"subscriptionStatus,omitempty"`
}
