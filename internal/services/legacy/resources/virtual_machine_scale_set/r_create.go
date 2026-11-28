// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/legacy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineScaleSetCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualmachinescalesets.NewVirtualMachineScaleSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, virtualmachinescalesets.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_virtual_machine_scale_set", id.ID())
			}
		}
	}

	storageProfile := virtualmachinescalesets.VirtualMachineScaleSetStorageProfile{}
	osDisk, err := expandAzureRMVirtualMachineScaleSetsStorageProfileOsDisk(d)
	if err != nil {
		return err
	}
	storageProfile.OsDisk = osDisk

	if _, ok := d.GetOk("storage_profile_data_disk"); ok {
		storageProfile.DataDisks = expandAzureRMVirtualMachineScaleSetsStorageProfileDataDisk(d)
	}

	if _, ok := d.GetOk("storage_profile_image_reference"); ok {
		imageRef, err2 := expandAzureRmVirtualMachineScaleSetStorageProfileImageReference(d)
		if err2 != nil {
			return err2
		}
		storageProfile.ImageReference = imageRef
	}

	osProfile := expandAzureRMVirtualMachineScaleSetsOsProfile(d)

	extensions, err := expandAzureRMVirtualMachineScaleSetExtensions(d)
	if err != nil {
		return err
	}

	upgradePolicy := d.Get("upgrade_policy_mode").(string)
	automaticOsUpgrade := d.Get("automatic_os_upgrade").(bool)
	overprovision := d.Get("overprovision").(bool)
	singlePlacementGroup := d.Get("single_placement_group").(bool)
	priority := d.Get("priority").(string)
	evictionPolicy := d.Get("eviction_policy").(string)

	scaleSetProps := virtualmachinescalesets.VirtualMachineScaleSetProperties{
		UpgradePolicy: &virtualmachinescalesets.UpgradePolicy{
			Mode: pointer.ToEnum[virtualmachinescalesets.UpgradeMode](upgradePolicy),
			AutomaticOSUpgradePolicy: &virtualmachinescalesets.AutomaticOSUpgradePolicy{
				EnableAutomaticOSUpgrade: pointer.To(automaticOsUpgrade),
			},
			RollingUpgradePolicy: expandAzureRmRollingUpgradePolicy(d),
		},
		VirtualMachineProfile: &virtualmachinescalesets.VirtualMachineScaleSetVMProfile{
			NetworkProfile:   expandAzureRmVirtualMachineScaleSetNetworkProfile(d),
			StorageProfile:   &storageProfile,
			OsProfile:        osProfile,
			ExtensionProfile: extensions,
			Priority:         pointer.ToEnum[virtualmachinescalesets.VirtualMachinePriorityTypes](priority),
		},
		// OrchestrationMode needs to be hardcoded to Uniform, for the
		// standard VMSS resource, since virtualMachineProfile is now supported
		// in both VMSS and Orchestrated VMSS...
		OrchestrationMode:    pointer.To(virtualmachinescalesets.OrchestrationModeUniform),
		Overprovision:        &overprovision,
		SinglePlacementGroup: &singlePlacementGroup,
	}

	if strings.EqualFold(priority, string(virtualmachinescalesets.VirtualMachinePriorityTypesLow)) {
		scaleSetProps.VirtualMachineProfile.EvictionPolicy = pointer.ToEnum[virtualmachinescalesets.VirtualMachineEvictionPolicyTypes](evictionPolicy)
	}

	if _, ok := d.GetOk("boot_diagnostics"); ok {
		diagnosticProfile := expandAzureRMVirtualMachineScaleSetsDiagnosticProfile(d)
		scaleSetProps.VirtualMachineProfile.DiagnosticsProfile = &diagnosticProfile
	}

	if v, ok := d.GetOk("health_probe_id"); ok {
		scaleSetProps.VirtualMachineProfile.NetworkProfile.HealthProbe = &virtualmachinescalesets.ApiEntityReference{
			Id: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		scaleSetProps.ProximityPlacementGroup = &virtualmachinescalesets.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	payload := virtualmachinescalesets.VirtualMachineScaleSet{
		Name:       &id.VirtualMachineScaleSetName,
		Location:   location.Normalize(d.Get("location").(string)),
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
		Sku:        expandVirtualMachineScaleSetSku(d),
		Properties: &scaleSetProps,
		Zones:      helpers.ExpandZones(d.Get("zones").([]any)),
	}

	if _, ok := d.GetOk("identity"); ok {
		identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = identityExpanded
	}

	if v, ok := d.GetOk("license_type"); ok {
		payload.Properties.VirtualMachineProfile.LicenseType = pointer.To(v.(string))
	}

	if _, ok := d.GetOk("plan"); ok {
		payload.Plan = expandAzureRmVirtualMachineScaleSetPlan(d)
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, virtualmachinescalesets.DefaultCreateOrUpdateOperationOptions(), sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, payload, virtualmachinescalesets.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceVirtualMachineScaleSetRead(d, meta)
}
