// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package trusted_signing_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/codesigningaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.CodeSigning.Client.CodeSigningAccounts

			id, err := codesigningaccounts.ParseCodeSigningAccountID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}
