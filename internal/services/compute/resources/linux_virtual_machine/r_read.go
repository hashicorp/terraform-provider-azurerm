// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package linux_virtual_machine

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLinuxVirtualMachineRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachines.ParseVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	options := virtualmachines.DefaultGetOperationOptions()
	options.Expand = pointer.To(virtualmachines.InstanceViewTypesUserData)
	resp, err := client.Get(ctx, *id, options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Linux %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Linux %s: %+v", id, err)
	}

	return resourceLinuxVirtualMachineFlatten(ctx, meta.(*clients.Client), d, id, resp.Model, true)
}

func resourceLinuxVirtualMachineFlatten(ctx context.Context, clientsClient *clients.Client, d *pluginsdk.ResourceData, id *virtualmachines.VirtualMachineId, model *virtualmachines.VirtualMachine, includeResource bool) error {
	disksClient := clientsClient.Compute.DisksClient
	networkInterfacesClient := clientsClient.Network.NetworkInterfacesClient
	publicIPAddressesClient := clientsClient.Network.PublicIPAddresses

	d.Set("name", id.VirtualMachineName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("edge_zone", helpers.FlattenEdgeZone(model.ExtendedLocation))

		zone := ""
		if model.Zones != nil {
			if zones := *model.Zones; len(zones) > 0 {
				zone = zones[0]
			}
		}
		d.Set("zone", zone)

		identityFlattened, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identityFlattened); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
		if err := d.Set("plan", helpers.FlattenPlan(model.Plan)); err != nil {
			return fmt.Errorf("setting `plan`: %+v", err)
		}
		if props := model.Properties; props != nil {
			if err := d.Set("additional_capabilities", helpers.FlattenVirtualMachineAdditionalCapabilities(props.AdditionalCapabilities)); err != nil {
				return fmt.Errorf("setting `additional_capabilities`: %+v", err)
			}

			availabilitySetId := ""
			if props.AvailabilitySet != nil && props.AvailabilitySet.Id != nil {
				availabilitySetId = *props.AvailabilitySet.Id
			}
			d.Set("availability_set_id", availabilitySetId)

			capacityReservationGroupId := ""
			if props.CapacityReservation != nil && props.CapacityReservation.CapacityReservationGroup != nil && props.CapacityReservation.CapacityReservationGroup.Id != nil {
				capacityReservationGroupId = *props.CapacityReservation.CapacityReservationGroup.Id
			}
			d.Set("capacity_reservation_group_id", capacityReservationGroupId)

			if props.ApplicationProfile != nil && props.ApplicationProfile.GalleryApplications != nil {
				d.Set("gallery_application", helpers.FlattenVirtualMachineGalleryApplication(props.ApplicationProfile.GalleryApplications))
			}

			licenseType := ""
			if v := props.LicenseType; v != nil && *v != "None" {
				licenseType = *v
			}
			d.Set("license_type", licenseType)

			if err := d.Set("boot_diagnostics", helpers.FlattenBootDiagnostics(props.DiagnosticsProfile)); err != nil {
				return fmt.Errorf("setting `boot_diagnostics`: %+v", err)
			}

			d.Set("eviction_policy", pointer.From(props.EvictionPolicy))
			if profile := props.HardwareProfile; profile != nil {
				d.Set("size", pointer.From(profile.VMSize))
			}

			extensionsTimeBudget := "PT1H30M"
			if props.ExtensionsTimeBudget != nil {
				extensionsTimeBudget = *props.ExtensionsTimeBudget
			}
			d.Set("extensions_time_budget", extensionsTimeBudget)

			// defaulted since BillingProfile isn't returned if it's unset
			maxBidPrice := float64(-1.0)
			if props.BillingProfile != nil && props.BillingProfile.MaxPrice != nil {
				maxBidPrice = *props.BillingProfile.MaxPrice
			}
			d.Set("max_bid_price", maxBidPrice)

			if profile := props.NetworkProfile; profile != nil {
				if err := d.Set("network_interface_ids", helpers.FlattenVirtualMachineNetworkInterfaceIDs(props.NetworkProfile.NetworkInterfaces)); err != nil {
					return fmt.Errorf("setting `network_interface_ids`: %+v", err)
				}
			}

			dedicatedHostId := ""
			if props.Host != nil && props.Host.Id != nil {
				dedicatedHostId = *props.Host.Id
			}
			d.Set("dedicated_host_id", dedicatedHostId)

			dedicatedHostGroupId := ""
			if props.HostGroup != nil && props.HostGroup.Id != nil {
				dedicatedHostGroupId = *props.HostGroup.Id
			}
			d.Set("dedicated_host_group_id", dedicatedHostGroupId)

			virtualMachineScaleSetId := ""
			if props.VirtualMachineScaleSet != nil && props.VirtualMachineScaleSet.Id != nil {
				virtualMachineScaleSetId = *props.VirtualMachineScaleSet.Id
			}
			d.Set("virtual_machine_scale_set_id", virtualMachineScaleSetId)
			platformFaultDomain := -1
			if props.PlatformFaultDomain != nil {
				platformFaultDomain = int(*props.PlatformFaultDomain)
			}
			d.Set("platform_fault_domain", platformFaultDomain)

			patchMode := string(virtualmachines.LinuxVMGuestPatchModeImageDefault)
			assessmentMode := string(virtualmachines.LinuxPatchAssessmentModeImageDefault)
			bypassPlatformSafetyChecksOnUserScheduleEnabled := false
			rebootSetting := ""

			if profile := props.OsProfile; profile != nil {
				d.Set("admin_username", profile.AdminUsername)
				d.Set("allow_extension_operations", profile.AllowExtensionOperations)
				d.Set("computer_name", profile.ComputerName)

				if config := profile.LinuxConfiguration; config != nil {
					d.Set("disable_password_authentication", config.DisablePasswordAuthentication)
					d.Set("provision_vm_agent", config.ProvisionVMAgent)
					d.Set("vm_agent_platform_updates_enabled", config.EnableVMAgentPlatformUpdates)

					flattenedSSHKeys, err := helpers.FlattenSSHKeys(config.Ssh)
					if err != nil {
						return fmt.Errorf("flattening `admin_ssh_key`: %+v", err)
					}
					if err := d.Set("admin_ssh_key", pluginsdk.NewSet(helpers.SSHKeySchemaHash, *flattenedSSHKeys)); err != nil {
						return fmt.Errorf("setting `admin_ssh_key`: %+v", err)
					}

					if patchSettings := config.PatchSettings; patchSettings != nil {
						if patchSettings.PatchMode != nil {
							patchMode = string(*patchSettings.PatchMode)
						}
						if patchSettings.AssessmentMode != nil {
							assessmentMode = string(*patchSettings.AssessmentMode)
						}
						if patchSettings.AutomaticByPlatformSettings != nil {
							bypassPlatformSafetyChecksOnUserScheduleEnabled = pointer.From(patchSettings.AutomaticByPlatformSettings.BypassPlatformSafetyChecksOnUserSchedule)
							rebootSetting = pointer.FromEnum(patchSettings.AutomaticByPlatformSettings.RebootSetting)
						}
					}
				}

				if err := d.Set("secret", helpers.FlattenLinuxSecrets(profile.Secrets)); err != nil {
					return fmt.Errorf("setting `secret`: %+v", err)
				}
			}

			d.Set("patch_mode", patchMode)
			d.Set("patch_assessment_mode", assessmentMode)
			d.Set("bypass_platform_safety_checks_on_user_schedule_enabled", bypassPlatformSafetyChecksOnUserScheduleEnabled)
			d.Set("reboot_setting", rebootSetting)

			// Resources created with azurerm_virtual_machine have priority set to ""
			// We need to treat "" as equal to "Regular" to allow migration azurerm_virtual_machine -> azurerm_linux_virtual_machine
			priority := string(virtualmachines.VirtualMachinePriorityTypesRegular)
			if props.Priority != nil && *props.Priority != "" {
				priority = string(*props.Priority)
			}
			d.Set("priority", priority)
			proximityPlacementGroupId := ""
			if props.ProximityPlacementGroup != nil && props.ProximityPlacementGroup.Id != nil {
				proximityPlacementGroupId = *props.ProximityPlacementGroup.Id
			}
			d.Set("proximity_placement_group_id", proximityPlacementGroupId)

			if profile := props.StorageProfile; profile != nil {
				d.Set("disk_controller_type", pointer.FromEnum(props.StorageProfile.DiskControllerType))

				if includeResource {
					// the storage_account_type isn't returned so we need to look it up
					flattenedOSDisk, err := helpers.FlattenVirtualMachineOSDisk(ctx, disksClient, profile.OsDisk)
					if err != nil {
						return fmt.Errorf("flattening `os_disk`: %+v", err)
					}
					if err := d.Set("os_disk", flattenedOSDisk); err != nil {
						return fmt.Errorf("settings `os_disk`: %+v", err)
					}
				}
				osManagedDiskId := ""
				if profile.OsDisk != nil && profile.OsDisk.ManagedDisk != nil && profile.OsDisk.ManagedDisk.Id != nil {
					osDiskId, err := commonids.ParseManagedDiskIDInsensitively(*profile.OsDisk.ManagedDisk.Id)
					if err != nil {
						return err
					}
					osManagedDiskId = osDiskId.ID()
				}
				d.Set("os_managed_disk_id", osManagedDiskId)
				var storageImageId string
				if profile.ImageReference != nil && profile.ImageReference.Id != nil {
					storageImageId = *profile.ImageReference.Id
				}
				if profile.ImageReference != nil && profile.ImageReference.CommunityGalleryImageId != nil {
					storageImageId = *profile.ImageReference.CommunityGalleryImageId
				}
				if profile.ImageReference != nil && profile.ImageReference.SharedGalleryImageId != nil {
					storageImageId = *profile.ImageReference.SharedGalleryImageId
				}

				d.Set("source_image_id", storageImageId)

				if err := d.Set("source_image_reference", helpers.FlattenSourceImageReference(profile.ImageReference, storageImageId != "")); err != nil {
					return fmt.Errorf("setting `source_image_reference`: %+v", err)
				}
			}

			if scheduleProfile := props.ScheduledEventsProfile; scheduleProfile != nil {
				if err := d.Set("os_image_notification", helpers.FlattenOsImageNotificationProfile(scheduleProfile.OsImageNotificationProfile)); err != nil {
					return fmt.Errorf("setting `termination_notification`: %+v", err)
				}

				if err := d.Set("termination_notification", helpers.FlattenTerminateNotificationProfile(scheduleProfile.TerminateNotificationProfile)); err != nil {
					return fmt.Errorf("setting `termination_notification`: %+v", err)
				}
			}

			encryptionAtHostEnabled := false
			vtpmEnabled := false
			secureBootEnabled := false

			if secprofile := props.SecurityProfile; secprofile != nil {
				if secprofile.EncryptionAtHost != nil {
					encryptionAtHostEnabled = *secprofile.EncryptionAtHost
				}
				if uefi := props.SecurityProfile.UefiSettings; uefi != nil {
					if uefi.VTpmEnabled != nil {
						vtpmEnabled = *uefi.VTpmEnabled
					}
					if uefi.SecureBootEnabled != nil {
						secureBootEnabled = *uefi.SecureBootEnabled
					}
				}
			}

			d.Set("encryption_at_host_enabled", encryptionAtHostEnabled)
			d.Set("vtpm_enabled", vtpmEnabled)
			d.Set("secure_boot_enabled", secureBootEnabled)
			d.Set("virtual_machine_id", props.VMId)
			d.Set("user_data", props.UserData)

			if includeResource {
				connectionInfo := helpers.RetrieveConnectionInformation(ctx, networkInterfacesClient, publicIPAddressesClient, props)
				d.Set("private_ip_address", connectionInfo.PrimaryPrivateAddress)
				d.Set("private_ip_addresses", connectionInfo.PrivateAddresses)
				d.Set("public_ip_address", connectionInfo.PrimaryPublicAddress)
				d.Set("public_ip_addresses", connectionInfo.PublicAddresses)
				helpers.SetConnectionInformation(d, connectionInfo, false)
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}
