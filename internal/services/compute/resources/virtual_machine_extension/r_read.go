// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_extension

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineextensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineExtensionsRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineExtensionsClient
	vmClient := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachineextensions.ParseExtensionID(d.Id())
	if err != nil {
		return err
	}

	virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineName)

	virtualMachine, err := vmClient.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(virtualMachine.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", virtualMachineId, err)
	}

	d.Set("virtual_machine_id", virtualMachineId.ID())

	resp, err := client.Get(ctx, *id, virtualmachineextensions.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %s", id.ExtensionName, err)
	}

	d.Set("name", id.ExtensionName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("publisher", props.Publisher)
			d.Set("type", props.Type)
			d.Set("type_handler_version", props.TypeHandlerVersion)
			d.Set("auto_upgrade_minor_version", props.AutoUpgradeMinorVersion)
			d.Set("automatic_upgrade_enabled", props.EnableAutomaticUpgrade)
			d.Set("protected_settings_from_key_vault", helpers.FlattenProtectedSettingsFromKeyVault(props.ProtectedSettingsFromKeyVault))
			d.Set("provision_after_extensions", pointer.From(props.ProvisionAfterExtensions))

			d.Set("failure_suppression_enabled", pointer.From(props.SuppressFailures))

			if props.Settings != nil {
				settings, err := json.Marshal(props.Settings)
				if err != nil {
					return fmt.Errorf("unmarshaling `settings`: %+v", err)
				}
				d.Set("settings", string(settings))
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
