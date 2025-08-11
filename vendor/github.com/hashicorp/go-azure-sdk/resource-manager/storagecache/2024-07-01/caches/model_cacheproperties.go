package caches

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheProperties struct {
	CacheSizeGB               *int64                          `json:"cacheSizeGB,omitempty"`
	DirectoryServicesSettings *CacheDirectorySettings         `json:"directoryServicesSettings,omitempty"`
	EncryptionSettings        *CacheEncryptionSettings        `json:"encryptionSettings,omitempty"`
	Health                    *CacheHealth                    `json:"health,omitempty"`
	MountAddresses            *[]string                       `json:"mountAddresses,omitempty"`
	NetworkSettings           *CacheNetworkSettings           `json:"networkSettings,omitempty"`
	PrimingJobs               *[]PrimingJob                   `json:"primingJobs,omitempty"`
	ProvisioningState         *ProvisioningStateType          `json:"provisioningState,omitempty"`
	SecuritySettings          *CacheSecuritySettings          `json:"securitySettings,omitempty"`
	SpaceAllocation           *[]StorageTargetSpaceAllocation `json:"spaceAllocation,omitempty"`
	Subnet                    *string                         `json:"subnet,omitempty"`
	UpgradeSettings           *CacheUpgradeSettings           `json:"upgradeSettings,omitempty"`
	UpgradeStatus             *CacheUpgradeStatus             `json:"upgradeStatus,omitempty"`
	Zones                     *zones.Schema                   `json:"zones,omitempty"`
}
