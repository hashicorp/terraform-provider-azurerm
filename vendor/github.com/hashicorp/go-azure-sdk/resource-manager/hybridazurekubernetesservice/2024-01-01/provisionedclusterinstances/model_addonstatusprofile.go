package provisionedclusterinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AddonStatusProfile struct {
	ErrorMessage *string     `json:"errorMessage,omitempty"`
	Name         *string     `json:"name,omitempty"`
	Phase        *AddonPhase `json:"phase,omitempty"`
	Ready        *bool       `json:"ready,omitempty"`
}
