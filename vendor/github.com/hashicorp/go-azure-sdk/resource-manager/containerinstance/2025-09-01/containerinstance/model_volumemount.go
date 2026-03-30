package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeMount struct {
	MountPath string `json:"mountPath"`
	Name      string `json:"name"`
	ReadOnly  *bool  `json:"readOnly,omitempty"`
}
