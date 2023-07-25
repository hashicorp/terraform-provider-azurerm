package labs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LabProperties struct {
	Announcement                         *LabAnnouncementProperties `json:"announcement,omitempty"`
	ArtifactsStorageAccount              *string                    `json:"artifactsStorageAccount,omitempty"`
	CreatedDate                          *string                    `json:"createdDate,omitempty"`
	DefaultPremiumStorageAccount         *string                    `json:"defaultPremiumStorageAccount,omitempty"`
	DefaultStorageAccount                *string                    `json:"defaultStorageAccount,omitempty"`
	EnvironmentPermission                *EnvironmentPermission     `json:"environmentPermission,omitempty"`
	ExtendedProperties                   *map[string]string         `json:"extendedProperties,omitempty"`
	LabStorageType                       *StorageType               `json:"labStorageType,omitempty"`
	LoadBalancerId                       *string                    `json:"loadBalancerId,omitempty"`
	MandatoryArtifactsResourceIdsLinux   *[]string                  `json:"mandatoryArtifactsResourceIdsLinux,omitempty"`
	MandatoryArtifactsResourceIdsWindows *[]string                  `json:"mandatoryArtifactsResourceIdsWindows,omitempty"`
	NetworkSecurityGroupId               *string                    `json:"networkSecurityGroupId,omitempty"`
	PremiumDataDiskStorageAccount        *string                    `json:"premiumDataDiskStorageAccount,omitempty"`
	PremiumDataDisks                     *PremiumDataDisk           `json:"premiumDataDisks,omitempty"`
	ProvisioningState                    *string                    `json:"provisioningState,omitempty"`
	PublicIPId                           *string                    `json:"publicIpId,omitempty"`
	Support                              *LabSupportProperties      `json:"support,omitempty"`
	UniqueIdentifier                     *string                    `json:"uniqueIdentifier,omitempty"`
	VMCreationResourceGroup              *string                    `json:"vmCreationResourceGroup,omitempty"`
	VaultName                            *string                    `json:"vaultName,omitempty"`
}

func (o *LabProperties) GetCreatedDateAsTime() (*time.Time, error) {
	if o.CreatedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *LabProperties) SetCreatedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDate = &formatted
}
