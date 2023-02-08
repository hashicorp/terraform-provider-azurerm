package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataBoxEdgeDeviceProperties struct {
	ConfiguredRoleTypes     *[]RoleTypes             `json:"configuredRoleTypes,omitempty"`
	Culture                 *string                  `json:"culture,omitempty"`
	DataBoxEdgeDeviceStatus *DataBoxEdgeDeviceStatus `json:"dataBoxEdgeDeviceStatus,omitempty"`
	Description             *string                  `json:"description,omitempty"`
	DeviceHcsVersion        *string                  `json:"deviceHcsVersion,omitempty"`
	DeviceLocalCapacity     *int64                   `json:"deviceLocalCapacity,omitempty"`
	DeviceModel             *string                  `json:"deviceModel,omitempty"`
	DeviceSoftwareVersion   *string                  `json:"deviceSoftwareVersion,omitempty"`
	DeviceType              *DeviceType              `json:"deviceType,omitempty"`
	EdgeProfile             *EdgeProfile             `json:"edgeProfile,omitempty"`
	FriendlyName            *string                  `json:"friendlyName,omitempty"`
	ModelDescription        *string                  `json:"modelDescription,omitempty"`
	NodeCount               *int64                   `json:"nodeCount,omitempty"`
	ResourceMoveDetails     *ResourceMoveDetails     `json:"resourceMoveDetails,omitempty"`
	SerialNumber            *string                  `json:"serialNumber,omitempty"`
	TimeZone                *string                  `json:"timeZone,omitempty"`
}
