package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProperties struct {
	AcceptedAudiences                  *[]AcceptedAudiences         `json:"acceptedAudiences,omitempty"`
	AllowedFqdnList                    *[]string                    `json:"allowedFqdnList,omitempty"`
	AllowedIPRangeList                 *[]string                    `json:"allowedIpRangeList,omitempty"`
	DataIngestionUri                   *string                      `json:"dataIngestionUri,omitempty"`
	EnableAutoStop                     *bool                        `json:"enableAutoStop,omitempty"`
	EnableDiskEncryption               *bool                        `json:"enableDiskEncryption,omitempty"`
	EnableDoubleEncryption             *bool                        `json:"enableDoubleEncryption,omitempty"`
	EnablePurge                        *bool                        `json:"enablePurge,omitempty"`
	EnableStreamingIngest              *bool                        `json:"enableStreamingIngest,omitempty"`
	EngineType                         *EngineType                  `json:"engineType,omitempty"`
	KeyVaultProperties                 *KeyVaultProperties          `json:"keyVaultProperties,omitempty"`
	LanguageExtensions                 *LanguageExtensionsList      `json:"languageExtensions,omitempty"`
	MigrationCluster                   *MigrationClusterProperties  `json:"migrationCluster,omitempty"`
	OptimizedAutoscale                 *OptimizedAutoscale          `json:"optimizedAutoscale,omitempty"`
	PrivateEndpointConnections         *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState                  *ProvisioningState           `json:"provisioningState,omitempty"`
	PublicIPType                       *PublicIPType                `json:"publicIPType,omitempty"`
	PublicNetworkAccess                *PublicNetworkAccess         `json:"publicNetworkAccess,omitempty"`
	RestrictOutboundNetworkAccess      *ClusterNetworkAccessFlag    `json:"restrictOutboundNetworkAccess,omitempty"`
	State                              *State                       `json:"state,omitempty"`
	StateReason                        *string                      `json:"stateReason,omitempty"`
	TrustedExternalTenants             *[]TrustedExternalTenant     `json:"trustedExternalTenants,omitempty"`
	Uri                                *string                      `json:"uri,omitempty"`
	VirtualClusterGraduationProperties *string                      `json:"virtualClusterGraduationProperties,omitempty"`
	VirtualNetworkConfiguration        *VirtualNetworkConfiguration `json:"virtualNetworkConfiguration,omitempty"`
}
