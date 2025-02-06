package images

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageProperties struct {
	Description                     *string                          `json:"description,omitempty"`
	HibernateSupport                *HibernateSupport                `json:"hibernateSupport,omitempty"`
	Offer                           *string                          `json:"offer,omitempty"`
	ProvisioningState               *ProvisioningState               `json:"provisioningState,omitempty"`
	Publisher                       *string                          `json:"publisher,omitempty"`
	RecommendedMachineConfiguration *RecommendedMachineConfiguration `json:"recommendedMachineConfiguration,omitempty"`
	Sku                             *string                          `json:"sku,omitempty"`
}
