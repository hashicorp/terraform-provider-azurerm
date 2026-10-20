// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_virtual_machine

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWindowsVirtualMachineUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachines.ParseVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualMachineName, helpers.VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, helpers.VirtualMachineResourceName)

	options := virtualmachines.DefaultGetOperationOptions()
	options.Expand = pointer.To(virtualmachines.InstanceViewTypesUserData)
	existing, err := client.Get(ctx, *id, options)
	if err != nil {
		return fmt.Errorf("retrieving Windows %s: %+v", id, err)
	}

	instanceView, err := client.InstanceView(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving InstanceView for Windows %s: %+v", id, err)
	}

	shouldTurnBackOn := helpers.VirtualMachineShouldBeStarted(instanceView.Model)
	hasEphemeralOSDisk := false

	if existing.Model == nil {
		return fmt.Errorf("`model` was nil for %s", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("`properties` was nil for %s", id)
	}

	props := existing.Model.Properties
	if storage := props.StorageProfile; storage != nil {
		if disk := storage.OsDisk; disk != nil {
			if settings := disk.DiffDiskSettings; settings != nil && settings.Option != nil {
				hasEphemeralOSDisk = *settings.Option == virtualmachines.DiffDiskOptionsLocal
			}
		}
	}

	shouldUpdate := false
	shouldShutDown := false
	shouldDeallocate := false

	update := virtualmachines.VirtualMachineUpdate{
		Properties: &virtualmachines.VirtualMachineProperties{},
	}

	if d.HasChange("boot_diagnostics") {
		shouldUpdate = true

		bootDiagnosticsRaw := d.Get("boot_diagnostics").([]any)
		update.Properties.DiagnosticsProfile = helpers.ExpandBootDiagnostics(bootDiagnosticsRaw)
	}

	if d.HasChange("secret") {
		shouldUpdate = true

		profile := virtualmachines.OSProfile{}

		if d.HasChange("secret") {
			secretsRaw := d.Get("secret").([]any)
			profile.Secrets = helpers.ExpandWindowsSecrets(secretsRaw)
		}

		update.Properties.OsProfile = &profile
	}

	if d.HasChange("allow_extension_operations") {
		allowExtensionOperations := d.Get("allow_extension_operations").(bool)

		shouldUpdate = true

		if update.Properties.OsProfile == nil {
			update.Properties.OsProfile = &virtualmachines.OSProfile{}
		}

		update.Properties.OsProfile.AllowExtensionOperations = pointer.To(allowExtensionOperations)
	}

	if d.HasChange("patch_mode") {
		shouldUpdate = true

		if update.Properties.OsProfile == nil {
			update.Properties.OsProfile = &virtualmachines.OSProfile{}
		}

		if update.Properties.OsProfile.WindowsConfiguration == nil {
			update.Properties.OsProfile.WindowsConfiguration = &virtualmachines.WindowsConfiguration{}
		}

		if update.Properties.OsProfile.WindowsConfiguration.PatchSettings == nil {
			update.Properties.OsProfile.WindowsConfiguration.PatchSettings = &virtualmachines.PatchSettings{}
		}

		update.Properties.OsProfile.WindowsConfiguration.PatchSettings.PatchMode = pointer.ToEnum[virtualmachines.WindowsVMGuestPatchMode](d.Get("patch_mode").(string))
	}

	if d.HasChange("patch_assessment_mode") {
		assessmentMode := d.Get("patch_assessment_mode").(string)
		if assessmentMode == string(virtualmachines.WindowsPatchAssessmentModeAutomaticByPlatform) && !d.Get("provision_vm_agent").(bool) {
			return fmt.Errorf("`provision_vm_agent` must be set to `true` when `patch_assessment_mode` is set to `AutomaticByPlatform`")
		}

		shouldUpdate = true

		if update.Properties.OsProfile == nil {
			update.Properties.OsProfile = &virtualmachines.OSProfile{}
		}

		if update.Properties.OsProfile.WindowsConfiguration == nil {
			update.Properties.OsProfile.WindowsConfiguration = &virtualmachines.WindowsConfiguration{}
		}

		if update.Properties.OsProfile.WindowsConfiguration.PatchSettings == nil {
			update.Properties.OsProfile.WindowsConfiguration.PatchSettings = &virtualmachines.PatchSettings{}
		}

		update.Properties.OsProfile.WindowsConfiguration.PatchSettings.AssessmentMode = pointer.ToEnum[virtualmachines.WindowsPatchAssessmentMode](assessmentMode)
	}

	isPatchModeAutomaticByPlatform := d.Get("patch_mode") == string(virtualmachines.WindowsVMGuestPatchModeAutomaticByPlatform)
	bypassPlatformSafetyChecksOnUserScheduleEnabled := d.Get("bypass_platform_safety_checks_on_user_schedule_enabled").(bool)
	if bypassPlatformSafetyChecksOnUserScheduleEnabled && !isPatchModeAutomaticByPlatform {
		return fmt.Errorf("`patch_mode` must be set to `AutomaticByPlatform` when `bypass_platform_safety_checks_on_user_schedule_enabled` is set to `true`")
	}
	if d.HasChange("bypass_platform_safety_checks_on_user_schedule_enabled") {
		shouldUpdate = true

		if update.Properties.OsProfile == nil {
			update.Properties.OsProfile = &virtualmachines.OSProfile{}
		}

		if update.Properties.OsProfile.WindowsConfiguration == nil {
			update.Properties.OsProfile.WindowsConfiguration = &virtualmachines.WindowsConfiguration{}
		}

		if update.Properties.OsProfile.WindowsConfiguration.PatchSettings == nil {
			update.Properties.OsProfile.WindowsConfiguration.PatchSettings = &virtualmachines.PatchSettings{}
		}

		if isPatchModeAutomaticByPlatform {
			if update.Properties.OsProfile.WindowsConfiguration.PatchSettings.AutomaticByPlatformSettings == nil {
				update.Properties.OsProfile.WindowsConfiguration.PatchSettings.AutomaticByPlatformSettings = &virtualmachines.WindowsVMGuestPatchAutomaticByPlatformSettings{}
			}

			update.Properties.OsProfile.WindowsConfiguration.PatchSettings.AutomaticByPlatformSettings.BypassPlatformSafetyChecksOnUserSchedule = pointer.To(bypassPlatformSafetyChecksOnUserScheduleEnabled)
		}
	}

	rebootSetting := d.Get("reboot_setting").(string)
	if rebootSetting != "" && !isPatchModeAutomaticByPlatform {
		return fmt.Errorf("`patch_mode` must be set to `AutomaticByPlatform` when `reboot_setting` is specified")
	}
	if d.HasChange("reboot_setting") {
		shouldUpdate = true

		if update.Properties.OsProfile == nil {
			update.Properties.OsProfile = &virtualmachines.OSProfile{}
		}

		if update.Properties.OsProfile.WindowsConfiguration == nil {
			update.Properties.OsProfile.WindowsConfiguration = &virtualmachines.WindowsConfiguration{}
		}

		if update.Properties.OsProfile.WindowsConfiguration.PatchSettings == nil {
			update.Properties.OsProfile.WindowsConfiguration.PatchSettings = &virtualmachines.PatchSettings{}
		}

		if isPatchModeAutomaticByPlatform {
			if update.Properties.OsProfile.WindowsConfiguration.PatchSettings.AutomaticByPlatformSettings == nil {
				update.Properties.OsProfile.WindowsConfiguration.PatchSettings.AutomaticByPlatformSettings = &virtualmachines.WindowsVMGuestPatchAutomaticByPlatformSettings{}
			}

			update.Properties.OsProfile.WindowsConfiguration.PatchSettings.AutomaticByPlatformSettings.RebootSetting = pointer.ToEnum[virtualmachines.WindowsVMGuestPatchAutomaticByPlatformRebootSetting](rebootSetting)
		}
	}

	if d.HasChange("hotpatching_enabled") {
		shouldUpdate = true

		if update.Properties.OsProfile == nil {
			update.Properties.OsProfile = &virtualmachines.OSProfile{}
		}

		if update.Properties.OsProfile.WindowsConfiguration == nil {
			update.Properties.OsProfile.WindowsConfiguration = &virtualmachines.WindowsConfiguration{}
		}

		if update.Properties.OsProfile.WindowsConfiguration.PatchSettings == nil {
			update.Properties.OsProfile.WindowsConfiguration.PatchSettings = &virtualmachines.PatchSettings{}
		}

		update.Properties.OsProfile.WindowsConfiguration.PatchSettings.EnableHotpatching = pointer.To(d.Get("hotpatching_enabled").(bool))
	}

	if d.HasChange("identity") {
		shouldUpdate = true

		identityExpanded, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		update.Identity = identityExpanded
	}

	if d.HasChange("capacity_reservation_group_id") {
		shouldUpdate = true
		shouldDeallocate = true

		if v, ok := d.GetOk("capacity_reservation_group_id"); ok {
			update.Properties.CapacityReservation = &virtualmachines.CapacityReservationProfile{
				CapacityReservationGroup: &virtualmachines.SubResource{
					Id: pointer.To(v.(string)),
				},
			}
		} else {
			update.Properties.CapacityReservation = &virtualmachines.CapacityReservationProfile{
				CapacityReservationGroup: &virtualmachines.SubResource{},
			}
		}
	}

	if d.HasChange("dedicated_host_id") {
		shouldUpdate = true

		// Code="PropertyChangeNotAllowed" Message="Updating Host of VM 'VMNAME' is not allowed as the VM is currently allocated. Please Deallocate the VM and retry the operation."
		shouldDeallocate = true

		if v, ok := d.GetOk("dedicated_host_id"); ok {
			update.Properties.Host = &virtualmachines.SubResource{
				Id: pointer.To(v.(string)),
			}
		} else {
			update.Properties.Host = &virtualmachines.SubResource{}
		}
	}

	if d.HasChange("dedicated_host_group_id") {
		shouldUpdate = true

		// Code="PropertyChangeNotAllowed" Message="Updating Host of VM 'VMNAME' is not allowed as the VM is currently allocated. Please Deallocate the VM and retry the operation."
		shouldDeallocate = true

		if v, ok := d.GetOk("dedicated_host_group_id"); ok {
			update.Properties.HostGroup = &virtualmachines.SubResource{
				Id: pointer.To(v.(string)),
			}
		} else {
			update.Properties.HostGroup = &virtualmachines.SubResource{}
		}
	}

	if d.HasChange("extensions_time_budget") {
		shouldUpdate = true
		update.Properties.ExtensionsTimeBudget = pointer.To(d.Get("extensions_time_budget").(string))
	}

	if d.HasChange("gallery_application") {
		shouldUpdate = true
		update.Properties.ApplicationProfile = &virtualmachines.ApplicationProfile{
			GalleryApplications: helpers.ExpandVirtualMachineGalleryApplication(d.Get("gallery_application").([]any)),
		}
	}

	if d.HasChange("max_bid_price") {
		shouldUpdate = true

		// Code="OperationNotAllowed" Message="Max price change is not allowed. For more information, see http://aka.ms/AzureSpot/errormessages"
		shouldShutDown = true

		// "code":"OperationNotAllowed"
		// "message": "Max price change is not allowed when the VM [name] is currently allocated.
		//			   Please deallocate and try again. For more information, see http://aka.ms/AzureSpot/errormessages"
		shouldDeallocate = true

		maxBidPrice := d.Get("max_bid_price").(float64)
		update.Properties.BillingProfile = &virtualmachines.BillingProfile{
			MaxPrice: pointer.To(maxBidPrice),
		}
	}

	if d.HasChange("network_interface_ids") {
		shouldUpdate = true

		// Code="CannotAddOrRemoveNetworkInterfacesFromARunningVirtualMachine"
		// Message="Secondary network interfaces cannot be added or removed from a running virtual machine.
		shouldShutDown = true

		// @tombuildsstuff: after testing shutting it down isn't sufficient - we need a full deallocation
		shouldDeallocate = true

		networkInterfaceIdsRaw := d.Get("network_interface_ids").([]any)
		networkInterfaceIds := helpers.ExpandVirtualMachineNetworkInterfaceIDs(networkInterfaceIdsRaw)

		update.Properties.NetworkProfile = &virtualmachines.NetworkProfile{
			NetworkInterfaces: &networkInterfaceIds,
		}
	}

	if d.HasChange("disk_controller_type") {
		shouldUpdate = true
		shouldDeallocate = true

		if update.Properties.StorageProfile == nil {
			update.Properties.StorageProfile = &virtualmachines.StorageProfile{}
		}

		update.Properties.StorageProfile.DiskControllerType = pointer.ToEnum[virtualmachines.DiskControllerTypes](d.Get("disk_controller_type").(string))
	}

	if d.HasChange("os_disk") {
		shouldUpdate = true

		// Code="Conflict" Message="Disk resizing is allowed only when creating a VM or when the VM is deallocated." Target="disk.diskSizeGB"
		shouldShutDown = true
		shouldDeallocate = true

		osDiskRaw := d.Get("os_disk").([]any)
		osDisk, err := helpers.ExpandVirtualMachineOSDisk(osDiskRaw, virtualmachines.OperatingSystemTypesWindows)
		if err != nil {
			return fmt.Errorf("expanding `os_disk`: %+v", err)
		}

		if v, _ := pluginsdk.GoValueFromTerraformValue[string](d.GetRawConfig().AsValueMap()["os_managed_disk_id"]); pointer.From(v) != "" {
			osDisk.CreateOption = virtualmachines.DiskCreateOptionTypesAttach
		}

		if update.Properties.StorageProfile == nil {
			update.Properties.StorageProfile = &virtualmachines.StorageProfile{}
		}

		update.Properties.StorageProfile.OsDisk = osDisk
	}

	if d.HasChange("virtual_machine_scale_set_id") {
		shouldUpdate = true

		if vmssIDRaw, ok := d.GetOk("virtual_machine_scale_set_id"); ok {
			update.Properties.VirtualMachineScaleSet = &virtualmachines.SubResource{
				Id: pointer.To(vmssIDRaw.(string)),
			}
		} else {
			update.Properties.VirtualMachineScaleSet = &virtualmachines.SubResource{}
		}
	}

	if d.HasChange("proximity_placement_group_id") {
		shouldUpdate = true

		// Code="OperationNotAllowed" Message="Updating proximity placement group of VM is not allowed while the VM is running. Please stop/deallocate the VM and retry the operation."
		shouldShutDown = true
		shouldDeallocate = true

		if ppgIDRaw, ok := d.GetOk("proximity_placement_group_id"); ok {
			update.Properties.ProximityPlacementGroup = &virtualmachines.SubResource{
				Id: pointer.To(ppgIDRaw.(string)),
			}
		} else {
			update.Properties.ProximityPlacementGroup = &virtualmachines.SubResource{}
		}
	}

	if d.HasChange("size") {
		shouldUpdate = true

		// this is kind of superflurious since Azure can do this for us, but if we do this we can subsequently
		// deallocate the VM to switch hosts if required
		shouldShutDown = true
		vmSize := d.Get("size").(string)

		// Azure will auto-reboot this for us, providing this machine will fit on this host
		// otherwise we need to shut down the VM to move it to another host to be able to use this size
		availableOnThisHost := false
		sizes, err := client.ListAvailableSizes(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving available sizes for Windows %s: %+v", id, err)
		}

		if sizes.Model != nil && sizes.Model.Value != nil {
			for _, size := range *sizes.Model.Value {
				if size.Name == nil {
					continue
				}

				if strings.EqualFold(*size.Name, vmSize) {
					availableOnThisHost = true
					break
				}
			}
		}

		if !availableOnThisHost {
			log.Printf("[DEBUG] Requested VM Size isn't available on the Host - must switch host to resize..")
			// Code="OperationNotAllowed"
			// Message="Unable to resize the VM [name] because the requested size Standard_F4s_v2 is not available in the current hardware cluster.
			//         The available sizes in this cluster are: [list]. The requested size might be available in other clusters of this region.
			//         Read more on VM resizing strategy at https://aka.ms/azure-resizevm."
			shouldDeallocate = true
		}

		update.Properties.HardwareProfile = &virtualmachines.HardwareProfile{
			VMSize: pointer.ToEnum[virtualmachines.VirtualMachineSizeTypes](vmSize),
		}
	}

	if d.HasChange("tags") {
		shouldUpdate = true

		tagsRaw := d.Get("tags").(map[string]any)
		update.Tags = tags.Expand(tagsRaw)
	}

	var osImageNotificationProfile *virtualmachines.OSImageNotificationProfile
	var terminateNotificationProfile *virtualmachines.TerminateNotificationProfile

	if d.HasChange("os_image_notification") {
		shouldUpdate = true
		osImageNotificationProfile = helpers.ExpandOsImageNotificationProfile(d.Get("os_image_notification").([]any))
	}

	if d.HasChange("termination_notification") {
		shouldUpdate = true
		terminateNotificationProfile = helpers.ExpandTerminateNotificationProfile(d.Get("termination_notification").([]any))
	}

	if osImageNotificationProfile != nil || terminateNotificationProfile != nil {
		update.Properties.ScheduledEventsProfile = &virtualmachines.ScheduledEventsProfile{
			OsImageNotificationProfile:   osImageNotificationProfile,
			TerminateNotificationProfile: terminateNotificationProfile,
		}
	}

	if d.HasChange("additional_capabilities") {
		shouldUpdate = true

		n, _ := d.GetChange("additional_capabilities")
		if len(n.([]any)) == 0 || d.HasChange("additional_capabilities.0.ultra_ssd_enabled") {
			shouldShutDown = true
			shouldDeallocate = true
		}

		additionalCapabilitiesRaw := d.Get("additional_capabilities").([]any)
		update.Properties.AdditionalCapabilities = helpers.ExpandVirtualMachineAdditionalCapabilities(additionalCapabilitiesRaw)
	}

	if d.HasChange("encryption_at_host_enabled") {
		shouldUpdate = true
		shouldDeallocate = true // API returns the following error if not deallocate: 'securityProfile.encryptionAtHost' can be updated only when VM is in deallocated state
		if update.Properties.SecurityProfile == nil {
			update.Properties.SecurityProfile = &virtualmachines.SecurityProfile{}
		}
		update.Properties.SecurityProfile.EncryptionAtHost = pointer.To(d.Get("encryption_at_host_enabled").(bool))
	}

	if d.HasChange("license_type") {
		shouldUpdate = true

		license := d.Get("license_type").(string)
		if license == "" {
			// Only for create no specification is possible in the API. API does not allow empty string in update.
			// So removing attribute license_type from Terraform configuration if it was set to value other than 'None' would lead to an endless loop in apply.
			// To allow updating in this case set value explicitly to 'None'.
			license = "None"
		}
		update.Properties.LicenseType = &license
	}

	if d.HasChange("user_data") {
		shouldUpdate = true
		update.Properties.UserData = pointer.To(d.Get("user_data").(string))
	}

	if instanceView.Model != nil && instanceView.Model.Statuses != nil {
		for _, status := range *instanceView.Model.Statuses {
			if status.Code == nil {
				continue
			}

			// could also be the provisioning state which we're not bothered with here
			state := strings.ToLower(*status.Code)
			if !strings.HasPrefix(state, "powerstate/") {
				continue
			}

			state = strings.TrimPrefix(state, "powerstate/")
			switch strings.ToLower(state) {
			case "deallocated":
				// VM already deallocated, no shutdown and deallocation needed anymore
				shouldShutDown = false
				shouldDeallocate = false
			case "deallocating":
				// VM is deallocating
				// To make sure we do not start updating before this action has finished,
				// only skip the shutdown and send another deallocation request if shouldDeallocate == true
				shouldShutDown = false
			case "stopped":
				// VM already stopped, no shutdown needed anymore
				shouldShutDown = false
			}
		}
	}

	if shouldShutDown {
		log.Printf("[DEBUG] Shutting Down Windows %s", id)
		if err := client.PowerOffThenPoll(ctx, *id, virtualmachines.DefaultPowerOffOperationOptions()); err != nil {
			return fmt.Errorf("sending Power Off to Windows %s: %+v", id, err)
		}

		log.Printf("[DEBUG] Shut Down Windows %s", id)
	}

	if shouldDeallocate {
		if !hasEphemeralOSDisk {
			log.Printf("[DEBUG] Deallocating Windows %s", id)
			// Upgrading to the 2021-07-01 exposed a new hibernate parameter in the Deallocate method
			if err := client.DeallocateThenPoll(ctx, *id, virtualmachines.DefaultDeallocateOperationOptions()); err != nil {
				return fmt.Errorf("deallocating Windows %s: %+v", id, err)
			}
			log.Printf("[DEBUG] Deallocated Windows %s", id)
		} else {
			// Code="OperationNotAllowed" Message="Operation 'deallocate' is not supported for VMs or VM Scale Set instances using an ephemeral OS disk."
			log.Printf("[DEBUG] Skipping deallocation for Windows %s since cannot deallocate a Virtual Machine with an Ephemeral OS Disk", id)
		}
	}

	// now the VM's shutdown/deallocated we can update the disk which can't be done via the VM API:
	// Code="ResizeDiskError" Message="Managed disk resize via Virtual Machine [name] is not allowed. Please resize disk resource at [id]."
	// Portal: "Disks can be resized or account type changed only when they are unattached or the owner VM is deallocated."
	if d.HasChange("os_disk.0.disk_size_gb") {
		diskName := d.Get("os_disk.0.name").(string)
		newSize := d.Get("os_disk.0.disk_size_gb").(int)
		log.Printf("[DEBUG] Resizing OS Disk %q for Windows %s to %dGB..", diskName, id, newSize)

		disksClient := meta.(*clients.Client).Compute.DisksClient
		subscriptionId := meta.(*clients.Client).Account.SubscriptionId
		id := commonids.NewManagedDiskID(subscriptionId, id.ResourceGroupName, diskName)

		update := disks.DiskUpdate{
			Properties: &disks.DiskUpdateProperties{
				DiskSizeGB: pointer.To(int64(newSize)),
			},
		}

		if err := disksClient.UpdateThenPoll(ctx, id, update); err != nil {
			return fmt.Errorf("resizing OS Disk %q for Windows Virtual Machine %q (Resource Group %q): %+v", diskName, id.DiskName, id.ResourceGroupName, err)
		}

		log.Printf("[DEBUG] Resized OS Disk %q for Windows Virtual Machine %q (Resource Group %q) to %dGB.", diskName, id.DiskName, id.ResourceGroupName, newSize)
	}

	if d.HasChange("os_disk.0.disk_encryption_set_id") {
		if diskEncryptionSetId := d.Get("os_disk.0.disk_encryption_set_id").(string); diskEncryptionSetId != "" {
			diskName := d.Get("os_disk.0.name").(string)

			encryptionType, err := helpers.RetrieveDiskEncryptionSetEncryptionType(ctx, meta.(*clients.Client).Compute.DiskEncryptionSetsClient, diskEncryptionSetId)
			if err != nil {
				return err
			}

			disksClient := meta.(*clients.Client).Compute.DisksClient
			subscriptionId := meta.(*clients.Client).Account.SubscriptionId
			id := commonids.NewManagedDiskID(subscriptionId, id.ResourceGroupName, diskName)

			update := disks.DiskUpdate{
				Properties: &disks.DiskUpdateProperties{
					Encryption: &disks.Encryption{
						Type:                encryptionType,
						DiskEncryptionSetId: pointer.To(diskEncryptionSetId),
					},
				},
			}

			if err = disksClient.UpdateThenPoll(ctx, id, update); err != nil {
				return fmt.Errorf("updating encryption settings of OS Disk %q for Windows Virtual Machine %q (Resource Group %q): %+v", diskName, id.DiskName, id.ResourceGroupName, err)
			}
		} else {
			return fmt.Errorf("once a customer-managed key is used, you can’t change the selection back to a platform-managed key")
		}
	}

	if shouldUpdate {
		if err := client.UpdateThenPoll(ctx, *id, update, virtualmachines.DefaultUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating Windows %s: %+v", id, err)
		}
	}

	// if we've shut it down and it was turned off, let's boot it back up
	if shouldTurnBackOn && (shouldShutDown || shouldDeallocate) {
		log.Printf("[DEBUG] Starting Windows %s", id)
		if err := client.StartThenPoll(ctx, *id); err != nil {
			return fmt.Errorf("starting Windows %s: %+v", id, err)
		}
		log.Printf("[DEBUG] Started Windows %s", id)
	}

	return resourceWindowsVirtualMachineRead(d, meta)
}
