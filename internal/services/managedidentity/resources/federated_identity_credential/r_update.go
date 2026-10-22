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

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, rmd sdk.ResourceMetaData) error {
			client := rmd.Client.ManagedIdentity.V20241130.FederatedIdentityCredentials

			var config FederatedIdentityCredentialResourceSchema
			if err := rmd.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := federatedidentitycredentials.ParseFederatedIdentityCredentialID(rmd.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil`", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}
			props := existing.Model.Properties

			if rmd.ResourceData.HasChange("audience") {
				props.Audiences = config.Audience
			}

			if rmd.ResourceData.HasChange("issuer") {
				props.Issuer = config.Issuer
			}

			if rmd.ResourceData.HasChange("subject") {
				props.Subject = config.Subject
			}

			parentId := commonids.NewUserAssignedIdentityID(id.SubscriptionId, id.ResourceGroupName, id.UserAssignedIdentityName)
			locks.ByID(parentId.ID())
			defer locks.UnlockByID(parentId.ID())

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
