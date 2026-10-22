// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package federated_identity_credential

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/federatedidentitycredentials"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.V20241130.FederatedIdentityCredentials

			var config FederatedIdentityCredentialResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parentId, err := commonids.ParseUserAssignedIdentityID(config.UserAssignedIdentityId)
			if err != nil {
				return err
			}

			locks.ByID(parentId.ID())
			defer locks.UnlockByID(parentId.ID())

			id, err := federatedidentitycredentials.ParseFederatedIdentityCredentialID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
