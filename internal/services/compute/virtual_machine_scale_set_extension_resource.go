// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

// NOTE (also in the docs): this is not intended to be used with the `azurerm_virtual_machine_scale_set` resource

func resourceVirtualMachineScaleSetExtension() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineScaleSetExtensionCreate,
		Read:   resourceVirtualMachineScaleSetExtensionRead,
		Update: resourceVirtualMachineScaleSetExtensionUpdate,
		Delete: resourceVirtualMachineScaleSetExtensionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualMachineScaleSetExtensionID(id)
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
				ValidateFunc: validate.VirtualMachineScaleSetID,
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
	client := meta.(*clients.Client).Compute.VMScaleSetExtensionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachineScaleSetId, err := parse.VirtualMachineScaleSetID(d.Get("virtual_machine_scale_set_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewVirtualMachineScaleSetExtensionID(virtualMachineScaleSetId.SubscriptionId, virtualMachineScaleSetId.ResourceGroup, virtualMachineScaleSetId.Name, d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualMachineScaleSetName, id.ExtensionName, "")
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(resp.Response) {
		return tf.ImportAsExistsError("azurerm_virtual_machine_scale_set_extension", *resp.ID)
	}

	settings := map[string]interface{}{}
	if settingsString := d.Get("settings").(string); settingsString != "" {
		s, err := pluginsdk.ExpandJsonFromString(settingsString)
		if err != nil {
			return fmt.Errorf("unable to parse `settings`: %s", err)
		}
		settings = s
	}

	provisionAfterExtensionsRaw := d.Get("provision_after_extensions").([]interface{})
	provisionAfterExtensions := utils.ExpandStringSlice(provisionAfterExtensionsRaw)

	props := compute.VirtualMachineScaleSetExtension{
		Name: utils.String(id.ExtensionName),
		VirtualMachineScaleSetExtensionProperties: &compute.VirtualMachineScaleSetExtensionProperties{
			Publisher:                     utils.String(d.Get("publisher").(string)),
			Type:                          utils.String(d.Get("type").(string)),
			TypeHandlerVersion:            utils.String(d.Get("type_handler_version").(string)),
			AutoUpgradeMinorVersion:       utils.Bool(d.Get("auto_upgrade_minor_version").(bool)),
			EnableAutomaticUpgrade:        utils.Bool(d.Get("automatic_upgrade_enabled").(bool)),
			SuppressFailures:              utils.Bool(d.Get("failure_suppression_enabled").(bool)),
			ProtectedSettingsFromKeyVault: expandProtectedSettingsFromKeyVault(d.Get("protected_settings_from_key_vault").([]interface{})),
			ProvisionAfterExtensions:      provisionAfterExtensions,
			Settings:                      settings,
		},
	}
	if v, ok := d.GetOk("force_update_tag"); ok {
		props.VirtualMachineScaleSetExtensionProperties.ForceUpdateTag = utils.String(v.(string))
	}

	if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
		protectedSettings, err := pluginsdk.ExpandJsonFromString(protectedSettingsString)
		if err != nil {
			return fmt.Errorf("unable to parse `protected_settings`: %s", err)
		}
		props.VirtualMachineScaleSetExtensionProperties.ProtectedSettings = &protectedSettings
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualMachineScaleSetName, id.ExtensionName, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualMachineScaleSetExtensionRead(d, meta)
}

func resourceVirtualMachineScaleSetExtensionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetExtensionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetExtensionID(d.Id())
	if err != nil {
		return err
	}

	props := compute.VirtualMachineScaleSetExtensionProperties{
		// if this isn't specified it defaults to false
		AutoUpgradeMinorVersion: utils.Bool(d.Get("auto_upgrade_minor_version").(bool)),
		EnableAutomaticUpgrade:  utils.Bool(d.Get("automatic_upgrade_enabled").(bool)),
	}

	if d.HasChange("failure_suppression_enabled") {
		props.SuppressFailures = utils.Bool(d.Get("failure_suppression_enabled").(bool))
	}

	if d.HasChange("force_update_tag") {
		props.ForceUpdateTag = utils.String(d.Get("force_update_tag").(string))
	}

	if d.HasChange("protected_settings") {
		protectedSettings := map[string]interface{}{}
		if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
			ps, err := pluginsdk.ExpandJsonFromString(protectedSettingsString)
			if err != nil {
				return fmt.Errorf("unable to parse `protected_settings`: %s", err)
			}
			protectedSettings = ps
		}

		props.ProtectedSettings = protectedSettings
	}

	if d.HasChange("protected_settings_from_key_vault") {
		props.ProtectedSettingsFromKeyVault = expandProtectedSettingsFromKeyVault(d.Get("protected_settings_from_key_vault").([]interface{}))
	}

	if d.HasChange("provision_after_extensions") {
		provisionAfterExtensionsRaw := d.Get("provision_after_extensions").([]interface{})
		props.ProvisionAfterExtensions = utils.ExpandStringSlice(provisionAfterExtensionsRaw)
	}

	if d.HasChange("publisher") {
		props.Publisher = utils.String(d.Get("publisher").(string))
	}

	if d.HasChange("settings") {
		settings := map[string]interface{}{}

		if settingsString := d.Get("settings").(string); settingsString != "" {
			s, err := pluginsdk.ExpandJsonFromString(settingsString)
			if err != nil {
				return fmt.Errorf("unable to parse `settings`: %s", err)
			}
			settings = s
		}

		props.Settings = settings
	}

	if d.HasChange("type") {
		props.Type = utils.String(d.Get("type").(string))
	}

	if d.HasChange("type_handler_version") {
		props.TypeHandlerVersion = utils.String(d.Get("type_handler_version").(string))
	}

	extension := compute.VirtualMachineScaleSetExtension{
		Name: utils.String(id.ExtensionName),
		VirtualMachineScaleSetExtensionProperties: &props,
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualMachineScaleSetName, id.ExtensionName, extension)
	if err != nil {
		return fmt.Errorf("updating Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", id.ExtensionName, id.VirtualMachineScaleSetName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", id.ExtensionName, id.VirtualMachineScaleSetName, id.ResourceGroup, err)
	}

	return resourceVirtualMachineScaleSetExtensionRead(d, meta)
}

func resourceVirtualMachineScaleSetExtensionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	vmssClient := meta.(*clients.Client).Compute.VMScaleSetClient
	client := meta.(*clients.Client).Compute.VMScaleSetExtensionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetExtensionID(d.Id())
	if err != nil {
		return err
	}

	// Upgrading to the 2021-07-01 exposed a new expand parameter in the GET method
	vmss, err := vmssClient.Get(ctx, id.ResourceGroup, id.VirtualMachineScaleSetName, "")
	if err != nil {
		if utils.ResponseWasNotFound(vmss.Response) {
			log.Printf("Virtual Machine Scale Set %q was not found in Resource Group %q - removing Extension from state!", id.VirtualMachineScaleSetName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Virtual Machine Scale Set %q (Resource Group %q): %+v", id.VirtualMachineScaleSetName, id.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualMachineScaleSetName, id.ExtensionName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("Extension %q (Virtual Machine Scale Set %q / Resource Group %q) was not found - removing from state!", id.ExtensionName, id.VirtualMachineScaleSetName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", id.ExtensionName, id.VirtualMachineScaleSetName, id.ResourceGroup, err)
	}

	d.Set("name", id.ExtensionName)
	d.Set("virtual_machine_scale_set_id", vmss.ID)

	if props := resp.VirtualMachineScaleSetExtensionProperties; props != nil {
		d.Set("auto_upgrade_minor_version", props.AutoUpgradeMinorVersion)
		d.Set("automatic_upgrade_enabled", props.EnableAutomaticUpgrade)
		d.Set("force_update_tag", props.ForceUpdateTag)
		d.Set("protected_settings_from_key_vault", flattenProtectedSettingsFromKeyVault(props.ProtectedSettingsFromKeyVault))
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
			settingsVal, ok := props.Settings.(map[string]interface{})
			if ok {
				settingsJson, err := pluginsdk.FlattenJsonToString(settingsVal)
				if err != nil {
					return fmt.Errorf("unable to parse settings from response: %s", err)
				}
				settings = settingsJson
			}
		}
		d.Set("settings", settings)
	}

	return nil
}

func resourceVirtualMachineScaleSetExtensionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMScaleSetExtensionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineScaleSetExtensionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.VirtualMachineScaleSetName, id.ExtensionName)
	if err != nil {
		return fmt.Errorf("deleting Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", id.ExtensionName, id.VirtualMachineScaleSetName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Extension %q (Virtual Machine Scale Set %q / Resource Group %q): %+v", id.ExtensionName, id.VirtualMachineScaleSetName, id.ResourceGroup, err)
	}

	return nil
}
