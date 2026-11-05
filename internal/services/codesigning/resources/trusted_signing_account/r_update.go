// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package trusted_signing_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/codesigningaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.CodeSigning.Client.CodeSigningAccounts
			id, err := codesigningaccounts.ParseCodeSigningAccountID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			var model TrustedSigningAccountModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			var patch codesigningaccounts.CodeSigningAccountPatch
			if meta.ResourceData.HasChange("sku_name") {
				patch.Properties = pointer.To(codesigningaccounts.CodeSigningAccountPatchProperties{
					Sku: pointer.To(codesigningaccounts.AccountSkuPatch{
						Name: pointer.ToEnum[codesigningaccounts.SkuName](model.SkuName),
					}),
				})
			}
			if meta.ResourceData.HasChange("tags") {
				patch.Tags = pointer.To(model.Tags)
			}

			if err = client.UpdateThenPoll(ctx, *id, patch); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}
