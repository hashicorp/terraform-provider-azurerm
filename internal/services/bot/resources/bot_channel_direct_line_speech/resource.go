// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_channel_direct_line_speech

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_channel_direct_line_speech/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBotChannelDirectLineSpeechCreate,
		Read:   resourceBotChannelDirectLineSpeechRead,
		Delete: resourceBotChannelDirectLineSpeechDelete,
		Update: resourceBotChannelDirectLineSpeechUpdate,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseBotServiceChannelID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			// v0 -> v1 normalises the casing of IDs imported while this resource parsed them with the
			// case-insensitive legacy parser, so they can be parsed with the case-sensitive SDK parser
			0: migration.BotChannelDirectLineSpeechV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: botChannelDirectLineSpeechSchema(),

		CustomizeDiff: pluginsdk.ForceNewIf("cognitive_account_id", func(ctx context.Context, d *schema.ResourceDiff, _ any) bool {
			if d == nil {
				return false
			}

			o, n := d.GetChange("cognitive_account_id")
			if o.(string) != "" && n.(string) == "" {
				// return `IsKnown` in case of unknown values to prevent replacements on valid but unknown at plan time values
				return d.GetRawConfig().AsValueMap()["cognitive_account_id"].IsKnown()
			}

			return false
		}),
	}
}
