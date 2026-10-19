// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package orchestrated_virtual_machine_scale_set

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceOrchestratedVirtualMachineScaleSetRead(d *pluginsdk.ResourceData, meta any) error {
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
			log.Printf("[DEBUG] Orchestrated %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Orchestrated %s: %w", id, err)
	}

	d.Set("name", id.VirtualMachineScaleSetName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		var skuName *string
		var instances int
		if model.Sku != nil {
			skuName, err = flattenOrchestratedVirtualMachineScaleSetSku(model.Sku)
			if err != nil || skuName == nil {
				return fmt.Errorf("setting `sku_name`: %w", err)
			}

			if resp.Model.Sku.Capacity != nil {
				instances = int(*resp.Model.Sku.Capacity)
			}

			d.Set("sku_name", skuName)
			d.Set("instances", instances)
		}

		identityFlattened, err := helpers.FlattenOrchestratedVirtualMachineScaleSetIdentity(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %w", err)
		}
		if err := d.Set("identity", identityFlattened); err != nil {
			return fmt.Errorf("setting `identity`: %w", err)
		}

		if err := d.Set("plan", helpers.FlattenPlanVMSS(model.Plan)); err != nil {
			return fmt.Errorf("setting `plan`: %w", err)
		}

		if props := model.Properties; props != nil {
			if err := d.Set("additional_capabilities", helpers.FlattenOrchestratedVirtualMachineScaleSetAdditionalCapabilities(props.AdditionalCapabilities)); err != nil {
				return fmt.Errorf("setting `additional_capabilities`: %+v", props.AdditionalCapabilities)
			}

			if err := d.Set("automatic_instance_repair", helpers.FlattenVirtualMachineScaleSetAutomaticRepairsPolicy(props.AutomaticRepairsPolicy)); err != nil {
				return fmt.Errorf("setting `automatic_instance_repair`: %w", err)
			}

			d.Set("platform_fault_domain_count", props.PlatformFaultDomainCount)
			proximityPlacementGroupId := ""
			if props.ProximityPlacementGroup != nil && props.ProximityPlacementGroup.Id != nil {
				proximityPlacementGroupId = *props.ProximityPlacementGroup.Id
			}
			d.Set("proximity_placement_group_id", proximityPlacementGroupId)

			// only write state for single_placement_group if it is returned by the RP...
			if props.SinglePlacementGroup != nil {
				d.Set("single_placement_group", props.SinglePlacementGroup)
			}

			if err := d.Set("sku_profile", flattenOrchestratedVirtualMachineScaleSetSkuProfile(props.SkuProfile)); err != nil {
				return fmt.Errorf("setting `sku_profile`: %w", err)
			}

			d.Set("unique_id", props.UniqueId)
			d.Set("zone_balance", props.ZoneBalance)

			extensionOperationsEnabled := true
			// if `VirtualMachineProfile` is nil, `UpgradeMode` will not exist in the response
			upgradeMode := string(virtualmachinescalesets.UpgradeModeManual)
			if profile := props.VirtualMachineProfile; profile != nil {
				if err := d.Set("boot_diagnostics", helpers.FlattenBootDiagnosticsVMSS(profile.DiagnosticsProfile)); err != nil {
					return fmt.Errorf("setting `boot_diagnostics`: %w", err)
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

				d.Set("eviction_policy", pointer.From(profile.EvictionPolicy))
				d.Set("license_type", profile.LicenseType)

				// the service just return empty when this is not assigned when provisioned
				// See discussion on https://github.com/Azure/azure-rest-api-specs/issues/10971
				priority := virtualmachinescalesets.VirtualMachinePriorityTypesRegular
				if profile.Priority != nil {
					priority = pointer.From(profile.Priority)
				}
				d.Set("priority", priority)

				if storageProfile := profile.StorageProfile; storageProfile != nil {
					if err := d.Set("os_disk", helpers.FlattenOrchestratedVirtualMachineScaleSetOSDisk(storageProfile.OsDisk)); err != nil {
						return fmt.Errorf("setting `os_disk`: %w", err)
					}

					if err := d.Set("data_disk", helpers.FlattenOrchestratedVirtualMachineScaleSetDataDisk(storageProfile.DataDisks)); err != nil {
						return fmt.Errorf("setting `data_disk`: %w", err)
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
						return fmt.Errorf("setting `source_image_reference`: %w", err)
					}
				}

				if osProfile := profile.OsProfile; osProfile != nil {
					if err := d.Set("os_profile", helpers.FlattenOrchestratedVirtualMachineScaleSetOSProfile(osProfile, d)); err != nil {
						return fmt.Errorf("setting `os_profile`: %w", err)
					}

					if osProfile.AllowExtensionOperations != nil {
						extensionOperationsEnabled = *osProfile.AllowExtensionOperations
					}
				}

				if nwProfile := profile.NetworkProfile; nwProfile != nil {
					d.Set("network_api_version", pointer.From(nwProfile.NetworkApiVersion))

					if err := d.Set("network_interface", helpers.FlattenOrchestratedVirtualMachineScaleSetNetworkInterface(nwProfile.NetworkInterfaceConfigurations)); err != nil {
						return fmt.Errorf("setting `network_interface`: %w", err)
					}
				}

				if scheduleProfile := profile.ScheduledEventsProfile; scheduleProfile != nil {
					if err := d.Set("termination_notification", helpers.FlattenOrchestratedVirtualMachineScaleSetScheduledEventsProfile(scheduleProfile)); err != nil {
						return fmt.Errorf("setting `termination_notification`: %w", err)
					}
				}

				extensionProfile, err := helpers.FlattenOrchestratedVirtualMachineScaleSetExtensions(profile.ExtensionProfile, d)
				if err != nil {
					return fmt.Errorf("failed flattening `extension`: %w", err)
				}
				d.Set("extension", extensionProfile)

				extensionsTimeBudget := "PT1H30M"
				if profile.ExtensionProfile != nil && profile.ExtensionProfile.ExtensionsTimeBudget != nil {
					extensionsTimeBudget = *profile.ExtensionProfile.ExtensionsTimeBudget
				}
				d.Set("extensions_time_budget", extensionsTimeBudget)

				encryptionAtHostEnabled := false
				if profile.SecurityProfile != nil && profile.SecurityProfile.EncryptionAtHost != nil {
					encryptionAtHostEnabled = *profile.SecurityProfile.EncryptionAtHost
				}
				d.Set("encryption_at_host_enabled", encryptionAtHostEnabled)
				d.Set("user_data_base64", profile.UserData)

				if policy := props.UpgradePolicy; policy != nil {
					upgradeMode = pointer.FromEnum(policy.Mode)
					if err := d.Set("rolling_upgrade_policy", helpers.FlattenVirtualMachineScaleSetRollingUpgradePolicy(policy.RollingUpgradePolicy)); err != nil {
						return fmt.Errorf("setting `rolling_upgrade_policy`: %w", err)
					}
				}
			}

			if priorityMixPolicy := props.PriorityMixPolicy; priorityMixPolicy != nil {
				if err := d.Set("priority_mix", helpers.FlattenOrchestratedVirtualMachineScaleSetPriorityMixPolicy(priorityMixPolicy)); err != nil {
					return fmt.Errorf("setting `priority_mix`: %w", err)
				}
			}

			d.Set("extension_operations_enabled", extensionOperationsEnabled)
			d.Set("upgrade_mode", upgradeMode)
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

func flattenOrchestratedVirtualMachineScaleSetSkuProfile(input *virtualmachinescalesets.SkuProfile) []any {
	if input == nil {
		return []any{}
	}

	result := map[string]any{
		"allocation_strategy": pointer.FromEnum(input.AllocationStrategy),
	}

	output := make([]any, 0)
	if input.VMSizes != nil {
		for _, vmSize := range *input.VMSizes {
			results := map[string]any{
				"name": pointer.From(vmSize.Name),
				"rank": nil,
			}

			if vmSize.Rank != nil {
				results["rank"] = int(pointer.From(vmSize.Rank)) + 1
			}

			output = append(output, results)
		}
	}

	result["virtual_machine_size"] = output

	return []any{result}
}

func flattenOrchestratedVirtualMachineScaleSetSku(input *virtualmachinescalesets.Sku) (*string, error) {
	var skuName string
	if input != nil && input.Name != nil {
		if strings.HasPrefix(strings.ToLower(*input.Name), "standard") || *input.Name == SkuNameMix {
			skuName = *input.Name
		} else {
			skuName = fmt.Sprintf("Standard_%s", *input.Name)
		}

		return &skuName, nil
	}

	return nil, fmt.Errorf("sku struct `name` is nil")
}

const (
	SkuNameMix = "Mix"
)
