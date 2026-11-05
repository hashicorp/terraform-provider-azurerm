// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_extension

func (r Resource) ModelObject() any {
	return &StackHCIExtensionResourceModel{}
}

type StackHCIExtensionResourceModel struct {
	Name                           string `tfschema:"name"`
	ArcSettingId                   string `tfschema:"arc_setting_id"`
	AutoUpgradeMinorVersionEnabled bool   `tfschema:"auto_upgrade_minor_version_enabled"`
	AutomaticUpgradeEnabled        bool   `tfschema:"automatic_upgrade_enabled"`
	ProtectedSettings              string `tfschema:"protected_settings"`
	Publisher                      string `tfschema:"publisher"`
	Settings                       string `tfschema:"settings"`
	Type                           string `tfschema:"type"`
	TypeHandlerVersion             string `tfschema:"type_handler_version"`
}
