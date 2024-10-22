package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSProfileLinuxConfiguration struct {
	PatchSettings *PatchSettings `json:"patchSettings,omitempty"`
}
