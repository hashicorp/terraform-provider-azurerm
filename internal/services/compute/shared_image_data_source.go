package compute

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSharedImageID(subscriptionId, d.Get("resource_group_name").(string), d.Get("gallery_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.GalleryName, id.ImageName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())

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

		if err := d.Set("identifier", flattenGalleryImageDataSourceIdentifier(props.Identifier)); err != nil {
			return fmt.Errorf("setting `identifier`: %+v", err)
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
