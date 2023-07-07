package offers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OfferProperties struct {
	Content           *string        `json:"content,omitempty"`
	ContentVersion    *string        `json:"contentVersion,omitempty"`
	ProvisioningState *string        `json:"provisioningState,omitempty"`
	PublisherId       *string        `json:"publisherId,omitempty"`
	SkuMappings       *[]SkuMappings `json:"skuMappings,omitempty"`
}
