// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package linux_virtual_machine_scale_set

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLinuxVirtualMachineScaleSetUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	updateInstances := false
	options := virtualmachinescalesets.DefaultGetOperationOptions()
	options.Expand = pointer.To(virtualmachinescalesets.ExpandTypesForGetVMScaleSetsUserData)
	existing, err := client.Get(ctx, *id, options)
	if err != nil {
		return fmt.Errorf("retrieving Linux %s: %+v", id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving Linux %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving Linux %s: `properties` was nil", id)
	}
	if existing.Model.Properties.VirtualMachineProfile == nil {
		return fmt.Errorf("retrieving Linux %s: `properties.virtualMachineProfile` was nil", id)
	}
	if existing.Model.Properties.VirtualMachineProfile.StorageProfile == nil {
		return fmt.Errorf("retrieving Linux %s: `properties.virtualMachineProfile,storageProfile` was nil", id)
	}

	updateProps := virtualmachinescalesets.VirtualMachineScaleSetUpdateProperties{
		VirtualMachineProfile: &virtualmachinescalesets.VirtualMachineScaleSetUpdateVMProfile{
			// if an image reference has been configured previously (it has to be), we would better to include that in this
			// update request to avoid some circumstances that the API will complain ImageReference is null
			// issue tracking: https://github.com/Azure/azure-rest-api-specs/issues/10322
			StorageProfile: &virtualmachinescalesets.VirtualMachineScaleSetUpdateStorageProfile{
				ImageReference: existing.Model.Properties.VirtualMachineProfile.StorageProfile.ImageReference,
			},
		},
		// if an upgrade policy's been configured previously (which it will have) it must be threaded through
		// this doesn't matter for Manual - but breaks when updating anything on a Automatic and Rolling Mode Scale Set
		UpgradePolicy: existing.Model.Properties.UpgradePolicy,
	}
	update := virtualmachinescalesets.VirtualMachineScaleSetUpdate{}

	// first try and pull this from existing vm, which covers no changes being made to this block
	automaticOSUpgradeIsEnabled := false
	if policy := existing.Model.Properties.UpgradePolicy; policy != nil {
		if policy.AutomaticOSUpgradePolicy != nil && policy.AutomaticOSUpgradePolicy.EnableAutomaticOSUpgrade != nil {
			automaticOSUpgradeIsEnabled = *policy.AutomaticOSUpgradePolicy.EnableAutomaticOSUpgrade
		}
	}

	if d.HasChange("zones") {
		update.Zones = pointer.To(zones.ExpandUntyped(d.Get("zones").(*schema.Set).List()))
	}

	if d.HasChanges("automatic_os_upgrade_policy", "rolling_upgrade_policy") {
		upgradePolicy := virtualmachinescalesets.UpgradePolicy{}
		if existing.Model.Properties.UpgradePolicy == nil {
			upgradePolicy = virtualmachinescalesets.UpgradePolicy{
				Mode: pointer.ToEnum[virtualmachinescalesets.UpgradeMode](d.Get("upgrade_mode").(string)),
			}
		} else {
			upgradePolicy = *existing.Model.Properties.UpgradePolicy
			upgradePolicy.Mode = pointer.ToEnum[virtualmachinescalesets.UpgradeMode](d.Get("upgrade_mode").(string))
		}

		if d.HasChange("automatic_os_upgrade_policy") {
			automaticRaw := d.Get("automatic_os_upgrade_policy").([]any)
			rollingRaw := d.Get("rolling_upgrade_policy").([]any)
			useRollingUpgradePolicy := len(rollingRaw) > 0
			upgradePolicy.AutomaticOSUpgradePolicy = helpers.ExpandVirtualMachineScaleSetAutomaticUpgradePolicy(automaticRaw, useRollingUpgradePolicy)

			if upgradePolicy.AutomaticOSUpgradePolicy != nil {
				automaticOSUpgradeIsEnabled = *upgradePolicy.AutomaticOSUpgradePolicy.EnableAutomaticOSUpgrade
			}
		}

		if d.HasChange("rolling_upgrade_policy") {
			rollingRaw := d.Get("rolling_upgrade_policy").([]any)
			zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
			rollingUpgradePolicy, err := helpers.ExpandVirtualMachineScaleSetRollingUpgradePolicy(rollingRaw, len(zones) > 0, d.Get("overprovision").(bool))
			if err != nil {
				return err
			}
			upgradePolicy.RollingUpgradePolicy = rollingUpgradePolicy
		}

		updateProps.UpgradePolicy = &upgradePolicy
	}

	priority := virtualmachinescalesets.VirtualMachinePriorityTypes(d.Get("priority").(string))
	if d.HasChange("max_bid_price") {
		if priority != virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `Spot`")
		}

		updateProps.VirtualMachineProfile.BillingProfile = &virtualmachinescalesets.BillingProfile{
			MaxPrice: pointer.To(d.Get("max_bid_price").(float64)),
		}
	}

	if d.HasChange("single_placement_group") {
		singlePlacementGroup := d.Get("single_placement_group").(bool)
		if singlePlacementGroup {
			return fmt.Errorf("%q can not be set to %q once it has been set to %q", "single_placement_group", "true", "false")
		}
		updateProps.SinglePlacementGroup = pointer.To(singlePlacementGroup)
	}

	// lintignore:R019 // deliberate subset: only the fields feeding the OSProfile update payload
	if d.HasChanges("admin_ssh_key", "custom_data", "disable_password_authentication", "provision_vm_agent", "secret") {
		osProfile := virtualmachinescalesets.VirtualMachineScaleSetUpdateOSProfile{}

		if d.HasChanges("admin_ssh_key", "disable_password_authentication", "provision_vm_agent") {
			linuxConfig := virtualmachinescalesets.LinuxConfiguration{}

			if d.HasChange("admin_ssh_key") {
				sshKeysRaw := d.Get("admin_ssh_key").(*pluginsdk.Set).List()
				sshKeys := helpers.ExpandSSHKeysVMSS(sshKeysRaw)
				linuxConfig.Ssh = &virtualmachinescalesets.SshConfiguration{
					PublicKeys: &sshKeys,
				}
			}

			if d.HasChange("disable_password_authentication") {
				linuxConfig.DisablePasswordAuthentication = pointer.To(d.Get("disable_password_authentication").(bool))
			}

			if d.HasChange("provision_vm_agent") {
				linuxConfig.ProvisionVMAgent = pointer.To(d.Get("provision_vm_agent").(bool))
			}

			osProfile.LinuxConfiguration = &linuxConfig
		}

		if d.HasChange("custom_data") {
			updateInstances = true

			// customData can only be sent if it's a base64 encoded string,
			// so it's not possible to remove this without tainting the resource
			if v, ok := d.GetOk("custom_data"); ok {
				osProfile.CustomData = pointer.To(v.(string))
			}
		}

		if d.HasChange("secret") {
			secretsRaw := d.Get("secret").([]any)
			osProfile.Secrets = helpers.ExpandLinuxSecretsVMSS(secretsRaw)
		}

		updateProps.VirtualMachineProfile.OsProfile = &osProfile
	}

	if d.HasChanges("data_disk", "os_disk", "source_image_id", "source_image_reference") {
		updateInstances = true

		if updateProps.VirtualMachineProfile.StorageProfile == nil {
			updateProps.VirtualMachineProfile.StorageProfile = &virtualmachinescalesets.VirtualMachineScaleSetUpdateStorageProfile{}
		}

		if d.HasChange("data_disk") {
			ultraSSDEnabled := d.Get("additional_capabilities.0.ultra_ssd_enabled").(bool)
			dataDisks, err := helpers.ExpandVirtualMachineScaleSetDataDisk(d.Get("data_disk").([]any), ultraSSDEnabled)
			if err != nil {
				return fmt.Errorf("expanding `data_disk`: %+v", err)
			}
			updateProps.VirtualMachineProfile.StorageProfile.DataDisks = dataDisks
		}

		if d.HasChange("os_disk") {
			osDiskRaw := d.Get("os_disk").([]any)
			updateProps.VirtualMachineProfile.StorageProfile.OsDisk = helpers.ExpandVirtualMachineScaleSetOSDiskUpdate(osDiskRaw)
		}

		if d.HasChanges("source_image_id", "source_image_reference") {
			sourceImageReferenceRaw := d.Get("source_image_reference").([]any)
			sourceImageId := d.Get("source_image_id").(string)

			// Must include all storage profile properties when updating disk image.  See: https://github.com/hashicorp/terraform-provider-azurerm/issues/8273
			updateProps.VirtualMachineProfile.StorageProfile.DataDisks = existing.Model.Properties.VirtualMachineProfile.StorageProfile.DataDisks
			updateProps.VirtualMachineProfile.StorageProfile.ImageReference = helpers.ExpandSourceImageReferenceVMSS(sourceImageReferenceRaw, sourceImageId)
			updateProps.VirtualMachineProfile.StorageProfile.OsDisk = &virtualmachinescalesets.VirtualMachineScaleSetUpdateOSDisk{
				Caching:                 existing.Model.Properties.VirtualMachineProfile.StorageProfile.OsDisk.Caching,
				WriteAcceleratorEnabled: existing.Model.Properties.VirtualMachineProfile.StorageProfile.OsDisk.WriteAcceleratorEnabled,
				DiskSizeGB:              existing.Model.Properties.VirtualMachineProfile.StorageProfile.OsDisk.DiskSizeGB,
				Image:                   existing.Model.Properties.VirtualMachineProfile.StorageProfile.OsDisk.Image,
				VhdContainers:           existing.Model.Properties.VirtualMachineProfile.StorageProfile.OsDisk.VhdContainers,
				ManagedDisk:             existing.Model.Properties.VirtualMachineProfile.StorageProfile.OsDisk.ManagedDisk,
			}
		}
	}

	if d.HasChanges("network_interface", "health_probe_id") {
		networkInterfacesRaw := d.Get("network_interface").([]any)
		networkInterfaces, err := helpers.ExpandVirtualMachineScaleSetNetworkInterfaceUpdate(networkInterfacesRaw)
		if err != nil {
			return fmt.Errorf("expanding `network_interface`: %+v", err)
		}

		updateProps.VirtualMachineProfile.NetworkProfile = &virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkProfile{
			NetworkInterfaceConfigurations: networkInterfaces,
		}

		healthProbeId := d.Get("health_probe_id").(string)
		if healthProbeId != "" {
			updateProps.VirtualMachineProfile.NetworkProfile.HealthProbe = &virtualmachinescalesets.ApiEntityReference{
				Id: pointer.To(healthProbeId),
			}
		}
	}

	if d.HasChange("boot_diagnostics") {
		updateInstances = true

		bootDiagnosticsRaw := d.Get("boot_diagnostics").([]any)
		updateProps.VirtualMachineProfile.DiagnosticsProfile = helpers.ExpandBootDiagnosticsVMSS(bootDiagnosticsRaw)
	}

	if d.HasChange("do_not_run_extensions_on_overprovisioned_machines") {
		v := d.Get("do_not_run_extensions_on_overprovisioned_machines").(bool)
		updateProps.DoNotRunExtensionsOnOverprovisionedVMs = pointer.To(v)
	}

	if d.HasChange("overprovision") {
		v := d.Get("overprovision").(bool)
		updateProps.Overprovision = pointer.To(v)
	}

	if d.HasChange("scale_in") {
		if updateScaleInPolicy := helpers.ExpandVirtualMachineScaleSetScaleInPolicy(d.Get("scale_in").([]any)); updateScaleInPolicy != nil {
			updateProps.ScaleInPolicy = updateScaleInPolicy
		}
	}

	if d.HasChanges("resilient_vm_creation_enabled", "resilient_vm_deletion_enabled") {
		resilientVMCreationEnabled := d.Get("resilient_vm_creation_enabled").(bool)
		resilientVMDeletionEnabled := d.Get("resilient_vm_deletion_enabled").(bool)
		updateProps.ResiliencyPolicy = helpers.ExpandVirtualMachineScaleSetResiliency(resilientVMCreationEnabled, resilientVMDeletionEnabled)
	}

	if d.HasChange("termination_notification") {
		notificationRaw := d.Get("termination_notification").([]any)
		updateProps.VirtualMachineProfile.ScheduledEventsProfile = helpers.ExpandVirtualMachineScaleSetScheduledEventsProfile(notificationRaw)
	}

	if d.HasChange("encryption_at_host_enabled") {
		updateProps.VirtualMachineProfile.SecurityProfile = &virtualmachinescalesets.SecurityProfile{
			EncryptionAtHost: pointer.To(d.Get("encryption_at_host_enabled").(bool)),
		}
	}

	if d.HasChange("automatic_instance_repair") {
		automaticRepairsPolicyRaw := d.Get("automatic_instance_repair").([]any)
		automaticRepairsPolicy := helpers.ExpandVirtualMachineScaleSetAutomaticRepairsPolicy(automaticRepairsPolicyRaw)

		if automaticRepairsPolicy != nil {
			// we need to know if the VMSS has a health extension or not
			hasHealthExtension := false

			if v, ok := d.GetOk("extension"); ok {
				var err error
				_, hasHealthExtension, err = helpers.ExpandOrchestratedVirtualMachineScaleSetExtensions(v.(*pluginsdk.Set).List())
				if err != nil {
					return err
				}
			}

			_, hasHealthProbeId := d.GetOk("health_probe_id")

			if !hasHealthProbeId && !hasHealthExtension {
				return fmt.Errorf("`automatic_instance_repair` can only be set if there is an application Health extension or a `health_probe_id` defined")
			}
		}

		updateProps.AutomaticRepairsPolicy = automaticRepairsPolicy
	}

	if d.HasChange("identity") {
		identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		existing.Model.Identity = identityExpanded
		// Removing a user-assigned identity using PATCH requires setting it to `null` in the payload which
		// 1. The go-azure-sdk for resource manager doesn't support at the moment
		// 2. The expand identity function doesn't behave this way
		// For the moment updating the identity with the PUT circumvents this API behaviour
		// See https://github.com/hashicorp/terraform-provider-azurerm/issues/25058 for more details
		if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, virtualmachinescalesets.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating identity for Linux %s: %+v", id, err)
		}
	}

	if d.HasChange("plan") {
		planRaw := d.Get("plan").([]any)
		update.Plan = helpers.ExpandPlanVMSS(planRaw)
	}

	if d.HasChanges("sku", "instances") {
		// in-case ignore_changes is being used, since both fields are required
		// look up the current values and override them as needed
		sku := existing.Model.Sku

		if d.HasChange("sku") {
			updateInstances = true

			sku.Name = pointer.To(d.Get("sku").(string))
		}

		if d.HasChange("instances") {
			sku.Capacity = pointer.To(int64(d.Get("instances").(int)))
		}

		update.Sku = sku
	}

	if d.HasChanges("extension", "extensions_time_budget") {
		updateInstances = true

		extensionProfile, _, err := helpers.ExpandVirtualMachineScaleSetExtensions(d.Get("extension").(*pluginsdk.Set).List())
		if err != nil {
			return err
		}
		updateProps.VirtualMachineProfile.ExtensionProfile = extensionProfile
		updateProps.VirtualMachineProfile.ExtensionProfile.ExtensionsTimeBudget = pointer.To(d.Get("extensions_time_budget").(string))
	}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("user_data") {
		updateInstances = true
		updateProps.VirtualMachineProfile.UserData = pointer.To(d.Get("user_data").(string))
	}

	update.Properties = &updateProps

	metaData := helpers.VirtualMachineScaleSetUpdateMetaData{
		AutomaticOSUpgradeIsEnabled:  automaticOSUpgradeIsEnabled,
		CanReimageOnManualUpgrade:    meta.(*clients.Client).Features.VirtualMachineScaleSet.ReimageOnManualUpgrade,
		CanRollInstancesWhenRequired: meta.(*clients.Client).Features.VirtualMachineScaleSet.RollInstancesWhenRequired,
		UpdateInstances:              updateInstances,
		Client:                       meta.(*clients.Client).Compute,
		Existing:                     *existing.Model,
		ID:                           id,
		OSType:                       virtualmachinescalesets.OperatingSystemTypesLinux,
	}

	if err := metaData.PerformUpdate(ctx, update); err != nil {
		return err
	}

	return resourceLinuxVirtualMachineScaleSetRead(d, meta)
}
