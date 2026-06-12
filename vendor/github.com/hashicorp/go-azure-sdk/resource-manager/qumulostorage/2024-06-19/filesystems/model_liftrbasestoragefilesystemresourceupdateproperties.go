package filesystems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiftrBaseStorageFileSystemResourceUpdateProperties struct {
	DelegatedSubnetId  *string                      `json:"delegatedSubnetId,omitempty"`
	MarketplaceDetails *LiftrBaseMarketplaceDetails `json:"marketplaceDetails,omitempty"`
	UserDetails        *LiftrBaseUserDetails        `json:"userDetails,omitempty"`
}
