// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudAPIPortal() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSpringCloudAPIPortalCreate,
		Read:   resourceSpringCloudAPIPortalRead,
		Update: resourceSpringCloudAPIPortalUpdate,
		Delete: resourceSpringCloudAPIPortalDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudApiPortalV0ToV1{},
		}),

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

func resourceSpringCloudAPIPortalCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_spring_cloud_api_portal", id.ID())
	}

	service, err := servicesClient.Get(ctx, springId.ResourceGroup, springId.SpringName)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", springId, err)
	}
	if service.Sku == nil || service.Sku.Name == nil || service.Sku.Tier == nil {
		return fmt.Errorf("invalid `sku` for %s", springId)
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
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSpringCloudAPIPortalRead(d, meta)
}

func resourceSpringCloudAPIPortalUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).AppPlatform.APIPortalClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	springId, err := parse.SpringCloudServiceID(d.Get("spring_cloud_service_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewSpringCloudAPIPortalID(subscriptionId, springId.ResourceGroup, springId.SpringName, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
	}
	if utils.ResponseWasNotFound(existing.Response) {
		return fmt.Errorf("retrieving %s: resource was not found", id)
	}

	if existing.Properties == nil {
		return fmt.Errorf("retrieving %s: properties are nil", id)
	}
	properties := existing.Properties

	if existing.Sku == nil {
		return fmt.Errorf("retrieving %s: sku is nil", id)
	}
	sku := existing.Sku

	if d.HasChange("gateway_ids") {
		properties.GatewayIds = utils.ExpandStringSlice(d.Get("gateway_ids").(*pluginsdk.Set).List())
	}

	if d.HasChange("https_only_enabled") {
		properties.HTTPSOnly = pointer.To(d.Get("https_only_enabled").(bool))
	}

	if d.HasChange("public_network_access_enabled") {
		properties.Public = pointer.To(d.Get("public_network_access_enabled").(bool))
	}

	if d.HasChange("sso") {
		properties.SsoProperties = expandAPIPortalSsoProperties(d.Get("sso").([]interface{}))
	}

	if d.HasChange("instance_count") {
		sku.Capacity = pointer.To(int32(d.Get("instance_count").(int)))
	}

	apiPortalResource := appplatform.APIPortalResource{
		Properties: properties,
		Sku:        sku,
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.ApiPortalName, apiPortalResource)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
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
		d.Set("gateway_ids", flattenSpringCloudAPIPortalGatewayIds(props.GatewayIds))
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

func flattenSpringCloudAPIPortalGatewayIds(ids *[]string) []string {
	if ids == nil || len(*ids) == 0 {
		return nil
	}
	out := make([]string, 0)
	for _, id := range *ids {
		gatewayId, err := parse.SpringCloudGatewayIDInsensitively(id)
		if err == nil {
			out = append(out, gatewayId.ID())
		}
	}
	return out
}
