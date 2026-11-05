// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package trusted_signing_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/codesigningaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CodeSigning.Client.CodeSigningAccounts

			var model TrustedSigningAccountModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionID := metadata.Client.Account.SubscriptionId
			id := codesigningaccounts.NewCodeSigningAccountID(subscriptionID, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("retrieving %s: %v", id, err)
					}
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			req := codesigningaccounts.CodeSigningAccount{
				Name:     &model.Name,
				Location: location.Normalize(model.Location),
				Tags:     &model.Tags,
				Properties: &codesigningaccounts.CodeSigningAccountProperties{
					Sku: &codesigningaccounts.AccountSku{
						Name: codesigningaccounts.SkuName(model.SkuName),
					},
				},
			}

			if err := client.CreateCallbackThenPoll(ctx, id, req, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
