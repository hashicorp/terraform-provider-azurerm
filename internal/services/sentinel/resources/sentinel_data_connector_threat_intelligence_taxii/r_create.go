// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_threat_intelligence_taxii

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient
			wspClient := metadata.Client.LogAnalytics.WorkspaceClient

			var plan DataConnectorThreatIntelligenceTAXIIModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			workspaceId, err := workspaces.ParseWorkspaceID(plan.LogAnalyticsWorkspaceId)
			if err != nil {
				return err
			}

			wsp, err := wspClient.Get(ctx, *workspaceId)
			if err != nil {
				return fmt.Errorf("retrieving the workspace %s: %+v", workspaceId, err)
			}
			if wsp.Model == nil {
				return fmt.Errorf("nil model of the workspace %s", workspaceId)
			}
			if wsp.Model.Properties == nil {
				return fmt.Errorf("nil properties of the workspace %s", workspaceId)
			}
			if wsp.Model.Properties.CustomerId == nil {
				return fmt.Errorf("nil workspace id of the workspace %s", workspaceId)
			}

			id := dataconnectors.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, plan.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			tenantId := plan.TenantId
			if tenantId == "" {
				tenantId = metadata.Client.Account.TenantId
			}

			// Format is guaranteed by schema validation
			lookbackDate, _ := time.Parse(time.RFC3339, plan.LookbackDate)

			params := dataconnectors.TiTaxiiDataConnector{
				Name: &plan.Name,
				Properties: &dataconnectors.TiTaxiiDataConnectorProperties{
					WorkspaceId:         wsp.Model.Properties.CustomerId,
					FriendlyName:        &plan.DisplayName,
					TaxiiServer:         &plan.APIRootURL,
					CollectionId:        &plan.CollectionID,
					PollingFrequency:    dataconnectors.PollingFrequency(plan.PollingFrequency),
					TaxiiLookbackPeriod: pointer.To(lookbackDate.Format(time.RFC3339)),
					DataTypes: dataconnectors.TiTaxiiDataConnectorDataTypes{
						TaxiiClient: dataconnectors.DataConnectorDataTypeCommon{
							State: dataconnectors.DataTypeStateEnabled,
						},
					},
					TenantId: tenantId,
				},
				Kind: dataconnectors.DataConnectorKindThreatIntelligenceTaxii,
			}

			if plan.UserName != "" {
				params.Properties.UserName = &plan.UserName
			}

			if plan.Password != "" {
				params.Properties.Password = &plan.Password
			}

			if _, err = client.CreateOrUpdate(ctx, id, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
