// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_hyperv_network_mapping

type HyperVNetworkMappingModel struct {
	Name            string `tfschema:"name"`
	VaultId         string `tfschema:"recovery_vault_id"`
	SCVMMname       string `tfschema:"source_system_center_virtual_machine_manager_name"`
	NetworkName     string `tfschema:"source_network_name"`
	TargetNetworkId string `tfschema:"target_network_id"`
}

func (s Resource) ModelObject() any {
	return &HyperVNetworkMappingModel{}
}
