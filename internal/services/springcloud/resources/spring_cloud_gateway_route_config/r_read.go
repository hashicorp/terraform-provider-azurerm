// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_gateway_route_config

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudGatewayRouteConfigRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppPlatform.GatewayRouteConfigClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform_rm.ParseRouteConfigIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.GatewayName, id.RouteConfigName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] appplatform %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.Set("name", id.RouteConfigName)
	d.Set("spring_cloud_gateway_id", appplatform_rm.NewGatewayID(id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.GatewayName).ID())
	if props := resp.Properties; props != nil {
		// The returned value has inconsistent casing
		// TODO: Remove the normalization codes once the following issue is fixed.
		// Issue: https://github.com/Azure/azure-rest-api-specs/issues/22845
		if props.AppResourceID != nil {
			appId, err := appplatform_rm.ParseAppIDInsensitively(*props.AppResourceID)
			if err != nil {
				return err
			}
			d.Set("spring_cloud_app_id", appId.ID())
		}
		d.Set("protocol", props.Protocol)
		if err := d.Set("route", flattenGatewayRouteConfigGatewayAPIRouteArray(props.Routes)); err != nil {
			return fmt.Errorf("setting `route`: %+v", err)
		}

		if err := d.Set("open_api", flattenGatewayRouteConfigOpenApi(props.OpenAPI)); err != nil {
			return fmt.Errorf("setting `open_api`: %+v", err)
		}

		if props.Filters != nil {
			d.Set("filters", pluginsdk.FlattenSlice(props.Filters))
		}
		if props.Predicates != nil {
			d.Set("predicates", pluginsdk.FlattenSlice(props.Predicates))
		}
		d.Set("sso_validation_enabled", props.SsoEnabled)
	}
	return nil
}

func flattenGatewayRouteConfigGatewayAPIRouteArray(input *[]appplatform.GatewayAPIRoute) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"description":            pointer.From(item.Description),
			"filters":                pluginsdk.FlattenSlice(item.Filters),
			"order":                  pointer.From(item.Order),
			"predicates":             pluginsdk.FlattenSlice(item.Predicates),
			"sso_validation_enabled": pointer.From(item.SsoEnabled),
			"title":                  pointer.From(item.Title),
			"token_relay":            pointer.From(item.TokenRelay),
			"uri":                    pointer.From(item.URI),
			"classification_tags":    pluginsdk.FlattenSlice(item.Tags),
		})
	}
	return results
}

func flattenGatewayRouteConfigOpenApi(input *appplatform.GatewayRouteConfigOpenAPIProperties) any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"uri": pointer.From(input.URI),
		},
	}
}
