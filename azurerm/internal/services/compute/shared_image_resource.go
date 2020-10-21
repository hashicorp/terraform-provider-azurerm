package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-30/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSharedImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSharedImageCreateUpdate,
		Read:   resourceArmSharedImageRead,
		Update: resourceArmSharedImageCreateUpdate,
		Delete: resourceArmSharedImageDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SharedImageID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SharedImageName,
			},

			"gallery_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"os_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Linux),
					string(compute.Windows),
				}, false),
			},

			"hyper_v_generation": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(compute.HyperVGenerationTypesV1),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.V1),
					string(compute.V2),
				}, false),
			},

			"identifier": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"publisher": {
							Type:     schema.TypeString,
							Required: true,
						},
						"offer": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sku": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"eula": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"purchase_plan": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"publisher": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"product": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"privacy_statement_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"release_note_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"specialized": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmSharedImageCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Shared Image creation.")

	name := d.Get("name").(string)
	galleryName := d.Get("gallery_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, galleryName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Shared Image %q (Gallery %q / Resource Group %q): %+v", name, galleryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_shared_image", *existing.ID)
		}
	}

	image := compute.GalleryImage{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		GalleryImageProperties: &compute.GalleryImageProperties{
			Description:         utils.String(d.Get("description").(string)),
			Eula:                utils.String(d.Get("eula").(string)),
			Identifier:          expandGalleryImageIdentifier(d),
			PrivacyStatementURI: utils.String(d.Get("privacy_statement_uri").(string)),
			ReleaseNoteURI:      utils.String(d.Get("release_note_uri").(string)),
			OsType:              compute.OperatingSystemTypes(d.Get("os_type").(string)),
			HyperVGeneration:    compute.HyperVGeneration(d.Get("hyper_v_generation").(string)),
			PurchasePlan:        expandGalleryImagePurchasePlan(d.Get("purchase_plan").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if d.Get("specialized").(bool) {
		image.GalleryImageProperties.OsState = compute.Specialized
	} else {
		image.GalleryImageProperties.OsState = compute.Generalized
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, galleryName, name, image)
	if err != nil {
		return fmt.Errorf("Error creating/updating Shared Image %q (Gallery %q / Resource Group %q): %+v", name, galleryName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of Shared Image %q (Gallery %q / Resource Group %q): %+v", name, galleryName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, galleryName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Shared Image %q (Gallery %q / Resource Group %q): %+v", name, galleryName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Shared Image %q (Gallery %q / Resource Group %q) ID", name, galleryName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmSharedImageRead(d, meta)
}

func resourceArmSharedImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SharedImageID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Gallery, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Shared Image %q (Gallery %q / Resource Group %q) was not found - removing from state", id.Name, id.Gallery, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Shared Image %q (Gallery %q / Resource Group %q): %+v", id.Name, id.Gallery, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("gallery_name", id.Gallery)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.GalleryImageProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("eula", props.Eula)
		d.Set("os_type", string(props.OsType))
		d.Set("specialized", props.OsState == compute.Specialized)
		d.Set("hyper_v_generation", string(props.HyperVGeneration))
		d.Set("privacy_statement_uri", props.PrivacyStatementURI)
		d.Set("release_note_uri", props.ReleaseNoteURI)

		if err := d.Set("identifier", flattenGalleryImageIdentifier(props.Identifier)); err != nil {
			return fmt.Errorf("Error setting `identifier`: %+v", err)
		}

		if err := d.Set("purchase_plan", flattenGalleryImagePurchasePlan(props.PurchasePlan)); err != nil {
			return fmt.Errorf("Error setting `purchase_plan`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmSharedImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SharedImageID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Gallery, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Shared Image %q (Gallery %q / Resource Group %q): %+v", id.Name, id.Gallery, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed to wait for deleting Shared Image %q (Gallery %q / Resource Group %q): %+v", id.Name, id.Gallery, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for Shared Image %q (Gallery %q / Resource Group %q) to be eventually deleted", id.Name, id.Gallery, id.ResourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   sharedImageDeleteStateRefreshFunc(ctx, client, id.ResourceGroup, id.Gallery, id.Name),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("failed to wait for Shared Image %q (Gallery %q / Resource Group %q) to be deleted: %+v", id.Name, id.Gallery, id.ResourceGroup, err)
	}

	return nil
}

func sharedImageDeleteStateRefreshFunc(ctx context.Context, client *compute.GalleryImagesClient, resourceGroupName string, galleryName string, imageName string) resource.StateRefreshFunc {
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

func expandGalleryImageIdentifier(d *schema.ResourceData) *compute.GalleryImageIdentifier {
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
