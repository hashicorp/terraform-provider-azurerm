// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_virtual_machine

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWindowsVirtualMachineCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualmachines.NewVirtualMachineID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.VirtualMachineName, helpers.VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, helpers.VirtualMachineResourceName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.Get(ctx, id, virtualmachines.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing Windows %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError(azureWindowsVirtualMachineResourceName, id.ID())
		}
	}

	additionalCapabilitiesRaw := d.Get("additional_capabilities").([]any)
	additionalCapabilities := helpers.ExpandVirtualMachineAdditionalCapabilities(additionalCapabilitiesRaw)

	additionalUnattendContentRaw := d.Get("additional_unattend_content").([]any)
	additionalUnattendContent := helpers.ExpandAdditionalUnattendContent(additionalUnattendContentRaw)

	allowExtensionOperations := true
	if !d.GetRawConfig().AsValueMap()["allow_extension_operations"].IsNull() {
		allowExtensionOperations = d.Get("allow_extension_operations").(bool)
	}

	bootDiagnosticsRaw := d.Get("boot_diagnostics").([]any)
	bootDiagnostics := helpers.ExpandBootDiagnostics(bootDiagnosticsRaw)

	identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	planRaw := d.Get("plan").([]any)
	plan := helpers.ExpandPlan(planRaw)

	priority := virtualmachines.VirtualMachinePriorityTypes(d.Get("priority").(string))

	patchMode := d.Get("patch_mode").(string)
	assessmentMode := d.Get("patch_assessment_mode").(string)
	hotPatch := d.Get("hotpatching_enabled").(bool)
	size := d.Get("size").(string)
	t := d.Get("tags").(map[string]any)

	networkInterfaceIdsRaw := d.Get("network_interface_ids").([]any)
	networkInterfaceIds := helpers.ExpandVirtualMachineNetworkInterfaceIDs(networkInterfaceIdsRaw)

	managedDiskIdRaw := d.Get("os_managed_disk_id").(string)
	// Note: The API fails if OsProfile is anything but nil with CreateOption = "Attach"
	osDiskIsImported := managedDiskIdRaw != ""

	secretsRaw := d.Get("secret").([]any)
	secrets := helpers.ExpandWindowsSecrets(secretsRaw)

	params := virtualmachines.VirtualMachine{
		Name:             pointer.To(id.VirtualMachineName),
		ExtendedLocation: helpers.ExpandEdgeZone(d.Get("edge_zone").(string)),
		Location:         location.Normalize(d.Get("location").(string)),
		Identity:         identityExpanded,
		Plan:             plan,
		Properties: &virtualmachines.VirtualMachineProperties{
			ApplicationProfile: &virtualmachines.ApplicationProfile{
				GalleryApplications: helpers.ExpandVirtualMachineGalleryApplication(d.Get("gallery_application").([]any)),
			},
			HardwareProfile: &virtualmachines.HardwareProfile{
				VMSize: pointer.ToEnum[virtualmachines.VirtualMachineSizeTypes](size),
			},
			NetworkProfile: &virtualmachines.NetworkProfile{
				NetworkInterfaces: &networkInterfaceIds,
			},
			Priority: pointer.To(priority),
			StorageProfile: &virtualmachines.StorageProfile{
				// Data Disks are instead handled via the Association resource - as such we can send an empty value here
				// but for Updates this'll need to be nil, else any associations will be overwritten
				DataDisks: &[]virtualmachines.DataDisk{},
			},

			// Optional
			AdditionalCapabilities: additionalCapabilities,
			DiagnosticsProfile:     bootDiagnostics,
			ExtensionsTimeBudget:   pointer.To(d.Get("extensions_time_budget").(string)),
		},
		Tags: tags.Expand(t),
	}
	osDiskRaw := d.Get("os_disk").([]any)
	osDisk, err := helpers.ExpandVirtualMachineOSDisk(osDiskRaw, virtualmachines.OperatingSystemTypesWindows)
	if err != nil {
		return fmt.Errorf("expanding `os_disk`: %+v", err)
	}
	securityEncryptionType := ""
	if !osDiskIsImported {
		securityEncryptionType = osDiskRaw[0].(map[string]any)["security_encryption_type"].(string)
		var computerName string
		if v, ok := d.GetOk("computer_name"); ok && len(v.(string)) > 0 {
			computerName = v.(string)
		} else {
			_, errs := validate.WindowsComputerNameFull(d.Get("name"), "computer_name")
			if len(errs) > 0 {
				return fmt.Errorf("unable to assume default computer name %s. Please adjust the `name`, or specify an explicit `computer_name`", errs[0])
			}
			computerName = id.VirtualMachineName
		}

		provisionVMAgent := true
		if p, ok := d.GetRawConfig().AsValueMap()["provision_vm_agent"]; ok && !p.IsNull() {
			provisionVMAgent = d.Get("provision_vm_agent").(bool)
		}

		if !provisionVMAgent && allowExtensionOperations {
			return fmt.Errorf("`allow_extension_operations` cannot be set to `true` when `provision_vm_agent` is set to `false`")
		}

		autoUpdatesEnabled := true

		if !d.GetRawConfig().AsValueMap()["automatic_updates_enabled"].IsNull() {
			autoUpdatesEnabled = d.Get("automatic_updates_enabled").(bool)
		}

		params.Properties.OsProfile = &virtualmachines.OSProfile{
			AdminPassword:            pointer.To(d.Get("admin_password").(string)),
			AdminUsername:            pointer.To(d.Get("admin_username").(string)),
			ComputerName:             pointer.To(computerName),
			AllowExtensionOperations: pointer.To(allowExtensionOperations),
			WindowsConfiguration: &virtualmachines.WindowsConfiguration{
				ProvisionVMAgent:       pointer.To(provisionVMAgent),
				EnableAutomaticUpdates: pointer.To(autoUpdatesEnabled),
				WinRM:                  helpers.ExpandWinRMListener(d.Get("winrm_listener").(*pluginsdk.Set).List()),
			},
			Secrets: secrets,
		}

		if len(additionalUnattendContentRaw) > 0 {
			params.Properties.OsProfile.WindowsConfiguration.AdditionalUnattendContent = additionalUnattendContent
		}
		sourceImageReferenceRaw := d.Get("source_image_reference").([]any)
		sourceImageId := d.Get("source_image_id").(string)
		if len(sourceImageReferenceRaw) != 0 || sourceImageId != "" {
			params.Properties.StorageProfile.ImageReference = helpers.ExpandSourceImageReference(sourceImageReferenceRaw, sourceImageId)
		}

		isHotpatchImage := helpers.IsValidHotPatchSourceImageReference(sourceImageReferenceRaw, sourceImageId)

		// Validate VM Guest Patch Mode configuration
		if patchMode == string(virtualmachines.WindowsVMGuestPatchModeAutomaticByPlatform) && !provisionVMAgent {
			return fmt.Errorf("%q cannot be set to %q when %q is set to %q", "patch_mode", "AutomaticByPlatform", "provision_vm_agent", "false")
		}

		if assessmentMode == string(virtualmachines.WindowsPatchAssessmentModeAutomaticByPlatform) && !provisionVMAgent {
			return fmt.Errorf("`provision_vm_agent` must be set to `true` when `patch_assessment_mode` is set to `AutomaticByPlatform`")
		}

		if isHotpatchImage && patchMode != string(virtualmachines.WindowsVMGuestPatchModeAutomaticByPlatform) {
			return fmt.Errorf("%q must always be set to %q when %q points to a hotpatch enabled image", "patch_mode", "AutomaticByPlatform", "source_image_reference")
		}

		// hot patching can only be enabled if the patch_mode is set to "AutomaticByPlatform"
		// and if the image reference is using one of the following skus:
		// 2022-datacenter-azure-edition-core, 2022-datacenter-azure-edition-core-smalldisk, 2022-datacenter-azure-edition-hotpatch, 2022-datacenter-azure-edition-hotpatch-smalldisk,
		// 2025-datacenter-azure-edition, 2025-datacenter-azure-edition-smalldisk, 2025-datacenter-azure-edition-core, 2025-datacenter-azure-edition-core-smalldisk
		if hotPatch {
			if patchMode != string(virtualmachines.WindowsVMGuestPatchModeAutomaticByPlatform) {
				return fmt.Errorf("%q cannot be set to %q when %q is set to %q", "hotpatching_enabled", "true", "patch_mode", patchMode)
			}

			if !provisionVMAgent {
				return fmt.Errorf("%q cannot be set to %q when %q is set to %q", "hotpatching_enabled", "true", "provisionVMAgent", "false")
			}

			if !isHotpatchImage {
				if sourceImageId != "" {
					return fmt.Errorf("the %q field is not supported if referencing the image via the %q field", "hotpatching_enabled", "source_image_id")
				}

				return fmt.Errorf("%q is currently only supported on %q, %q, %q, %q, %q, %q, %q or %q image reference skus", "hotpatching_enabled", "2022-datacenter-azure-edition-core", "2022-datacenter-azure-edition-core-smalldisk", "2022-datacenter-azure-edition-hotpatch", "2022-datacenter-azure-edition-hotpatch-smalldisk", "2025-datacenter-azure-edition", "2025-datacenter-azure-edition-smalldisk", "2025-datacenter-azure-edition-core", "2025-datacenter-azure-edition-core-smalldisk")
			}
		}

		params.Properties.OsProfile.WindowsConfiguration.PatchSettings = &virtualmachines.PatchSettings{
			PatchMode:         pointer.ToEnum[virtualmachines.WindowsVMGuestPatchMode](patchMode),
			EnableHotpatching: pointer.To(hotPatch),
			AssessmentMode:    pointer.ToEnum[virtualmachines.WindowsPatchAssessmentMode](assessmentMode),
		}

		if d.Get("bypass_platform_safety_checks_on_user_schedule_enabled").(bool) {
			if patchMode != string(virtualmachines.WindowsVMGuestPatchModeAutomaticByPlatform) {
				return fmt.Errorf("`patch_mode` must be set to `AutomaticByPlatform` when `bypass_platform_safety_checks_on_user_schedule_enabled` is set to `true`")
			}

			if params.Properties.OsProfile.WindowsConfiguration.PatchSettings.AutomaticByPlatformSettings == nil {
				params.Properties.OsProfile.WindowsConfiguration.PatchSettings.AutomaticByPlatformSettings = &virtualmachines.WindowsVMGuestPatchAutomaticByPlatformSettings{}
			}

			params.Properties.OsProfile.WindowsConfiguration.PatchSettings.AutomaticByPlatformSettings.BypassPlatformSafetyChecksOnUserSchedule = pointer.To(true)
		}

		if v, ok := d.GetOk("reboot_setting"); ok {
			if patchMode != string(virtualmachines.WindowsVMGuestPatchModeAutomaticByPlatform) {
				return fmt.Errorf("`patch_mode` must be set to `AutomaticByPlatform` when `reboot_setting` is specified")
			}

			if params.Properties.OsProfile.WindowsConfiguration.PatchSettings.AutomaticByPlatformSettings == nil {
				params.Properties.OsProfile.WindowsConfiguration.PatchSettings.AutomaticByPlatformSettings = &virtualmachines.WindowsVMGuestPatchAutomaticByPlatformSettings{}
			}

			params.Properties.OsProfile.WindowsConfiguration.PatchSettings.AutomaticByPlatformSettings.RebootSetting = pointer.ToEnum[virtualmachines.WindowsVMGuestPatchAutomaticByPlatformRebootSetting](v.(string))
		}

		if v, ok := d.GetOk("custom_data"); ok {
			params.Properties.OsProfile.CustomData = pointer.To(v.(string))
		}

		if v, ok := d.GetOk("timezone"); ok {
			params.Properties.OsProfile.WindowsConfiguration.TimeZone = pointer.To(v.(string))
		}
	} else {
		diskId, err := commonids.ParseManagedDiskID(managedDiskIdRaw)
		if err != nil {
			return err
		}

		osDisk.ManagedDisk.Id = pointer.To(diskId.ID())
		osDisk.CreateOption = virtualmachines.DiskCreateOptionTypesAttach
	}

	params.Properties.StorageProfile.OsDisk = osDisk

	if diskControllerType, ok := d.GetOk("disk_controller_type"); ok {
		params.Properties.StorageProfile.DiskControllerType = pointer.ToEnum[virtualmachines.DiskControllerTypes](diskControllerType.(string))
	}

	if v, ok := d.GetOk("availability_set_id"); ok {
		params.Properties.AvailabilitySet = &virtualmachines.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("capacity_reservation_group_id"); ok {
		params.Properties.CapacityReservation = &virtualmachines.CapacityReservationProfile{
			CapacityReservationGroup: &virtualmachines.SubResource{
				Id: pointer.To(v.(string)),
			},
		}
	}

	if v, ok := d.GetOk("dedicated_host_id"); ok {
		params.Properties.Host = &virtualmachines.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("dedicated_host_group_id"); ok {
		params.Properties.HostGroup = &virtualmachines.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if encryptionAtHostEnabled, ok := d.GetOk("encryption_at_host_enabled"); ok {
		if params.Properties.SecurityProfile == nil {
			params.Properties.SecurityProfile = &virtualmachines.SecurityProfile{}
		}
		params.Properties.SecurityProfile.EncryptionAtHost = pointer.To(encryptionAtHostEnabled.(bool))
	}

	secureBootEnabled := d.Get("secure_boot_enabled").(bool)
	vtpmEnabled := d.Get("vtpm_enabled").(bool)
	if securityEncryptionType != "" {
		if virtualmachines.SecurityEncryptionTypesDiskWithVMGuestState == virtualmachines.SecurityEncryptionTypes(securityEncryptionType) && !secureBootEnabled {
			return fmt.Errorf("`secure_boot_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is set to `DiskWithVMGuestState`")
		}
		if !vtpmEnabled {
			return fmt.Errorf("`vtpm_enabled` must be set to `true` when `os_disk.0.security_encryption_type` is set")
		}

		if params.Properties.SecurityProfile == nil {
			params.Properties.SecurityProfile = &virtualmachines.SecurityProfile{}
		}
		params.Properties.SecurityProfile.SecurityType = pointer.To(virtualmachines.SecurityTypesConfidentialVM)

		if params.Properties.SecurityProfile.UefiSettings == nil {
			params.Properties.SecurityProfile.UefiSettings = &virtualmachines.UefiSettings{}
		}
		params.Properties.SecurityProfile.UefiSettings.SecureBootEnabled = pointer.To(secureBootEnabled)
		params.Properties.SecurityProfile.UefiSettings.VTpmEnabled = pointer.To(vtpmEnabled)
	} else {
		if secureBootEnabled {
			if params.Properties.SecurityProfile == nil {
				params.Properties.SecurityProfile = &virtualmachines.SecurityProfile{}
			}
			if params.Properties.SecurityProfile.UefiSettings == nil {
				params.Properties.SecurityProfile.UefiSettings = &virtualmachines.UefiSettings{}
			}
			params.Properties.SecurityProfile.SecurityType = pointer.To(virtualmachines.SecurityTypesTrustedLaunch)
			params.Properties.SecurityProfile.UefiSettings.SecureBootEnabled = pointer.To(secureBootEnabled)
		}

		if vtpmEnabled {
			if params.Properties.SecurityProfile == nil {
				params.Properties.SecurityProfile = &virtualmachines.SecurityProfile{}
			}
			if params.Properties.SecurityProfile.UefiSettings == nil {
				params.Properties.SecurityProfile.UefiSettings = &virtualmachines.UefiSettings{}
			}
			params.Properties.SecurityProfile.SecurityType = pointer.To(virtualmachines.SecurityTypesTrustedLaunch)
			params.Properties.SecurityProfile.UefiSettings.VTpmEnabled = pointer.To(vtpmEnabled)
		}
	}

	if evictionPolicyRaw, ok := d.GetOk("eviction_policy"); ok {
		if params.Properties.Priority != nil && *params.Properties.Priority != virtualmachines.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("an `eviction_policy` can only be specified when `priority` is set to `Spot`")
		}

		params.Properties.EvictionPolicy = pointer.ToEnum[virtualmachines.VirtualMachineEvictionPolicyTypes](evictionPolicyRaw.(string))
	} else if priority == virtualmachines.VirtualMachinePriorityTypesSpot {
		return fmt.Errorf("an `eviction_policy` must be specified when `priority` is set to `Spot`")
	}

	if v, ok := d.GetOk("license_type"); ok {
		params.Properties.LicenseType = pointer.To(v.(string))
	}

	if v, ok := d.Get("max_bid_price").(float64); ok && v > 0 {
		if priority != virtualmachines.VirtualMachinePriorityTypesSpot {
			return fmt.Errorf("`max_bid_price` can only be configured when `priority` is set to `Spot`")
		}

		params.Properties.BillingProfile = &virtualmachines.BillingProfile{
			MaxPrice: pointer.To(v),
		}
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		params.Properties.ProximityPlacementGroup = &virtualmachines.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("virtual_machine_scale_set_id"); ok {
		params.Properties.VirtualMachineScaleSet = &virtualmachines.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	platformFaultDomain := d.Get("platform_fault_domain").(int)
	if platformFaultDomain != -1 {
		params.Properties.PlatformFaultDomain = pointer.To(int64(platformFaultDomain))
	}

	var osImageNotificationProfile *virtualmachines.OSImageNotificationProfile
	var terminateNotificationProfile *virtualmachines.TerminateNotificationProfile

	if v, ok := d.GetOk("os_image_notification"); ok {
		osImageNotificationProfile = helpers.ExpandOsImageNotificationProfile(v.([]any))
	}

	if v, ok := d.GetOk("termination_notification"); ok {
		terminateNotificationProfile = helpers.ExpandTerminateNotificationProfile(v.([]any))
	}

	if terminateNotificationProfile != nil || osImageNotificationProfile != nil {
		params.Properties.ScheduledEventsProfile = &virtualmachines.ScheduledEventsProfile{
			OsImageNotificationProfile:   osImageNotificationProfile,
			TerminateNotificationProfile: terminateNotificationProfile,
		}
	}

	if v, ok := d.GetOk("user_data"); ok {
		params.Properties.UserData = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok {
		params.Zones = &[]string{
			v.(string),
		}
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, params, virtualmachines.DefaultCreateOrUpdateOperationOptions(), sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating Windows %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceWindowsVirtualMachineRead(d, meta)
}
