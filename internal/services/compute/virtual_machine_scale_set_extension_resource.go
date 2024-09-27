// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetextensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// NOTE (also in the docs): this is not intended to be used with the `azurerm_virtual_machine_scale_set` resource

func resourceVirtualMachineScaleSetExtension() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineScaleSetExtensionCreate,
		Read:   resourceVirtualMachineScaleSetExtensionRead,
		Update: resourceVirtualMachineScaleSetExtensionUpdate,
		Delete: resourceVirtualMachineScaleSetExtensionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualmachinescalesetextensions.ParseVirtualMachineScaleSetExtensionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringDoesNotContainAny("/"),
				),
			},

			"virtual_machine_scale_set_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateVirtualMachineScaleSetID,
			},

			"publisher": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"type": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"type_handler_version": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"auto_upgrade_minor_version": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"automatic_upgrade_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"failure_suppression_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"force_update_tag": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"protected_settings": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Sensitive:        true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			"protected_settings_from_key_vault": protectedSettingsFromKeyVaultSchema(true),

			"provision_after_extensions": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"settings": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},
		},
	}
}

func resourceVirtualMachineScaleSetExtensionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetExtensionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachineScaleSetId, err := commonids.ParseVirtualMachineScaleSetID(d.Get("virtual_machine_scale_set_id").(string))
	if err != nil {
		return err
	}
	id := virtualmachinescalesetextensions.NewVirtualMachineScaleSetExtensionID(virtualMachineScaleSetId.SubscriptionId, virtualMachineScaleSetId.ResourceGroupName, virtualMachineScaleSetId.VirtualMachineScaleSetName, d.Get("name").(string))

	resp, err := client.Get(ctx, id, virtualmachinescalesetextensions.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(resp.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_virtual_machine_scale_set_extension", id.ID())
	}

	var settings *interface{}
	if settingsString := d.Get("settings").(string); settingsString != "" {
		var result interface{}
		err := json.Unmarshal([]byte(settingsString), &result)
		if err != nil {
			return fmt.Errorf("unmarshaling `settings`: %+v", err)
		}
		settings = pointer.To(result)
	}

	provisionAfterExtensionsRaw := d.Get("provision_after_extensions").([]interface{})
	provisionAfterExtensions := utils.ExpandStringSlice(provisionAfterExtensionsRaw)

	props := virtualmachinescalesetextensions.VirtualMachineScaleSetExtension{
		Name: pointer.To(id.ExtensionName),
		Properties: &virtualmachinescalesetextensions.VirtualMachineScaleSetExtensionProperties{
			Publisher:                     pointer.To(d.Get("publisher").(string)),
			Type:                          pointer.To(d.Get("type").(string)),
			TypeHandlerVersion:            pointer.To(d.Get("type_handler_version").(string)),
			AutoUpgradeMinorVersion:       pointer.To(d.Get("auto_upgrade_minor_version").(bool)),
			EnableAutomaticUpgrade:        pointer.To(d.Get("automatic_upgrade_enabled").(bool)),
			SuppressFailures:              pointer.To(d.Get("failure_suppression_enabled").(bool)),
			ProtectedSettingsFromKeyVault: expandProtectedSettingsFromKeyVaultOldVMSSExtension(d.Get("protected_settings_from_key_vault").([]interface{})),
			ProvisionAfterExtensions:      provisionAfterExtensions,
			Settings:                      settings,
		},
	}
	if v, ok := d.GetOk("force_update_tag"); ok {
		props.Properties.ForceUpdateTag = pointer.To(v.(string))
	}

	if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
		var result interface{}
		err := json.Unmarshal([]byte(protectedSettingsString), &result)
		if err != nil {
			return fmt.Errorf("unmarshaling `protected_settings`: %+v", err)
		}
		props.Properties.ProtectedSettings = pointer.To(result)
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualMachineScaleSetExtensionRead(d, meta)
}

func resourceVirtualMachineScaleSetExtensionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		var protectedSettings interface{}
		if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
			var result interface{}
			err := json.Unmarshal([]byte(protectedSettingsString), &result)
			if err != nil {
				return fmt.Errorf("unmarshaling `protected_settings`: %+v", err)
			}

			protectedSettings = result
		}
		props.ProtectedSettings = pointer.To(protectedSettings)
	}

	if d.HasChange("protected_settings_from_key_vault") {
		props.ProtectedSettingsFromKeyVault = expandProtectedSettingsFromKeyVaultOldVMSSExtension(d.Get("protected_settings_from_key_vault").([]interface{}))
	}

	if d.HasChange("provision_after_extensions") {
		provisionAfterExtensionsRaw := d.Get("provision_after_extensions").([]interface{})
		props.ProvisionAfterExtensions = utils.ExpandStringSlice(provisionAfterExtensionsRaw)
	}

	if d.HasChange("publisher") {
		props.Publisher = pointer.To(d.Get("publisher").(string))
	}

	if d.HasChange("settings") {
		var settings interface{}
		if settingsString := d.Get("settings").(string); settingsString != "" {
			var result interface{}
			err := json.Unmarshal([]byte(settingsString), &result)
			if err != nil {
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

func resourceVirtualMachineScaleSetExtensionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			d.Set("protected_settings_from_key_vault", flattenProtectedSettingsFromKeyVaultOldVMSSExtension(props.ProtectedSettingsFromKeyVault))
			d.Set("provision_after_extensions", utils.FlattenStringSlice(props.ProvisionAfterExtensions))
			d.Set("publisher", props.Publisher)
			d.Set("type", props.Type)
			d.Set("type_handler_version", props.TypeHandlerVersion)

			suppressFailure := false
			if props.SuppressFailures != nil {
				suppressFailure = *props.SuppressFailures
			}
			d.Set("failure_suppression_enabled", suppressFailure)

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

func resourceVirtualMachineScaleSetExtensionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetExtensionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesetextensions.ParseVirtualMachineScaleSetExtensionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
