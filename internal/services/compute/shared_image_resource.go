// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSharedImage() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSharedImageCreate,
		Read:   resourceSharedImageRead,
		Update: resourceSharedImageUpdate,
		Delete: resourceSharedImageDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := galleryimages.ParseGalleryImageID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SharedImageName,
			},

			"gallery_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"architecture": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(galleryimages.ArchitectureXSixFour),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(galleryimages.ArchitectureXSixFour),
					string(galleryimages.ArchitectureArmSixFour),
				}, false),
			},

			"os_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(galleryimages.OperatingSystemTypesLinux),
					string(galleryimages.OperatingSystemTypesWindows),
				}, false),
			},

			"disk_types_not_allowed": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(galleryimageversions.StorageAccountTypeStandardLRS),
						string(galleryimageversions.StorageAccountTypePremiumLRS),
					}, false),
				},
			},

			"end_of_life_date": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time,
			},

			"hyper_v_generation": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(galleryimages.HyperVGenerationVOne),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(galleryimages.HyperVGenerationVOne),
					string(galleryimages.HyperVGenerationVTwo),
				}, false),
			},

			"identifier": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"publisher": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Required:     true,
							ValidateFunc: validate.SharedImageIdentifierAttribute(128),
						},
						"offer": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Required:     true,
							ValidateFunc: validate.SharedImageIdentifierAttribute(64),
						},
						"sku": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Required:     true,
							ValidateFunc: validate.SharedImageIdentifierAttribute(64),
						},
					},
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"eula": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"purchase_plan": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"publisher": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"product": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"privacy_statement_uri": {
				Type:     pluginsdk.TypeString,
				ForceNew: true,
				Optional: true,
			},

			"max_recommended_vcpu_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 80),
			},

			"min_recommended_vcpu_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 80),
			},

			"max_recommended_memory_in_gb": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 640),
			},

			"min_recommended_memory_in_gb": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 640),
			},

			"release_note_uri": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"specialized": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"trusted_launch_supported": {
				Type:          pluginsdk.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"trusted_launch_enabled", "confidential_vm_supported", "confidential_vm_enabled"},
			},

			"trusted_launch_enabled": {
				Type:          pluginsdk.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"trusted_launch_supported", "confidential_vm_supported", "confidential_vm_enabled"},
			},

			"confidential_vm_supported": {
				Type:          pluginsdk.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"trusted_launch_supported", "trusted_launch_enabled", "confidential_vm_enabled"},
			},

			"confidential_vm_enabled": {
				Type:          pluginsdk.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"trusted_launch_supported", "trusted_launch_enabled", "confidential_vm_supported"},
			},

			"accelerated_network_support_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"hibernation_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"disk_controller_type_nvme_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("end_of_life_date", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(string) != "" && new.(string) == ""
			}),
		),
	}
}

func resourceSharedImageCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Shared Image creation.")
	id := galleryimages.NewGalleryImageID(subscriptionId, d.Get("resource_group_name").(string), d.Get("gallery_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_shared_image", id.ID())
	}

	recommended, err := expandGalleryImageRecommended(d)
	if err != nil {
		return err
	}

	image := galleryimages.GalleryImage{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &galleryimages.GalleryImageProperties{
			Description:         pointer.To(d.Get("description").(string)),
			Disallowed:          expandGalleryImageDisallowed(d),
			Identifier:          expandGalleryImageIdentifier(d),
			PrivacyStatementUri: pointer.To(d.Get("privacy_statement_uri").(string)),
			ReleaseNoteUri:      pointer.To(d.Get("release_note_uri").(string)),
			Architecture:        pointer.To(galleryimages.Architecture(d.Get("architecture").(string))),
			OsType:              galleryimages.OperatingSystemTypes(d.Get("os_type").(string)),
			HyperVGeneration:    pointer.To(galleryimages.HyperVGeneration(d.Get("hyper_v_generation").(string))),
			PurchasePlan:        expandGalleryImagePurchasePlan(d.Get("purchase_plan").([]interface{})),
			Features:            expandSharedImageFeatures(d),
			Recommended:         recommended,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("end_of_life_date"); ok {
		endOfLifeDate, _ := time.Parse(time.RFC3339, v.(string))
		image.Properties.EndOfLifeDate = pointer.To(date.Time{
			Time: endOfLifeDate,
		}.String())
	}

	if v, ok := d.GetOk("eula"); ok {
		image.Properties.Eula = pointer.To(v.(string))
	}

	if d.Get("specialized").(bool) {
		image.Properties.OsState = galleryimages.OperatingSystemStateTypesSpecialized
	} else {
		image.Properties.OsState = galleryimages.OperatingSystemStateTypesGeneralized
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, image); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSharedImageRead(d, meta)
}

func resourceSharedImageUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := galleryimages.ParseGalleryImageID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	payload := existing.Model

	if payload == nil {
		return fmt.Errorf("model is nil for %s", id)
	}

	if d.HasChange("disk_types_not_allowed") {
		payload.Properties.Disallowed = expandGalleryImageDisallowed(d)
	}

	if d.HasChange("end_of_life_date") {
		endOfLifeDate, _ := time.Parse(time.RFC3339, d.Get("end_of_life_date").(string))
		payload.Properties.EndOfLifeDate = pointer.To(date.Time{
			Time: endOfLifeDate,
		}.String())
	}

	if d.HasChange("description") {
		payload.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("eula") {
		payload.Properties.Description = pointer.To(d.Get("eula").(string))
	}

	if d.HasChange("specialized") {
		if d.Get("specialized").(bool) {
			payload.Properties.OsState = galleryimages.OperatingSystemStateTypesSpecialized
		} else {
			payload.Properties.OsState = galleryimages.OperatingSystemStateTypesGeneralized
		}
	}

	if d.HasChange("release_note_uri") {
		payload.Properties.ReleaseNoteUri = pointer.To(d.Get("release_note_uri").(string))
	}

	if d.HasChanges("max_recommended_vcpu_count", "min_recommended_vcpu_count", "max_recommended_memory_in_gb", "min_recommended_memory_in_gb") {
		recommended, err := expandGalleryImageRecommended(d)
		if err != nil {
			return err
		}
		payload.Properties.Recommended = recommended
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSharedImageRead(d, meta)
}

func resourceSharedImageRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

			architecture := string((galleryimages.ArchitectureXSixFour))
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

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceSharedImageDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := galleryimages.ParseGalleryImageID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	log.Printf("[DEBUG] Waiting for %s to be eventually deleted", *id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   sharedImageDeleteStateRefreshFunc(ctx, client, *id),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
	}

	return nil
}

func sharedImageDeleteStateRefreshFunc(ctx context.Context, client *galleryimages.GalleryImagesClient, id galleryimages.GalleryImageId) pluginsdk.StateRefreshFunc {
	// The resource Shared Image depends on the resource Shared Image Gallery.
	// Although the delete API returns 404 which means the Shared Image resource has been deleted.
	// Then it tries to immediately delete Shared Image Gallery but it still throws error `Can not delete resource before nested resources are deleted.`
	// In this case we're going to try triggering the Deletion again, in-case it didn't work prior to this attempt.
	// For more details, see related Bug: https://github.com/Azure/azure-sdk-for-go/issues/8314
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("failed to poll to check if the Shared Image has been deleted: %+v", err)
		}

		return res, "Exists", nil
	}
}

func expandGalleryImageIdentifier(d *pluginsdk.ResourceData) galleryimages.GalleryImageIdentifier {
	vs := d.Get("identifier").([]interface{})
	v := vs[0].(map[string]interface{})

	offer := v["offer"].(string)
	publisher := v["publisher"].(string)
	sku := v["sku"].(string)

	return galleryimages.GalleryImageIdentifier{
		Sku:       sku,
		Publisher: publisher,
		Offer:     offer,
	}
}

func flattenGalleryImageIdentifier(input *galleryimages.GalleryImageIdentifier) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"offer":     input.Offer,
			"publisher": input.Publisher,
			"sku":       input.Sku,
		},
	}
}

func expandGalleryImagePurchasePlan(input []interface{}) *galleryimages.ImagePurchasePlan {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
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

func flattenGalleryImagePurchasePlan(input *galleryimages.ImagePurchasePlan) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	name := ""
	if input.Name != nil {
		name = *input.Name
	}

	publisher := ""
	if input.Publisher != nil {
		publisher = *input.Publisher
	}

	product := ""
	if input.Product != nil {
		product = *input.Product
	}

	return []interface{}{
		map[string]interface{}{
			"name":      name,
			"publisher": publisher,
			"product":   product,
		},
	}
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
