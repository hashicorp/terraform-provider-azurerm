package publicmaintenanceconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InputPatchConfiguration struct {
	LinuxParameters   *InputLinuxParameters             `json:"linuxParameters,omitempty"`
	RebootSetting     *RebootOptions                    `json:"rebootSetting,omitempty"`
	Tasks             *SoftwareUpdateConfigurationTasks `json:"tasks,omitempty"`
	WindowsParameters *InputWindowsParameters           `json:"windowsParameters,omitempty"`
}
