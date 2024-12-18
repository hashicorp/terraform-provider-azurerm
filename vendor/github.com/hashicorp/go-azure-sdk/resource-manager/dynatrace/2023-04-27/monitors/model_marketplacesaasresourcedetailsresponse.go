package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MarketplaceSaaSResourceDetailsResponse struct {
	MarketplaceSaaSResourceId     *string                        `json:"marketplaceSaaSResourceId,omitempty"`
	MarketplaceSubscriptionStatus *MarketplaceSubscriptionStatus `json:"marketplaceSubscriptionStatus,omitempty"`
	PlanId                        *string                        `json:"planId,omitempty"`
}
