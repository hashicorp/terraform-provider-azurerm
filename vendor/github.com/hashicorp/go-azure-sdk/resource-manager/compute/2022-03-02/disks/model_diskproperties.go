package disks

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskProperties struct {
	BurstingEnabled              *bool                         `json:"burstingEnabled,omitempty"`
	CompletionPercent            *float64                      `json:"completionPercent,omitempty"`
	CreationData                 CreationData                  `json:"creationData"`
	DataAccessAuthMode           *DataAccessAuthMode           `json:"dataAccessAuthMode,omitempty"`
	DiskAccessId                 *string                       `json:"diskAccessId,omitempty"`
	DiskIOPSReadOnly             *int64                        `json:"diskIOPSReadOnly,omitempty"`
	DiskIOPSReadWrite            *int64                        `json:"diskIOPSReadWrite,omitempty"`
	DiskMBpsReadOnly             *int64                        `json:"diskMBpsReadOnly,omitempty"`
	DiskMBpsReadWrite            *int64                        `json:"diskMBpsReadWrite,omitempty"`
	DiskSizeBytes                *int64                        `json:"diskSizeBytes,omitempty"`
	DiskSizeGB                   *int64                        `json:"diskSizeGB,omitempty"`
	DiskState                    *DiskState                    `json:"diskState,omitempty"`
	Encryption                   *Encryption                   `json:"encryption,omitempty"`
	EncryptionSettingsCollection *EncryptionSettingsCollection `json:"encryptionSettingsCollection,omitempty"`
	HyperVGeneration             *HyperVGeneration             `json:"hyperVGeneration,omitempty"`
	MaxShares                    *int64                        `json:"maxShares,omitempty"`
	NetworkAccessPolicy          *NetworkAccessPolicy          `json:"networkAccessPolicy,omitempty"`
	OsType                       *OperatingSystemTypes         `json:"osType,omitempty"`
	PropertyUpdatesInProgress    *PropertyUpdatesInProgress    `json:"propertyUpdatesInProgress,omitempty"`
	ProvisioningState            *string                       `json:"provisioningState,omitempty"`
	PublicNetworkAccess          *PublicNetworkAccess          `json:"publicNetworkAccess,omitempty"`
	PurchasePlan                 *PurchasePlan                 `json:"purchasePlan,omitempty"`
	SecurityProfile              *DiskSecurityProfile          `json:"securityProfile,omitempty"`
	ShareInfo                    *[]ShareInfoElement           `json:"shareInfo,omitempty"`
	SupportedCapabilities        *SupportedCapabilities        `json:"supportedCapabilities,omitempty"`
	SupportsHibernation          *bool                         `json:"supportsHibernation,omitempty"`
	Tier                         *string                       `json:"tier,omitempty"`
	TimeCreated                  *string                       `json:"timeCreated,omitempty"`
	UniqueId                     *string                       `json:"uniqueId,omitempty"`
}

func (o *DiskProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	if o.TimeCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *DiskProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = &formatted
}
