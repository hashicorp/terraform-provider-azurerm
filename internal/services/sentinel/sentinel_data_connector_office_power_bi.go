package sentinel

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataConnectorOfficePowerBIResource struct{}

var _ sdk.ResourceWithCustomImporter = DataConnectorOfficePowerBIResource{}

type DataConnectorOfficePowerBIModel struct {
	Name                    string `tfschema:"name"`
	LogAnalyticsWorkspaceId string `tfschema:"log_analytics_workspace_id"`
	TenantId                string `tfschema:"tenant_id"`
}

func (r DataConnectorOfficePowerBIResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: dataconnectors.ValidateWorkspaceID,
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

func (r DataConnectorOfficePowerBIResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataConnectorOfficePowerBIResource) ResourceType() string {
	return "azurerm_sentinel_data_connector_office_power_bi"
}

func (r DataConnectorOfficePowerBIResource) ModelObject() interface{} {
	return &DataConnectorOfficePowerBIModel{}
}

func (r DataConnectorOfficePowerBIResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dataconnectors.ValidateDataConnectorID
}

func (r DataConnectorOfficePowerBIResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		_, err := importSentinelDataConnector(dataconnectors.DataConnectorKindOfficePowerBI)(ctx, metadata.ResourceData, metadata.Client)
		return err
	}
}

func (r DataConnectorOfficePowerBIResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			var plan DataConnectorOfficePowerBIModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			workspaceId, err := dataconnectors.ParseWorkspaceID(plan.LogAnalyticsWorkspaceId)
			if err != nil {
				return err
			}

			id := dataconnectors.NewDataConnectorID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, plan.Name)
			existing, err := client.DataConnectorsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			tenantId := plan.TenantId
			if tenantId == "" {
				tenantId = metadata.Client.Account.TenantId
			}

			params := dataconnectors.OfficePowerBIDataConnector{
				Name: &plan.Name,
				Properties: &dataconnectors.OfficePowerBIDataConnectorProperties{
					TenantId: tenantId,
					DataTypes: dataconnectors.OfficePowerBIConnectorDataTypes{
						Logs: dataconnectors.DataConnectorDataTypeCommon{
							State: dataconnectors.DataTypeStateEnabled,
						},
					},
				},
			}

			if _, err = client.DataConnectorsCreateOrUpdate(ctx, id, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataConnectorOfficePowerBIResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient
			id, err := dataconnectors.ParseDataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			workspaceId := dataconnectors.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

			existing, err := client.DataConnectorsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("model was nil for %s", id)
			}

			modelPtr := *existing.Model
			dc, ok := modelPtr.(dataconnectors.OfficePowerBIDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an Office PowerBI Data Connector", id)
			}

			var tenantId string
			if props := dc.Properties; props != nil {
				tenantId = props.TenantId
			}

			model := DataConnectorOfficePowerBIModel{
				Name:                    id.DataConnectorId,
				LogAnalyticsWorkspaceId: workspaceId.ID(),
				TenantId:                tenantId,
			}

			return metadata.Encode(&model)
		},
	}
}

func (r DataConnectorOfficePowerBIResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			id, err := dataconnectors.ParseDataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.DataConnectorsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
