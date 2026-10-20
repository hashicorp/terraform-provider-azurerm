// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package route_filter

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routefilters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteFilterCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.RouteFilters
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := routefilters.NewRouteFilterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]any)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, routefilters.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_route_filter", id.ID())
		}
	}

	routeSet := routefilters.RouteFilter{
		Name:     &id.RouteFilterName,
		Location: pointer.To(location),
		Properties: &routefilters.RouteFilterPropertiesFormat{
			Rules: expandRouteFilterRules(d),
		},
		Tags: tags.Expand(t),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, routeSet, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceRouteFilterRead(d, meta)
}

func expandRouteFilterRules(d *pluginsdk.ResourceData) *[]routefilters.RouteFilterRule {
	configs := d.Get("rule").([]any)
	rules := make([]routefilters.RouteFilterRule, 0, len(configs))

	for _, configRaw := range configs {
		data := configRaw.(map[string]any)

		rule := routefilters.RouteFilterRule{
			Name: pointer.To(data["name"].(string)),
			Properties: &routefilters.RouteFilterRulePropertiesFormat{
				Access:              routefilters.Access(data["access"].(string)),
				RouteFilterRuleType: routefilters.RouteFilterRuleType(data["rule_type"].(string)),
				Communities:         *pluginsdk.ExpandStringSlice(data["communities"].([]any)),
			},
		}

		rules = append(rules, rule)
	}

	return &rules
}
