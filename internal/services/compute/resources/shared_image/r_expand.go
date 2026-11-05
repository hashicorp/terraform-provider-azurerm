// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package shared_image

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandGalleryImageIdentifier(d *pluginsdk.ResourceData) galleryimages.GalleryImageIdentifier {
	vs := d.Get("identifier").([]any)
	v := vs[0].(map[string]any)

	offer := v["offer"].(string)
	publisher := v["publisher"].(string)
	sku := v["sku"].(string)

	return galleryimages.GalleryImageIdentifier{
		Sku:       sku,
		Publisher: publisher,
		Offer:     offer,
	}
}

func expandGalleryImagePurchasePlan(input []any) *galleryimages.ImagePurchasePlan {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	result := galleryimages.ImagePurchasePlan{
		Name: pointer.To(v["name"].(string)),
	}

	if publisher := v["publisher"].(string); publisher != "" {
		result.Publisher = &publisher
	}

	if product := v["product"].(string); product != "" {
		result.Product = &product
	}

	return &result
}

func expandGalleryImageDisallowed(d *pluginsdk.ResourceData) *galleryimages.Disallowed {
	diskTypesNotAllowedRaw := d.Get("disk_types_not_allowed").(*pluginsdk.Set).List()

	diskTypesNotAllowed := make([]string, 0)
	for _, v := range diskTypesNotAllowedRaw {
		diskTypesNotAllowed = append(diskTypesNotAllowed, v.(string))
	}

	return &galleryimages.Disallowed{
		DiskTypes: &diskTypesNotAllowed,
	}
}

func expandGalleryImageRecommended(d *pluginsdk.ResourceData) (*galleryimages.RecommendedMachineConfiguration, error) {
	result := &galleryimages.RecommendedMachineConfiguration{
		VCPUs:  &galleryimages.ResourceRange{},
		Memory: &galleryimages.ResourceRange{},
	}

	maxVcpuCount := d.Get("max_recommended_vcpu_count").(int)
	minVcpuCount := d.Get("min_recommended_vcpu_count").(int)
	if maxVcpuCount != 0 && minVcpuCount != 0 && maxVcpuCount < minVcpuCount {
		return nil, fmt.Errorf("`max_recommended_vcpu_count` must be greater than or equal to `min_recommended_vcpu_count`")
	}
	if maxVcpuCount != 0 {
		result.VCPUs.Max = pointer.To(int64(maxVcpuCount))
	}
	if minVcpuCount != 0 {
		result.VCPUs.Min = pointer.To(int64(minVcpuCount))
	}

	maxMemory := d.Get("max_recommended_memory_in_gb").(int)
	minMemory := d.Get("min_recommended_memory_in_gb").(int)
	if maxMemory != 0 && minMemory != 0 && maxMemory < minMemory {
		return nil, fmt.Errorf("`max_recommended_memory_in_gb` must be greater than or equal to `min_recommended_memory_in_gb`")
	}
	if maxMemory != 0 {
		result.Memory.Max = pointer.To(int64(maxMemory))
	}
	if minMemory != 0 {
		result.Memory.Min = pointer.To(int64(minMemory))
	}

	return result, nil
}

func expandSharedImageFeatures(d *pluginsdk.ResourceData) *[]galleryimages.GalleryImageFeature {
	var features []galleryimages.GalleryImageFeature
	if d.Get("accelerated_network_support_enabled").(bool) {
		features = append(features, galleryimages.GalleryImageFeature{
			Name:  pointer.To("IsAcceleratedNetworkSupported"),
			Value: pointer.To("true"),
		})
	}

	if d.Get("disk_controller_type_nvme_enabled").(bool) {
		features = append(features, galleryimages.GalleryImageFeature{
			Name:  pointer.To("DiskControllerTypes"),
			Value: pointer.To("SCSI, NVMe"),
		})
	}

	if tvmSupported := d.Get("trusted_launch_supported").(bool); tvmSupported {
		features = append(features, galleryimages.GalleryImageFeature{
			Name:  pointer.To("SecurityType"),
			Value: pointer.To("TrustedLaunchSupported"),
		})
	}

	if tvmEnabled := d.Get("trusted_launch_enabled").(bool); tvmEnabled {
		features = append(features, galleryimages.GalleryImageFeature{
			Name:  pointer.To("SecurityType"),
			Value: pointer.To("TrustedLaunch"),
		})
	}

	if cvmSupported := d.Get("confidential_vm_supported").(bool); cvmSupported {
		features = append(features, galleryimages.GalleryImageFeature{
			Name:  pointer.To("SecurityType"),
			Value: pointer.To("ConfidentialVmSupported"),
		})
	}

	if cvmEnabled := d.Get("confidential_vm_enabled").(bool); cvmEnabled {
		features = append(features, galleryimages.GalleryImageFeature{
			Name:  pointer.To("SecurityType"),
			Value: pointer.To("ConfidentialVM"),
		})
	}

	if hibernationEnabled := d.Get("hibernation_enabled").(bool); hibernationEnabled {
		features = append(features, galleryimages.GalleryImageFeature{
			Name:  pointer.To("IsHibernateSupported"),
			Value: pointer.To(strconv.FormatBool(hibernationEnabled)),
		})
	}

	return &features
}
