// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_rule_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionruleassociations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DataCollectionRuleAssociationModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.Monitor.DataCollectionRuleAssociationsClient

			id := datacollectionruleassociations.NewScopedDataCollectionRuleAssociationID(model.TargetResourceId, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			input := datacollectionruleassociations.DataCollectionRuleAssociationProxyOnlyResource{
				Name: pointer.To(model.Name),
				Properties: &datacollectionruleassociations.DataCollectionRuleAssociation{
					Description: pointer.To(model.Description),
				},
			}

			if model.DataCollectionEndpointId != "" {
				input.Properties.DataCollectionEndpointId = pointer.To(model.DataCollectionEndpointId)
			}
			if model.DataCollectionRuleId != "" {
				input.Properties.DataCollectionRuleId = pointer.To(model.DataCollectionRuleId)
			}

			if _, err := client.Create(ctx, id, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
