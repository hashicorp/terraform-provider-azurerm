package filesystems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiftrBaseStorageFileSystemResourceProperties struct {
	AdminPassword      string                      `json:"adminPassword"`
	AvailabilityZone   *string                     `json:"availabilityZone,omitempty"`
	ClusterLoginURL    *string                     `json:"clusterLoginUrl,omitempty"`
	DelegatedSubnetId  string                      `json:"delegatedSubnetId"`
	MarketplaceDetails LiftrBaseMarketplaceDetails `json:"marketplaceDetails"`
	PrivateIPs         *[]string                   `json:"privateIPs,omitempty"`
	ProvisioningState  *ProvisioningState          `json:"provisioningState,omitempty"`
	StorageSku         string                      `json:"storageSku"`
	UserDetails        LiftrBaseUserDetails        `json:"userDetails"`
}
