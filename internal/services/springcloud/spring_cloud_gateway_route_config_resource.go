// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudGatewayRouteConfig() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		DeprecationMessage: "Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_gateway_route_config` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.",

		Create: resourceSpringCloudGatewayRouteConfigCreateUpdate,
		Read:   resourceSpringCloudGatewayRouteConfigRead,
		Update: resourceSpringCloudGatewayRouteConfigCreateUpdate,
		Delete: resourceSpringCloudGatewayRouteConfigDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudGatewayRouteConfigV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := appplatform_rm.ParseRouteConfigID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"spring_cloud_gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: appplatform_rm.ValidateGatewayID,
			},

			"open_api": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"uri": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},
					},
				},
			},

			"protocol": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(appplatform.GatewayRouteConfigProtocolHTTP),
					string(appplatform.GatewayRouteConfigProtocolHTTPS),
				}, false),
			},

			"spring_cloud_app_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: appplatform_rm.ValidateAppID,
			},

			"filters": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"predicates": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"route": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"order": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},

						"description": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"filters": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"predicates": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"sso_validation_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"title": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"token_relay": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"uri": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"classification_tags": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},

			"sso_validation_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceSpringCloudGatewayRouteConfigCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).AppPlatform.GatewayRouteConfigClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	gatewayId, err := appplatform_rm.ParseGatewayID(d.Get("spring_cloud_gateway_id").(string))
	if err != nil {
		return err
	}
	id := appplatform_rm.NewRouteConfigID(subscriptionId, gatewayId.ResourceGroupName, gatewayId.SpringName, gatewayId.GatewayName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.GatewayName, id.RouteConfigName)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
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
			OpenAPI:       expandGatewayRouteConfigOpenApi(d.Get("open_api").([]interface{})),
		},
	}

	filters := d.Get("filters").(*pluginsdk.Set).List()
	if len(filters) > 0 {
		gatewayRouteConfigResource.Properties.Filters = helpers.ExpandStringSlice(filters)
	}

	predicates := d.Get("predicates").(*pluginsdk.Set).List()
	if len(predicates) > 0 {
		gatewayRouteConfigResource.Properties.Predicates = helpers.ExpandStringSlice(predicates)
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

func resourceSpringCloudGatewayRouteConfigRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.GatewayRouteConfigClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := appplatform_rm.ParseRouteConfigID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.GatewayName, id.RouteConfigName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
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
				return fmt.Errorf("parsing `spring_cloud_app_id`: %+v", err)
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
			d.Set("filters", helpers.FlattenStringSlice(props.Filters))
		}
		if props.Predicates != nil {
			d.Set("predicates", helpers.FlattenStringSlice(props.Predicates))
		}
		d.Set("sso_validation_enabled", props.SsoEnabled)
	}
	return nil
}

func resourceSpringCloudGatewayRouteConfigDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.GatewayRouteConfigClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := appplatform_rm.ParseRouteConfigID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroupName, id.SpringName, id.GatewayName, id.RouteConfigName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}
	return nil
}

func expandGatewayRouteConfigGatewayAPIRouteArray(input []interface{}) *[]appplatform.GatewayAPIRoute {
	results := make([]appplatform.GatewayAPIRoute, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, appplatform.GatewayAPIRoute{
			Title:       pointer.To(v["title"].(string)),
			Description: pointer.To(v["description"].(string)),
			URI:         pointer.To(v["uri"].(string)),
			SsoEnabled:  pointer.To(v["sso_validation_enabled"].(bool)),
			TokenRelay:  pointer.To(v["token_relay"].(bool)),
			Predicates:  helpers.ExpandStringSlice(v["predicates"].(*pluginsdk.Set).List()),
			Filters:     helpers.ExpandStringSlice(v["filters"].(*pluginsdk.Set).List()),
			Order:       pointer.To(int32(v["order"].(int))),
			Tags:        helpers.ExpandStringSlice(v["classification_tags"].(*pluginsdk.Set).List()),
		})
	}
	return &results
}

func flattenGatewayRouteConfigGatewayAPIRouteArray(input *[]appplatform.GatewayAPIRoute) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"description":            pointer.From(item.Description),
			"filters":                helpers.FlattenStringSlice(item.Filters),
			"order":                  pointer.From(item.Order),
			"predicates":             helpers.FlattenStringSlice(item.Predicates),
			"sso_validation_enabled": pointer.From(item.SsoEnabled),
			"title":                  pointer.From(item.Title),
			"token_relay":            pointer.From(item.TokenRelay),
			"uri":                    pointer.From(item.URI),
			"classification_tags":    helpers.FlattenStringSlice(item.Tags),
		})
	}
	return results
}

func expandGatewayRouteConfigOpenApi(input []interface{}) *appplatform.GatewayRouteConfigOpenAPIProperties {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})
	return &appplatform.GatewayRouteConfigOpenAPIProperties{
		URI: pointer.To(config["uri"].(string)),
	}
}

func flattenGatewayRouteConfigOpenApi(input *appplatform.GatewayRouteConfigOpenAPIProperties) interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"uri": pointer.From(input.URI),
		},
	}
}
