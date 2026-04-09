package virtualmachines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LabVirtualMachineProperties struct {
	AllowClaim                   *bool                               `json:"allowClaim,omitempty"`
	ApplicableSchedule           *ApplicableSchedule                 `json:"applicableSchedule,omitempty"`
	ArtifactDeploymentStatus     *ArtifactDeploymentStatusProperties `json:"artifactDeploymentStatus,omitempty"`
	Artifacts                    *[]ArtifactInstallProperties        `json:"artifacts,omitempty"`
	ComputeId                    *string                             `json:"computeId,omitempty"`
	ComputeVM                    *ComputeVMProperties                `json:"computeVm,omitempty"`
	CreatedByUser                *string                             `json:"createdByUser,omitempty"`
	CreatedByUserId              *string                             `json:"createdByUserId,omitempty"`
	CreatedDate                  *string                             `json:"createdDate,omitempty"`
	CustomImageId                *string                             `json:"customImageId,omitempty"`
	DataDiskParameters           *[]DataDiskProperties               `json:"dataDiskParameters,omitempty"`
	DisallowPublicIPAddress      *bool                               `json:"disallowPublicIpAddress,omitempty"`
	EnvironmentId                *string                             `json:"environmentId,omitempty"`
	ExpirationDate               *string                             `json:"expirationDate,omitempty"`
	Fqdn                         *string                             `json:"fqdn,omitempty"`
	GalleryImageReference        *GalleryImageReference              `json:"galleryImageReference,omitempty"`
	IsAuthenticationWithSshKey   *bool                               `json:"isAuthenticationWithSshKey,omitempty"`
	LabSubnetName                *string                             `json:"labSubnetName,omitempty"`
	LabVirtualNetworkId          *string                             `json:"labVirtualNetworkId,omitempty"`
	LastKnownPowerState          *string                             `json:"lastKnownPowerState,omitempty"`
	NetworkInterface             *NetworkInterfaceProperties         `json:"networkInterface,omitempty"`
	Notes                        *string                             `json:"notes,omitempty"`
	OsType                       *string                             `json:"osType,omitempty"`
	OwnerObjectId                *string                             `json:"ownerObjectId,omitempty"`
	OwnerUserPrincipalName       *string                             `json:"ownerUserPrincipalName,omitempty"`
	Password                     *string                             `json:"password,omitempty"`
	PlanId                       *string                             `json:"planId,omitempty"`
	ProvisioningState            *string                             `json:"provisioningState,omitempty"`
	ScheduleParameters           *[]ScheduleCreationParameter        `json:"scheduleParameters,omitempty"`
	Size                         *string                             `json:"size,omitempty"`
	SshKey                       *string                             `json:"sshKey,omitempty"`
	StorageType                  *string                             `json:"storageType,omitempty"`
	UniqueIdentifier             *string                             `json:"uniqueIdentifier,omitempty"`
	UserName                     *string                             `json:"userName,omitempty"`
	VirtualMachineCreationSource *VirtualMachineCreationSource       `json:"virtualMachineCreationSource,omitempty"`
}

func (o *LabVirtualMachineProperties) GetCreatedDateAsTime() (*time.Time, error) {
	if o.CreatedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *LabVirtualMachineProperties) SetCreatedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDate = &formatted
}

func (o *LabVirtualMachineProperties) GetExpirationDateAsTime() (*time.Time, error) {
	if o.ExpirationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *LabVirtualMachineProperties) SetExpirationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationDate = &formatted
}
