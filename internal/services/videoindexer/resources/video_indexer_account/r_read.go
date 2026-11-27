// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package video_indexer_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/videoindexer/2025-04-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	videoindexerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/videoindexer/client"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := videoindexerClient.From(metadata.Client).AccountClient
			id, err := accounts.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			account, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(account.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, account.Model)
		},
	}
}
