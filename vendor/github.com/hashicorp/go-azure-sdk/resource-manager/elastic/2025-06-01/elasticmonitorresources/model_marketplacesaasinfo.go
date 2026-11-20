package elasticmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MarketplaceSaaSInfo struct {
	BilledAzureSubscriptionId *string                                     `json:"billedAzureSubscriptionId,omitempty"`
	MarketplaceName           *string                                     `json:"marketplaceName,omitempty"`
	MarketplaceResourceId     *string                                     `json:"marketplaceResourceId,omitempty"`
	MarketplaceStatus         *string                                     `json:"marketplaceStatus,omitempty"`
	MarketplaceSubscription   *MarketplaceSaaSInfoMarketplaceSubscription `json:"marketplaceSubscription,omitempty"`
	Subscribed                *bool                                       `json:"subscribed,omitempty"`
}
