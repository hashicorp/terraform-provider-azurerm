package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PatchSettings struct {
	AssessmentMode    *AssessmentModeTypes `json:"assessmentMode,omitempty"`
	EnableHotpatching *bool                `json:"enableHotpatching,omitempty"`
	PatchMode         *PatchModeTypes      `json:"patchMode,omitempty"`
	Status            *PatchSettingsStatus `json:"status,omitempty"`
}
