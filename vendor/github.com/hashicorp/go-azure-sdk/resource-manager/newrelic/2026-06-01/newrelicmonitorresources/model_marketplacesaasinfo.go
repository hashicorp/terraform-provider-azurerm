package newrelicmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MarketplaceSaaSInfo struct {
	BilledAzureSubscriptionId   *string `json:"billedAzureSubscriptionId,omitempty"`
	MarketplaceResourceId       *string `json:"marketplaceResourceId,omitempty"`
	MarketplaceStatus           *string `json:"marketplaceStatus,omitempty"`
	MarketplaceSubscriptionId   *string `json:"marketplaceSubscriptionId,omitempty"`
	MarketplaceSubscriptionName *string `json:"marketplaceSubscriptionName,omitempty"`
	OfferId                     *string `json:"offerId,omitempty"`
	PublisherId                 *string `json:"publisherId,omitempty"`
}
