// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_rule_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionruleassociations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := datacollectionruleassociations.ParseScopedDataCollectionRuleAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Monitor.DataCollectionRuleAssociationsClient
			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}
			existing := resp.Model
			if existing.Properties == nil {
				return fmt.Errorf("unexpected null properties of %s", *id)
			}

			var model DataCollectionRuleAssociationModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("data_collection_endpoint_id") {
				if model.DataCollectionEndpointId != "" {
					existing.Properties.DataCollectionEndpointId = pointer.To(model.DataCollectionEndpointId)
				} else {
					existing.Properties.DataCollectionEndpointId = nil
				}
			}

			if metadata.ResourceData.HasChange("data_collection_rule_id") {
				if model.DataCollectionRuleId != "" {
					existing.Properties.DataCollectionRuleId = pointer.To(model.DataCollectionRuleId)
				} else {
					existing.Properties.DataCollectionRuleId = nil
				}
			}

			if metadata.ResourceData.HasChange("description") {
				existing.Properties.Description = pointer.To(model.Description)
			}

			if _, err := client.Create(ctx, *id, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
