// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_microsoft_threat_intelligence

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := sentinelClient.From(metadata.Client).DataConnectorsClient
			var metaModel DataConnectorMicrosoftThreatIntelligenceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			workSpaceId, err := workspaces.ParseWorkspaceID(metaModel.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing workspace id %+v", err)
			}

			id := dataconnectors.NewDataConnectorID(workSpaceId.SubscriptionId, workSpaceId.ResourceGroupName, workSpaceId.WorkspaceName, metaModel.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(s.ResourceType(), id)
				}
			}

			tenantId := metaModel.TenantId
			if tenantId == "" {
				tenantId = account.From(metadata.Client).TenantId
			}

			dataConnector := dataconnectors.MSTIDataConnector{
				Name: &id.DataConnectorId,
				Kind: dataconnectors.DataConnectorKindMicrosoftThreatIntelligence,
				Properties: &dataconnectors.MSTIDataConnectorProperties{
					DataTypes: dataconnectors.MSTIDataConnectorDataTypes{
						MicrosoftEmergingThreatFeed: expandSentinelDataConnectorMicrosoftThreatIntelligenceMicrosoftEmergingThreatFeed(metaModel),
					},
					TenantId: tenantId,
				},
			}
			if _, err = client.CreateOrUpdate(ctx, id, dataConnector); err != nil {
				return fmt.Errorf("creating %+v", err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandSentinelDataConnectorMicrosoftThreatIntelligenceMicrosoftEmergingThreatFeed(input DataConnectorMicrosoftThreatIntelligenceModel) dataconnectors.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed {
	if input.MicrosoftEmergingThreatFeedLookBackDate == "" {
		return dataconnectors.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed{
			State: pointer.To(dataconnectors.DataTypeStateDisabled),
		}
	}

	return dataconnectors.MSTIDataConnectorDataTypesMicrosoftEmergingThreatFeed{
		LookbackPeriod: input.MicrosoftEmergingThreatFeedLookBackDate,
		State:          pointer.To(dataconnectors.DataTypeStateEnabled),
	}
}
