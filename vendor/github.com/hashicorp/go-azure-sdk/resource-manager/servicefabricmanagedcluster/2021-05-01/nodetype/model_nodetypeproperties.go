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
	VmExtensions            *[]VMSSExtension                  `json:"vmExtensions,omitempty"`
	VmImageOffer            *string                           `json:"vmImageOffer,omitempty"`
	VmImagePublisher        *string                           `json:"vmImagePublisher,omitempty"`
	VmImageSku              *string                           `json:"vmImageSku,omitempty"`
	VmImageVersion          *string                           `json:"vmImageVersion,omitempty"`
	VmInstanceCount         int64                             `json:"vmInstanceCount"`
	VmManagedIdentity       *identity.UserAssignedList        `json:"vmManagedIdentity,omitempty"`
	VmSecrets               *[]VaultSecretGroup               `json:"vmSecrets,omitempty"`
	VmSize                  *string                           `json:"vmSize,omitempty"`
}
