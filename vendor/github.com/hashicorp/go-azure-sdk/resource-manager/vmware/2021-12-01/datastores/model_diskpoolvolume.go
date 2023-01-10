package datastores

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskPoolVolume struct {
	LunName     string           `json:"lunName"`
	MountOption *MountOptionEnum `json:"mountOption,omitempty"`
	Path        *string          `json:"path,omitempty"`
	TargetId    string           `json:"targetId"`
}
