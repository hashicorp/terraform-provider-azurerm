// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_connection

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceArmBotConnectionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Bot.ConnectionClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewBotConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, resourceId.ConnectionName)
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing Bot Connection %q (Bot %q / Resource Group %q): %+v", resourceId.ConnectionName, resourceId.BotServiceName, resourceId.ResourceGroup, err)
			}
		}
		if !response.WasNotFound(existing.Response.Response) {
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
		connection.Properties.Parameters = expandBotConnectionParameters(v.(map[string]any))
	}

	if _, err := client.Create(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, resourceId.ConnectionName, connection); err != nil {
		return fmt.Errorf("creating Bot Connection %q (Bot %q / Resource Group %q): %+v", resourceId.ConnectionName, resourceId.BotServiceName, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())
	return resourceArmBotConnectionRead(d, meta)
}

func expandBotConnectionParameters(input map[string]any) *[]botservice.ConnectionSettingParameter {
	output := make([]botservice.ConnectionSettingParameter, 0)

	for k, v := range input {
		output = append(output, botservice.ConnectionSettingParameter{
			Key:   pointer.To(k),
			Value: pointer.To(v.(string)),
		})
	}
	return &output
}
