// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_rai_blocklist

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/raiblocklists"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (c Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiBlocklistsClient

			var model cognitiveRaiBlocklistModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountId, err := raiblocklists.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			id, err := raiblocklists.ParseRaiBlocklistID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := existing.Model

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
