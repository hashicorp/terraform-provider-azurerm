package virtualmachinescalesetvms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PatchSettings struct {
	AssessmentMode              *WindowsPatchAssessmentMode                     `json:"assessmentMode,omitempty"`
	AutomaticByPlatformSettings *WindowsVMGuestPatchAutomaticByPlatformSettings `json:"automaticByPlatformSettings,omitempty"`
	EnableHotpatching           *bool                                           `json:"enableHotpatching,omitempty"`
	PatchMode                   *WindowsVMGuestPatchMode                        `json:"patchMode,omitempty"`
}
