// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/routefilters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceRouteFilter() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRouteFilterCreate,
		Read:   resourceRouteFilterRead,
		Update: resourceRouteFilterUpdate,
		Delete: resourceRouteFilterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := routefilters.ParseRouteFilterID(id)
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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"rule": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				Computed:   true,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				MaxItems:   1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"access": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(routefilters.AccessAllow),
							}, false),
						},

						"rule_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Community",
							}, false),
						},

						"communities": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceRouteFilterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RouteFilters
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Route Filter create.")

	id := routefilters.NewRouteFilterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	existing, err := client.Get(ctx, id, routefilters.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_route_filter", id.ID())
	}

	routeSet := routefilters.RouteFilter{
		Name:     &id.RouteFilterName,
		Location: pointer.To(location),
		Properties: &routefilters.RouteFilterPropertiesFormat{
			Rules: expandRouteFilterRules(d),
		},
		Tags: tags.Expand(t),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, routeSet); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceRouteFilterRead(d, meta)
}

func resourceRouteFilterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RouteFilters
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Route Filter update.")

	id, err := routefilters.ParseRouteFilterID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, routefilters.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	payload := existing.Model

	if d.HasChange("rule") {
		payload.Properties.Rules = expandRouteFilterRules(d)
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceRouteFilterRead(d, meta)
}

func resourceRouteFilterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RouteFilters
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routefilters.ParseRouteFilterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, routefilters.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.RouteFilterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("rule", flattenRouteFilterRules(props.Rules)); err != nil {
				return err
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceRouteFilterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RouteFilters
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routefilters.ParseRouteFilterID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandRouteFilterRules(d *pluginsdk.ResourceData) *[]routefilters.RouteFilterRule {
	configs := d.Get("rule").([]interface{})
	rules := make([]routefilters.RouteFilterRule, 0, len(configs))

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		rule := routefilters.RouteFilterRule{
			Name: utils.String(data["name"].(string)),
			Properties: &routefilters.RouteFilterRulePropertiesFormat{
				Access:              routefilters.Access(data["access"].(string)),
				RouteFilterRuleType: routefilters.RouteFilterRuleType(data["rule_type"].(string)),
				Communities:         *utils.ExpandStringSlice(data["communities"].([]interface{})),
			},
		}

		rules = append(rules, rule)
	}

	return &rules
}

func flattenRouteFilterRules(input *[]routefilters.RouteFilterRule) []interface{} {
	results := make([]interface{}, 0)

	if rules := input; rules != nil {
		for _, rule := range *rules {
			r := make(map[string]interface{})

			r["name"] = *rule.Name
			if props := rule.Properties; props != nil {
				r["access"] = string(props.Access)
				r["rule_type"] = props.RouteFilterRuleType
				r["communities"] = props.Communities
			}

			results = append(results, r)
		}
	}

	return results
}
