// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_microsoft_threat_intelligence

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := sentinelClient.From(metadata.Client).DataConnectorsClient
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

			dc, ok := existing.Model.(dataconnectors.MSTIDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an Microsoft Threat Protection Data Connector", id)
			}

			state := DataConnectorMicrosoftThreatIntelligenceModel{
				Name:        id.DataConnectorId,
				WorkspaceId: workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID(),
			}

			if props := dc.Properties; props != nil {
				state.TenantId = props.TenantId

				dt := props.DataTypes
				if metf := dt.MicrosoftEmergingThreatFeed; metf.State != nil {
					if strings.EqualFold(string(*metf.State), string(dataconnectors.DataTypeStateEnabled)) {
						state.MicrosoftEmergingThreatFeedLookBackDate, err = flattenSentinelDataConnectorMicrosoftThreatIntelligenceTime(metf.LookbackPeriod)
						if err != nil {
							return fmt.Errorf("flattening `microsoft_emerging_threat_feed`: %+v", err)
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenSentinelDataConnectorMicrosoftThreatIntelligenceTime(input string) (string, error) {
	// TODO: check if this workaround could be removed in the future
	t, err := time.Parse(time.RFC3339, input)
	if err != nil {
		t, err = time.Parse("01/02/2006 15:04:05", input)
		if err != nil {
			return "", err
		}
	}

	return t.Format(time.RFC3339), nil
}
