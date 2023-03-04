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

type DataConnectorMicrosoftThreatProtectionResource struct{}

var _ sdk.ResourceWithCustomImporter = DataConnectorMicrosoftThreatProtectionResource{}

type DataConnectorMicrosoftThreatProtectionModel struct {
	Name                    string `tfschema:"name"`
	LogAnalyticsWorkspaceId string `tfschema:"log_analytics_workspace_id"`
	TenantId                string `tfschema:"tenant_id"`
}

func (r DataConnectorMicrosoftThreatProtectionResource) Arguments() map[string]*pluginsdk.Schema {
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

func (r DataConnectorMicrosoftThreatProtectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataConnectorMicrosoftThreatProtectionResource) ResourceType() string {
	return "azurerm_sentinel_data_connector_microsoft_threat_protection"
}

func (r DataConnectorMicrosoftThreatProtectionResource) ModelObject() interface{} {
	return &DataConnectorMicrosoftThreatProtectionModel{}
}

func (r DataConnectorMicrosoftThreatProtectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dataconnectors.ValidateDataConnectorID
}

func (r DataConnectorMicrosoftThreatProtectionResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		_, err := importSentinelDataConnector(dataconnectors.DataConnectorKindMicrosoftThreatProtection)(ctx, metadata.ResourceData, metadata.Client)
		return err
	}
}

func (r DataConnectorMicrosoftThreatProtectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			var plan DataConnectorMicrosoftThreatProtectionModel
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

			params := dataconnectors.MTPDataConnector{
				Name: &plan.Name,
				Properties: &dataconnectors.MTPDataConnectorProperties{
					TenantId: tenantId,
					DataTypes: dataconnectors.MTPDataConnectorDataTypes{
						Incidents: dataconnectors.DataConnectorDataTypeCommon{
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

func (r DataConnectorMicrosoftThreatProtectionResource) Read() sdk.ResourceFunc {
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
			dc, ok := modelPtr.(dataconnectors.MTPDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an Microsoft Threat Protection Data Connector", id)
			}

			var tenantId string
			if props := dc.Properties; props != nil {
				tenantId = props.TenantId
			}

			model := DataConnectorMicrosoftThreatProtectionModel{
				Name:                    id.DataConnectorId,
				LogAnalyticsWorkspaceId: workspaceId.ID(),
				TenantId:                tenantId,
			}

			return metadata.Encode(&model)
		},
	}
}

func (r DataConnectorMicrosoftThreatProtectionResource) Delete() sdk.ResourceFunc {
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
