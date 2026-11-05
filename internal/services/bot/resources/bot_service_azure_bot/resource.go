// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package bot_service_azure_bot

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/resources/bot_service_azure_bot/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/botservice/2021-05-01-preview/botservice"
)

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
	_ sdk.ResourceWithStateMigration = Resource{}
)

type Resource struct{}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateBotServiceID
}

func (r Resource) ResourceType() string {
	return "azurerm_bot_service_azure_bot"
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			// v0 -> v1 normalises the casing of IDs imported while this resource parsed them with the
			// case-insensitive legacy parser, so they can be parsed with the case-sensitive SDK parser
			0: migration.BotServiceAzureBotV0ToV1{},
		},
	}
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.Bot.BotClient

		id, err := commonids.ParseBotServiceID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroupName, id.BotServiceName)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if actualKind := string(resp.Kind); actualKind != string(botservice.KindAzurebot) {
			return fmt.Errorf("bot has mismatched type, expected: %q, got %q", string(botservice.KindAzurebot), actualKind)
		}

		return nil
	}
}
