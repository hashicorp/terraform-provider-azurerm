package compute

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceSharedImage() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSharedImageRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageName,
			},

			"gallery_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.SharedImageGalleryName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"os_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"specialized": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"hyper_v_generation": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identifier": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"publisher": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"offer": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"sku": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"eula": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"privacy_statement_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"release_note_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceSharedImageRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.GalleryImagesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
		d.Set("specialized", props.OsState == compute.Specialized)
		d.Set("hyper_v_generation", string(props.HyperVGeneration))
		d.Set("privacy_statement_uri", props.PrivacyStatementURI)
		d.Set("release_note_uri", props.ReleaseNoteURI)

		if err := d.Set("identifier", flattenGalleryImageDataSourceIdentifier(props.Identifier)); err != nil {
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
