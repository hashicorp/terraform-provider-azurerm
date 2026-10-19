// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_rai_blocklist

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/raiblocklists"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (c Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiBlocklistsClient

			id, err := raiblocklists.ParseRaiBlocklistID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := cognitiveRaiBlocklistModel{
				Name:               id.RaiBlocklistName,
				CognitiveAccountId: raiblocklists.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
			}

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Description = pointer.From(props.Description)
				}
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}
