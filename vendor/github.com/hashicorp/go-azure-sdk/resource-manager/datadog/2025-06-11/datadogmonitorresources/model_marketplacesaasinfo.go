package datadogmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MarketplaceSaaSInfo struct {
	BilledAzureSubscriptionId *string `json:"billedAzureSubscriptionId,omitempty"`
	MarketplaceName           *string `json:"marketplaceName,omitempty"`
	MarketplaceStatus         *string `json:"marketplaceStatus,omitempty"`
	MarketplaceSubscriptionId *string `json:"marketplaceSubscriptionId,omitempty"`
	Subscribed                *bool   `json:"subscribed,omitempty"`
}
