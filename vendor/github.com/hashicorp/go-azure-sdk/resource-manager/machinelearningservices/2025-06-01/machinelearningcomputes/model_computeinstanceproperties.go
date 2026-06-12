package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeInstanceProperties struct {
	ApplicationSharingPolicy         *ApplicationSharingPolicy             `json:"applicationSharingPolicy,omitempty"`
	Applications                     *[]ComputeInstanceApplication         `json:"applications,omitempty"`
	ComputeInstanceAuthorizationType *ComputeInstanceAuthorizationType     `json:"computeInstanceAuthorizationType,omitempty"`
	ConnectivityEndpoints            *ComputeInstanceConnectivityEndpoints `json:"connectivityEndpoints,omitempty"`
	Containers                       *[]ComputeInstanceContainer           `json:"containers,omitempty"`
	CreatedBy                        *ComputeInstanceCreatedBy             `json:"createdBy,omitempty"`
	CustomServices                   *[]CustomService                      `json:"customServices,omitempty"`
	DataDisks                        *[]ComputeInstanceDataDisk            `json:"dataDisks,omitempty"`
	DataMounts                       *[]ComputeInstanceDataMount           `json:"dataMounts,omitempty"`
	EnableNodePublicIP               *bool                                 `json:"enableNodePublicIp,omitempty"`
	Errors                           *[]ErrorResponse                      `json:"errors,omitempty"`
	LastOperation                    *ComputeInstanceLastOperation         `json:"lastOperation,omitempty"`
	OsImageMetadata                  *ImageMetadata                        `json:"osImageMetadata,omitempty"`
	PersonalComputeInstanceSettings  *PersonalComputeInstanceSettings      `json:"personalComputeInstanceSettings,omitempty"`
	Schedules                        *ComputeSchedules                     `json:"schedules,omitempty"`
	SetupScripts                     *SetupScripts                         `json:"setupScripts,omitempty"`
	SshSettings                      *ComputeInstanceSshSettings           `json:"sshSettings,omitempty"`
	State                            *ComputeInstanceState                 `json:"state,omitempty"`
	Subnet                           *ResourceId                           `json:"subnet,omitempty"`
	VMSize                           *string                               `json:"vmSize,omitempty"`
	Versions                         *ComputeInstanceVersion               `json:"versions,omitempty"`
}
