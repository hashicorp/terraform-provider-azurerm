package nodetype

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodeTypeProperties struct {
	ApplicationPorts        *EndpointRangeDescription         `json:"applicationPorts,omitempty"`
	Capacities              *map[string]string                `json:"capacities,omitempty"`
	DataDiskSizeGB          int64                             `json:"dataDiskSizeGB"`
	DataDiskType            *DiskType                         `json:"dataDiskType,omitempty"`
	EphemeralPorts          *EndpointRangeDescription         `json:"ephemeralPorts,omitempty"`
	IsPrimary               bool                              `json:"isPrimary"`
	IsStateless             *bool                             `json:"isStateless,omitempty"`
	MultiplePlacementGroups *bool                             `json:"multiplePlacementGroups,omitempty"`
	PlacementProperties     *map[string]string                `json:"placementProperties,omitempty"`
	ProvisioningState       *ManagedResourceProvisioningState `json:"provisioningState,omitempty"`
	VMExtensions            *[]VMSSExtension                  `json:"vmExtensions,omitempty"`
	VMImageOffer            *string                           `json:"vmImageOffer,omitempty"`
	VMImagePublisher        *string                           `json:"vmImagePublisher,omitempty"`
	VMImageSku              *string                           `json:"vmImageSku,omitempty"`
	VMImageVersion          *string                           `json:"vmImageVersion,omitempty"`
	VMInstanceCount         int64                             `json:"vmInstanceCount"`
	VMManagedIdentity       *identity.UserAssignedList        `json:"vmManagedIdentity,omitempty"`
	VMSecrets               *[]VaultSecretGroup               `json:"vmSecrets,omitempty"`
	VMSize                  *string                           `json:"vmSize,omitempty"`
}
