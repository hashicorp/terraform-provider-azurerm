package sentinel

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type DataConnectorGenericUIResourceModel struct {
	Name                    string `tfschema:"name"`
	LogAnalyticsWorkspaceId string `tfschema:"log_analytics_workspace_id"`
	// below are copied from `DataConnectorGenericUIConfigModel`
	Title                 string                                            `tfschema:"title"`
	Publisher             string                                            `tfschema:"publisher"`
	DescriptionMarkdown   string                                            `tfschema:"description_markdown"`
	CustomImage           string                                            `tfschema:"custom_image"`
	GraphQueriesTableName string                                            `tfschema:"graph_queries_table_name"`
	GraphQueries          []DataConnectorGenericUIGraphQueryModel           `tfschema:"graph_query"`
	SampleQueries         []DataConnectorGenericUISampleQueryModel          `tfschema:"sample_query"`
	DataTypes             []DataConnectorGenericUIDataTypeModel             `tfschema:"data_type"`
	ConnectivityCriteria  []DataConnectorGenericUIConnectivityCriteriaModel `tfschema:"connectivity_criteria"`
	Availability          []DataConnectorGenericUIAvailabilityModel         `tfschema:"availability"`
	Permissions           []DataConnectorGenericUIPermissionsModel          `tfschema:"permission"`
	InstructionSteps      []DataConnectorGenericUIInstructionStepModel      `tfschema:"instruction"`
}

type DataConnectorGenericUIResource struct{}

var _ sdk.ResourceWithUpdate = DataConnectorGenericUIResource{}

func (r DataConnectorGenericUIResource) Arguments() map[string]*schema.Schema {
	arg := map[string]*schema.Schema{
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
	}

	for k, v := range DataConnectorGenericUIConfigSchema() {
		arg[k] = v
	}
	return arg
}

func (r DataConnectorGenericUIResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r DataConnectorGenericUIResource) ModelObject() interface{} {
	return &DataConnectorGenericUIConfigModel{}
}

func (r DataConnectorGenericUIResource) ResourceType() string {
	return "azurerm_sentinel_data_connector_generic_ui"
}

func (r DataConnectorGenericUIResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataConnectorID
}

func (r DataConnectorGenericUIResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			var plan DataConnectorGenericUIResourceModel
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

			uiConfig, err := expandDataConnectorGenericUIResourceModel(plan)
			if err != nil {
				return fmt.Errorf("expanding: %+v", err)
			}

			params := securityinsight.CodelessUIDataConnector{
				Name: &plan.Name,
				CodelessParameters: &securityinsight.CodelessParameters{
					ConnectorUIConfig: uiConfig,
				},
				Kind: securityinsight.KindBasicDataConnectorKindGenericUI,
			}

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataConnectorGenericUIResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient
			id, err := parse.DataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			dc, ok := existing.Value.(securityinsight.CodelessUIDataConnector)
			if !ok {
				return fmt.Errorf("%s was not an IoT Data Connector", id)
			}

			if dc.ConnectorUIConfig == nil {
				return fmt.Errorf("retrieving %s: `ConnectorUIConfig` was nil", id)
			}

			model, err := flattenDataConnectorGenericUIResourceModel(dc.ConnectorUIConfig)
			if err != nil {
				return fmt.Errorf("flattening: %+v", err)
			}

			if dc.Name != nil {
				model.Name = *dc.Name
			}

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
			model.LogAnalyticsWorkspaceId = workspaceId.ID()

			return metadata.Encode(&model)
		},
	}
}

func (r DataConnectorGenericUIResource) Delete() sdk.ResourceFunc {
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

func (r DataConnectorGenericUIResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient
			id, err := parse.DataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan DataConnectorGenericUIResourceModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			_, err = client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			uiConfig, err := expandDataConnectorGenericUIResourceModel(plan)
			if err != nil {
				return fmt.Errorf("expanding : %+v", err)
			}

			params := securityinsight.CodelessUIDataConnector{
				Name: &plan.Name,
				CodelessParameters: &securityinsight.CodelessParameters{
					ConnectorUIConfig: uiConfig,
				},
				Kind: securityinsight.KindBasicDataConnectorKindGenericUI,
			}

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, params); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandDataConnectorGenericUIResourceModel(input DataConnectorGenericUIResourceModel) (*securityinsight.CodelessUIConnectorConfigProperties, error) {
	uiModel := DataConnectorGenericUIConfigModel{
		Title:                 input.Title,
		Publisher:             input.Publisher,
		DescriptionMarkdown:   input.DescriptionMarkdown,
		CustomImage:           input.CustomImage,
		GraphQueriesTableName: input.GraphQueriesTableName,
		GraphQueries:          input.GraphQueries,
		SampleQueries:         input.SampleQueries,
		DataTypes:             input.DataTypes,
		ConnectivityCriteria:  input.ConnectivityCriteria,
		Availability:          input.Availability,
		Permissions:           input.Permissions,
		InstructionSteps:      input.InstructionSteps,
	}
	return expandDataConnectorGenericUIConfigModel(append([]DataConnectorGenericUIConfigModel{}, uiModel))
}

func flattenDataConnectorGenericUIResourceModel(input *securityinsight.CodelessUIConnectorConfigProperties) (DataConnectorGenericUIResourceModel, error) {
	model := DataConnectorGenericUIResourceModel{}
	uiConfig, err := flattenDataConnectorGenericUIConfigModel(input)
	if err != nil {
		return model, err
	}

	model.Title = uiConfig[0].Title
	model.Publisher = uiConfig[0].Publisher
	model.DescriptionMarkdown = uiConfig[0].DescriptionMarkdown
	model.CustomImage = uiConfig[0].CustomImage
	model.GraphQueriesTableName = uiConfig[0].GraphQueriesTableName
	model.GraphQueries = uiConfig[0].GraphQueries
	model.SampleQueries = uiConfig[0].SampleQueries
	model.DataTypes = uiConfig[0].DataTypes
	model.ConnectivityCriteria = uiConfig[0].ConnectivityCriteria
	model.Availability = uiConfig[0].Availability
	model.Permissions = uiConfig[0].Permissions
	model.InstructionSteps = uiConfig[0].InstructionSteps

	return model, nil
}
