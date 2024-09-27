// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineextensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualMachineExtension() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineExtensionsCreateUpdate,
		Read:   resourceVirtualMachineExtensionsRead,
		Update: resourceVirtualMachineExtensionsCreateUpdate,
		Delete: resourceVirtualMachineExtensionsDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualmachineextensions.ParseExtensionID(id)
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
				ValidateFunc: commonids.ValidateVirtualMachineID,
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

			"provision_after_extensions": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceVirtualMachineExtensionsCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	publisher := d.Get("publisher").(string)
	extensionType := d.Get("type").(string)
	typeHandlerVersion := d.Get("type_handler_version").(string)
	autoUpgradeMinor := d.Get("auto_upgrade_minor_version").(bool)
	enableAutomaticUpgrade := d.Get("automatic_upgrade_enabled").(bool)
	suppressFailure := d.Get("failure_suppression_enabled").(bool)
	t := d.Get("tags").(map[string]interface{})

	extension := virtualmachineextensions.VirtualMachineExtension{
		Location: &location,
		Properties: &virtualmachineextensions.VirtualMachineExtensionProperties{
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
		var result interface{}
		err := json.Unmarshal([]byte(settingsString), &result)
		if err != nil {
			return fmt.Errorf("unmarshaling `settings`: %+v", err)
		}
		extension.Properties.Settings = pointer.To(result)
	}

	if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
		var result interface{}
		err := json.Unmarshal([]byte(protectedSettingsString), &result)
		if err != nil {
			return fmt.Errorf("unmarshaling `protected_settings`: %+v", err)
		}
		extension.Properties.ProtectedSettings = pointer.To(result)
	}

	if provisionAfterExtensionsValue, exists := d.GetOk("provision_after_extensions"); exists {
		extension.Properties.ProvisionAfterExtensions = utils.ExpandStringSlice(provisionAfterExtensionsValue.([]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, extension); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualMachineExtensionsRead(d, meta)
}

func resourceVirtualMachineExtensionsRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			d.Set("protected_settings_from_key_vault", flattenProtectedSettingsFromKeyVault(props.ProtectedSettingsFromKeyVault))
			d.Set("provision_after_extensions", pointer.From(props.ProvisionAfterExtensions))

			suppressFailure := false
			if props.SuppressFailures != nil {
				suppressFailure = *props.SuppressFailures
			}
			d.Set("failure_suppression_enabled", suppressFailure)

			if props.Settings != nil {
				settings, err := json.Marshal(props.Settings)
				if err != nil {
					return fmt.Errorf("unmarshaling `settings`: %+v", err)
				}
				d.Set("settings", string(settings))
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceVirtualMachineExtensionsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineExtensionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachineextensions.ParseExtensionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
