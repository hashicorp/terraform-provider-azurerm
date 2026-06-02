package resources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Plan struct {
	Name          *string `json:"name,omitempty"`
	Product       *string `json:"product,omitempty"`
	PromotionCode *string `json:"promotionCode,omitempty"`
	Publisher     *string `json:"publisher,omitempty"`
	Version       *string `json:"version,omitempty"`
}
