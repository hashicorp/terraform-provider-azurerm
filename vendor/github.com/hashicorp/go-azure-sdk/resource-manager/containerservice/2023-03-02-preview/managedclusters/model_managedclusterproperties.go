package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterProperties struct {
	AadProfile                *ManagedClusterAADProfile                  `json:"aadProfile,omitempty"`
	AddonProfiles             *map[string]ManagedClusterAddonProfile     `json:"addonProfiles,omitempty"`
	AgentPoolProfiles         *[]ManagedClusterAgentPoolProfile          `json:"agentPoolProfiles,omitempty"`
	ApiServerAccessProfile    *ManagedClusterAPIServerAccessProfile      `json:"apiServerAccessProfile,omitempty"`
	AutoScalerProfile         *ManagedClusterPropertiesAutoScalerProfile `json:"autoScalerProfile,omitempty"`
	AutoUpgradeProfile        *ManagedClusterAutoUpgradeProfile          `json:"autoUpgradeProfile,omitempty"`
	AzureMonitorProfile       *ManagedClusterAzureMonitorProfile         `json:"azureMonitorProfile,omitempty"`
	AzurePortalFQDN           *string                                    `json:"azurePortalFQDN,omitempty"`
	CreationData              *CreationData                              `json:"creationData,omitempty"`
	CurrentKubernetesVersion  *string                                    `json:"currentKubernetesVersion,omitempty"`
	DisableLocalAccounts      *bool                                      `json:"disableLocalAccounts,omitempty"`
	DiskEncryptionSetID       *string                                    `json:"diskEncryptionSetID,omitempty"`
	DnsPrefix                 *string                                    `json:"dnsPrefix,omitempty"`
	EnableNamespaceResources  *bool                                      `json:"enableNamespaceResources,omitempty"`
	EnablePodSecurityPolicy   *bool                                      `json:"enablePodSecurityPolicy,omitempty"`
	EnableRBAC                *bool                                      `json:"enableRBAC,omitempty"`
	Fqdn                      *string                                    `json:"fqdn,omitempty"`
	FqdnSubdomain             *string                                    `json:"fqdnSubdomain,omitempty"`
	GuardrailsProfile         *GuardrailsProfile                         `json:"guardrailsProfile,omitempty"`
	HTTPProxyConfig           *ManagedClusterHTTPProxyConfig             `json:"httpProxyConfig,omitempty"`
	IdentityProfile           *map[string]UserAssignedIdentity           `json:"identityProfile,omitempty"`
	IngressProfile            *ManagedClusterIngressProfile              `json:"ingressProfile,omitempty"`
	KubernetesVersion         *string                                    `json:"kubernetesVersion,omitempty"`
	LinuxProfile              *ContainerServiceLinuxProfile              `json:"linuxProfile,omitempty"`
	MaxAgentPools             *int64                                     `json:"maxAgentPools,omitempty"`
	NetworkProfile            *ContainerServiceNetworkProfile            `json:"networkProfile,omitempty"`
	NodeResourceGroup         *string                                    `json:"nodeResourceGroup,omitempty"`
	NodeResourceGroupProfile  *ManagedClusterNodeResourceGroupProfile    `json:"nodeResourceGroupProfile,omitempty"`
	OidcIssuerProfile         *ManagedClusterOIDCIssuerProfile           `json:"oidcIssuerProfile,omitempty"`
	PodIdentityProfile        *ManagedClusterPodIdentityProfile          `json:"podIdentityProfile,omitempty"`
	PowerState                *PowerState                                `json:"powerState,omitempty"`
	PrivateFQDN               *string                                    `json:"privateFQDN,omitempty"`
	PrivateLinkResources      *[]PrivateLinkResource                     `json:"privateLinkResources,omitempty"`
	ProvisioningState         *string                                    `json:"provisioningState,omitempty"`
	PublicNetworkAccess       *PublicNetworkAccess                       `json:"publicNetworkAccess,omitempty"`
	SecurityProfile           *ManagedClusterSecurityProfile             `json:"securityProfile,omitempty"`
	ServiceMeshProfile        *ServiceMeshProfile                        `json:"serviceMeshProfile,omitempty"`
	ServicePrincipalProfile   *ManagedClusterServicePrincipalProfile     `json:"servicePrincipalProfile,omitempty"`
	StorageProfile            *ManagedClusterStorageProfile              `json:"storageProfile,omitempty"`
	SupportPlan               *KubernetesSupportPlan                     `json:"supportPlan,omitempty"`
	UpgradeSettings           *ClusterUpgradeSettings                    `json:"upgradeSettings,omitempty"`
	WindowsProfile            *ManagedClusterWindowsProfile              `json:"windowsProfile,omitempty"`
	WorkloadAutoScalerProfile *ManagedClusterWorkloadAutoScalerProfile   `json:"workloadAutoScalerProfile,omitempty"`
}
