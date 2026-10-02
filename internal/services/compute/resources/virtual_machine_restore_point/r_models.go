// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_restore_point

func (r Resource) ModelObject() any {
	return &VirtualMachineRestorePointResourceModel{}
}

type VirtualMachineRestorePointResourceModel struct {
	Name                                   string   `tfschema:"name"`
	VirtualMachineRestorePointCollectionId string   `tfschema:"virtual_machine_restore_point_collection_id"`
	CrashConsistencyModeEnabled            bool     `tfschema:"crash_consistency_mode_enabled"`
	ExcludedDisks                          []string `tfschema:"excluded_disks"`
}
