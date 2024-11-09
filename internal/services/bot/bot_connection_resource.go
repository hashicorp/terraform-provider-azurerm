// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceArmBotConnection() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceArmBotConnectionCreate,
		Read:   resourceArmBotConnectionRead,
		Update: resourceArmBotConnectionUpdate,
		Delete: resourceArmBotConnectionDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.BotConnectionID(id)
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
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"bot_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"service_provider_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"client_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"client_secret": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"scopes": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}

	return resource
}

func resourceArmBotConnectionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ConnectionClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewBotConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, resourceId.ConnectionName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Bot Connection %q (Bot %q / Resource Group %q): %+v", resourceId.ConnectionName, resourceId.BotServiceName, resourceId.ResourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_bot_connection", resourceId.ID())
		}
	}

	serviceProviderName := d.Get("service_provider_name").(string)
	var serviceProviderId *string

	serviceProviders, err := client.ListServiceProviders(ctx)
	if err != nil {
		return fmt.Errorf("listing Bot Connection service provider: %+v", err)
	}

	if serviceProviders.Value == nil {
		return errors.New("no service providers were returned from the Azure API")
	}

	availableProviders := make([]string, 0, len(*serviceProviders.Value))
	for _, provider := range *serviceProviders.Value {
		if provider.Properties == nil || provider.Properties.ServiceProviderName == nil {
			continue
		}
		name := provider.Properties.ServiceProviderName
		if strings.EqualFold(serviceProviderName, *name) {
			serviceProviderId = provider.Properties.ID
			break
		}
		availableProviders = append(availableProviders, *name)
	}

	if serviceProviderId == nil {
		return fmt.Errorf("the Service Provider %q was not found. The available service providers are %s", serviceProviderName, strings.Join(availableProviders, ","))
	}

	connection := botservice.ConnectionSetting{
		Properties: &botservice.ConnectionSettingProperties{
			ServiceProviderID: serviceProviderId,
			ClientID:          pointer.To(d.Get("client_id").(string)),
			ClientSecret:      pointer.To(d.Get("client_secret").(string)),
			Scopes:            pointer.To(d.Get("scopes").(string)),
		},
		Kind:     botservice.KindBot,
		Location: pointer.To(d.Get("location").(string)),
	}

	if v, ok := d.GetOk("parameters"); ok {
		connection.Properties.Parameters = expandBotConnectionParameters(v.(map[string]interface{}))
	}

	if _, err := client.Create(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, resourceId.ConnectionName, connection); err != nil {
		return fmt.Errorf("creating Bot Connection %q (Bot %q / Resource Group %q): %+v", resourceId.ConnectionName, resourceId.BotServiceName, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())
	return resourceArmBotConnectionRead(d, meta)
}

func resourceArmBotConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ConnectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, id.ConnectionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Bot Connection %q (Bot %q / Resource Group %q) was not found - removing from state!", id.ConnectionName, id.BotServiceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Bot Connection %q (Bot %q / Resource Group %q): %+v", id.ConnectionName, id.BotServiceName, id.ResourceGroup, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		d.Set("client_id", props.ClientID)
		d.Set("scopes", props.Scopes)
		if err := d.Set("parameters", flattenBotConnectionParameters(props.Parameters)); err != nil {
			return fmt.Errorf("setting `parameters`: %+v", err)
		}
	}

	return nil
}

func resourceArmBotConnectionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ConnectionClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotConnectionID(d.Id())
	if err != nil {
		return err
	}

	connection := botservice.ConnectionSetting{
		Properties: &botservice.ConnectionSettingProperties{
			ServiceProviderDisplayName: utils.String(d.Get("service_provider_name").(string)),
			ClientID:                   utils.String(d.Get("client_id").(string)),
			ClientSecret:               utils.String(d.Get("client_secret").(string)),
			Scopes:                     utils.String(d.Get("scopes").(string)),
		},
		Kind:     botservice.KindBot,
		Location: utils.String(d.Get("location").(string)),
	}

	if v, ok := d.GetOk("parameters"); ok {
		connection.Properties.Parameters = expandBotConnectionParameters(v.(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, id.ConnectionName, connection); err != nil {
		return fmt.Errorf("updating Bot Connection %q (Bot %q / Resource Group %q): %+v", id.ConnectionName, id.BotServiceName, id.ResourceGroup, err)
	}

	return resourceArmBotConnectionRead(d, meta)
}

func resourceArmBotConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ConnectionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, id.ConnectionName)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting Bot Connection %q (Bot %q / Resource Group %q): %+v", id.ConnectionName, id.BotServiceName, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandBotConnectionParameters(input map[string]interface{}) *[]botservice.ConnectionSettingParameter {
	output := make([]botservice.ConnectionSettingParameter, 0)

	for k, v := range input {
		output = append(output, botservice.ConnectionSettingParameter{
			Key:   utils.String(k),
			Value: utils.String(v.(string)),
		})
	}
	return &output
}

func flattenBotConnectionParameters(input *[]botservice.ConnectionSettingParameter) map[string]interface{} {
	output := make(map[string]interface{})
	if input == nil {
		return output
	}

	for _, parameter := range *input {
		if key := parameter.Key; key != nil {
			// We disregard the clientSecret and clientId as one is sensitive and the other is returned in the ClientId attribute.
			if *key != "clientSecret" && *key != "clientId" && *key != "scopes" {
				if value := parameter.Value; value != nil {
					output[*key] = *value
				}
			}
		}
	}

	return output
}
