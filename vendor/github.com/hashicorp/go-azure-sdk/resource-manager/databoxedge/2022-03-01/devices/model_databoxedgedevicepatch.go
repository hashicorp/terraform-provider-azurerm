package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataBoxEdgeDevicePatch struct {
	Identity   *ResourceIdentity                 `json:"identity,omitempty"`
	Properties *DataBoxEdgeDevicePropertiesPatch `json:"properties,omitempty"`
	Tags       *map[string]string                `json:"tags,omitempty"`
}
