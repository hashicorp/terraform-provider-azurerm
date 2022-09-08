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

func resourceSpringCloudAPIPortal() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSpringCloudAPIPortalCreateUpdate,
		Read:   resourceSpringCloudAPIPortalRead,
		Update: resourceSpringCloudAPIPortalCreateUpdate,
		Delete: resourceSpringCloudAPIPortalDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SpringCloudAPIPortalID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"default",
				}, false),
			},

			"spring_cloud_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudServiceID,
			},

			"gateway_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.SpringCloudGatewayID,
				},
			},

			"https_only_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"instance_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 500),
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"sso": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"client_id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"client_secret": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"issuer_uri": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"scope": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceSpringCloudAPIPortalCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).AppPlatform.APIPortalClient
	servicesClient := meta.(*clients.Client).AppPlatform.ServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	springId, err := parse.SpringCloudServiceID(d.Get("spring_cloud_service_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewSpringCloudAPIPortalID(subscriptionId, springId.ResourceGroup, springId.SpringName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_api_portal", id.ID())
		}
	}

	service, err := servicesClient.Get(ctx, springId.ResourceGroup, springId.SpringName)
	if err != nil {
		return fmt.Errorf("checking for presence of existing Spring Cloud Service %q (Resource Group %q): %+v", springId.SpringName, springId.ResourceGroup, err)
	}
	if service.Sku == nil || service.Sku.Name == nil || service.Sku.Tier == nil {
		return fmt.Errorf("invalid `sku` for Spring Cloud Service %q (Resource Group %q)", springId.SpringName, springId.ResourceGroup)
	}

	apiPortalResource := appplatform.APIPortalResource{
		Properties: &appplatform.APIPortalProperties{
			GatewayIds:    utils.ExpandStringSlice(d.Get("gateway_ids").(*pluginsdk.Set).List()),
			HTTPSOnly:     utils.Bool(d.Get("https_only_enabled").(bool)),
			Public:        utils.Bool(d.Get("public_network_access_enabled").(bool)),
			SsoProperties: expandAPIPortalSsoProperties(d.Get("sso").([]interface{})),
		},
		Sku: &appplatform.Sku{
			Name:     service.Sku.Name,
			Tier:     service.Sku.Tier,
			Capacity: utils.Int32(int32(d.Get("instance_count").(int))),
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName, apiPortalResource)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSpringCloudAPIPortalRead(d, meta)
}

func resourceSpringCloudAPIPortalRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.APIPortalClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAPIPortalID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] appplatform %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.Set("name", id.ApiPortalName)
	d.Set("spring_cloud_service_id", parse.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroup, id.SpringName).ID())
	if resp.Sku != nil {
		d.Set("instance_count", resp.Sku.Capacity)
	}
	if props := resp.Properties; props != nil {
		d.Set("gateway_ids", utils.FlattenStringSlice(props.GatewayIds))
		d.Set("https_only_enabled", props.HTTPSOnly)
		d.Set("public_network_access_enabled", props.Public)
		if err := d.Set("sso", flattenAPIPortalSsoProperties(props.SsoProperties, d.Get("sso").([]interface{}))); err != nil {
			return fmt.Errorf("setting `sso`: %+v", err)
		}
		d.Set("url", props.URL)
	}
	return nil
}

func resourceSpringCloudAPIPortalDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.APIPortalClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAPIPortalID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}
	return nil
}

func expandAPIPortalSsoProperties(input []interface{}) *appplatform.SsoProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &appplatform.SsoProperties{
		Scope:        utils.ExpandStringSlice(v["scope"].(*pluginsdk.Set).List()),
		ClientID:     utils.String(v["client_id"].(string)),
		ClientSecret: utils.String(v["client_secret"].(string)),
		IssuerURI:    utils.String(v["issuer_uri"].(string)),
	}
}

func flattenAPIPortalSsoProperties(input *appplatform.SsoProperties, old []interface{}) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	oldItems := make(map[string]map[string]interface{})
	for _, item := range old {
		v := item.(map[string]interface{})
		if name, ok := v["issuer_uri"]; ok {
			oldItems[name.(string)] = v
		}
	}

	var issuerUri string
	if input.IssuerURI != nil {
		issuerUri = *input.IssuerURI
	}
	var clientId string
	var clientSecret string
	if oldItem, ok := oldItems[issuerUri]; ok {
		if value, ok := oldItem["client_id"]; ok {
			clientId = value.(string)
		}
		if value, ok := oldItem["client_secret"]; ok {
			clientSecret = value.(string)
		}
	}
	return []interface{}{
		map[string]interface{}{
			"client_id":     clientId,
			"client_secret": clientSecret,
			"issuer_uri":    issuerUri,
			"scope":         utils.FlattenStringSlice(input.Scope),
		},
	}
}
