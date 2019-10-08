package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceArmResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmResourceRead,
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name", "resource_group_name", "type", "required_tags"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"resource_id"},
			},
			"resource_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"resource_id"},
			},
			"type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"resource_id"},
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

func dataSourceArmResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Resource.ResourcesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)
	resourceName := d.Get("name").(string)
	resourceType := d.Get("type").(string)
	resourceID := d.Get("resource_id").(string)

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

	requiredTags := d.Get("required_tags").(map[string]interface{})

	resources := make([]map[string]interface{}, 0)

	resource, err := client.ListComplete(ctx, filter, "", nil)
	if err != nil {
		return fmt.Errorf("Error getting resources: %+v", err)
	}

	for resource.NotDone() {

		res := resource.Value()

		// currently the Azure-Go-SDK method "GetByID" does not work for some resources, as the
		// API Version is hard coded, therefore we use ListComplete and look for the ResourceID 'manually'
		if resourceID != "" && *res.ID != resourceID {
			err = resource.NextWithContext(ctx)
			if err != nil {
				return fmt.Errorf("Error loading Resource List: %s", err)
			}
			continue
		}

		// currently its not supported to use a other filters together with the tags filter
		// therefore we need to filter the resources manually.
		tagMatches := 0
		for requiredTagName, requiredTagVal := range requiredTags {
			for tagName, tagVal := range res.Tags {
				if requiredTagName == tagName && requiredTagVal == *tagVal {
					tagMatches++
				}
			}
		}

		if tagMatches == len(requiredTags) {

			s := make(map[string]interface{})

			if v := *res.Name; v != "" {
				s["name"] = v
			}

			if v := *res.ID; v != "" {
				s["id"] = v
			}

			if v := *res.Type; v != "" {
				s["type"] = v
			}

			if v := *res.Location; v != "" {
				s["location"] = v
			}

			tags := make(map[string]interface{}, len(res.Tags))
			for key, value := range res.Tags {
				tags[key] = *value
			}

			s["tags"] = tags

			resources = append(resources, s)
		} else {
			log.Printf("[DEBUG] azurerm_resource - resources %q (id: %q) skipped as a required tag is not set or has the wrong value.", *res.Name, *res.ID)
		}

		err = resource.NextWithContext(ctx)
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
