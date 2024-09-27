package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BillingResponseListResult struct {
	BillingResources            *[]BillingResources            `json:"billingResources,omitempty"`
	VMSizeFilters               *[]VMSizeCompatibilityFilterV2 `json:"vmSizeFilters,omitempty"`
	VMSizeProperties            *[]VMSizeProperty              `json:"vmSizeProperties,omitempty"`
	VMSizes                     *[]string                      `json:"vmSizes,omitempty"`
	VMSizesWithEncryptionAtHost *[]string                      `json:"vmSizesWithEncryptionAtHost,omitempty"`
}
