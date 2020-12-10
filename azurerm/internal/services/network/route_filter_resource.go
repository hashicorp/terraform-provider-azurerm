package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRouteFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRouteFilterCreateUpdate,
		Read:   resourceArmRouteFilterRead,
		Update: resourceArmRouteFilterCreateUpdate,
		Delete: resourceArmRouteFilterDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.RouteFilterID(id)
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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"rule": {
				Type:       schema.TypeList,
				ConfigMode: schema.SchemaConfigModeAttr,
				Optional:   true,
				Computed:   true,
				MaxItems:   1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"access": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Allow),
							}, false),
						},

						"rule_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Community",
							}, false),
						},

						"communities": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmRouteFilterCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RouteFiltersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Route Filter create/update.")

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Route Filter %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_route_filter", *existing.ID)
		}
	}

	routeSet := network.RouteFilter{
		Name:     &name,
		Location: &location,
		RouteFilterPropertiesFormat: &network.RouteFilterPropertiesFormat{
			Rules: expandRouteFilterRules(d),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, routeSet)
	if err != nil {
		return fmt.Errorf("creating/updating Route Filter %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of Route Filter %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("retrieving Route Filter %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("ID was nil for Route Filter %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmRouteFilterRead(d, meta)
}

func resourceArmRouteFilterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RouteFiltersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RouteFilterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Route Table %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.RouteFilterPropertiesFormat; props != nil {
		if err := d.Set("rule", flattenRouteFilterRules(props.Rules)); err != nil {
			return err
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmRouteFilterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RouteFiltersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RouteFilterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("deleting Route Filter %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Route Filter %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandRouteFilterRules(d *schema.ResourceData) *[]network.RouteFilterRule {
	configs := d.Get("rule").([]interface{})
	rules := make([]network.RouteFilterRule, 0, len(configs))

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		rule := network.RouteFilterRule{
			Name: utils.String(data["name"].(string)),
			RouteFilterRulePropertiesFormat: &network.RouteFilterRulePropertiesFormat{
				Access:              network.Access(data["access"].(string)),
				RouteFilterRuleType: utils.String(data["rule_type"].(string)),
				Communities:         utils.ExpandStringSlice(data["communities"].([]interface{})),
			},
		}

		rules = append(rules, rule)
	}

	return &rules
}

func flattenRouteFilterRules(input *[]network.RouteFilterRule) []interface{} {
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
