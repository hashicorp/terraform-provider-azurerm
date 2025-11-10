// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/jackofallops/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type DataConnectorOfficeIRMResource struct{}

var _ sdk.ResourceWithCustomImporter = DataConnectorOfficeIRMResource{}

type DataConnectorOfficeIRMModel struct {
	Name                    string `tfschema:"name"`
	LogAnalyticsWorkspaceId string `tfschema:"log_analytics_workspace_id"`
	TenantId                string `tfschema:"tenant_id"`
}

func (r DataConnectorOfficeIRMResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"tenant_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}

func (r DataConnectorOfficeIRMResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataConnectorOfficeIRMResource) ResourceType() string {
	return "azurerm_sentinel_data_connector_office_irm"
}

func (r DataConnectorOfficeIRMResource) ModelObject() interface{} {
	return &DataConnectorOfficeIRMModel{}
}

func (r DataConnectorOfficeIRMResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataConnectorID
}

func (r DataConnectorOfficeIRMResource) CustomImporter() sdk.ResourceRunFunc {
	return importDataConnectorTyped(securityinsight.DataConnectorKindOfficeIRM)
}

func (r DataConnectorOfficeIRMResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			var plan DataConnectorOfficeIRMModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			workspaceId, err := workspaces.ParseWorkspaceID(plan.LogAnalyticsWorkspaceId)
			if err != nil {
				return err
			}

			id := parse.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, plan.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			tenantId := plan.TenantId
			if tenantId == "" {
				tenantId = metadata.Client.Account.TenantId
			}

			params := securityinsight.OfficeIRMDataConnector{
				Name: &plan.Name,
				OfficeIRMDataConnectorProperties: &securityinsight.OfficeIRMDataConnectorProperties{
					TenantID: &tenantId,
					DataTypes: &securityinsight.AlertsDataTypeOfDataConnector{
						Alerts: &securityinsight.DataConnectorDataTypeCommon{
							State: securityinsight.DataTypeStateEnabled,
						},
					},
				},
				Kind: securityinsight.KindBasicDataConnectorKindOfficeIRM,
			}

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataConnectorOfficeIRMResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient
			id, err := parse.DataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)

			existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			dc, ok := existing.Value.(securityinsight.OfficeIRMDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an Office IRM Data Connector", id)
			}

			var tenantId string
			if props := dc.OfficeIRMDataConnectorProperties; props != nil {
				if props.TenantID != nil {
					tenantId = *props.TenantID
				}
			}

			model := DataConnectorOfficeIRMModel{
				Name:                    id.Name,
				LogAnalyticsWorkspaceId: workspaceId.ID(),
				TenantId:                tenantId,
			}

			return metadata.Encode(&model)
		},
	}
}

func (r DataConnectorOfficeIRMResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			id, err := parse.DataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
