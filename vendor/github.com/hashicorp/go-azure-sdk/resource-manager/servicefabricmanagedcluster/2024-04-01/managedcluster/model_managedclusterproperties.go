package managedcluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterProperties struct {
	AddonFeatures                        *[]ManagedClusterAddOnFeature         `json:"addonFeatures,omitempty"`
	AdminPassword                        *string                               `json:"adminPassword,omitempty"`
	AdminUserName                        string                                `json:"adminUserName"`
	AllowRdpAccess                       *bool                                 `json:"allowRdpAccess,omitempty"`
	ApplicationTypeVersionsCleanupPolicy *ApplicationTypeVersionsCleanupPolicy `json:"applicationTypeVersionsCleanupPolicy,omitempty"`
	AuxiliarySubnets                     *[]Subnet                             `json:"auxiliarySubnets,omitempty"`
	AzureActiveDirectory                 *AzureActiveDirectory                 `json:"azureActiveDirectory,omitempty"`
	ClientConnectionPort                 *int64                                `json:"clientConnectionPort,omitempty"`
	Clients                              *[]ClientCertificate                  `json:"clients,omitempty"`
	ClusterCertificateThumbprints        *[]string                             `json:"clusterCertificateThumbprints,omitempty"`
	ClusterCodeVersion                   *string                               `json:"clusterCodeVersion,omitempty"`
	ClusterId                            *string                               `json:"clusterId,omitempty"`
	ClusterState                         *ClusterState                         `json:"clusterState,omitempty"`
	ClusterUpgradeCadence                *ClusterUpgradeCadence                `json:"clusterUpgradeCadence,omitempty"`
	ClusterUpgradeMode                   *ClusterUpgradeMode                   `json:"clusterUpgradeMode,omitempty"`
	DdosProtectionPlanId                 *string                               `json:"ddosProtectionPlanId,omitempty"`
	DnsName                              string                                `json:"dnsName"`
	EnableAutoOSUpgrade                  *bool                                 `json:"enableAutoOSUpgrade,omitempty"`
	EnableHTTPGatewayExclusiveAuthMode   *bool                                 `json:"enableHttpGatewayExclusiveAuthMode,omitempty"`
	EnableIPv6                           *bool                                 `json:"enableIpv6,omitempty"`
	EnableServicePublicIP                *bool                                 `json:"enableServicePublicIP,omitempty"`
	FabricSettings                       *[]SettingsSectionDescription         `json:"fabricSettings,omitempty"`
	Fqdn                                 *string                               `json:"fqdn,omitempty"`
	HTTPGatewayConnectionPort            *int64                                `json:"httpGatewayConnectionPort,omitempty"`
	HTTPGatewayTokenAuthConnectionPort   *int64                                `json:"httpGatewayTokenAuthConnectionPort,omitempty"`
	IPTags                               *[]IPTag                              `json:"ipTags,omitempty"`
	IPv4Address                          *string                               `json:"ipv4Address,omitempty"`
	IPv6Address                          *string                               `json:"ipv6Address,omitempty"`
	LoadBalancingRules                   *[]LoadBalancingRule                  `json:"loadBalancingRules,omitempty"`
	NetworkSecurityRules                 *[]NetworkSecurityRule                `json:"networkSecurityRules,omitempty"`
	ProvisioningState                    *ManagedResourceProvisioningState     `json:"provisioningState,omitempty"`
	PublicIPPrefixId                     *string                               `json:"publicIPPrefixId,omitempty"`
	PublicIPv6PrefixId                   *string                               `json:"publicIPv6PrefixId,omitempty"`
	ServiceEndpoints                     *[]ServiceEndpoint                    `json:"serviceEndpoints,omitempty"`
	SubnetId                             *string                               `json:"subnetId,omitempty"`
	UpgradeDescription                   *ClusterUpgradePolicy                 `json:"upgradeDescription,omitempty"`
	UseCustomVnet                        *bool                                 `json:"useCustomVnet,omitempty"`
	ZonalResiliency                      *bool                                 `json:"zonalResiliency,omitempty"`
	ZonalUpdateMode                      *ZonalUpdateMode                      `json:"zonalUpdateMode,omitempty"`
}
