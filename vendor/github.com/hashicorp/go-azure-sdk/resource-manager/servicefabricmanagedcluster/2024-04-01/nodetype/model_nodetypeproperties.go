package nodetype

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodeTypeProperties struct {
	AdditionalDataDisks                      *[]VMSSDataDisk                            `json:"additionalDataDisks,omitempty"`
	AdditionalNetworkInterfaceConfigurations *[]AdditionalNetworkInterfaceConfiguration `json:"additionalNetworkInterfaceConfigurations,omitempty"`
	ApplicationPorts                         *EndpointRangeDescription                  `json:"applicationPorts,omitempty"`
	Capacities                               *map[string]string                         `json:"capacities,omitempty"`
	ComputerNamePrefix                       *string                                    `json:"computerNamePrefix,omitempty"`
	DataDiskLetter                           *string                                    `json:"dataDiskLetter,omitempty"`
	DataDiskSizeGB                           *int64                                     `json:"dataDiskSizeGB,omitempty"`
	DataDiskType                             *DiskType                                  `json:"dataDiskType,omitempty"`
	DscpConfigurationId                      *string                                    `json:"dscpConfigurationId,omitempty"`
	EnableAcceleratedNetworking              *bool                                      `json:"enableAcceleratedNetworking,omitempty"`
	EnableEncryptionAtHost                   *bool                                      `json:"enableEncryptionAtHost,omitempty"`
	EnableNodePublicIP                       *bool                                      `json:"enableNodePublicIP,omitempty"`
	EnableNodePublicIPv6                     *bool                                      `json:"enableNodePublicIPv6,omitempty"`
	EnableOverProvisioning                   *bool                                      `json:"enableOverProvisioning,omitempty"`
	EphemeralPorts                           *EndpointRangeDescription                  `json:"ephemeralPorts,omitempty"`
	EvictionPolicy                           *EvictionPolicyType                        `json:"evictionPolicy,omitempty"`
	FrontendConfigurations                   *[]FrontendConfiguration                   `json:"frontendConfigurations,omitempty"`
	HostGroupId                              *string                                    `json:"hostGroupId,omitempty"`
	IsPrimary                                bool                                       `json:"isPrimary"`
	IsSpotVM                                 *bool                                      `json:"isSpotVM,omitempty"`
	IsStateless                              *bool                                      `json:"isStateless,omitempty"`
	MultiplePlacementGroups                  *bool                                      `json:"multiplePlacementGroups,omitempty"`
	NatConfigurations                        *[]NodeTypeNatConfig                       `json:"natConfigurations,omitempty"`
	NatGatewayId                             *string                                    `json:"natGatewayId,omitempty"`
	NetworkSecurityRules                     *[]NetworkSecurityRule                     `json:"networkSecurityRules,omitempty"`
	PlacementProperties                      *map[string]string                         `json:"placementProperties,omitempty"`
	ProvisioningState                        *ManagedResourceProvisioningState          `json:"provisioningState,omitempty"`
	SecureBootEnabled                        *bool                                      `json:"secureBootEnabled,omitempty"`
	SecurityType                             *SecurityType                              `json:"securityType,omitempty"`
	ServiceArtifactReferenceId               *string                                    `json:"serviceArtifactReferenceId,omitempty"`
	SpotRestoreTimeout                       *string                                    `json:"spotRestoreTimeout,omitempty"`
	SubnetId                                 *string                                    `json:"subnetId,omitempty"`
	UseDefaultPublicLoadBalancer             *bool                                      `json:"useDefaultPublicLoadBalancer,omitempty"`
	UseEphemeralOSDisk                       *bool                                      `json:"useEphemeralOSDisk,omitempty"`
	UseTempDataDisk                          *bool                                      `json:"useTempDataDisk,omitempty"`
	VMExtensions                             *[]VMSSExtension                           `json:"vmExtensions,omitempty"`
	VMImageOffer                             *string                                    `json:"vmImageOffer,omitempty"`
	VMImagePlan                              *VMImagePlan                               `json:"vmImagePlan,omitempty"`
	VMImagePublisher                         *string                                    `json:"vmImagePublisher,omitempty"`
	VMImageResourceId                        *string                                    `json:"vmImageResourceId,omitempty"`
	VMImageSku                               *string                                    `json:"vmImageSku,omitempty"`
	VMImageVersion                           *string                                    `json:"vmImageVersion,omitempty"`
	VMInstanceCount                          int64                                      `json:"vmInstanceCount"`
	VMManagedIdentity                        *identity.UserAssignedList                 `json:"vmManagedIdentity,omitempty"`
	VMSecrets                                *[]VaultSecretGroup                        `json:"vmSecrets,omitempty"`
	VMSetupActions                           *[]VMSetupAction                           `json:"vmSetupActions,omitempty"`
	VMSharedGalleryImageId                   *string                                    `json:"vmSharedGalleryImageId,omitempty"`
	VMSize                                   *string                                    `json:"vmSize,omitempty"`
	Zones                                    *zones.Schema                              `json:"zones,omitempty"`
}
