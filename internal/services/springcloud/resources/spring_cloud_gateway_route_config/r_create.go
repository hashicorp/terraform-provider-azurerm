// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_gateway_route_config

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudGatewayRouteConfigCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).AppPlatform.GatewayRouteConfigClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	gatewayId, err := appplatform_rm.ParseGatewayIDInsensitively(d.Get("spring_cloud_gateway_id").(string))
	if err != nil {
		return err
	}
	id := appplatform_rm.NewRouteConfigID(subscriptionId, gatewayId.ResourceGroupName, gatewayId.SpringName, gatewayId.GatewayName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.GatewayName, id.RouteConfigName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_spring_cloud_gateway_route_config", id.ID())
			}
		}
	}

	gatewayRouteConfigResource := appplatform.GatewayRouteConfigResource{
		Properties: &appplatform.GatewayRouteConfigProperties{
			AppResourceID: pointer.To(d.Get("spring_cloud_app_id").(string)),
			Protocol:      appplatform.GatewayRouteConfigProtocol(d.Get("protocol").(string)),
			Routes:        expandGatewayRouteConfigGatewayAPIRouteArray(d.Get("route").(*pluginsdk.Set).List()),
			SsoEnabled:    pointer.To(d.Get("sso_validation_enabled").(bool)),
			OpenAPI:       expandGatewayRouteConfigOpenApi(d.Get("open_api").([]any)),
		},
	}

	filters := d.Get("filters").(*pluginsdk.Set).List()
	if len(filters) > 0 {
		gatewayRouteConfigResource.Properties.Filters = pluginsdk.ExpandStringSlice(filters)
	}

	predicates := d.Get("predicates").(*pluginsdk.Set).List()
	if len(predicates) > 0 {
		gatewayRouteConfigResource.Properties.Predicates = pluginsdk.ExpandStringSlice(predicates)
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, id.GatewayName, id.RouteConfigName, gatewayRouteConfigResource)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	return resourceSpringCloudGatewayRouteConfigRead(d, meta)
}

func expandGatewayRouteConfigGatewayAPIRouteArray(input []any) *[]appplatform.GatewayAPIRoute {
	results := make([]appplatform.GatewayAPIRoute, 0)
	for _, item := range input {
		v := item.(map[string]any)
		results = append(results, appplatform.GatewayAPIRoute{
			Title:       pointer.To(v["title"].(string)),
			Description: pointer.To(v["description"].(string)),
			URI:         pointer.To(v["uri"].(string)),
			SsoEnabled:  pointer.To(v["sso_validation_enabled"].(bool)),
			TokenRelay:  pointer.To(v["token_relay"].(bool)),
			Predicates:  pluginsdk.ExpandStringSlice(v["predicates"].(*pluginsdk.Set).List()),
			Filters:     pluginsdk.ExpandStringSlice(v["filters"].(*pluginsdk.Set).List()),
			Order:       pointer.To(int32(v["order"].(int))),
			Tags:        pluginsdk.ExpandStringSlice(v["classification_tags"].(*pluginsdk.Set).List()),
		})
	}
	return &results
}

func expandGatewayRouteConfigOpenApi(input []any) *appplatform.GatewayRouteConfigOpenAPIProperties {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]any)
	return &appplatform.GatewayRouteConfigOpenAPIProperties{
		URI: pointer.To(config["uri"].(string)),
	}
}
