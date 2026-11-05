// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package trusted_signing_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/codesigningaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.CodeSigning.Client.CodeSigningAccounts

			id, err := codesigningaccounts.ParseCodeSigningAccountID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			output := TrustedSigningAccountModel{
				Name:              id.CodeSigningAccountName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := result.Model; model != nil {
				output.Location = location.Normalize(model.Location)
				output.Tags = pointer.From(model.Tags)

				if prop := model.Properties; prop != nil {
					output.AccountUri = pointer.From(prop.AccountUri)
					if sku := prop.Sku; sku != nil {
						output.SkuName = string(sku.Name)
					}
				}
			}

			return meta.Encode(&output)
		},
	}
}
