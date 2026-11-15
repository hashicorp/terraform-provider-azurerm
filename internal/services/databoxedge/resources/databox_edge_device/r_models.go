// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databox_edge_device

type DevicePropertiesModel struct {
	Capacity            int64    `tfschema:"capacity"`
	ConfiguredRoleTypes []string `tfschema:"configured_role_types"`
	Culture             string   `tfschema:"culture"`
	HcsVersion          string   `tfschema:"hcs_version"`
	Model               string   `tfschema:"model"`
	NodeCount           int64    `tfschema:"node_count"`
	SerialNumber        string   `tfschema:"serial_number"`
	SoftwareVersion     string   `tfschema:"software_version"`
	Status              string   `tfschema:"status"`
	TimeZone            string   `tfschema:"time_zone"`
	Type                string   `tfschema:"type"`
}

type EdgeDeviceModel struct {
	DeviceProperties  []DevicePropertiesModel `tfschema:"device_properties"`
	Location          string                  `tfschema:"location"`
	Name              string                  `tfschema:"name"`
	ResourceGroupName string                  `tfschema:"resource_group_name"`
	SkuName           string                  `tfschema:"sku_name"`
	Tags              map[string]string       `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &EdgeDeviceModel{}
}
