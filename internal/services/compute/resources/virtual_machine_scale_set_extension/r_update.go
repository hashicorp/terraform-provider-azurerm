// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_extension

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesetextensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineScaleSetExtensionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetExtensionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesetextensions.ParseVirtualMachineScaleSetExtensionID(d.Id())
	if err != nil {
		return err
	}

	props := virtualmachinescalesetextensions.VirtualMachineScaleSetExtensionProperties{
		// if this isn't specified it defaults to false
		AutoUpgradeMinorVersion: pointer.To(d.Get("auto_upgrade_minor_version").(bool)),
		EnableAutomaticUpgrade:  pointer.To(d.Get("automatic_upgrade_enabled").(bool)),
	}

	if d.HasChange("failure_suppression_enabled") {
		props.SuppressFailures = pointer.To(d.Get("failure_suppression_enabled").(bool))
	}

	if d.HasChange("force_update_tag") {
		props.ForceUpdateTag = pointer.To(d.Get("force_update_tag").(string))
	}

	if d.HasChange("protected_settings") {
		var protectedSettings any
		if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
			var result any
			if err := json.Unmarshal([]byte(protectedSettingsString), &result); err != nil {
				return fmt.Errorf("unmarshaling `protected_settings`: %+v", err)
			}

			protectedSettings = result
		}
		props.ProtectedSettings = pointer.To(protectedSettings)
	}

	if d.HasChange("protected_settings_from_key_vault") {
		props.ProtectedSettingsFromKeyVault = helpers.ExpandProtectedSettingsFromKeyVaultOldVMSSExtension(d.Get("protected_settings_from_key_vault").([]any))
	}

	if d.HasChange("provision_after_extensions") {
		provisionAfterExtensionsRaw := d.Get("provision_after_extensions").([]any)
		props.ProvisionAfterExtensions = pluginsdk.ExpandStringSlice(provisionAfterExtensionsRaw)
	}

	if d.HasChange("publisher") {
		props.Publisher = pointer.To(d.Get("publisher").(string))
	}

	if d.HasChange("settings") {
		var settings any
		if settingsString := d.Get("settings").(string); settingsString != "" {
			var result any
			if err := json.Unmarshal([]byte(settingsString), &result); err != nil {
				return fmt.Errorf("unmarshaling `settings`: %+v", err)
			}
			settings = result
		}

		props.Settings = pointer.To(settings)
	}

	if d.HasChange("type") {
		props.Type = pointer.To(d.Get("type").(string))
	}

	if d.HasChange("type_handler_version") {
		props.TypeHandlerVersion = pointer.To(d.Get("type_handler_version").(string))
	}

	extension := virtualmachinescalesetextensions.VirtualMachineScaleSetExtension{
		Name:       pointer.To(id.ExtensionName),
		Properties: &props,
	}
	if err := client.CreateOrUpdateThenPoll(ctx, *id, extension); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceVirtualMachineScaleSetExtensionRead(d, meta)
}
