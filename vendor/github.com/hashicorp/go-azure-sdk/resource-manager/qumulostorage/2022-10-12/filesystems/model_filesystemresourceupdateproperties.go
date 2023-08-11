package filesystems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileSystemResourceUpdateProperties struct {
	ClusterLoginUrl    *string             `json:"clusterLoginUrl,omitempty"`
	DelegatedSubnetId  *string             `json:"delegatedSubnetId,omitempty"`
	MarketplaceDetails *MarketplaceDetails `json:"marketplaceDetails,omitempty"`
	PrivateIPs         *[]string           `json:"privateIPs,omitempty"`
	UserDetails        *UserDetails        `json:"userDetails,omitempty"`
}
