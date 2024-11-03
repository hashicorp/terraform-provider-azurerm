package organizations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiftrBaseDataOrganizationProperties struct {
	Marketplace                   LiftrBaseMarketplaceDetails                 `json:"marketplace"`
	PartnerOrganizationProperties *LiftrBaseDataPartnerOrganizationProperties `json:"partnerOrganizationProperties,omitempty"`
	ProvisioningState             *ResourceProvisioningState                  `json:"provisioningState,omitempty"`
	User                          LiftrBaseUserDetails                        `json:"user"`
}
