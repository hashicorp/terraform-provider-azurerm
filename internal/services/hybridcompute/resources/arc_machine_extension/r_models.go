// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_machine_extension

type MachineExtensionModel struct {
	Name                   string            `tfschema:"name"`
	HybridComputeMachineId string            `tfschema:"arc_machine_id"`
	EnableAutomaticUpgrade bool              `tfschema:"automatic_upgrade_enabled"`
	ForceUpdateTag         string            `tfschema:"force_update_tag"`
	Location               string            `tfschema:"location"`
	ProtectedSettings      string            `tfschema:"protected_settings"`
	Publisher              string            `tfschema:"publisher"`
	Settings               string            `tfschema:"settings"`
	Tags                   map[string]string `tfschema:"tags"`
	Type                   string            `tfschema:"type"`
	TypeHandlerVersion     string            `tfschema:"type_handler_version"`
}

func (r Resource) ModelObject() any {
	return &MachineExtensionModel{}
}
