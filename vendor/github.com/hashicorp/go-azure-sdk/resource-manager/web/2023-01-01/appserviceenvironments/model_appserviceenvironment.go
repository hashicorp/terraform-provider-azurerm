package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServiceEnvironment struct {
	ClusterSettings              *[]NameValuePair              `json:"clusterSettings,omitempty"`
	CustomDnsSuffixConfiguration *CustomDnsSuffixConfiguration `json:"customDnsSuffixConfiguration,omitempty"`
	DedicatedHostCount           *int64                        `json:"dedicatedHostCount,omitempty"`
	DnsSuffix                    *string                       `json:"dnsSuffix,omitempty"`
	FrontEndScaleFactor          *int64                        `json:"frontEndScaleFactor,omitempty"`
	HasLinuxWorkers              *bool                         `json:"hasLinuxWorkers,omitempty"`
	IPsslAddressCount            *int64                        `json:"ipsslAddressCount,omitempty"`
	InternalLoadBalancingMode    *LoadBalancingMode            `json:"internalLoadBalancingMode,omitempty"`
	MaximumNumberOfMachines      *int64                        `json:"maximumNumberOfMachines,omitempty"`
	MultiRoleCount               *int64                        `json:"multiRoleCount,omitempty"`
	MultiSize                    *string                       `json:"multiSize,omitempty"`
	NetworkingConfiguration      *AseV3NetworkingConfiguration `json:"networkingConfiguration,omitempty"`
	ProvisioningState            *ProvisioningState            `json:"provisioningState,omitempty"`
	Status                       *HostingEnvironmentStatus     `json:"status,omitempty"`
	Suspended                    *bool                         `json:"suspended,omitempty"`
	UpgradeAvailability          *UpgradeAvailability          `json:"upgradeAvailability,omitempty"`
	UpgradePreference            *UpgradePreference            `json:"upgradePreference,omitempty"`
	UserWhitelistedIPRanges      *[]string                     `json:"userWhitelistedIpRanges,omitempty"`
	VirtualNetwork               VirtualNetworkProfile         `json:"virtualNetwork"`
	ZoneRedundant                *bool                         `json:"zoneRedundant,omitempty"`
}
