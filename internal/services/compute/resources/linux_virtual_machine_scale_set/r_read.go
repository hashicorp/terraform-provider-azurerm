// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package linux_virtual_machine_scale_set

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLinuxVirtualMachineScaleSetRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachineScaleSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachinescalesets.ParseVirtualMachineScaleSetID(d.Id())
	if err != nil {
		return err
	}

	options := virtualmachinescalesets.DefaultGetOperationOptions()
	options.Expand = pointer.To(virtualmachinescalesets.ExpandTypesForGetVMScaleSetsUserData)
	resp, err := client.Get(ctx, *id, options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Linux %s - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Linux %s: %+v", id, err)
	}

	d.Set("name", id.VirtualMachineScaleSetName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("edge_zone", helpers.FlattenEdgeZone(model.ExtendedLocation))
		d.Set("zones", zones.FlattenUntyped(model.Zones))
		var skuName *string
		var instances int
		if model.Sku != nil {
			skuName = model.Sku.Name
			if model.Sku.Capacity != nil {
				instances = int(*model.Sku.Capacity)
			}
		}
		d.Set("instances", instances)
		d.Set("sku", skuName)

		identityFlattened, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return err
		}
		if err := d.Set("identity", identityFlattened); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := d.Set("plan", helpers.FlattenPlanVMSS(model.Plan)); err != nil {
			return fmt.Errorf("setting `plan`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if err := d.Set("additional_capabilities", helpers.FlattenVirtualMachineScaleSetAdditionalCapabilities(props.AdditionalCapabilities)); err != nil {
				return fmt.Errorf("setting `additional_capabilities`: %+v", props.AdditionalCapabilities)
			}

			if err := d.Set("automatic_instance_repair", helpers.FlattenVirtualMachineScaleSetAutomaticRepairsPolicy(props.AutomaticRepairsPolicy)); err != nil {
				return fmt.Errorf("setting `automatic_instance_repair`: %+v", err)
			}

			d.Set("do_not_run_extensions_on_overprovisioned_machines", props.DoNotRunExtensionsOnOverprovisionedVMs)
			if props.HostGroup != nil && props.HostGroup.Id != nil {
				d.Set("host_group_id", props.HostGroup.Id)
			}
			d.Set("overprovision", props.Overprovision)
			proximityPlacementGroupId := ""
			if props.ProximityPlacementGroup != nil && props.ProximityPlacementGroup.Id != nil {
				proximityPlacementGroupId = *props.ProximityPlacementGroup.Id
			}
			d.Set("platform_fault_domain_count", props.PlatformFaultDomainCount)
			d.Set("proximity_placement_group_id", proximityPlacementGroupId)
			d.Set("single_placement_group", props.SinglePlacementGroup)
			d.Set("unique_id", props.UniqueId)
			d.Set("zone_balance", props.ZoneBalance)
			d.Set("scale_in", helpers.FlattenVirtualMachineScaleSetScaleInPolicy(props.ScaleInPolicy))

			if props.SpotRestorePolicy != nil {
				d.Set("spot_restore", helpers.FlattenVirtualMachineScaleSetSpotRestorePolicy(props.SpotRestorePolicy))
			}

			resilientVMCreationEnabled, resilientVMDeletionEnabled := helpers.FlattenVirtualMachineScaleSetResiliency(props.ResiliencyPolicy)
			d.Set("resilient_vm_creation_enabled", resilientVMCreationEnabled)
			d.Set("resilient_vm_deletion_enabled", resilientVMDeletionEnabled)

			if profile := props.VirtualMachineProfile; profile != nil {
				if err := d.Set("boot_diagnostics", helpers.FlattenBootDiagnosticsVMSS(profile.DiagnosticsProfile)); err != nil {
					return fmt.Errorf("setting `boot_diagnostics`: %+v", err)
				}

				capacityReservationGroupId := ""
				if profile.CapacityReservation != nil && profile.CapacityReservation.CapacityReservationGroup != nil && profile.CapacityReservation.CapacityReservationGroup.Id != nil {
					capacityReservationGroupId = *profile.CapacityReservation.CapacityReservationGroup.Id
				}
				d.Set("capacity_reservation_group_id", capacityReservationGroupId)

				// defaulted since BillingProfile isn't returned if it's unset
				maxBidPrice := float64(-1.0)
				if profile.BillingProfile != nil && profile.BillingProfile.MaxPrice != nil {
					maxBidPrice = *profile.BillingProfile.MaxPrice
				}
				d.Set("max_bid_price", maxBidPrice)

				d.Set("eviction_policy", pointer.FromEnum(profile.EvictionPolicy))

				if profile.ApplicationProfile != nil && profile.ApplicationProfile.GalleryApplications != nil {
					d.Set("gallery_application", helpers.FlattenVirtualMachineScaleSetGalleryApplication(profile.ApplicationProfile.GalleryApplications))
				}

				// the service just return empty when this is not assigned when provisioned
				// See discussion on https://github.com/Azure/azure-rest-api-specs/issues/10971
				priority := virtualmachinescalesets.VirtualMachinePriorityTypesRegular
				if pointer.From(profile.Priority) != "" {
					priority = pointer.From(profile.Priority)
				}
				d.Set("priority", priority)

				if storageProfile := profile.StorageProfile; storageProfile != nil {
					if err := d.Set("os_disk", helpers.FlattenVirtualMachineScaleSetOSDisk(storageProfile.OsDisk)); err != nil {
						return fmt.Errorf("setting `os_disk`: %+v", err)
					}

					if err := d.Set("data_disk", helpers.FlattenVirtualMachineScaleSetDataDisk(storageProfile.DataDisks)); err != nil {
						return fmt.Errorf("setting `data_disk`: %+v", err)
					}

					var storageImageId string
					if storageProfile.ImageReference != nil && storageProfile.ImageReference.Id != nil {
						storageImageId = *storageProfile.ImageReference.Id
					}
					if storageProfile.ImageReference != nil && storageProfile.ImageReference.CommunityGalleryImageId != nil {
						storageImageId = *storageProfile.ImageReference.CommunityGalleryImageId
					}
					if storageProfile.ImageReference != nil && storageProfile.ImageReference.SharedGalleryImageId != nil {
						storageImageId = *storageProfile.ImageReference.SharedGalleryImageId
					}
					d.Set("source_image_id", storageImageId)

					if err := d.Set("source_image_reference", helpers.FlattenSourceImageReferenceVMSS(storageProfile.ImageReference, storageImageId != "")); err != nil {
						return fmt.Errorf("setting `source_image_reference`: %+v", err)
					}
				}

				extensionOperationsEnabled := true
				if osProfile := profile.OsProfile; osProfile != nil {
					// admin_password isn't returned, but it's a top level field so we can ignore it without consequence
					d.Set("admin_username", osProfile.AdminUsername)
					d.Set("computer_name_prefix", osProfile.ComputerNamePrefix)

					if osProfile.AllowExtensionOperations != nil {
						extensionOperationsEnabled = *osProfile.AllowExtensionOperations
					}

					if linux := osProfile.LinuxConfiguration; linux != nil {
						d.Set("disable_password_authentication", linux.DisablePasswordAuthentication)
						d.Set("provision_vm_agent", linux.ProvisionVMAgent)

						flattenedSshKeys, err := helpers.FlattenSSHKeysVMSS(linux.Ssh)
						if err != nil {
							return fmt.Errorf("flattening `admin_ssh_key`: %+v", err)
						}
						if err := d.Set("admin_ssh_key", pluginsdk.NewSet(helpers.SSHKeySchemaHash, *flattenedSshKeys)); err != nil {
							return fmt.Errorf("setting `admin_ssh_key`: %+v", err)
						}
					}

					if err := d.Set("secret", helpers.FlattenLinuxSecretsVMSS(osProfile.Secrets)); err != nil {
						return fmt.Errorf("setting `secret`: %+v", err)
					}
				}
				d.Set("extension_operations_enabled", extensionOperationsEnabled)

				if nwProfile := profile.NetworkProfile; nwProfile != nil {
					if err := d.Set("network_interface", helpers.FlattenVirtualMachineScaleSetNetworkInterface(nwProfile.NetworkInterfaceConfigurations)); err != nil {
						return fmt.Errorf("setting `network_interface`: %+v", err)
					}

					healthProbeId := ""
					if nwProfile.HealthProbe != nil && nwProfile.HealthProbe.Id != nil {
						healthProbeId = *nwProfile.HealthProbe.Id
					}
					d.Set("health_probe_id", healthProbeId)
				}

				if scheduleProfile := profile.ScheduledEventsProfile; scheduleProfile != nil {
					if err := d.Set("termination_notification", helpers.FlattenVirtualMachineScaleSetScheduledEventsProfile(scheduleProfile)); err != nil {
						return fmt.Errorf("setting `termination_notification`: %+v", err)
					}
				}

				extensionProfile, err := helpers.FlattenVirtualMachineScaleSetExtensions(profile.ExtensionProfile, d)
				if err != nil {
					return fmt.Errorf("failed flattening `extension`: %+v", err)
				}
				d.Set("extension", extensionProfile)

				extensionsTimeBudget := "PT1H30M"
				if profile.ExtensionProfile != nil && profile.ExtensionProfile.ExtensionsTimeBudget != nil {
					extensionsTimeBudget = *profile.ExtensionProfile.ExtensionsTimeBudget
				}
				d.Set("extensions_time_budget", extensionsTimeBudget)

				encryptionAtHostEnabled := false
				vtpmEnabled := false
				secureBootEnabled := false

				if secprofile := profile.SecurityProfile; secprofile != nil {
					if secprofile.EncryptionAtHost != nil {
						encryptionAtHostEnabled = *secprofile.EncryptionAtHost
					}
					if uefi := profile.SecurityProfile.UefiSettings; uefi != nil {
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
				d.Set("user_data", profile.UserData)
			}

			if policy := props.UpgradePolicy; policy != nil {
				d.Set("upgrade_mode", pointer.FromEnum(policy.Mode))

				if err := d.Set("automatic_os_upgrade_policy", helpers.FlattenVirtualMachineScaleSetAutomaticOSUpgradePolicy(policy.AutomaticOSUpgradePolicy)); err != nil {
					return fmt.Errorf("setting `automatic_os_upgrade_policy`: %+v", err)
				}

				if err := d.Set("rolling_upgrade_policy", helpers.FlattenVirtualMachineScaleSetRollingUpgradePolicy(policy.RollingUpgradePolicy)); err != nil {
					return fmt.Errorf("setting `rolling_upgrade_policy`: %+v", err)
				}
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}
