package tenants

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TenantProperties struct {
	BillingConfig *BillingConfig `json:"billingConfig,omitempty"`
	CountryCode   *string        `json:"countryCode,omitempty"`
	DisplayName   *string        `json:"displayName,omitempty"`
	TenantId      *string        `json:"tenantId,omitempty"`
}
