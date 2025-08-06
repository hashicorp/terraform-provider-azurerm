package restorepointcollections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinuxPatchSettings struct {
	AssessmentMode              *LinuxPatchAssessmentMode                     `json:"assessmentMode,omitempty"`
	AutomaticByPlatformSettings *LinuxVMGuestPatchAutomaticByPlatformSettings `json:"automaticByPlatformSettings,omitempty"`
	PatchMode                   *LinuxVMGuestPatchMode                        `json:"patchMode,omitempty"`
}
