// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_hybrid_runbook_worker_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/hybridrunbookworkergroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	automationClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/client"
)

func (m Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := automationClient.From(metadata.Client).HybridRunbookWorkerGroup

			var model HybridRunbookWorkerGroupModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			subscriptionID := account.From(metadata.Client).SubscriptionId
			id := hybridrunbookworkergroup.NewHybridRunbookWorkerGroupID(subscriptionID, model.ResourceGroupName,
				model.AutomationAccountName, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("retrieving %s: %v", id, err)
					}
					return metadata.ResourceRequiresImport(m.ResourceType(), id)
				}
			}
			req := hybridrunbookworkergroup.HybridRunbookWorkerGroupCreateOrUpdateParameters{
				Name: pointer.To(model.Name),
			}
			if model.CredentialName != "" {
				req.Properties = &hybridrunbookworkergroup.HybridRunbookWorkerGroupCreateOrUpdateProperties{}
				req.Properties.Credential = &hybridrunbookworkergroup.RunAsCredentialAssociationProperty{
					Name: pointer.To(model.CredentialName),
				}
			}
			// return 201 cause err in autorest sdk
			// Workaround swagger issue https://github.com/Azure/azure-rest-api-specs/issues/19741
			if _, err := client.Create(ctx, id, req); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
