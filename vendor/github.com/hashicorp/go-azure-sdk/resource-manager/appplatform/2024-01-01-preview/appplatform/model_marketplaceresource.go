package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MarketplaceResource struct {
	Plan      *string `json:"plan,omitempty"`
	Product   *string `json:"product,omitempty"`
	Publisher *string `json:"publisher,omitempty"`
}
