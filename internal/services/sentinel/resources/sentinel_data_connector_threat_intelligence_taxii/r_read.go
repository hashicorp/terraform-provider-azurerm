// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_threat_intelligence_taxii

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := sentinelClient.From(metadata.Client).DataConnectorsClient

			var state DataConnectorThreatIntelligenceTAXIIModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, err := dataconnectors.ParseDataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			dc, ok := existing.Model.(dataconnectors.TiTaxiiDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an Threat Intelligence TAXII Data Connector", id)
			}

			model := DataConnectorThreatIntelligenceTAXIIModel{
				Name:                    id.DataConnectorId,
				LogAnalyticsWorkspaceId: workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID(),
				UserName:                state.UserName, // setting the user name from state, as it is not returned from API
				Password:                state.Password, // setting the password from state, as it is not returned from API
			}

			if props := dc.Properties; props != nil {
				if props.FriendlyName != nil {
					model.DisplayName = *props.FriendlyName
				}

				if props.TaxiiServer != nil {
					model.APIRootURL = *props.TaxiiServer
				}

				if props.CollectionId != nil {
					model.CollectionID = *props.CollectionId
				}

				model.PollingFrequency = string(props.PollingFrequency)

				if lookback, err := props.GetTaxiiLookbackPeriodAsTime(); err == nil && lookback != nil {
					model.LookbackDate = lookback.Format(time.RFC3339)
				}

				model.TenantId = props.TenantId
			}

			return metadata.Encode(&model)
		},
	}
}
