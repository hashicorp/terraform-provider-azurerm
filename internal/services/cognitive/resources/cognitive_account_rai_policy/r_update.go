// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_rai_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/raipolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiPoliciesClient

			id, err := raipolicies.ParseAccountRaiPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model AccountRaiPolicyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			cognitiveAccountId := raipolicies.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)

			locks.ByID(cognitiveAccountId.ID())
			defer locks.UnlockByID(cognitiveAccountId.ID())

			payload := existing.Model

			if metadata.ResourceData.HasChange("content_filter") {
				payload.Properties.ContentFilters = expandRaiPolicyContentFilters(model.ContentFilter)
			}

			if metadata.ResourceData.HasChange("mode") {
				payload.Properties.Mode = pointer.ToEnum[raipolicies.RaiPolicyMode](model.Mode)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
