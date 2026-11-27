// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover_agent

type StorageMoverAgentResourceModel struct {
	Name                string `tfschema:"name"`
	StorageMoverId      string `tfschema:"storage_mover_id"`
	ArcVirtualMachineId string `tfschema:"arc_virtual_machine_id"`
	ArcVmUuid           string `tfschema:"arc_virtual_machine_uuid"`
	Description         string `tfschema:"description"`
}

func (r Resource) ModelObject() any {
	return &StorageMoverAgentResourceModel{}
}
