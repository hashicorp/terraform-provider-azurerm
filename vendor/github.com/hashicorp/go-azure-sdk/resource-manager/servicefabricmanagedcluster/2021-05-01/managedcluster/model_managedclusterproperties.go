package managedcluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterProperties struct {
	AddonFeatures                        *[]AddonFeatures                      `json:"addonFeatures,omitempty"`
	AdminPassword                        *string                               `json:"adminPassword,omitempty"`
	AdminUserName                        string                                `json:"adminUserName"`
	AllowRdpAccess                       *bool                                 `json:"allowRdpAccess,omitempty"`
	ApplicationTypeVersionsCleanupPolicy *ApplicationTypeVersionsCleanupPolicy `json:"applicationTypeVersionsCleanupPolicy,omitempty"`
	AzureActiveDirectory                 *AzureActiveDirectory                 `json:"azureActiveDirectory,omitempty"`
	ClientConnectionPort                 *int64                                `json:"clientConnectionPort,omitempty"`
	Clients                              *[]ClientCertificate                  `json:"clients,omitempty"`
	ClusterCertificateThumbprints        *[]string                             `json:"clusterCertificateThumbprints,omitempty"`
	ClusterCodeVersion                   *string                               `json:"clusterCodeVersion,omitempty"`
	ClusterId                            *string                               `json:"clusterId,omitempty"`
	ClusterState                         *ClusterState                         `json:"clusterState,omitempty"`
	ClusterUpgradeCadence                *ClusterUpgradeCadence                `json:"clusterUpgradeCadence,omitempty"`
	ClusterUpgradeMode                   *ClusterUpgradeMode                   `json:"clusterUpgradeMode,omitempty"`
	DnsName                              string                                `json:"dnsName"`
	EnableAutoOSUpgrade                  *bool                                 `json:"enableAutoOSUpgrade,omitempty"`
	FabricSettings                       *[]SettingsSectionDescription         `json:"fabricSettings,omitempty"`
	Fqdn                                 *string                               `json:"fqdn,omitempty"`
	HTTPGatewayConnectionPort            *int64                                `json:"httpGatewayConnectionPort,omitempty"`
	IPv4Address                          *string                               `json:"ipv4Address,omitempty"`
	LoadBalancingRules                   *[]LoadBalancingRule                  `json:"loadBalancingRules,omitempty"`
	NetworkSecurityRules                 *[]NetworkSecurityRule                `json:"networkSecurityRules,omitempty"`
	ProvisioningState                    *ManagedResourceProvisioningState     `json:"provisioningState,omitempty"`
	ZonalResiliency                      *bool                                 `json:"zonalResiliency,omitempty"`
}
