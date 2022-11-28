package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeInstanceProperties struct {
	ApplicationSharingPolicy         *ApplicationSharingPolicy             `json:"applicationSharingPolicy,omitempty"`
	Applications                     *[]ComputeInstanceApplication         `json:"applications,omitempty"`
	ComputeInstanceAuthorizationType *ComputeInstanceAuthorizationType     `json:"computeInstanceAuthorizationType,omitempty"`
	ConnectivityEndpoints            *ComputeInstanceConnectivityEndpoints `json:"connectivityEndpoints"`
	Containers                       *[]ComputeInstanceContainer           `json:"containers,omitempty"`
	CreatedBy                        *ComputeInstanceCreatedBy             `json:"createdBy"`
	DataDisks                        *[]ComputeInstanceDataDisk            `json:"dataDisks,omitempty"`
	DataMounts                       *[]ComputeInstanceDataMount           `json:"dataMounts,omitempty"`
	EnableNodePublicIP               *bool                                 `json:"enableNodePublicIp,omitempty"`
	Errors                           *[]ErrorResponse                      `json:"errors,omitempty"`
	LastOperation                    *ComputeInstanceLastOperation         `json:"lastOperation"`
	PersonalComputeInstanceSettings  *PersonalComputeInstanceSettings      `json:"personalComputeInstanceSettings"`
	Schedules                        *ComputeSchedules                     `json:"schedules"`
	SetupScripts                     *SetupScripts                         `json:"setupScripts"`
	SshSettings                      *ComputeInstanceSshSettings           `json:"sshSettings"`
	State                            *ComputeInstanceState                 `json:"state,omitempty"`
	Subnet                           *ResourceId                           `json:"subnet"`
	Versions                         *ComputeInstanceVersion               `json:"versions"`
	VmSize                           *string                               `json:"vmSize,omitempty"`
}
