// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func resourceSharedImage() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSharedImageCreateUpdate,
		Read:   resourceSharedImageRead,
		Update: resourceSharedImageCreateUpdate,
		Delete: resourceSharedImageDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SharedImageID(id)
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
				Default:  string(compute.ArchitectureTypesX64),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.ArchitectureTypesX64),
					string(compute.ArchitectureTypesArm64),
				}, false),
			},

			"os_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.OperatingSystemTypesLinux),
					string(compute.OperatingSystemTypesWindows),
				}, false),
			},

			"disk_types_not_allowed": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.DiskStorageAccountTypesStandardLRS),
						string(compute.DiskStorageAccountTypesPremiumLRS),
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
				Default:  string(compute.HyperVGenerationTypesV1),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.HyperVGenerationV1),
					string(compute.HyperVGenerationV2),
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

			"trusted_launch_enabled": {
				Type:          pluginsdk.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"confidential_vm_supported", "confidential_vm_enabled"},
			},

			"confidential_vm_supported": {
				Type:          pluginsdk.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"trusted_launch_enabled", "confidential_vm_enabled"},
			},

			"confidential_vm_enabled": {
				Type:          pluginsdk.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"trusted_launch_enabled", "confidential_vm_supported"},
			},

			"accelerated_network_support_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("end_of_life_date", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(string) != "" && new.(string) == ""
			}),
		),
	}
}

func resourceSharedImageCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Shared Image creation.")
	id := parse.NewSharedImageID(subscriptionId, d.Get("resource_group_name").(string), d.Get("gallery_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.GalleryName, id.ImageName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_shared_image", id.ID())
		}
	}

	recommended, err := expandGalleryImageRecommended(d)
	if err != nil {
		return err
	}

	image := compute.GalleryImage{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		GalleryImageProperties: &compute.GalleryImageProperties{
			Description:         utils.String(d.Get("description").(string)),
			Disallowed:          expandGalleryImageDisallowed(d),
			Identifier:          expandGalleryImageIdentifier(d),
			PrivacyStatementURI: utils.String(d.Get("privacy_statement_uri").(string)),
			ReleaseNoteURI:      utils.String(d.Get("release_note_uri").(string)),
			Architecture:        compute.Architecture(d.Get("architecture").(string)),
			OsType:              compute.OperatingSystemTypes(d.Get("os_type").(string)),
			HyperVGeneration:    compute.HyperVGeneration(d.Get("hyper_v_generation").(string)),
			PurchasePlan:        expandGalleryImagePurchasePlan(d.Get("purchase_plan").([]interface{})),
			Features:            expandSharedImageFeatures(d),
			Recommended:         recommended,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("end_of_life_date"); ok {
		endOfLifeDate, _ := time.Parse(time.RFC3339, v.(string))
		image.GalleryImageProperties.EndOfLifeDate = &date.Time{
			Time: endOfLifeDate,
		}
	}

	if v, ok := d.GetOk("eula"); ok {
		image.GalleryImageProperties.Eula = utils.String(v.(string))
	}

	if d.Get("specialized").(bool) {
		image.GalleryImageProperties.OsState = compute.OperatingSystemStateTypesSpecialized
	} else {
		image.GalleryImageProperties.OsState = compute.OperatingSystemStateTypesGeneralized
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.GalleryName, id.ImageName, image)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSharedImageRead(d, meta)
}

func resourceSharedImageRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SharedImageID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.GalleryName, id.ImageName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Shared Image %q (Gallery %q / Resource Group %q) was not found - removing from state", id.ImageName, id.GalleryName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Shared Image %q (Gallery %q / Resource Group %q): %+v", id.ImageName, id.GalleryName, id.ResourceGroup, err)
	}

	d.Set("name", id.ImageName)
	d.Set("gallery_name", id.GalleryName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.GalleryImageProperties; props != nil {
		d.Set("description", props.Description)

		diskTypesNotAllowed := make([]string, 0)
		if disallowed := props.Disallowed; disallowed != nil {
			if disallowed.DiskTypes != nil {
				diskTypesNotAllowed = append(diskTypesNotAllowed, *disallowed.DiskTypes...)
			}
		}
		d.Set("disk_types_not_allowed", diskTypesNotAllowed)

		if v := props.EndOfLifeDate; v != nil {
			d.Set("end_of_life_date", props.EndOfLifeDate.Format(time.RFC3339))
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

		architecture := string((compute.ArchitectureTypesX64))
		if props.Architecture != "" {
			architecture = string(props.Architecture)
		}
		d.Set("architecture", architecture)

		d.Set("specialized", props.OsState == compute.OperatingSystemStateTypesSpecialized)
		d.Set("hyper_v_generation", string(props.HyperVGeneration))
		d.Set("privacy_statement_uri", props.PrivacyStatementURI)
		d.Set("release_note_uri", props.ReleaseNoteURI)

		if err := d.Set("identifier", flattenGalleryImageIdentifier(props.Identifier)); err != nil {
			return fmt.Errorf("setting `identifier`: %+v", err)
		}

		if err := d.Set("purchase_plan", flattenGalleryImagePurchasePlan(props.PurchasePlan)); err != nil {
			return fmt.Errorf("setting `purchase_plan`: %+v", err)
		}

		trustedLaunchEnabled := false
		cvmEnabled := false
		cvmSupported := false
		acceleratedNetworkSupportEnabled := false
		if features := props.Features; features != nil {
			for _, feature := range *features {
				if feature.Name == nil || feature.Value == nil {
					continue
				}

				if strings.EqualFold(*feature.Name, "SecurityType") {
					trustedLaunchEnabled = strings.EqualFold(*feature.Value, "TrustedLaunch")
					cvmSupported = strings.EqualFold(*feature.Value, "ConfidentialVmSupported")
					cvmEnabled = strings.EqualFold(*feature.Value, "ConfidentialVm")
				}

				if strings.EqualFold(*feature.Name, "IsAcceleratedNetworkSupported") {
					acceleratedNetworkSupportEnabled = strings.EqualFold(*feature.Value, "true")
				}
			}
		}
		d.Set("confidential_vm_supported", cvmSupported)
		d.Set("confidential_vm_enabled", cvmEnabled)
		d.Set("trusted_launch_enabled", trustedLaunchEnabled)
		d.Set("accelerated_network_support_enabled", acceleratedNetworkSupportEnabled)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSharedImageDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SharedImageID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.GalleryName, id.ImageName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of: %s: %+v", *id, err)
	}

	log.Printf("[DEBUG] Waiting for %s to be eventually deleted", *id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   sharedImageDeleteStateRefreshFunc(ctx, client, id.ResourceGroup, id.GalleryName, id.ImageName),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
	}

	return nil
}

func sharedImageDeleteStateRefreshFunc(ctx context.Context, client *compute.GalleryImagesClient, resourceGroupName string, galleryName string, imageName string) pluginsdk.StateRefreshFunc {
	// The resource Shared Image depends on the resource Shared Image Gallery.
	// Although the delete API returns 404 which means the Shared Image resource has been deleted.
	// Then it tries to immediately delete Shared Image Gallery but it still throws error `Can not delete resource before nested resources are deleted.`
	// In this case we're going to try triggering the Deletion again, in-case it didn't work prior to this attempt.
	// For more details, see related Bug: https://github.com/Azure/azure-sdk-for-go/issues/8314
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, galleryName, imageName)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("failed to poll to check if the Shared Image has been deleted: %+v", err)
		}

		return res, "Exists", nil
	}
}

func expandGalleryImageIdentifier(d *pluginsdk.ResourceData) *compute.GalleryImageIdentifier {
	vs := d.Get("identifier").([]interface{})
	v := vs[0].(map[string]interface{})

	offer := v["offer"].(string)
	publisher := v["publisher"].(string)
	sku := v["sku"].(string)

	return &compute.GalleryImageIdentifier{
		Sku:       utils.String(sku),
		Publisher: utils.String(publisher),
		Offer:     utils.String(offer),
	}
}

func flattenGalleryImageIdentifier(input *compute.GalleryImageIdentifier) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	offer := ""
	if input.Offer != nil {
		offer = *input.Offer
	}

	publisher := ""
	if input.Publisher != nil {
		publisher = *input.Publisher
	}

	sku := ""
	if input.Sku != nil {
		sku = *input.Sku
	}

	return []interface{}{
		map[string]interface{}{
			"offer":     offer,
			"publisher": publisher,
			"sku":       sku,
		},
	}
}

func expandGalleryImagePurchasePlan(input []interface{}) *compute.ImagePurchasePlan {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	result := compute.ImagePurchasePlan{
		Name: utils.String(v["name"].(string)),
	}

	if publisher := v["publisher"].(string); publisher != "" {
		result.Publisher = &publisher
	}

	if product := v["product"].(string); product != "" {
		result.Product = &product
	}

	return &result
}

func flattenGalleryImagePurchasePlan(input *compute.ImagePurchasePlan) []interface{} {
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

func expandGalleryImageDisallowed(d *pluginsdk.ResourceData) *compute.Disallowed {
	diskTypesNotAllowedRaw := d.Get("disk_types_not_allowed").(*pluginsdk.Set).List()

	diskTypesNotAllowed := make([]string, 0)
	for _, v := range diskTypesNotAllowedRaw {
		diskTypesNotAllowed = append(diskTypesNotAllowed, v.(string))
	}

	return &compute.Disallowed{
		DiskTypes: &diskTypesNotAllowed,
	}
}

func expandGalleryImageRecommended(d *pluginsdk.ResourceData) (*compute.RecommendedMachineConfiguration, error) {
	result := &compute.RecommendedMachineConfiguration{
		VCPUs:  &compute.ResourceRange{},
		Memory: &compute.ResourceRange{},
	}

	maxVcpuCount := d.Get("max_recommended_vcpu_count").(int)
	minVcpuCount := d.Get("min_recommended_vcpu_count").(int)
	if maxVcpuCount != 0 && minVcpuCount != 0 && maxVcpuCount < minVcpuCount {
		return nil, fmt.Errorf("`max_recommended_vcpu_count` must be greater than or equal to `min_recommended_vcpu_count`")
	}
	if maxVcpuCount != 0 {
		result.VCPUs.Max = utils.Int32(int32(maxVcpuCount))
	}
	if minVcpuCount != 0 {
		result.VCPUs.Min = utils.Int32(int32(minVcpuCount))
	}

	maxMemory := d.Get("max_recommended_memory_in_gb").(int)
	minMemory := d.Get("min_recommended_memory_in_gb").(int)
	if maxMemory != 0 && minMemory != 0 && maxMemory < minMemory {
		return nil, fmt.Errorf("`max_recommended_memory_in_gb` must be greater than or equal to `min_recommended_memory_in_gb`")
	}
	if maxMemory != 0 {
		result.Memory.Max = utils.Int32(int32(maxMemory))
	}
	if minMemory != 0 {
		result.Memory.Min = utils.Int32(int32(minMemory))
	}

	return result, nil
}

func expandSharedImageFeatures(d *pluginsdk.ResourceData) *[]compute.GalleryImageFeature {
	var features []compute.GalleryImageFeature
	if d.Get("accelerated_network_support_enabled").(bool) {
		features = append(features, compute.GalleryImageFeature{
			Name:  utils.String("IsAcceleratedNetworkSupported"),
			Value: utils.String("true"),
		})
	}

	if tvmEnabled := d.Get("trusted_launch_enabled").(bool); tvmEnabled {
		features = append(features, compute.GalleryImageFeature{
			Name:  utils.String("SecurityType"),
			Value: utils.String("TrustedLaunch"),
		})
	}

	if cvmSupported := d.Get("confidential_vm_supported").(bool); cvmSupported {
		features = append(features, compute.GalleryImageFeature{
			Name:  utils.String("SecurityType"),
			Value: utils.String("ConfidentialVmSupported"),
		})
	}

	if cvmEnabled := d.Get("confidential_vm_enabled").(bool); cvmEnabled {
		features = append(features, compute.GalleryImageFeature{
			Name:  utils.String("SecurityType"),
			Value: utils.String("ConfidentialVM"),
		})
	}

	return &features
}
