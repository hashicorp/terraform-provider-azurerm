// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_restore_point_collection

func (r Resource) ModelObject() any {
	return &VirtualMachineRestorePointCollectionResourceModel{}
}

type VirtualMachineRestorePointCollectionResourceModel struct {
	Name                   string         `tfschema:"name"`
	ResourceGroup          string         `tfschema:"resource_group_name"`
	Location               string         `tfschema:"location"`
	SourceVirtualMachineId string         `tfschema:"source_virtual_machine_id"`
	Tags                   map[string]any `tfschema:"tags"`
}
