// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_implicit_data_disk_from_source

func (r Resource) ModelObject() any {
	return &VirtualMachineImplicitDataDiskFromSourceResourceModel{}
}

type VirtualMachineImplicitDataDiskFromSourceResourceModel struct {
	Name                    string `tfschema:"name"`
	VirtualMachineId        string `tfschema:"virtual_machine_id"`
	Lun                     int64  `tfschema:"lun"`
	Caching                 string `tfschema:"caching"`
	CreateOption            string `tfschema:"create_option"`
	DiskSizeGb              int64  `tfschema:"disk_size_gb"`
	SourceResourceId        string `tfschema:"source_resource_id"`
	WriteAcceleratorEnabled bool   `tfschema:"write_accelerator_enabled"`
}
