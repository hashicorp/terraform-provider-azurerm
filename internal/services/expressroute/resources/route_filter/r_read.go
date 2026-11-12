// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route_filter

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routefilters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteFilterRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.RouteFilters
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

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenRouteFilterRules(input *[]routefilters.RouteFilterRule) []any {
	results := make([]any, 0)

	if rules := input; rules != nil {
		for _, rule := range *rules {
			r := make(map[string]any)

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
