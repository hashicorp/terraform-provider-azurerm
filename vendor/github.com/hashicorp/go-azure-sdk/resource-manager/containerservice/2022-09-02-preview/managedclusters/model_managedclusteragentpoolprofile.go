package managedclusters

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterAgentPoolProfile struct {
	AvailabilityZones          *zones.Schema             `json:"availabilityZones,omitempty"`
	CapacityReservationGroupID *string                   `json:"capacityReservationGroupID,omitempty"`
	Count                      *int64                    `json:"count,omitempty"`
	CreationData               *CreationData             `json:"creationData,omitempty"`
	CurrentOrchestratorVersion *string                   `json:"currentOrchestratorVersion,omitempty"`
	EnableAutoScaling          *bool                     `json:"enableAutoScaling,omitempty"`
	EnableCustomCATrust        *bool                     `json:"enableCustomCATrust,omitempty"`
	EnableEncryptionAtHost     *bool                     `json:"enableEncryptionAtHost,omitempty"`
	EnableFIPS                 *bool                     `json:"enableFIPS,omitempty"`
	EnableNodePublicIP         *bool                     `json:"enableNodePublicIP,omitempty"`
	EnableUltraSSD             *bool                     `json:"enableUltraSSD,omitempty"`
	GpuInstanceProfile         *GPUInstanceProfile       `json:"gpuInstanceProfile,omitempty"`
	HostGroupID                *string                   `json:"hostGroupID,omitempty"`
	KubeletConfig              *KubeletConfig            `json:"kubeletConfig,omitempty"`
	KubeletDiskType            *KubeletDiskType          `json:"kubeletDiskType,omitempty"`
	LinuxOSConfig              *LinuxOSConfig            `json:"linuxOSConfig,omitempty"`
	MaxCount                   *int64                    `json:"maxCount,omitempty"`
	MaxPods                    *int64                    `json:"maxPods,omitempty"`
	MessageOfTheDay            *string                   `json:"messageOfTheDay,omitempty"`
	MinCount                   *int64                    `json:"minCount,omitempty"`
	Mode                       *AgentPoolMode            `json:"mode,omitempty"`
	Name                       string                    `json:"name"`
	NetworkProfile             *AgentPoolNetworkProfile  `json:"networkProfile,omitempty"`
	NodeImageVersion           *string                   `json:"nodeImageVersion,omitempty"`
	NodeLabels                 *map[string]string        `json:"nodeLabels,omitempty"`
	NodePublicIPPrefixID       *string                   `json:"nodePublicIPPrefixID,omitempty"`
	NodeTaints                 *[]string                 `json:"nodeTaints,omitempty"`
	OrchestratorVersion        *string                   `json:"orchestratorVersion,omitempty"`
	OsDiskSizeGB               *int64                    `json:"osDiskSizeGB,omitempty"`
	OsDiskType                 *OSDiskType               `json:"osDiskType,omitempty"`
	OsSKU                      *OSSKU                    `json:"osSKU,omitempty"`
	OsType                     *OSType                   `json:"osType,omitempty"`
	PodSubnetID                *string                   `json:"podSubnetID,omitempty"`
	PowerState                 *PowerState               `json:"powerState,omitempty"`
	ProvisioningState          *string                   `json:"provisioningState,omitempty"`
	ProximityPlacementGroupID  *string                   `json:"proximityPlacementGroupID,omitempty"`
	ScaleDownMode              *ScaleDownMode            `json:"scaleDownMode,omitempty"`
	ScaleSetEvictionPolicy     *ScaleSetEvictionPolicy   `json:"scaleSetEvictionPolicy,omitempty"`
	ScaleSetPriority           *ScaleSetPriority         `json:"scaleSetPriority,omitempty"`
	SpotMaxPrice               *float64                  `json:"spotMaxPrice,omitempty"`
	Tags                       *map[string]string        `json:"tags,omitempty"`
	Type                       *AgentPoolType            `json:"type,omitempty"`
	UpgradeSettings            *AgentPoolUpgradeSettings `json:"upgradeSettings,omitempty"`
	VMSize                     *string                   `json:"vmSize,omitempty"`
	VnetSubnetID               *string                   `json:"vnetSubnetID,omitempty"`
	WindowsProfile             *AgentPoolWindowsProfile  `json:"windowsProfile,omitempty"`
	WorkloadRuntime            *WorkloadRuntime          `json:"workloadRuntime,omitempty"`
}
