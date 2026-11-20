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

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionRuleAssociationsClient
			id, err := datacollectionruleassociations.ParseScopedDataCollectionRuleAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var description, dataCollectionEndpointId, dataCollectionRuleId string

			if model := resp.Model; model != nil {
				if prop := model.Properties; prop != nil {
					dataCollectionEndpointId = pointer.From(prop.DataCollectionEndpointId)
					dataCollectionRuleId = pointer.From(prop.DataCollectionRuleId)
					description = pointer.From(prop.Description)
				}
			}

			return metadata.Encode(&DataCollectionRuleAssociationModel{
				Name:                     id.DataCollectionRuleAssociationName,
				TargetResourceId:         id.ResourceUri,
				DataCollectionEndpointId: dataCollectionEndpointId,
				DataCollectionRuleId:     dataCollectionRuleId,
				Description:              description,
			})
		},
		Timeout: 5 * time.Minute,
	}
}
