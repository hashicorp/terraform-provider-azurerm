// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_extension

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineextensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Properties struct {
	Output Output `tfschema:"output"`
}

type Output struct {
	AccessSAS string `tfschema:"accessSAS"`
}

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name virtual_machine_extension -service-package-name compute -properties "name" -compare-values "virtual_machine_name:virtual_machine_id,resource_group_name:virtual_machine_id,subscription_id:virtual_machine_id"

func RegisterResource() *pluginsdk.Resource {
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

		Schema: virtualMachineExtensionSchema(),
	}
}

// resolvePreflightVMLocation looks up the location of the parent VM for use in preflight
// validation. Returns skip=true if the virtual_machine_id is not yet known or the VM does not
// exist and no location fallback is configured.
func resolvePreflightVMLocation(ctx context.Context, client *clients.Client, d *schema.ResourceDiff) (loc string, skip bool, err error) {
	vmIdRaw := d.Get("virtual_machine_id").(string)
	if vmIdRaw == "" {
		return "", true, nil
	}

	vmId, err := virtualmachines.ParseVirtualMachineID(vmIdRaw)
	if err != nil {
		return "", false, fmt.Errorf("parsing virtual_machine_id for preflight validation: %w", err)
	}

	vm, err := client.Compute.VirtualMachinesClient.Get(ctx, *vmId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(vm.HttpResponse) {
			if fallback := client.Features.EnhancedValidation.LocationFallback; fallback != nil {
				return *fallback, false, nil
			}

			return "", true, nil
		}

		return "", false, fmt.Errorf("retrieving %s for preflight validation: %+v", vmId, err)
	}

	if vm.Model == nil || vm.Model.Location == "" {
		return "", false, fmt.Errorf("determining location of %s for preflight validation: location was missing", vmId)
	}

	return vm.Model.Location, false, nil
}

// resourceVirtualMachineExtensionCustomizeDiff implements preflight validation for
// azurerm_virtual_machine_extension. The resource uses the same complete PUT payload for
// create and update, so validation can run for both operations.
func resourceVirtualMachineExtensionCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	client := meta.(*clients.Client)

	if d == nil || !client.Features.EnhancedValidation.PreflightEnabled {
		return nil
	}

	if len(d.GetChangedKeysPrefix("")) == 0 && d.Id() != "" {
		return nil
	}

	loc, skip, err := resolvePreflightVMLocation(ctx, client, d)
	if err != nil {
		return err
	}
	if skip {
		return nil
	}

	vmIdRaw := d.Get("virtual_machine_id").(string)
	vmId, err := virtualmachines.ParseVirtualMachineID(vmIdRaw)
	if err != nil {
		return fmt.Errorf("parsing virtual_machine_id for preflight validation: %w", err)
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
