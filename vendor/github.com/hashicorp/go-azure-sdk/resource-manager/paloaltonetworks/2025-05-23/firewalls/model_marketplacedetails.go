package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MarketplaceDetails struct {
	MarketplaceSubscriptionId     *string                        `json:"marketplaceSubscriptionId,omitempty"`
	MarketplaceSubscriptionStatus *MarketplaceSubscriptionStatus `json:"marketplaceSubscriptionStatus,omitempty"`
	OfferId                       string                         `json:"offerId"`
	PublisherId                   string                         `json:"publisherId"`
}
