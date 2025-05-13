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

type DataConnectorIOTResource struct{}

var _ sdk.ResourceWithCustomImporter = DataConnectorIOTResource{}

type DataConnectorIOTModel struct {
	Name                    string `tfschema:"name"`
	LogAnalyticsWorkspaceId string `tfschema:"log_analytics_workspace_id"`
	SubscriptionId          string `tfschema:"subscription_id"`
}

func (r DataConnectorIOTResource) Arguments() map[string]*pluginsdk.Schema {
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

		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}

func (r DataConnectorIOTResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataConnectorIOTResource) ResourceType() string {
	return "azurerm_sentinel_data_connector_iot"
}

func (r DataConnectorIOTResource) ModelObject() interface{} {
	return &DataConnectorIOTModel{}
}

func (r DataConnectorIOTResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataConnectorID
}

func (r DataConnectorIOTResource) CustomImporter() sdk.ResourceRunFunc {
	return importDataConnectorTyped(securityinsight.DataConnectorKindIOT)
}

func (r DataConnectorIOTResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			var plan DataConnectorIOTModel
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

			subscriptionId := plan.SubscriptionId
			if subscriptionId == "" {
				subscriptionId = metadata.Client.Account.SubscriptionId
			}

			params := securityinsight.IoTDataConnector{
				Name: &plan.Name,
				IoTDataConnectorProperties: &securityinsight.IoTDataConnectorProperties{
					SubscriptionID: &subscriptionId,
					DataTypes: &securityinsight.AlertsDataTypeOfDataConnector{
						Alerts: &securityinsight.DataConnectorDataTypeCommon{
							State: securityinsight.DataTypeStateEnabled,
						},
					},
				},
				Kind: securityinsight.KindBasicDataConnectorKindIOT,
			}

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataConnectorIOTResource) Read() sdk.ResourceFunc {
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

			dc, ok := existing.Value.(securityinsight.IoTDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an IoT Data Connector", id)
			}

			var subscriptionId string
			if props := dc.IoTDataConnectorProperties; props != nil {
				if props.SubscriptionID != nil {
					subscriptionId = *props.SubscriptionID
				}
			}

			model := DataConnectorIOTModel{
				Name:                    id.Name,
				LogAnalyticsWorkspaceId: workspaceId.ID(),
				SubscriptionId:          subscriptionId,
			}

			return metadata.Encode(&model)
		},
	}
}

func (r DataConnectorIOTResource) Delete() sdk.ResourceFunc {
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
