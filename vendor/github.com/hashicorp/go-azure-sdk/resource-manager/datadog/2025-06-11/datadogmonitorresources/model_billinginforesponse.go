package datadogmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BillingInfoResponse struct {
	MarketplaceSaasInfo  *MarketplaceSaaSInfo  `json:"marketplaceSaasInfo,omitempty"`
	PartnerBillingEntity *PartnerBillingEntity `json:"partnerBillingEntity,omitempty"`
}
