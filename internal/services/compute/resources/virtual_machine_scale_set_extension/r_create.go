// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_extension

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesetextensions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineScaleSetExtensionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetExtensionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachineScaleSetId, err := commonids.ParseVirtualMachineScaleSetID(d.Get("virtual_machine_scale_set_id").(string))
	if err != nil {
		return err
	}
	id := virtualmachinescalesetextensions.NewVirtualMachineScaleSetExtensionID(virtualMachineScaleSetId.SubscriptionId, virtualMachineScaleSetId.ResourceGroupName, virtualMachineScaleSetId.VirtualMachineScaleSetName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.Get(ctx, id, virtualmachinescalesetextensions.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_machine_scale_set_extension", id.ID())
		}
	}

	var settings *any
	if settingsString := d.Get("settings").(string); settingsString != "" {
		var result any
		if err := json.Unmarshal([]byte(settingsString), &result); err != nil {
			return fmt.Errorf("unmarshaling `settings`: %+v", err)
		}
		settings = pointer.To(result)
	}

	provisionAfterExtensionsRaw := d.Get("provision_after_extensions").([]any)
	provisionAfterExtensions := pluginsdk.ExpandStringSlice(provisionAfterExtensionsRaw)

	props := virtualmachinescalesetextensions.VirtualMachineScaleSetExtension{
		Name: pointer.To(id.ExtensionName),
		Properties: &virtualmachinescalesetextensions.VirtualMachineScaleSetExtensionProperties{
			Publisher:                     pointer.To(d.Get("publisher").(string)),
			Type:                          pointer.To(d.Get("type").(string)),
			TypeHandlerVersion:            pointer.To(d.Get("type_handler_version").(string)),
			AutoUpgradeMinorVersion:       pointer.To(d.Get("auto_upgrade_minor_version").(bool)),
			EnableAutomaticUpgrade:        pointer.To(d.Get("automatic_upgrade_enabled").(bool)),
			SuppressFailures:              pointer.To(d.Get("failure_suppression_enabled").(bool)),
			ProtectedSettingsFromKeyVault: helpers.ExpandProtectedSettingsFromKeyVaultOldVMSSExtension(d.Get("protected_settings_from_key_vault").([]any)),
			ProvisionAfterExtensions:      provisionAfterExtensions,
			Settings:                      settings,
		},
	}
	if v, ok := d.GetOk("force_update_tag"); ok {
		props.Properties.ForceUpdateTag = pointer.To(v.(string))
	}

	if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
		var result any
		if err := json.Unmarshal([]byte(protectedSettingsString), &result); err != nil {
			return fmt.Errorf("unmarshaling `protected_settings`: %+v", err)
		}
		props.Properties.ProtectedSettings = pointer.To(result)
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualMachineScaleSetExtensionRead(d, meta)
}
