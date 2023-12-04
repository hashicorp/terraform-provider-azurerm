package grafanaresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnterpriseConfigurations struct {
	MarketplaceAutoRenew *MarketplaceAutoRenew `json:"marketplaceAutoRenew,omitempty"`
	MarketplacePlanId    *string               `json:"marketplacePlanId,omitempty"`
}
