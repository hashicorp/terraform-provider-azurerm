// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_cloud

type SystemCenterVirtualMachineManagerCloudModel struct {
	Name                                                   string            `tfschema:"name"`
	Location                                               string            `tfschema:"location"`
	ResourceGroupName                                      string            `tfschema:"resource_group_name"`
	CustomLocationId                                       string            `tfschema:"custom_location_id"`
	SystemCenterVirtualMachineManagerServerInventoryItemId string            `tfschema:"system_center_virtual_machine_manager_server_inventory_item_id"`
	Tags                                                   map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &SystemCenterVirtualMachineManagerCloudModel{}
}
