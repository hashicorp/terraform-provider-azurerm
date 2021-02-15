package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceRouteFilter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRouteFilterRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"rule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"access": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"rule_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"communities": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceRouteFilterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RouteFiltersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Route Filter %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("making Read request on Azure Route Filter %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.RouteFilterPropertiesFormat; props != nil {
		if err := d.Set("rule", flattenRouteFilterDataSourceRules(props.Rules)); err != nil {
			return err
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenRouteFilterDataSourceRules(input *[]network.RouteFilterRule) []interface{} {
	results := make([]interface{}, 0)

	if rules := input; rules != nil {
		for _, rule := range *rules {
			r := make(map[string]interface{})

			r["name"] = *rule.Name
			if props := rule.RouteFilterRulePropertiesFormat; props != nil {
				r["access"] = string(props.Access)
				r["rule_type"] = *props.RouteFilterRuleType
				r["communities"] = utils.FlattenStringSlice(props.Communities)
			}

			results = append(results, r)
		}
	}

	return results
}
