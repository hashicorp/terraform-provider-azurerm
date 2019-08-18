package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSharedImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSharedImageRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageName,
			},

			"gallery_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"identifier": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"publisher": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"offer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sku": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"eula": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"privacy_statement_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"release_note_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}
func dataSourceArmSharedImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.GalleryImagesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	galleryName := d.Get("gallery_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, galleryName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Shared Image %q was not found in Gallery %q / Resource Group %q", name, galleryName, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on Shared Image %q (Gallery %q / Resource Group %q): %+v", name, galleryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("gallery_name", galleryName)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.GalleryImageProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("eula", props.Eula)
		d.Set("os_type", string(props.OsType))
		d.Set("privacy_statement_uri", props.PrivacyStatementURI)
		d.Set("release_note_uri", props.ReleaseNoteURI)

		flattenedIdentifier := flattenGalleryImageDataSourceIdentifier(props.Identifier)
		if err := d.Set("identifier", flattenedIdentifier); err != nil {
			return fmt.Errorf("Error setting `identifier`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenGalleryImageDataSourceIdentifier(input *compute.GalleryImageIdentifier) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if input.Offer != nil {
		result["offer"] = *input.Offer
	}

	if input.Publisher != nil {
		result["publisher"] = *input.Publisher
	}

	if input.Sku != nil {
		result["sku"] = *input.Sku
	}

	return []interface{}{result}
}
