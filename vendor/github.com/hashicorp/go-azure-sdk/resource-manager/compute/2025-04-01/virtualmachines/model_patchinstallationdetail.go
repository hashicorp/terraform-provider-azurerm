package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PatchInstallationDetail struct {
	Classifications   *[]string               `json:"classifications,omitempty"`
	InstallationState *PatchInstallationState `json:"installationState,omitempty"`
	KbId              *string                 `json:"kbId,omitempty"`
	Name              *string                 `json:"name,omitempty"`
	PatchId           *string                 `json:"patchId,omitempty"`
	Version           *string                 `json:"version,omitempty"`
}
