// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSharedImageRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := galleryimages.ParseGalleryImageID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ImageName)
	d.Set("gallery_name", id.GalleryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("description", props.Description)

			diskTypesNotAllowed := make([]string, 0)
			if disallowed := props.Disallowed; disallowed != nil {
				if disallowed.DiskTypes != nil {
					diskTypesNotAllowed = append(diskTypesNotAllowed, *disallowed.DiskTypes...)
				}
			}
			d.Set("disk_types_not_allowed", diskTypesNotAllowed)

			if v := props.EndOfLifeDate; v != nil {
				d.Set("end_of_life_date", props.EndOfLifeDate)
			}

			d.Set("eula", props.Eula)

			maxRecommendedVcpuCount := 0
			minRecommendedVcpuCount := 0
			maxRecommendedMemoryInGB := 0
			minRecommendedMemoryInGB := 0
			if recommended := props.Recommended; recommended != nil {
				if vcpus := recommended.VCPUs; vcpus != nil {
					if vcpus.Max != nil {
						maxRecommendedVcpuCount = int(*vcpus.Max)
					}
					if vcpus.Min != nil {
						minRecommendedVcpuCount = int(*vcpus.Min)
					}
				}
				if memory := recommended.Memory; memory != nil {
					if memory.Max != nil {
						maxRecommendedMemoryInGB = int(*memory.Max)
					}
					if memory.Min != nil {
						minRecommendedMemoryInGB = int(*memory.Min)
					}
				}
			}
			d.Set("max_recommended_vcpu_count", maxRecommendedVcpuCount)
			d.Set("min_recommended_vcpu_count", minRecommendedVcpuCount)
			d.Set("max_recommended_memory_in_gb", maxRecommendedMemoryInGB)
			d.Set("min_recommended_memory_in_gb", minRecommendedMemoryInGB)

			d.Set("os_type", string(props.OsType))

			architecture := string(galleryimages.ArchitectureXSixFour)
			if props.Architecture != nil {
				architecture = string(*props.Architecture)
			}
			d.Set("architecture", architecture)

			d.Set("specialized", props.OsState == galleryimages.OperatingSystemStateTypesSpecialized)

			hyperVGeneration := string(galleryimages.HyperVGenerationVOne)
			if props.HyperVGeneration != nil {
				hyperVGeneration = string(*props.HyperVGeneration)
			}
			d.Set("hyper_v_generation", hyperVGeneration)
			d.Set("privacy_statement_uri", props.PrivacyStatementUri)
			d.Set("release_note_uri", props.ReleaseNoteUri)

			if err := d.Set("identifier", flattenGalleryImageIdentifier(&props.Identifier)); err != nil {
				return fmt.Errorf("setting `identifier`: %+v", err)
			}

			if err := d.Set("purchase_plan", flattenGalleryImagePurchasePlan(props.PurchasePlan)); err != nil {
				return fmt.Errorf("setting `purchase_plan`: %+v", err)
			}

			trustedLaunchSupported := false
			trustedLaunchEnabled := false
			cvmEnabled := false
			cvmSupported := false
			acceleratedNetworkSupportEnabled := false
			hibernationEnabled := false
			diskControllerTypeNVMEEnabled := false
			if features := props.Features; features != nil {
				for _, feature := range *features {
					if feature.Name == nil || feature.Value == nil {
						continue
					}

					if strings.EqualFold(*feature.Name, "SecurityType") {
						trustedLaunchSupported = strings.EqualFold(*feature.Value, "TrustedLaunchSupported")
						trustedLaunchEnabled = strings.EqualFold(*feature.Value, "TrustedLaunch")
						cvmSupported = strings.EqualFold(*feature.Value, "ConfidentialVmSupported")
						cvmEnabled = strings.EqualFold(*feature.Value, "ConfidentialVm")
					}

					if strings.EqualFold(*feature.Name, "IsAcceleratedNetworkSupported") {
						acceleratedNetworkSupportEnabled = strings.EqualFold(*feature.Value, "true")
					}

					if strings.EqualFold(*feature.Name, "IsHibernateSupported") {
						hibernationEnabled = strings.EqualFold(*feature.Value, "true")
					}

					if strings.EqualFold(*feature.Name, "DiskControllerTypes") {
						diskControllerTypeNVMEEnabled = strings.Contains(*feature.Value, "NVMe")
					}
				}
			}
			d.Set("confidential_vm_supported", cvmSupported)
			d.Set("confidential_vm_enabled", cvmEnabled)
			d.Set("trusted_launch_supported", trustedLaunchSupported)
			d.Set("trusted_launch_enabled", trustedLaunchEnabled)
			d.Set("accelerated_network_support_enabled", acceleratedNetworkSupportEnabled)
			d.Set("hibernation_enabled", hibernationEnabled)
			d.Set("disk_controller_type_nvme_enabled", diskControllerTypeNVMEEnabled)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("flattening `tags`: %+v", err)
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenGalleryImageIdentifier(input *galleryimages.GalleryImageIdentifier) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"offer":     input.Offer,
			"publisher": input.Publisher,
			"sku":       input.Sku,
		},
	}
}

func flattenGalleryImagePurchasePlan(input *galleryimages.ImagePurchasePlan) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"name":      pointer.From(input.Name),
			"publisher": pointer.From(input.Publisher),
			"product":   pointer.From(input.Product),
		},
	}
}
