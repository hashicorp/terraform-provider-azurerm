package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeMount struct {
	ContainerMountPath string  `json:"containerMountPath"`
	Data               *string `json:"data,omitempty"`
	ReadOnly           *bool   `json:"readOnly,omitempty"`
	VolumeSubPath      string  `json:"volumeSubPath"`
}
