package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceArmResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmResourcesRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"required_tags": tags.Schema(),

			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": azure.SchemaLocationForDataSource(),
						"tags":     tags.SchemaDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceArmResourcesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Resource.ResourcesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)
	resourceName := d.Get("name").(string)
	resourceType := d.Get("type").(string)
	requiredTags := d.Get("required_tags").(map[string]interface{})

	if resourceGroupName == "" && resourceName == "" && resourceType == "" {
		return fmt.Errorf("At least one of `name`, `resource_group_name` or `type` must be specified")
	}

	var filter string

	if resourceGroupName != "" {
		v := fmt.Sprintf("resourceGroup eq '%s'", resourceGroupName)
		filter = filter + v
	}

	if resourceName != "" {
		if strings.Contains(filter, "eq") {
			filter = filter + " and "
		}
		v := fmt.Sprintf("name eq '%s'", resourceName)
		filter = filter + v
	}

	if resourceType != "" {
		if strings.Contains(filter, "eq") {
			filter = filter + " and "
		}
		v := fmt.Sprintf("resourceType eq '%s'", resourceType)
		filter = filter + v
	}

	resources := make([]map[string]interface{}, 0)
	resourcesResp, err := client.ListComplete(ctx, filter, "", nil)
	if err != nil {
		return fmt.Errorf("Error getting resources: %+v", err)
	}

	for resourcesResp.NotDone() {
		res := resourcesResp.Value()
		if res.ID == nil {
			continue
		}

		// currently its not supported to use tags filter with other filters
		// therefore we need to filter the resources manually.
		tagMatches := 0
		if res.Tags != nil {
			for requiredTagName, requiredTagVal := range requiredTags {
				for tagName, tagVal := range res.Tags {
					if requiredTagName == tagName && tagVal != nil && requiredTagVal == *tagVal {
						tagMatches++
					}
				}
			}
		}

		if tagMatches == len(requiredTags) {
			resName := ""
			if res.Name != nil {
				resName = *res.Name
			}

			resID := ""
			if res.ID != nil {
				resID = *res.ID
			}

			resType := ""
			if res.Type != nil {
				resType = *res.Type
			}

			resLocation := ""
			if res.Location != nil {
				resLocation = *res.Location
			}

			resTags := make(map[string]interface{})
			if res.Tags != nil {
				resTags = make(map[string]interface{}, len(res.Tags))
				for key, value := range res.Tags {
					if value != nil {
						resTags[key] = *value
					}
				}
			}

			resources = append(resources, map[string]interface{}{
				"name":     resName,
				"id":       resID,
				"type":     resType,
				"location": resLocation,
				"tags":     resTags,
			})
		} else {
			log.Printf("[DEBUG] azurerm_resources - resources %q (id: %q) skipped as a required tag is not set or has the wrong value.", *res.Name, *res.ID)
		}

		err = resourcesResp.NextWithContext(ctx)
		if err != nil {
			return fmt.Errorf("Error loading Resource List: %s", err)
		}
	}

	d.SetId("resource-" + uuid.New().String())
	if err := d.Set("resources", resources); err != nil {
		return fmt.Errorf("Error setting `resources`: %+v", err)
	}

	return nil
}
