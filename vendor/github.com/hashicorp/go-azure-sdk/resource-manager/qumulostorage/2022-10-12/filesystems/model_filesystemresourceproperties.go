package filesystems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileSystemResourceProperties struct {
	AdminPassword      string             `json:"adminPassword"`
	AvailabilityZone   *string            `json:"availabilityZone,omitempty"`
	ClusterLoginUrl    *string            `json:"clusterLoginUrl,omitempty"`
	DelegatedSubnetId  string             `json:"delegatedSubnetId"`
	InitialCapacity    int64              `json:"initialCapacity"`
	MarketplaceDetails MarketplaceDetails `json:"marketplaceDetails"`
	PrivateIPs         *[]string          `json:"privateIPs,omitempty"`
	ProvisioningState  *ProvisioningState `json:"provisioningState,omitempty"`
	StorageSku         StorageSku         `json:"storageSku"`
	UserDetails        UserDetails        `json:"userDetails"`
}
