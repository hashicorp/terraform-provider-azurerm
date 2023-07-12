// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func resourceVirtualMachineExtension() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineExtensionsCreateUpdate,
		Read:   resourceVirtualMachineExtensionsRead,
		Update: resourceVirtualMachineExtensionsCreateUpdate,
		Delete: resourceVirtualMachineExtensionsDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualMachineExtensionID(id)
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

			"virtual_machine_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VirtualMachineID,
			},

			"publisher": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"type_handler_version": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"auto_upgrade_minor_version": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
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

			"settings": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			// due to the sensitive nature, these are not returned by the API
			"protected_settings": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Sensitive:        true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				ConflictsWith:    []string{"protected_settings_from_key_vault"},
			},

			"protected_settings_from_key_vault": protectedSettingsFromKeyVaultSchema(true),

			"tags": tags.Schema(),
		},
	}
}

func resourceVirtualMachineExtensionsCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	vmExtensionClient := meta.(*clients.Client).Compute.VMExtensionClient
	vmClient := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachineId, err := parse.VirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Virtual Machine ID %q: %+v", virtualMachineId, err)
	}
	id := parse.NewVirtualMachineExtensionID(virtualMachineId.SubscriptionId, virtualMachineId.ResourceGroup, virtualMachineId.Name, d.Get("name").(string))

	virtualMachine, err := vmClient.Get(ctx, id.ResourceGroup, id.VirtualMachineName, "")
	if err != nil {
		return fmt.Errorf("getting %s: %+v", virtualMachineId, err)
	}

	location := *virtualMachine.Location
	if location == "" {
		return fmt.Errorf("reading location of %s", virtualMachineId)
	}

	if d.IsNewResource() {
		existing, err := vmExtensionClient.Get(ctx, id.ResourceGroup, id.VirtualMachineName, id.ExtensionName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_virtual_machine_extension", id.ID())
		}
	}

	publisher := d.Get("publisher").(string)
	extensionType := d.Get("type").(string)
	typeHandlerVersion := d.Get("type_handler_version").(string)
	autoUpgradeMinor := d.Get("auto_upgrade_minor_version").(bool)
	enableAutomaticUpgrade := d.Get("automatic_upgrade_enabled").(bool)
	suppressFailure := d.Get("failure_suppression_enabled").(bool)
	t := d.Get("tags").(map[string]interface{})

	extension := compute.VirtualMachineExtension{
		Location: &location,
		VirtualMachineExtensionProperties: &compute.VirtualMachineExtensionProperties{
			Publisher:                     &publisher,
			Type:                          &extensionType,
			TypeHandlerVersion:            &typeHandlerVersion,
			AutoUpgradeMinorVersion:       &autoUpgradeMinor,
			EnableAutomaticUpgrade:        &enableAutomaticUpgrade,
			ProtectedSettingsFromKeyVault: expandProtectedSettingsFromKeyVault(d.Get("protected_settings_from_key_vault").([]interface{})),
			SuppressFailures:              &suppressFailure,
		},
		Tags: tags.Expand(t),
	}

	if settingsString := d.Get("settings").(string); settingsString != "" {
		settings, err := pluginsdk.ExpandJsonFromString(settingsString)
		if err != nil {
			return fmt.Errorf("unable to parse settings: %s", err)
		}
		extension.VirtualMachineExtensionProperties.Settings = &settings
	}

	if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
		protectedSettings, err := pluginsdk.ExpandJsonFromString(protectedSettingsString)
		if err != nil {
			return fmt.Errorf("unable to parse protected_settings: %s", err)
		}
		extension.VirtualMachineExtensionProperties.ProtectedSettings = &protectedSettings
	}

	future, err := vmExtensionClient.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualMachineName, id.ExtensionName, extension)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, vmExtensionClient.Client); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceVirtualMachineExtensionsRead(d, meta)
}

func resourceVirtualMachineExtensionsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	vmExtensionClient := meta.(*clients.Client).Compute.VMExtensionClient
	vmClient := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineExtensionID(d.Id())
	if err != nil {
		return err
	}

	virtualMachine, err := vmClient.Get(ctx, id.ResourceGroup, id.VirtualMachineName, "")
	if err != nil {
		if utils.ResponseWasNotFound(virtualMachine.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Virtual Machine %s: %s", id.ExtensionName, err)
	}

	d.Set("virtual_machine_id", virtualMachine.ID)

	resp, err := vmExtensionClient.Get(ctx, id.ResourceGroup, id.VirtualMachineName, id.ExtensionName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Virtual Machine Extension %s: %s", id.ExtensionName, err)
	}

	d.Set("name", resp.Name)

	if props := resp.VirtualMachineExtensionProperties; props != nil {
		d.Set("publisher", props.Publisher)
		d.Set("type", props.Type)
		d.Set("type_handler_version", props.TypeHandlerVersion)
		d.Set("auto_upgrade_minor_version", props.AutoUpgradeMinorVersion)
		d.Set("automatic_upgrade_enabled", props.EnableAutomaticUpgrade)
		d.Set("protected_settings_from_key_vault", flattenProtectedSettingsFromKeyVault(props.ProtectedSettingsFromKeyVault))

		suppressFailure := false
		if props.SuppressFailures != nil {
			suppressFailure = *props.SuppressFailures
		}
		d.Set("failure_suppression_enabled", suppressFailure)

		if settings := props.Settings; settings != nil {
			settingsVal := settings.(map[string]interface{})
			settingsJson, err := pluginsdk.FlattenJsonToString(settingsVal)
			if err != nil {
				return fmt.Errorf("unable to parse settings from response: %s", err)
			}
			d.Set("settings", settingsJson)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVirtualMachineExtensionsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMExtensionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualMachineExtensionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.VirtualMachineName, id.ExtensionName)
	if err != nil {
		return err
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}
