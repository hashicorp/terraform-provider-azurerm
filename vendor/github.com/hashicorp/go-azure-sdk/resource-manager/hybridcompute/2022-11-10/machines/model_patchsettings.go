package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PatchSettings struct {
	AssessmentMode *AssessmentModeTypes `json:"assessmentMode,omitempty"`
	PatchMode      *PatchModeTypes      `json:"patchMode,omitempty"`
}
