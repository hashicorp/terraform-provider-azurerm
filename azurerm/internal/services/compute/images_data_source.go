package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmImagesRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"tags_filter": tags.Schema(),

			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"location": location.SchemaComputed(),

						"zone_resilient": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"os_disk": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"blob_uri": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"caching": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"managed_disk_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size_gb": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"data_disk": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"blob_uri": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"caching": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"lun": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"managed_disk_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size_gb": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"tags": tags.SchemaDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceArmImagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ImagesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	filterTags := tags.Expand(d.Get("tags_filter").(map[string]interface{}))

	resp, err := client.ListByResourceGroupComplete(ctx, resourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response().Response) {
			return fmt.Errorf("Images (Resource Group %q) was not found", resourceGroup)
		}
		return fmt.Errorf("retrieving Images (Resource Group %q): %+v", resourceGroup, err)
	}

	images, err := flattenImages(ctx, resp, filterTags)
	if err != nil {
		return fmt.Errorf("retrieving Images (Resource Group %q): %+v", resourceGroup, err)
	}
	if len(images) == 0 {
		return fmt.Errorf("unable to find any images")
	}

	d.SetId(time.Now().UTC().String())

	d.Set("resource_group_name", resourceGroup)

	if err := d.Set("images", images); err != nil {
		return fmt.Errorf("setting `images`: %+v", err)
	}

	return nil
}

func flattenImages(ctx context.Context, iterator compute.ImageListResultIterator, filterTags map[string]*string) ([]interface{}, error) {
	results := make([]interface{}, 0)

	var err error
	for ; iterator.NotDone(); err = iterator.NextWithContext(ctx) {
		if err != nil {
			return nil, fmt.Errorf("loading Images list: %+v", err)
		}
		image := iterator.Value()
		found := true
		// Loop through our filter tags and see if they match
		for k, v := range filterTags {
			if v != nil {
				// If the tags do not match return false
				if image.Tags[k] == nil || *v != *image.Tags[k] {
					found = false
				}
			}
		}

		if found {
			results = append(results, flattenImage(image))
		}
	}

	return results, nil
}

func flattenImage(input compute.Image) map[string]interface{} {
	output := make(map[string]interface{})

	output["name"] = input.Name
	output["location"] = location.NormalizeNilable(input.Location)

	if input.ImageProperties != nil {
		if storageProfile := input.ImageProperties.StorageProfile; storageProfile != nil {
			output["zone_resilient"] = storageProfile.ZoneResilient

			output["os_disk"] = flattenAzureRmImageOSDisk(storageProfile.OsDisk)

			output["data_disk"] = flattenAzureRmImageDataDisks(storageProfile.DataDisks)
		}
	}

	output["tags"] = tags.Flatten(input.Tags)

	return output
}
