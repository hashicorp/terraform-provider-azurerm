package labs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LabVirtualMachineCreationParameterProperties struct {
	AllowClaim                 *bool                        `json:"allowClaim,omitempty"`
	Artifacts                  *[]ArtifactInstallProperties `json:"artifacts,omitempty"`
	BulkCreationParameters     *BulkCreationParameters      `json:"bulkCreationParameters,omitempty"`
	CreatedDate                *string                      `json:"createdDate,omitempty"`
	CustomImageId              *string                      `json:"customImageId,omitempty"`
	DataDiskParameters         *[]DataDiskProperties        `json:"dataDiskParameters,omitempty"`
	DisallowPublicIPAddress    *bool                        `json:"disallowPublicIpAddress,omitempty"`
	EnvironmentId              *string                      `json:"environmentId,omitempty"`
	ExpirationDate             *string                      `json:"expirationDate,omitempty"`
	GalleryImageReference      *GalleryImageReference       `json:"galleryImageReference,omitempty"`
	IsAuthenticationWithSshKey *bool                        `json:"isAuthenticationWithSshKey,omitempty"`
	LabSubnetName              *string                      `json:"labSubnetName,omitempty"`
	LabVirtualNetworkId        *string                      `json:"labVirtualNetworkId,omitempty"`
	NetworkInterface           *NetworkInterfaceProperties  `json:"networkInterface,omitempty"`
	Notes                      *string                      `json:"notes,omitempty"`
	OwnerObjectId              *string                      `json:"ownerObjectId,omitempty"`
	OwnerUserPrincipalName     *string                      `json:"ownerUserPrincipalName,omitempty"`
	Password                   *string                      `json:"password,omitempty"`
	PlanId                     *string                      `json:"planId,omitempty"`
	ScheduleParameters         *[]ScheduleCreationParameter `json:"scheduleParameters,omitempty"`
	Size                       *string                      `json:"size,omitempty"`
	SshKey                     *string                      `json:"sshKey,omitempty"`
	StorageType                *string                      `json:"storageType,omitempty"`
	UserName                   *string                      `json:"userName,omitempty"`
}

func (o *LabVirtualMachineCreationParameterProperties) GetCreatedDateAsTime() (*time.Time, error) {
	if o.CreatedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *LabVirtualMachineCreationParameterProperties) SetCreatedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDate = &formatted
}

func (o *LabVirtualMachineCreationParameterProperties) GetExpirationDateAsTime() (*time.Time, error) {
	if o.ExpirationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *LabVirtualMachineCreationParameterProperties) SetExpirationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationDate = &formatted
}
