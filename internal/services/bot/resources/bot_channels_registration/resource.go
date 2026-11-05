// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channels_registration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_channels_registration/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
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

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			// v0 -> v1 normalises the casing of IDs imported while this resource parsed them with the
			// case-insensitive legacy parser, so they can be parsed with the case-sensitive SDK parser
			0: migration.BotChannelsRegistrationV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := commonids.ParseBotServiceID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).Bot.BotClient

			id, err := commonids.ParseBotServiceID(d.Id())
			if err != nil {
				return nil, err
			}

			resp, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName)
			if err != nil {
				if response.WasNotFound(resp.Response.Response) {
					return nil, fmt.Errorf("the Bot Channels Registration %q was not found in Resource Group %q", id.BotServiceName, id.ResourceGroupName)
				}

				return nil, fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if resp.Kind != botservice.KindBot {
				return nil, fmt.Errorf(" %s was not a Channel Registration - got %q", id, string(resp.Kind))
			}

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Schema: botChannelsRegistrationSchema(),
	}
}
