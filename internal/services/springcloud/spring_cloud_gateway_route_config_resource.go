package springcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2022-05-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSpringCloudGatewayRouteConfig() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSpringCloudGatewayRouteConfigCreateUpdate,
		Read:   resourceSpringCloudGatewayRouteConfigRead,
		Update: resourceSpringCloudGatewayRouteConfigCreateUpdate,
		Delete: resourceSpringCloudGatewayRouteConfigDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SpringCloudGatewayRouteConfigID(id)
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
				ValidateFunc: validate.SpringCloudGatewayID,
			},

			"spring_cloud_app_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.SpringCloudAppID,
			},

			"route": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
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

						"order": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
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
		},
	}
}
func resourceSpringCloudGatewayRouteConfigCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).AppPlatform.GatewayRouteConfigClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	gatewayId, err := parse.SpringCloudGatewayID(d.Get("spring_cloud_gateway_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewSpringCloudGatewayRouteConfigID(subscriptionId, gatewayId.ResourceGroup, gatewayId.SpringName, gatewayId.GatewayName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.GatewayName, id.RouteConfigName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_gateway_route_config", id.ID())
		}
	}

	gatewayRouteConfigResource := appplatform.GatewayRouteConfigResource{
		Properties: &appplatform.GatewayRouteConfigProperties{
			AppResourceID: utils.String(d.Get("spring_cloud_app_id").(string)),
			Routes:        expandGatewayRouteConfigGatewayAPIRouteArray(d.Get("route").(*pluginsdk.Set).List()),
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.GatewayName, id.RouteConfigName, gatewayRouteConfigResource)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSpringCloudGatewayRouteConfigRead(d, meta)
}

func resourceSpringCloudGatewayRouteConfigRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.GatewayRouteConfigClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudGatewayRouteConfigID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.GatewayName, id.RouteConfigName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] appplatform %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.Set("name", id.RouteConfigName)
	d.Set("spring_cloud_gateway_id", parse.NewSpringCloudGatewayID(id.SubscriptionId, id.ResourceGroup, id.SpringName, id.GatewayName).ID())
	if props := resp.Properties; props != nil {
		d.Set("spring_cloud_app_id", props.AppResourceID)
		if err := d.Set("route", flattenGatewayRouteConfigGatewayAPIRouteArray(props.Routes)); err != nil {
			return fmt.Errorf("setting `route`: %+v", err)
		}
	}
	return nil
}

func resourceSpringCloudGatewayRouteConfigDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.GatewayRouteConfigClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudGatewayRouteConfigID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.GatewayName, id.RouteConfigName)
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
			Title:       utils.String(v["title"].(string)),
			Description: utils.String(v["description"].(string)),
			URI:         utils.String(v["uri"].(string)),
			SsoEnabled:  utils.Bool(v["sso_validation_enabled"].(bool)),
			TokenRelay:  utils.Bool(v["token_relay"].(bool)),
			Predicates:  utils.ExpandStringSlice(v["predicates"].(*pluginsdk.Set).List()),
			Filters:     utils.ExpandStringSlice(v["filters"].(*pluginsdk.Set).List()),
			Order:       utils.Int32(int32(v["order"].(int))),
			Tags:        utils.ExpandStringSlice(v["classification_tags"].(*pluginsdk.Set).List()),
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
		var description string
		if item.Description != nil {
			description = *item.Description
		}
		var order int32
		if item.Order != nil {
			order = *item.Order
		}
		var ssoEnabled bool
		if item.SsoEnabled != nil {
			ssoEnabled = *item.SsoEnabled
		}
		var title string
		if item.Title != nil {
			title = *item.Title
		}
		var tokenRelay bool
		if item.TokenRelay != nil {
			tokenRelay = *item.TokenRelay
		}
		var uri string
		if item.URI != nil {
			uri = *item.URI
		}
		results = append(results, map[string]interface{}{
			"description":            description,
			"filters":                utils.FlattenStringSlice(item.Filters),
			"order":                  order,
			"predicates":             utils.FlattenStringSlice(item.Predicates),
			"sso_validation_enabled": ssoEnabled,
			"title":                  title,
			"token_relay":            tokenRelay,
			"uri":                    uri,
			"classification_tags":    utils.FlattenStringSlice(item.Tags),
		})
	}
	return results
}
