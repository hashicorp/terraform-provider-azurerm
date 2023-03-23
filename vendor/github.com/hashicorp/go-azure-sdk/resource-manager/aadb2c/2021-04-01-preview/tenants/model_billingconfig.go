package tenants

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BillingConfig struct {
	BillingType           *BillingType `json:"billingType,omitempty"`
	EffectiveStartDateUtc *string      `json:"effectiveStartDateUtc,omitempty"`
}
