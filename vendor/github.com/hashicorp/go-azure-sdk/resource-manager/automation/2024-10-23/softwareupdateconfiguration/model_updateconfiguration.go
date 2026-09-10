package softwareupdateconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateConfiguration struct {
	AzureVirtualMachines  *[]string           `json:"azureVirtualMachines,omitempty"`
	Duration              *string             `json:"duration,omitempty"`
	Linux                 *LinuxProperties    `json:"linux,omitempty"`
	NonAzureComputerNames *[]string           `json:"nonAzureComputerNames,omitempty"`
	OperatingSystem       OperatingSystemType `json:"operatingSystem"`
	Targets               *TargetProperties   `json:"targets,omitempty"`
	Windows               *WindowsProperties  `json:"windows,omitempty"`
}
