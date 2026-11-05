// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_extension

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesetextensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineScaleSetExtensionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetExtensionsClient
	vmssClient := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesetextensions.ParseVirtualMachineScaleSetExtensionID(d.Id())
	if err != nil {
		return err
	}

	virtualMachineScaleSetId := virtualmachinescalesets.NewVirtualMachineScaleSetID(id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineScaleSetName)

	vmss, err := vmssClient.Get(ctx, virtualMachineScaleSetId, virtualmachinescalesets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(vmss.HttpResponse) {
			log.Printf("%s was not found - removing Extension from state!", virtualMachineScaleSetId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", virtualMachineScaleSetId, err)
	}

	resp, err := client.Get(ctx, *id, virtualmachinescalesetextensions.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("%s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ExtensionName)
	d.Set("virtual_machine_scale_set_id", virtualMachineScaleSetId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("auto_upgrade_minor_version", props.AutoUpgradeMinorVersion)
			d.Set("automatic_upgrade_enabled", props.EnableAutomaticUpgrade)
			d.Set("force_update_tag", props.ForceUpdateTag)
			d.Set("protected_settings_from_key_vault", helpers.FlattenProtectedSettingsFromKeyVaultOldVMSSExtension(props.ProtectedSettingsFromKeyVault))
			d.Set("provision_after_extensions", pluginsdk.FlattenSlice(props.ProvisionAfterExtensions))
			d.Set("publisher", props.Publisher)
			d.Set("type", props.Type)
			d.Set("type_handler_version", props.TypeHandlerVersion)

			d.Set("failure_suppression_enabled", pointer.From(props.SuppressFailures))

			settings := ""
			if props.Settings != nil {
				settingsRaw, err := json.Marshal(props.Settings)
				if err != nil {
					return fmt.Errorf("unmarshaling `settings`: %+v", err)
				}
				settings = string(settingsRaw)
			}
			d.Set("settings", settings)
		}
	}

	return nil
}
