// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_availability_set

type SystemCenterVirtualMachineManagerAvailabilitySetModel struct {
	Name                                      string            `tfschema:"name"`
	ResourceGroupName                         string            `tfschema:"resource_group_name"`
	Location                                  string            `tfschema:"location"`
	CustomLocationId                          string            `tfschema:"custom_location_id"`
	SystemCenterVirtualMachineManagerServerId string            `tfschema:"system_center_virtual_machine_manager_server_id"`
	Tags                                      map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &SystemCenterVirtualMachineManagerAvailabilitySetModel{}
}
