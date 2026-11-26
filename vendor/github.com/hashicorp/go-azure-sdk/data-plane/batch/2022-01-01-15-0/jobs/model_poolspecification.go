package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolSpecification struct {
	ApplicationLicenses          *[]string                      `json:"applicationLicenses,omitempty"`
	ApplicationPackageReferences *[]ApplicationPackageReference `json:"applicationPackageReferences,omitempty"`
	AutoScaleEvaluationInterval  *string                        `json:"autoScaleEvaluationInterval,omitempty"`
	AutoScaleFormula             *string                        `json:"autoScaleFormula,omitempty"`
	CertificateReferences        *[]CertificateReference        `json:"certificateReferences,omitempty"`
	CloudServiceConfiguration    *CloudServiceConfiguration     `json:"cloudServiceConfiguration,omitempty"`
	DisplayName                  *string                        `json:"displayName,omitempty"`
	EnableAutoScale              *bool                          `json:"enableAutoScale,omitempty"`
	EnableInterNodeCommunication *bool                          `json:"enableInterNodeCommunication,omitempty"`
	Metadata                     *[]MetadataItem                `json:"metadata,omitempty"`
	MountConfiguration           *[]MountConfiguration          `json:"mountConfiguration,omitempty"`
	NetworkConfiguration         *NetworkConfiguration          `json:"networkConfiguration,omitempty"`
	ResizeTimeout                *string                        `json:"resizeTimeout,omitempty"`
	StartTask                    *StartTask                     `json:"startTask,omitempty"`
	TargetDedicatedNodes         *int64                         `json:"targetDedicatedNodes,omitempty"`
	TargetLowPriorityNodes       *int64                         `json:"targetLowPriorityNodes,omitempty"`
	TaskSchedulingPolicy         *TaskSchedulingPolicy          `json:"taskSchedulingPolicy,omitempty"`
	TaskSlotsPerNode             *int64                         `json:"taskSlotsPerNode,omitempty"`
	UserAccounts                 *[]UserAccount                 `json:"userAccounts,omitempty"`
	VMSize                       string                         `json:"vmSize"`
	VirtualMachineConfiguration  *VirtualMachineConfiguration   `json:"virtualMachineConfiguration,omitempty"`
}
