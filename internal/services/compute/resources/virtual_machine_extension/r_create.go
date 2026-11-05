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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineExtensionsCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineExtensionsClient
	vmClient := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachineId, err := virtualmachines.ParseVirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return err
	}
	id := virtualmachineextensions.NewExtensionID(virtualMachineId.SubscriptionId, virtualMachineId.ResourceGroupName, virtualMachineId.VirtualMachineName, d.Get("name").(string))

	virtualMachine, err := vmClient.Get(ctx, *virtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", virtualMachineId, err)
	}

	if virtualMachine.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", virtualMachineId)
	}

	location := virtualMachine.Model.Location
	if location == "" {
		return fmt.Errorf("reading location of %s", virtualMachineId)
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, virtualmachineextensions.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_virtual_machine_extension", id.ID())
			}
		}
	}

	publisher := d.Get("publisher").(string)
	extensionType := d.Get("type").(string)
	typeHandlerVersion := d.Get("type_handler_version").(string)
	autoUpgradeMinor := d.Get("auto_upgrade_minor_version").(bool)
	enableAutomaticUpgrade := d.Get("automatic_upgrade_enabled").(bool)
	suppressFailure := d.Get("failure_suppression_enabled").(bool)
	t := d.Get("tags").(map[string]any)

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
		Tags: tags.Expand(t),
	}

	if settingsString := d.Get("settings").(string); settingsString != "" {
		var result any
		if err := json.Unmarshal([]byte(settingsString), &result); err != nil {
			return fmt.Errorf("unmarshaling `settings`: %+v", err)
		}
		extension.Properties.Settings = pointer.To(result)
	}

	if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
		var result any
		if err := json.Unmarshal([]byte(protectedSettingsString), &result); err != nil {
			return fmt.Errorf("unmarshaling `protected_settings`: %+v", err)
		}
		extension.Properties.ProtectedSettings = pointer.To(result)
	}

	if provisionAfterExtensionsValue, exists := d.GetOk("provision_after_extensions"); exists {
		extension.Properties.ProvisionAfterExtensions = pluginsdk.ExpandStringSlice(provisionAfterExtensionsValue.([]any))
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, extension, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, extension); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceVirtualMachineExtensionsRead(d, meta)
}
