// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/validate"
	kvValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelsRegistration() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceBotChannelsRegistrationCreate,
		Read:   resourceBotChannelsRegistrationRead,
		Delete: resourceBotChannelsRegistrationDelete,
		Update: resourceBotChannelsRegistrationUpdate,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.BotServiceID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).Bot.BotClient

			id, err := parse.BotServiceID(d.Id())
			if err != nil {
				return nil, err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return nil, fmt.Errorf("Bot Channels Registration %q was not found in Resource Group %q", id.Name, id.ResourceGroup)
				}

				return nil, fmt.Errorf("retrieving Bot Channels Registration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}
			if resp.Kind != botservice.KindBot {
				return nil, fmt.Errorf("Bot %q (Resource Group %q) was not a Channel Registration - got %q", id.Name, id.ResourceGroup, string(resp.Kind))
			}

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(botservice.SkuNameF0),
					string(botservice.SkuNameS1),
				}, false),
			},

			"microsoft_app_id": {
				Type:         pluginsdk.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"cmk_key_vault_url": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: kvValidate.NestedItemIdWithOptionalVersion,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.BotChannelRegistrationDescription,
			},

			"display_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"endpoint": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"developer_app_insights_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},

			"developer_app_insights_api_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"developer_app_insights_application_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},

			"icon_url": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.BotChannelRegistrationIconUrl,
			},

			"streaming_endpoint_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: !features.FourPointOhBeta(),
				ConflictsWith: func() []string {
					if !features.FourPointOhBeta() {
						return []string{"isolated_network_enabled"}
					}
					return []string{}
				}(),
			},

			"tags": tags.Schema(),
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["isolated_network_enabled"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			Deprecated:    "`isolated_network_enabled` will be removed in favour of the property `public_network_access_enabled` in version 4.0 of the AzureRM Provider.",
			ConflictsWith: []string{"public_network_access_enabled"},
		}
	}

	return resource
}

func resourceBotChannelsRegistrationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.BotClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewBotServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Bot Channels Registration %q (Resource Group %q): %+v", resourceId.Name, resourceId.ResourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_bot_channels_registration", resourceId.ID())
		}
	}

	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = resourceId.Name
	}

	bot := botservice.Bot{
		Properties: &botservice.BotProperties{
			DisplayName:                       utils.String(displayName),
			Endpoint:                          utils.String(d.Get("endpoint").(string)),
			MsaAppID:                          utils.String(d.Get("microsoft_app_id").(string)),
			CmekKeyVaultURL:                   utils.String(d.Get("cmk_key_vault_url").(string)),
			Description:                       utils.String(d.Get("description").(string)),
			DeveloperAppInsightKey:            utils.String(d.Get("developer_app_insights_key").(string)),
			DeveloperAppInsightsAPIKey:        utils.String(d.Get("developer_app_insights_api_key").(string)),
			DeveloperAppInsightsApplicationID: utils.String(d.Get("developer_app_insights_application_id").(string)),
			IconURL:                           utils.String(d.Get("icon_url").(string)),
			IsCmekEnabled:                     utils.Bool(false),
			IsStreamingSupported:              utils.Bool(d.Get("streaming_endpoint_enabled").(bool)),
		},
		Location: utils.String(d.Get("location").(string)),
		Sku: &botservice.Sku{
			Name: botservice.SkuName(d.Get("sku").(string)),
		},
		Kind: botservice.KindBot,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, ok := d.GetOk("cmk_key_vault_url"); ok {
		bot.Properties.IsCmekEnabled = utils.Bool(true)
	}

	if _, err := client.Create(ctx, resourceId.ResourceGroup, resourceId.Name, bot); err != nil {
		return fmt.Errorf("creating Bot Channels Registration %q (Resource Group %q): %+v", resourceId.Name, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())

	// d.GetOk cannot identify whether user sets the property that is bool type and `public_network_access_enabled` is set as `false`. So it has to identify it using `d.GetRawConfig()`
	publicNetworkAccessEnabled := d.GetRawConfig().AsValueMap()["public_network_access_enabled"]
	if !features.FourPointOhBeta() {
		// d.GetOk cannot identify whether user sets the property that is bool type and `isolated_network_enabled` is set as `false`. So it has to identify it using `d.GetRawConfig()`
		isolatedNetworkEnabled := d.GetRawConfig().AsValueMap()["isolated_network_enabled"]
		if !isolatedNetworkEnabled.IsNull() || !publicNetworkAccessEnabled.IsNull() {
			return resourceBotChannelsRegistrationUpdate(d, meta)
		} else {
			return resourceBotChannelsRegistrationRead(d, meta)
		}
	} else {
		if !publicNetworkAccessEnabled.IsNull() {
			return resourceBotChannelsRegistrationUpdate(d, meta)
		} else {
			return resourceBotChannelsRegistrationRead(d, meta)
		}
	}
}

func resourceBotChannelsRegistrationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.BotClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Bot Channels Registration %q (Resource Group %q) was not found - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Bot Channels Registration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := resp.Properties; props != nil {
		d.Set("cmk_key_vault_url", props.CmekKeyVaultURL)
		d.Set("microsoft_app_id", props.MsaAppID)
		d.Set("endpoint", props.Endpoint)
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("developer_app_insights_key", props.DeveloperAppInsightKey)
		d.Set("developer_app_insights_application_id", props.DeveloperAppInsightsApplicationID)
		d.Set("icon_url", props.IconURL)
		d.Set("streaming_endpoint_enabled", props.IsStreamingSupported)

		// `PublicNetworkAccess` is empty string when `public_network_access_enabled` or `isolated_network_enabled` isn't specified. So `public_network_access_enabled` and `isolated_network_enabled` shouldn't be set at this time to avoid diff
		if props.PublicNetworkAccess != "" {
			d.Set("public_network_access_enabled", props.PublicNetworkAccess == botservice.PublicNetworkAccessEnabled)

			if !features.FourPointOhBeta() {
				d.Set("isolated_network_enabled", props.PublicNetworkAccess == botservice.PublicNetworkAccessDisabled)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceBotChannelsRegistrationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.BotClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotServiceID(d.Id())
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]interface{})
	displayName := d.Get("display_name").(string)
	if displayName == "" {
		displayName = id.Name
	}

	bot := botservice.Bot{
		Properties: &botservice.BotProperties{
			DisplayName:                       utils.String(displayName),
			Endpoint:                          utils.String(d.Get("endpoint").(string)),
			MsaAppID:                          utils.String(d.Get("microsoft_app_id").(string)),
			CmekKeyVaultURL:                   utils.String(d.Get("cmk_key_vault_url").(string)),
			Description:                       utils.String(d.Get("description").(string)),
			DeveloperAppInsightKey:            utils.String(d.Get("developer_app_insights_key").(string)),
			DeveloperAppInsightsAPIKey:        utils.String(d.Get("developer_app_insights_api_key").(string)),
			DeveloperAppInsightsApplicationID: utils.String(d.Get("developer_app_insights_application_id").(string)),
			IconURL:                           utils.String(d.Get("icon_url").(string)),
			IsCmekEnabled:                     utils.Bool(false),
			IsStreamingSupported:              utils.Bool(d.Get("streaming_endpoint_enabled").(bool)),
		},
		Location: utils.String(d.Get("location").(string)),
		Sku: &botservice.Sku{
			Name: botservice.SkuName(d.Get("sku").(string)),
		},
		Kind: botservice.KindBot,
		Tags: tags.Expand(t),
	}

	if _, ok := d.GetOk("cmk_key_vault_url"); ok {
		bot.Properties.IsCmekEnabled = utils.Bool(true)
	}

	if !features.FourPointOhBeta() {
		// d.GetOk cannot identify whether user sets the property that is bool type and `isolated_network_enabled` is set as `false`. So it has to identify it using `d.GetRawConfig()`
		if v := d.GetRawConfig().AsValueMap()["isolated_network_enabled"]; !v.IsNull() {
			publicNetworkAccessEnabled := botservice.PublicNetworkAccessEnabled
			if v.True() {
				publicNetworkAccessEnabled = botservice.PublicNetworkAccessDisabled
			}
			bot.Properties.PublicNetworkAccess = publicNetworkAccessEnabled
		}
	}

	// d.GetOk cannot identify whether user sets the property that is bool type and `public_network_access_enabled` is set as `false`. So it has to identify it using `d.GetRawConfig()`
	if v := d.GetRawConfig().AsValueMap()["public_network_access_enabled"]; !v.IsNull() {
		publicNetworkAccessEnabled := botservice.PublicNetworkAccessEnabled
		if v.False() {
			publicNetworkAccessEnabled = botservice.PublicNetworkAccessDisabled
		}
		bot.Properties.PublicNetworkAccess = publicNetworkAccessEnabled
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, bot); err != nil {
		return fmt.Errorf("updating Bot Channels Registration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceBotChannelsRegistrationRead(d, meta)
}

func resourceBotChannelsRegistrationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.BotClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting Bot Channels Registration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
