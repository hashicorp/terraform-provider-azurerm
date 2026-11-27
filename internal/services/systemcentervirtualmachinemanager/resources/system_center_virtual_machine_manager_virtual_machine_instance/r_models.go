// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_virtual_machine_instance

type SystemCenterVirtualMachineManagerVirtualMachineInstanceModel struct {
	ScopedResourceId                                    string             `tfschema:"scoped_resource_id"`
	CustomLocationId                                    string             `tfschema:"custom_location_id"`
	Hardware                                            []Hardware         `tfschema:"hardware"`
	Infrastructure                                      []Infrastructure   `tfschema:"infrastructure"`
	NetworkInterfaces                                   []NetworkInterface `tfschema:"network_interface"`
	OperatingSystem                                     []OperatingSystem  `tfschema:"operating_system"`
	StorageDisks                                        []StorageDisk      `tfschema:"storage_disk"`
	SystemCenterVirtualMachineManagerAvailabilitySetIds []string           `tfschema:"system_center_virtual_machine_manager_availability_set_ids"`
}

func (r Resource) ModelObject() any {
	return &SystemCenterVirtualMachineManagerVirtualMachineInstanceModel{}
}
