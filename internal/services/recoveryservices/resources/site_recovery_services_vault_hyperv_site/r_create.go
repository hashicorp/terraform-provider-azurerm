// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_services_vault_hyperv_site

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel HyperVSiteModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding %s", err)
			}

			client := metadata.Client.RecoveryServices.FabricClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			vaultId, err := replicationfabrics.ParseVaultID(metaModel.RecoveryVaultId)
			if err != nil {
				return err
			}

			id := replicationfabrics.NewReplicationFabricID(subscriptionId, vaultId.ResourceGroupName, vaultId.VaultName, metaModel.Name)

			// the instance type `HyperVSite` is not exposed in Swagger, tracked on https://github.com/Azure/azure-rest-api-specs/issues/22016
			parameters := replicationfabrics.FabricCreationInput{
				Properties: &replicationfabrics.FabricCreationInputProperties{
					CustomDetails: replicationfabrics.BaseFabricSpecificCreationInputImpl{
						InstanceType: "HyperVSite",
					},
				},
			}

			if err := client.CreateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
