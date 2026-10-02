// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package orchestrated_virtual_machine_scale_set

import (
	"fmt"
	"log"
	"strconv"

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

func resourceOrchestratedVirtualMachineScaleSetUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	isLegacy := true
	updateInstances := false
	isHotpatchEnabledImage := false
	linuxAutomaticVMGuestPatchingEnabled := false

	options := virtualmachinescalesets.DefaultGetOperationOptions()
	options.Expand = pointer.To(virtualmachinescalesets.ExpandTypesForGetVMScaleSetsUserData)
	existing, err := client.Get(ctx, *id, options)
	if err != nil {
		return fmt.Errorf("retrieving Orchestrated %s: %w", id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving Orchestrated %s: `model` was nil", id)
	}
	if existing.Model.Sku != nil {
		isLegacy = false
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving Orchestrated %s: `properties` was nil", id)
	}

	if !isLegacy {
		if existing.Model.Properties.VirtualMachineProfile == nil {
			return fmt.Errorf("retrieving Orchestrated %s: `Properties.VirtualMachineProfile` was nil", id)
		}
		if existing.Model.Properties.VirtualMachineProfile.StorageProfile == nil {
			return fmt.Errorf("retrieving Orchestrated %s: `Properties.VirtualMachineProfile.StorageProfile` was nil", id)
		}
	}

	updateProps := virtualmachinescalesets.VirtualMachineScaleSetUpdateProperties{}
	update := virtualmachinescalesets.VirtualMachineScaleSetUpdate{}
	var osType virtualmachinescalesets.OperatingSystemTypes

	if !isLegacy {
		updateProps = virtualmachinescalesets.VirtualMachineScaleSetUpdateProperties{
			VirtualMachineProfile: &virtualmachinescalesets.VirtualMachineScaleSetUpdateVMProfile{
				// if an image reference has been configured previously (it has to be), we would better to include that in this
				// update request to avoid some circumstances that the API will complain ImageReference is null
				// issue tracking: https://github.com/Azure/azure-rest-api-specs/issues/10322
				StorageProfile: &virtualmachinescalesets.VirtualMachineScaleSetUpdateStorageProfile{
					ImageReference: existing.Model.Properties.VirtualMachineProfile.StorageProfile.ImageReference,
				},
			},
			// Currently not supported in orchestrated VMSS
			// if an upgrade policy's been configured previously (which it will have) it must be threaded through
			// this doesn't matter for Manual - but breaks when updating anything on a Automatic and Rolling Mode Scale Set
			// UpgradePolicy: existing.Properties.UpgradePolicy,
		}

		priority := virtualmachinescalesets.VirtualMachinePriorityTypes(d.Get("priority").(string))

		if d.HasChange("single_placement_group") {
			// Since null is now a valid value for single_placement_group
			// make sure it is in the config file before you set the value
			// on the update props...
			if !pluginsdk.IsExplicitlyNullInConfig(d, "single_placement_group") {
				singlePlacementGroup := d.Get("single_placement_group").(bool)
				if singlePlacementGroup {
					return fmt.Errorf("`single_placement_group` cannot be changed from `false` to `true`")
				}
				updateProps.SinglePlacementGroup = pointer.To(singlePlacementGroup)
			}
		}

		if d.HasChange("sku_profile") {
			updateInstances = true
			updateProps.SkuProfile = expandOrchestratedVirtualMachineScaleSetSkuProfile(d.Get("sku_profile").([]any))
		}

		if d.HasChange("max_bid_price") {
			if priority != virtualmachinescalesets.VirtualMachinePriorityTypesSpot {
				return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `%s`", string(virtualmachinescalesets.VirtualMachinePriorityTypesSpot))
			}

			updateProps.VirtualMachineProfile.BillingProfile = &virtualmachinescalesets.BillingProfile{
				MaxPrice: pointer.To(d.Get("max_bid_price").(float64)),
			}
		}

		osProfileRaw := d.Get("os_profile").([]any)
		vmssOsProfile := virtualmachinescalesets.VirtualMachineScaleSetUpdateOSProfile{}
		windowsConfig := virtualmachinescalesets.WindowsConfiguration{}
		windowsConfig.PatchSettings = &virtualmachinescalesets.PatchSettings{}
		linuxConfig := virtualmachinescalesets.LinuxConfiguration{}

		if len(osProfileRaw) > 0 && osProfileRaw[0] != nil {
			osProfile := osProfileRaw[0].(map[string]any)
			winConfigRaw := osProfile["windows_configuration"].([]any)
			linConfigRaw := osProfile["linux_configuration"].([]any)

			if d.HasChange("os_profile.0.custom_data") {
				updateInstances = true

				// customData can only be sent if it's a base64 encoded string,
				// so it's not possible to remove this without tainting the resource
				vmssOsProfile.CustomData = pointer.To(osProfile["custom_data"].(string))
			}

			if len(winConfigRaw) > 0 && winConfigRaw[0] != nil {
				osType = virtualmachinescalesets.OperatingSystemTypesWindows
				winConfig := winConfigRaw[0].(map[string]any)
				provisionVMAgent := winConfig["provision_vm_agent"].(bool)
				patchAssessmentMode := winConfig["patch_assessment_mode"].(string)
				patchMode := winConfig["patch_mode"].(string)
				hotpatchingEnabled := winConfig["hotpatching_enabled"].(bool)

				// If the image allows hotpatching the patch mode can only ever be AutomaticByPlatform.
				sourceImageReferenceRaw := d.Get("source_image_reference").([]any)
				sourceImageId := d.Get("source_image_id").(string)
				isHotpatchEnabledImage = helpers.IsValidHotPatchSourceImageReference(sourceImageReferenceRaw, sourceImageId)

				// PatchSettings is required by PATCH API when running Hotpatch-compatible images.
				if isHotpatchEnabledImage {
					windowsConfig.PatchSettings.AssessmentMode = pointer.ToEnum[virtualmachinescalesets.WindowsPatchAssessmentMode](patchAssessmentMode)
					windowsConfig.PatchSettings.PatchMode = pointer.ToEnum[virtualmachinescalesets.WindowsVMGuestPatchMode](patchMode)
					windowsConfig.PatchSettings.EnableHotpatching = pointer.To(hotpatchingEnabled)
				}

				// lintignore:R019 // deliberate subset: the windows_configuration fields that require rolling the instances to take effect
				if d.HasChanges(
					"os_profile.0.windows_configuration.0.enable_automatic_updates",
					"os_profile.0.windows_configuration.0.provision_vm_agent",
					"os_profile.0.windows_configuration.0.timezone",
					"os_profile.0.windows_configuration.0.secret",
					"os_profile.0.windows_configuration.0.winrm_listener",
				) {
					updateInstances = true
				}

				if d.HasChange("os_profile.0.windows_configuration.0.enable_automatic_updates") {
					windowsConfig.EnableAutomaticUpdates = pointer.To(winConfig["enable_automatic_updates"].(bool))
				}

				if d.HasChange("os_profile.0.windows_configuration.0.provision_vm_agent") {
					if isHotpatchEnabledImage && !provisionVMAgent {
						return fmt.Errorf("when using a hotpatching enabled image, `provision_vm_agent` must be set to `true`, got `%s`", strconv.FormatBool(provisionVMAgent))
					}
					windowsConfig.ProvisionVMAgent = pointer.To(provisionVMAgent)
				}

				if d.HasChange("os_profile.0.windows_configuration.0.patch_assessment_mode") {
					if !provisionVMAgent && (patchAssessmentMode == string(virtualmachinescalesets.WindowsPatchAssessmentModeAutomaticByPlatform)) {
						return fmt.Errorf("when `patch_assessment_mode` is set to `%s`, `provision_vm_agent` must be set to `true`", virtualmachinescalesets.WindowsPatchAssessmentModeAutomaticByPlatform)
					}
					windowsConfig.PatchSettings.AssessmentMode = pointer.ToEnum[virtualmachinescalesets.WindowsPatchAssessmentMode](patchAssessmentMode)
				}

				if d.HasChange("os_profile.0.windows_configuration.0.patch_mode") {
					if isHotpatchEnabledImage && (patchMode != string(virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByPlatform)) {
						return fmt.Errorf("when using a hotpatching enabled image, `patch_mode` must be set to `%s`, got `%s`", virtualmachinescalesets.WindowsVMGuestPatchModeAutomaticByPlatform, patchMode)
					}
					windowsConfig.PatchSettings.PatchMode = pointer.ToEnum[virtualmachinescalesets.WindowsVMGuestPatchMode](patchMode)
				}

				// Disabling hotpatching is not supported in images that support hotpatching
				// so while the attribute is exposed in VMSS it is hardcoded inside the images that
				// support hotpatching to always be enabled and cannot be set to false, ever.
				if d.HasChange("os_profile.0.windows_configuration.0.hotpatching_enabled") {
					if isHotpatchEnabledImage && !hotpatchingEnabled {
						return fmt.Errorf("when using a hotpatching enabled image, `hotpatching_enabled` must be set to `true`, got `%s`", strconv.FormatBool(hotpatchingEnabled))
					}
					windowsConfig.PatchSettings.EnableHotpatching = pointer.To(hotpatchingEnabled)
				}

				if d.HasChange("os_profile.0.windows_configuration.0.secret") {
					vmssOsProfile.Secrets = helpers.ExpandWindowsSecretsVMSS(winConfig["secret"].([]any))
				}

				if d.HasChange("os_profile.0.windows_configuration.0.timezone") {
					windowsConfig.TimeZone = pointer.To(winConfig["timezone"].(string))
				}

				if d.HasChange("os_profile.0.windows_configuration.0.winrm_listener") {
					winRmListenersRaw := winConfig["winrm_listener"].(*pluginsdk.Set).List()
					vmssOsProfile.WindowsConfiguration.WinRM = helpers.ExpandWinRMListenerVMSS(winRmListenersRaw)
				}

				vmssOsProfile.WindowsConfiguration = &windowsConfig
			}

			if len(linConfigRaw) > 0 && linConfigRaw[0] != nil {
				osType = virtualmachinescalesets.OperatingSystemTypesLinux
				linConfig := linConfigRaw[0].(map[string]any)
				provisionVMAgent := linConfig["provision_vm_agent"].(bool)
				patchAssessmentMode := linConfig["patch_assessment_mode"].(string)
				patchMode := linConfig["patch_mode"].(string)

				if d.HasChanges(
					"os_profile.0.linux_configuration.0.provision_vm_agent",
					"os_profile.0.linux_configuration.0.disable_password_authentication",
					"os_profile.0.linux_configuration.0.admin_ssh_key",
				) {
					updateInstances = true
				}

				if d.HasChange("os_profile.0.linux_configuration.0.provision_vm_agent") {
					linuxConfig.ProvisionVMAgent = pointer.To(provisionVMAgent)
				}

				if d.HasChange("os_profile.0.linux_configuration.0.disable_password_authentication") {
					linuxConfig.DisablePasswordAuthentication = pointer.To(linConfig["disable_password_authentication"].(bool))
				}

				if d.HasChange("os_profile.0.linux_configuration.0.admin_ssh_key") {
					sshPublicKeys := helpers.ExpandSSHKeysVMSS(linConfig["admin_ssh_key"].(*pluginsdk.Set).List())
					if linuxConfig.Ssh == nil {
						linuxConfig.Ssh = &virtualmachinescalesets.SshConfiguration{}
					}
					linuxConfig.Ssh.PublicKeys = &sshPublicKeys
				}

				if d.HasChange("os_profile.0.linux_configuration.0.patch_assessment_mode") {
					if !provisionVMAgent && (patchAssessmentMode == string(virtualmachinescalesets.LinuxPatchAssessmentModeAutomaticByPlatform)) {
						return fmt.Errorf("when the `patch_assessment_mode` field is set to `%s` the `provision_vm_agent` must always be set to `true`", virtualmachinescalesets.LinuxPatchAssessmentModeAutomaticByPlatform)
					}

					if linuxConfig.PatchSettings == nil {
						linuxConfig.PatchSettings = &virtualmachinescalesets.LinuxPatchSettings{}
					}
					linuxConfig.PatchSettings.AssessmentMode = pointer.ToEnum[virtualmachinescalesets.LinuxPatchAssessmentMode](patchAssessmentMode)
				}

				if d.HasChange("os_profile.0.linux_configuration.0.patch_mode") {
					if patchMode == string(virtualmachinescalesets.LinuxPatchAssessmentModeAutomaticByPlatform) {
						if !provisionVMAgent {
							return fmt.Errorf("when the `patch_mode` field is set to `%s` the `provision_vm_agent` field must always be set to `true`, got `%s`", patchMode, strconv.FormatBool(provisionVMAgent))
						}

						linuxAutomaticVMGuestPatchingEnabled = true
					}

					if linuxConfig.PatchSettings == nil {
						linuxConfig.PatchSettings = &virtualmachinescalesets.LinuxPatchSettings{}
					}
					linuxConfig.PatchSettings.PatchMode = pointer.ToEnum[virtualmachinescalesets.LinuxVMGuestPatchMode](patchMode)
				}

				vmssOsProfile.LinuxConfiguration = &linuxConfig
			}

			updateProps.VirtualMachineProfile.OsProfile = &vmssOsProfile
		}

		if d.HasChanges("data_disk", "os_disk", "source_image_id", "source_image_reference") {
			updateInstances = true

			if updateProps.VirtualMachineProfile.StorageProfile == nil {
				updateProps.VirtualMachineProfile.StorageProfile = &virtualmachinescalesets.VirtualMachineScaleSetUpdateStorageProfile{}
			}

			if d.HasChange("data_disk") {
				ultraSSDEnabled := false // Currently not supported in orchestrated vmss
				dataDisks, err := helpers.ExpandOrchestratedVirtualMachineScaleSetDataDisk(d.Get("data_disk").([]any), ultraSSDEnabled)
				if err != nil {
					return fmt.Errorf("expanding `data_disk`: %w", err)
				}
				updateProps.VirtualMachineProfile.StorageProfile.DataDisks = dataDisks
			}

			if d.HasChange("os_disk") {
				osDiskRaw := d.Get("os_disk").([]any)
				updateProps.VirtualMachineProfile.StorageProfile.OsDisk = helpers.ExpandOrchestratedVirtualMachineScaleSetOSDiskUpdate(osDiskRaw)
			}

			if d.HasChanges("source_image_id", "source_image_reference") {
				sourceImageReferenceRaw := d.Get("source_image_reference").([]any)
				sourceImageId := d.Get("source_image_id").(string)

				if len(sourceImageReferenceRaw) != 0 || sourceImageId != "" {
					updateProps.VirtualMachineProfile.StorageProfile.ImageReference = helpers.ExpandSourceImageReferenceVMSS(sourceImageReferenceRaw, sourceImageId)
				}

				// Must include all storage profile properties when updating disk image.  See: https://github.com/hashicorp/terraform-provider-azurerm/issues/8273
				updateProps.VirtualMachineProfile.StorageProfile.DataDisks = existing.Model.Properties.VirtualMachineProfile.StorageProfile.DataDisks
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

		if d.HasChanges("network_api_version", "network_interface") {
			if updateProps.VirtualMachineProfile.NetworkProfile == nil {
				updateProps.VirtualMachineProfile.NetworkProfile = &virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkProfile{}
			}

			updateProps.VirtualMachineProfile.NetworkProfile.NetworkApiVersion = pointer.ToEnum[virtualmachinescalesets.NetworkApiVersion](d.Get("network_api_version").(string))

			networkInterfacesRaw := d.Get("network_interface").([]any)
			networkInterfaces, err := helpers.ExpandOrchestratedVirtualMachineScaleSetNetworkInterfaceUpdate(networkInterfacesRaw)
			if err != nil {
				return fmt.Errorf("expanding `network_interface`: %w", err)
			}

			updateProps.VirtualMachineProfile.NetworkProfile.NetworkInterfaceConfigurations = networkInterfaces
		}

		if d.HasChange("boot_diagnostics") {
			updateInstances = true

			bootDiagnosticsRaw := d.Get("boot_diagnostics").([]any)
			updateProps.VirtualMachineProfile.DiagnosticsProfile = helpers.ExpandBootDiagnosticsVMSS(bootDiagnosticsRaw)
		}

		if d.HasChange("termination_notification") {
			notificationRaw := d.Get("termination_notification").([]any)
			updateProps.VirtualMachineProfile.ScheduledEventsProfile = helpers.ExpandOrchestratedVirtualMachineScaleSetScheduledEventsProfile(notificationRaw)
		}

		if d.HasChange("encryption_at_host_enabled") {
			updateProps.VirtualMachineProfile.SecurityProfile = &virtualmachinescalesets.SecurityProfile{
				EncryptionAtHost: pointer.To(d.Get("encryption_at_host_enabled").(bool)),
			}
		}

		if d.HasChange("license_type") {
			license := d.Get("license_type").(string)
			if license == "" {
				// Only for create no specification is possible in the API. API does not allow empty string in update.
				// So removing attribute license_type from Terraform configuration if it was set to value other than 'None' would lead to an endless loop in apply.
				// To allow updating in this case set value explicitly to 'None'.
				license = "None"
			}
			updateProps.VirtualMachineProfile.LicenseType = &license
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

				if !hasHealthExtension {
					return fmt.Errorf("`automatic_instance_repair` can only be enabled when an application health extension is configured")
				}
			}
			updateProps.AutomaticRepairsPolicy = automaticRepairsPolicy
		}

		if d.HasChange("identity") {
			identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %w", err)
			}

			update.Identity = identityExpanded
		}

		if d.HasChange("plan") {
			planRaw := d.Get("plan").([]any)
			update.Plan = helpers.ExpandPlanVMSS(planRaw)
		}

		if d.HasChanges("sku_name", "instances") {
			// in-case ignore_changes is being used, since both fields are required
			// look up the current values and override them as needed
			sku := existing.Model.Sku
			instances := int(*sku.Capacity)
			skuName := d.Get("sku_name").(string)

			if d.HasChange("instances") {
				instances = d.Get("instances").(int)

				sku, err = expandOrchestratedVirtualMachineScaleSetSku(skuName, instances)
				if err != nil {
					return err
				}
			}

			if d.HasChange("sku_name") {
				updateInstances = true

				sku, err = expandOrchestratedVirtualMachineScaleSetSku(skuName, instances)
				if err != nil {
					return err
				}
			}

			update.Sku = sku
		}

		if d.HasChanges("extension", "extensions_time_budget") {
			updateInstances = true

			extensionProfile, hasHealthExtension, err := helpers.ExpandOrchestratedVirtualMachineScaleSetExtensions(d.Get("extension").(*pluginsdk.Set).List())
			if err != nil {
				return err
			}

			if isHotpatchEnabledImage && !hasHealthExtension {
				return fmt.Errorf("when using a hotpatching enabled image, an application health extension must be configured")
			}

			if linuxAutomaticVMGuestPatchingEnabled && !hasHealthExtension {
				return fmt.Errorf("when `patch_mode` is set to `%s`, at least one application health extension must be configured, got 0", virtualmachinescalesets.LinuxPatchAssessmentModeAutomaticByPlatform)
			}

			updateProps.VirtualMachineProfile.ExtensionProfile = extensionProfile
			updateProps.VirtualMachineProfile.ExtensionProfile.ExtensionsTimeBudget = pointer.To(d.Get("extensions_time_budget").(string))
		}
	}

	if d.HasChange("zones") {
		update.Zones = pointer.To(zones.ExpandUntyped(d.Get("zones").(*schema.Set).List()))
	}

	if d.HasChange("rolling_upgrade_policy") {
		upgradePolicy := virtualmachinescalesets.UpgradePolicy{}

		if existing.Model.Properties.UpgradePolicy != nil {
			upgradePolicy = *existing.Model.Properties.UpgradePolicy
		}

		upgradePolicy.Mode = pointer.ToEnum[virtualmachinescalesets.UpgradeMode](d.Get("upgrade_mode").(string))

		rollingRaw := d.Get("rolling_upgrade_policy").([]any)
		rollingUpgradePolicy, err := helpers.ExpandVirtualMachineScaleSetRollingUpgradePolicy(rollingRaw, len(zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())) > 0, false)
		if err != nil {
			return fmt.Errorf("expanding `rolling_upgrade_policy`: %w", err)
		}

		upgradePolicy.RollingUpgradePolicy = rollingUpgradePolicy
		updateProps.UpgradePolicy = &upgradePolicy
	}

	// Only two fields that can change in legacy mode
	if d.HasChange("proximity_placement_group_id") {
		if v, ok := d.GetOk("proximity_placement_group_id"); ok {
			updateInstances = true
			updateProps.ProximityPlacementGroup = &virtualmachinescalesets.SubResource{
				Id: pointer.To(v.(string)),
			}
		}
	}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("user_data_base64") {
		updateInstances = true
		updateProps.VirtualMachineProfile.UserData = pointer.To(d.Get("user_data_base64").(string))
	}

	update.Properties = &updateProps

	if updateInstances {
		log.Printf("[DEBUG] Orchestrated %s - updateInstances is true", id)
	}

	// AutomaticOSUpgradeIsEnabled currently is not supported in orchestrated VMSS flex
	metaData := helpers.VirtualMachineScaleSetUpdateMetaData{
		Client:   meta.(*clients.Client).Compute,
		Existing: pointer.From(existing.Model),
		ID:       id,
		OSType:   osType,
	}

	if err := metaData.PerformUpdate(ctx, update); err != nil {
		return err
	}

	return resourceOrchestratedVirtualMachineScaleSetRead(d, meta)
}
