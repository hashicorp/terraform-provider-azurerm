// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_extension

import (
	"context"
	"fmt"
	"log"
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
// validation. Returns skip=true if the VM does not exist and no location fallback is configured.
func resolvePreflightVMLocation(ctx context.Context, client *clients.Client, vmId virtualmachines.VirtualMachineId) (loc string, skip bool, err error) {
	vm, err := client.Compute.VirtualMachinesClient.Get(ctx, vmId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(vm.HttpResponse) {
			return "", false, fmt.Errorf("retrieving %s for preflight validation: %+v", vmId, err)
		}

		if fallback := client.Features.EnhancedValidation.LocationFallback; fallback != nil {
			log.Printf("[DEBUG] %s was not found, using the configured location fallback for preflight validation", vmId)
			return *fallback, false, nil
		}

		log.Printf("[DEBUG] skipping preflight validation: %s was not found and no location fallback is configured", vmId)
		return "", true, nil
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

	if !d.NewValueKnown("virtual_machine_id") || !d.NewValueKnown("name") {
		log.Printf("[DEBUG] skipping preflight validation for `azurerm_virtual_machine_extension`: `virtual_machine_id` or `name` is not known until apply")
		return nil
	}

	// Untyped CustomizeDiff contexts have no deadline, which the SDK clients require for polling.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	vmId, err := virtualmachines.ParseVirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `virtual_machine_id` for preflight validation: %w", err)
	}

	loc, skip, err := resolvePreflightVMLocation(ctx, client, *vmId)
	if err != nil {
		return err
	}
	if skip {
		return nil
	}

	id := virtualmachineextensions.NewExtensionID(vmId.SubscriptionId, vmId.ResourceGroupName, vmId.VirtualMachineName, d.Get("name").(string))

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
