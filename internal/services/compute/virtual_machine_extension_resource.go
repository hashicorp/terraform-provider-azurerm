// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name virtual_machine_extension -service-package-name compute -properties "name" -compare-values "virtual_machine_name:virtual_machine_id,resource_group_name:virtual_machine_id,subscription_id:virtual_machine_id"

func resourceVirtualMachineExtension() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceVirtualMachineExtensionsCreateUpdate,
		Read:     resourceVirtualMachineExtensionsRead,
		Update:   resourceVirtualMachineExtensionsCreateUpdate,
		Delete:   resourceVirtualMachineExtensionsDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&virtualmachineextensions.ExtensionId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&virtualmachineextensions.ExtensionId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(resourceVirtualMachineExtensionCustomizeDiff),

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
			ProtectedSettingsFromKeyVault: expandProtectedSettingsFromKeyVault(d.Get("protected_settings_from_key_vault").([]interface{})),
			SuppressFailures:              &suppressFailure,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if settingsString := d.Get("settings").(string); settingsString != "" {
		var result interface{}
		if err := json.Unmarshal([]byte(settingsString), &result); err != nil {
			return extension, fmt.Errorf("unmarshaling `settings`: %+v", err)
		}
		extension.Properties.Settings = pointer.To(result)
	}

	if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
		var result interface{}
		if err := json.Unmarshal([]byte(protectedSettingsString), &result); err != nil {
			return extension, fmt.Errorf("unmarshaling `protected_settings`: %+v", err)
		}
		extension.Properties.ProtectedSettings = pointer.To(result)
	}

	if provisionAfterExtensionsValue, exists := d.GetOk("provision_after_extensions"); exists {
		extension.Properties.ProvisionAfterExtensions = helpers.ExpandStringSlice(provisionAfterExtensionsValue.([]interface{}))
	}

	return extension, nil
}

// resolvePreflightVMLocation looks up the location of the parent VM for use in preflight
// validation. Returns skip=true if the virtual_machine_id is not yet known (cross-resource
// computed value) or if the VM cannot be found, so that validation is gracefully skipped
// rather than failing the plan.
func resolvePreflightVMLocation(ctx context.Context, client *clients.Client, d *schema.ResourceDiff) (loc string, skip bool) {
	vmIdRaw := d.Get("virtual_machine_id").(string)
	if vmIdRaw == "" {
		// virtual_machine_id is (known after apply) — skip gracefully.
		return "", true
	}

	vmId, err := virtualmachines.ParseVirtualMachineID(vmIdRaw)
	if err != nil {
		return "", true
	}

	vm, err := client.Compute.VirtualMachinesClient.Get(ctx, *vmId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(vm.HttpResponse) {
			// VM doesn't exist yet — gracefully skip.
			return "", true
		}
		// For any other error also skip rather than failing plan.
		return "", true
	}

	if vm.Model == nil || vm.Model.Location == "" {
		return "", true
	}

	return vm.Model.Location, false
}

// resourceVirtualMachineExtensionCustomizeDiff implements preflight validation for
// azurerm_virtual_machine_extension. It uses Pattern 2 (create and ForceNew replacement
// only) because in-place updates do not change the extension type or its parent VM, so
// running the full ARM PUT payload against the preflight API on every update would produce
// false positives from immutable fields (publisher, virtual_machine_id) that cannot change
// without a ForceNew destroy+create.
func resourceVirtualMachineExtensionCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	client := meta.(*clients.Client)

	if !client.Features.EnhancedValidation.PreflightEnabled {
		return nil
	}

	isNewResource := d.Id() == ""
	isForceNewReplacement := false

	if !isNewResource {
		forceNewKeys := []string{"name", "virtual_machine_id", "publisher"}
		for _, key := range forceNewKeys {
			if d.HasChange(key) {
				isForceNewReplacement = true
				break
			}
		}
	}

	if !isNewResource && !isForceNewReplacement {
		return nil
	}

	loc, skip := resolvePreflightVMLocation(ctx, client, d)
	if skip {
		return nil
	}

	vmIdRaw := d.Get("virtual_machine_id").(string)
	vmId, err := virtualmachines.ParseVirtualMachineID(vmIdRaw)
	if err != nil {
		return nil
	}

	extensionName := d.Get("name").(string)
	id := virtualmachineextensions.NewExtensionID(vmId.SubscriptionId, vmId.ResourceGroupName, vmId.VirtualMachineName, extensionName)

	req, err := expandCreateForVirtualMachineExtension(d, loc)
	if err != nil {
		return err
	}

	preflightValidate, err := preflight.NewValidationRequest(pointer.To(loc), pointer.To(id), "2024-03-01", req)
	if err != nil {
		return fmt.Errorf("constructing preflight validation request: %w", err)
	}

	metadata := sdk.ResourceMetaData{
		Client:       client,
		ResourceDiff: d,
		Logger:       sdk.ConsoleLogger{},
	}

	return preflightValidate.ValidateResource(ctx, metadata)
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
		if err := json.Unmarshal([]byte(settingsString), &result); err != nil {
			return fmt.Errorf("unmarshaling `settings`: %+v", err)
		}
		extension.Properties.Settings = pointer.To(result)
	}

	if protectedSettingsString := d.Get("protected_settings").(string); protectedSettingsString != "" {
		var result interface{}
		if err := json.Unmarshal([]byte(protectedSettingsString), &result); err != nil {
			return fmt.Errorf("unmarshaling `protected_settings`: %+v", err)
		}
		extension.Properties.ProtectedSettings = pointer.To(result)
	}

	if provisionAfterExtensionsValue, exists := d.GetOk("provision_after_extensions"); exists {
		extension.Properties.ProvisionAfterExtensions = helpers.ExpandStringSlice(provisionAfterExtensionsValue.([]interface{}))
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
