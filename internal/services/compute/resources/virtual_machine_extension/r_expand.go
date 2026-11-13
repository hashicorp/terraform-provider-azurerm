// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_extension

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineextensions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// expandCreateForVirtualMachineExtension builds the full ARM PUT body for a VM extension from
// a ResourceDiff at plan time. location must be resolved beforehand from the parent VM.
func expandCreateForVirtualMachineExtension(d *schema.ResourceDiff, location string) (virtualmachineextensions.VirtualMachineExtension, error) {
	publisher := d.Get("publisher").(string)
	extensionType := d.Get("type").(string)
	typeHandlerVersion := d.Get("type_handler_version").(string)
	autoUpgradeMinor := d.Get("auto_upgrade_minor_version").(bool)
	enableAutomaticUpgrade := d.Get("automatic_upgrade_enabled").(bool)
	suppressFailure := d.Get("failure_suppression_enabled").(bool)

	extension := virtualmachineextensions.VirtualMachineExtension{
		Location: &location,
		Properties: &virtualmachineextensions.VirtualMachineExtensionProperties{
			Publisher:                     &publisher,
			Type:                          &extensionType,
			TypeHandlerVersion:            &typeHandlerVersion,
			AutoUpgradeMinorVersion:       &autoUpgradeMinor,
			EnableAutomaticUpgrade:        &enableAutomaticUpgrade,
			ProtectedSettingsFromKeyVault: helpers.ExpandProtectedSettingsFromKeyVault(d.Get("protected_settings_from_key_vault").([]any)),
			SuppressFailures:              &suppressFailure,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if settingsString := d.Get("settings").(string); settingsString != "" {
		var result any
		if err := json.Unmarshal([]byte(settingsString), &result); err != nil {
			return extension, fmt.Errorf("unmarshaling `settings`: %+v", err)
		}
		extension.Properties.Settings = pointer.To(result)
	}

	if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
		var result any
		if err := json.Unmarshal([]byte(protectedSettingsString), &result); err != nil {
			return extension, fmt.Errorf("unmarshaling `protected_settings`: %+v", err)
		}
		extension.Properties.ProtectedSettings = pointer.To(result)
	}

	if provisionAfterExtensionsValue, exists := d.GetOk("provision_after_extensions"); exists {
		extension.Properties.ProvisionAfterExtensions = pluginsdk.ExpandStringSlice(provisionAfterExtensionsValue.([]any))
	}

	return extension, nil
}
