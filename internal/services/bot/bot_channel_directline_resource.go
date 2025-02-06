// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func resourceBotChannelDirectline() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBotChannelDirectlineCreate,
		Read:   resourceBotChannelDirectlineRead,
		Delete: resourceBotChannelDirectlineDelete,
		Update: resourceBotChannelDirectlineUpdate,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.BotChannelID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"bot_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"site": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"user_upload_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"endpoint_parameters_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"storage_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"v1_allowed": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"v3_allowed": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"enhanced_authentication_enabled": {
							Type:     pluginsdk.TypeBool,
							Default:  false,
							Optional: true,
						},

						"trusted_origins": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"key": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},

						"key2": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceBotChannelDirectlineCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewBotChannelID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), string(botservice.ChannelNameBasicChannelChannelNameDirectLineChannel))
	existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, resourceId.ChannelName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Directline Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroup, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		// a "Default Site" site gets created and returned.. so let's check it's not just that
		if props := existing.Properties; props != nil {
			directLineChannel, ok := props.AsDirectLineChannel()
			if ok && directLineChannel.Properties != nil {
				sites := filterSites(directLineChannel.Properties.Sites)
				if len(sites) != 0 {
					return tf.ImportAsExistsError("azurerm_bot_channel_directline", resourceId.ID())
				}
			}
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.DirectLineChannel{
			Properties: &botservice.DirectLineChannelProperties{
				Sites: expandDirectlineSites(d.Get("site").(*pluginsdk.Set).List()),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameDirectLineChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Create(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, botservice.ChannelNameDirectLineChannel, channel); err != nil {
		return fmt.Errorf("creating Directline Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroup, err)
	}
	d.SetId(resourceId.ID())

	// Unable to add a new site with enhanced_authentication_enabled, user_upload_enabled, endpoint_parameters_enabled, storage_enabled in the same operation, so we need to make two calls
	if _, err := client.Update(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, botservice.ChannelNameDirectLineChannel, channel); err != nil {
		return fmt.Errorf("updating Directline Channel for Bot %q (Resource Group %q): %+v", resourceId.BotServiceName, resourceId.ResourceGroup, err)
	}

	return resourceBotChannelDirectlineRead(d, meta)
}

func resourceBotChannelDirectlineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameBasicChannelChannelNameDirectLineChannel))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Directline Channel for Bot %q (Resource Group %q) was not found - removing from state!", id.ResourceGroup, id.BotServiceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Channel Directline for Bot %q (Resource Group %q): %+v", id.ResourceGroup, id.BotServiceName, err)
	}

	channelsResp, err := client.ListWithKeys(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameDirectLineChannel)
	if err != nil {
		return fmt.Errorf("listing Keys for Directline Channel for Bot %q (Resource Group %q): %+v", id.ResourceGroup, id.BotServiceName, err)
	}

	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := channelsResp.Properties; props != nil {
		if channel, ok := props.AsDirectLineChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("site", flattenDirectlineSites(filterSites(channelProps.Sites)))
			}
		}
	}

	return nil
}

func resourceBotChannelDirectlineUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	channel := botservice.BotChannel{
		Properties: botservice.DirectLineChannel{
			Properties: &botservice.DirectLineChannelProperties{
				Sites: expandDirectlineSites(d.Get("site").(*pluginsdk.Set).List()),
			},
			ChannelName: botservice.ChannelNameBasicChannelChannelNameDirectLineChannel,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameDirectLineChannel, channel); err != nil {
		return fmt.Errorf("updating Directline Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
	}

	// Unable to add a new site with enhanced_authentication_enabled, user_upload_enabled, endpoint_parameters_enabled, storage_enabled in the same operation, so we need to make two calls
	// Once this issue https://github.com/Azure/azure-rest-api-specs/issues/25758 is fixed, this update will be removed
	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, botservice.ChannelNameDirectLineChannel, channel); err != nil {
		return fmt.Errorf("updating Directline Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
	}

	return resourceBotChannelDirectlineRead(d, meta)
}

func resourceBotChannelDirectlineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotChannelID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameBasicChannelChannelNameDirectLineChannel))
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting Directline Channel for Bot %q (Resource Group %q): %+v", id.BotServiceName, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandDirectlineSites(input []interface{}) *[]botservice.DirectLineSite {
	sites := make([]botservice.DirectLineSite, 0)

	for _, element := range input {
		if element == nil {
			continue
		}

		site := element.(map[string]interface{})
		expanded := botservice.DirectLineSite{
			IsBlockUserUploadEnabled:    utils.Bool(!site["user_upload_enabled"].(bool)),
			IsEndpointParametersEnabled: utils.Bool(site["endpoint_parameters_enabled"].(bool)),
			IsNoStorageEnabled:          utils.Bool(!site["storage_enabled"].(bool)),
		}

		if v, ok := site["name"].(string); ok {
			expanded.SiteName = &v
		}
		if v, ok := site["enabled"].(bool); ok {
			expanded.IsEnabled = &v
		}
		if v, ok := site["v1_allowed"].(bool); ok {
			expanded.IsV1Enabled = &v
		}
		if v, ok := site["v3_allowed"].(bool); ok {
			expanded.IsV3Enabled = &v
		}
		if v, ok := site["enhanced_authentication_enabled"].(bool); ok {
			expanded.IsSecureSiteEnabled = &v
		}
		if v, ok := site["trusted_origins"].(*pluginsdk.Set); ok {
			origins := v.List()
			items := make([]string, len(origins))
			for i, raw := range origins {
				items[i] = raw.(string)
			}
			expanded.TrustedOrigins = &items
		}

		sites = append(sites, expanded)
	}

	return &sites
}

func flattenDirectlineSites(input []botservice.DirectLineSite) []interface{} {
	sites := make([]interface{}, len(input))

	for i, element := range input {
		site := make(map[string]interface{})

		if v := element.SiteName; v != nil {
			site["name"] = *v
		}

		userUploadEnabled := true
		if v := element.IsBlockUserUploadEnabled; v != nil {
			userUploadEnabled = !*v
		}
		site["user_upload_enabled"] = userUploadEnabled

		var endpointParametersEnabled bool
		if v := element.IsEndpointParametersEnabled; v != nil {
			endpointParametersEnabled = *v
		}
		site["endpoint_parameters_enabled"] = endpointParametersEnabled

		storageEnabled := true
		if v := element.IsNoStorageEnabled; v != nil {
			storageEnabled = !*v
		}
		site["storage_enabled"] = storageEnabled

		if element.Key != nil {
			site["key"] = *element.Key
		}

		if element.Key2 != nil {
			site["key2"] = *element.Key2
		}

		if element.IsEnabled != nil {
			site["enabled"] = *element.IsEnabled
		}

		if element.IsV1Enabled != nil {
			site["v1_allowed"] = *element.IsV1Enabled
		}

		if element.IsV3Enabled != nil {
			site["v3_allowed"] = *element.IsV3Enabled
		}

		if element.IsSecureSiteEnabled != nil {
			site["enhanced_authentication_enabled"] = *element.IsSecureSiteEnabled
		}

		if element.TrustedOrigins != nil {
			site["trusted_origins"] = *element.TrustedOrigins
		}

		sites[i] = site
	}

	return sites
}

// When creating a new directline channel, a Default Site is created
// There is a race condition where this site is not removed before the create request is completed
func filterSites(sites *[]botservice.DirectLineSite) []botservice.DirectLineSite {
	filtered := make([]botservice.DirectLineSite, 0)
	for _, site := range *sites {
		if *site.SiteName == "Default Site" {
			continue
		}
		filtered = append(filtered, site)
	}
	return filtered
}
