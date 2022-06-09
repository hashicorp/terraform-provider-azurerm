package compute

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"os_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.OperatingSystemTypesLinux),
					string(compute.OperatingSystemTypesWindows),
				}, false),
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
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"offer": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"sku": {
							Type:     pluginsdk.TypeString,
							Required: true,
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
				Optional: true,
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
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"accelerated_network_support_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"tags": tags.Schema(),
		},
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

	var features []compute.GalleryImageFeature
	if d.Get("trusted_launch_enabled").(bool) {
		features = append(features, compute.GalleryImageFeature{
			Name:  utils.String("SecurityType"),
			Value: utils.String("TrustedLaunch"),
		})
	}
	if d.Get("accelerated_network_support_enabled").(bool) {
		features = append(features, compute.GalleryImageFeature{
			Name:  utils.String("IsAcceleratedNetworkSupported"),
			Value: utils.String("true"),
		})
	}

	image := compute.GalleryImage{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		GalleryImageProperties: &compute.GalleryImageProperties{
			Description:         utils.String(d.Get("description").(string)),
			Identifier:          expandGalleryImageIdentifier(d),
			PrivacyStatementURI: utils.String(d.Get("privacy_statement_uri").(string)),
			ReleaseNoteURI:      utils.String(d.Get("release_note_uri").(string)),
			OsType:              compute.OperatingSystemTypes(d.Get("os_type").(string)),
			HyperVGeneration:    compute.HyperVGeneration(d.Get("hyper_v_generation").(string)),
			PurchasePlan:        expandGalleryImagePurchasePlan(d.Get("purchase_plan").([]interface{})),
			Features:            &features,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
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
		d.Set("eula", props.Eula)
		d.Set("os_type", string(props.OsType))
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
		acceleratedNetworkSupportEnabled := false
		if features := props.Features; features != nil {
			for _, feature := range *features {
				if feature.Name == nil || feature.Value == nil {
					continue
				}

				if strings.EqualFold(*feature.Name, "SecurityType") {
					trustedLaunchEnabled = strings.EqualFold(*feature.Value, "TrustedLaunch")
				}

				if strings.EqualFold(*feature.Name, "IsAcceleratedNetworkSupported") {
					acceleratedNetworkSupportEnabled = strings.EqualFold(*feature.Value, "true")
				}
			}
		}
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
