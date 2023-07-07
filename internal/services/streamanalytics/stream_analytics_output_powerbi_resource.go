// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/streamingjobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type OutputPowerBIResource struct{}

var (
	_ sdk.ResourceWithCustomImporter = OutputPowerBIResource{}
	_ sdk.ResourceWithStateMigration = OutputPowerBIResource{}
)

type OutputPowerBIResourceModel struct {
	Name                   string `tfschema:"name"`
	StreamAnalyticsJob     string `tfschema:"stream_analytics_job_id"`
	DataSet                string `tfschema:"dataset"`
	Table                  string `tfschema:"table"`
	GroupID                string `tfschema:"group_id"`
	GroupName              string `tfschema:"group_name"`
	TokenUserPrincipalName string `tfschema:"token_user_principal_name"`
	TokenUserDisplayName   string `tfschema:"token_user_display_name"`
}

func (r OutputPowerBIResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"stream_analytics_job_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: streamingjobs.ValidateStreamingJobID,
		},

		"dataset": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"table": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},

		"group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"token_user_principal_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"token_user_display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r OutputPowerBIResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r OutputPowerBIResource) ModelObject() interface{} {
	return &OutputPowerBIResourceModel{}
}

func (r OutputPowerBIResource) ResourceType() string {
	return "azurerm_stream_analytics_output_powerbi"
}

func (r OutputPowerBIResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model OutputPowerBIResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.StreamAnalytics.OutputsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			streamingJobId, err := streamingjobs.ParseStreamingJobID(model.StreamAnalyticsJob)
			if err != nil {
				return err
			}
			id := outputs.NewOutputID(subscriptionId, streamingJobId.ResourceGroupName, streamingJobId.StreamingJobName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			powerBIOutputProps := &outputs.PowerBIOutputDataSourceProperties{
				Dataset:            utils.String(model.DataSet),
				Table:              utils.String(model.Table),
				GroupId:            utils.String(model.GroupID),
				GroupName:          utils.String(model.GroupName),
				RefreshToken:       utils.String("someRefreshToken"),               // A valid refresh token is currently only obtainable via the Azure Portal. Put a dummy string value here when creating the data source and then going to the Azure Portal to authenticate the data source which will update this property with a valid refresh token.
				AuthenticationMode: utils.ToPtr(outputs.AuthenticationMode("Msi")), // Set authentication mode as "Msi" here since other modes requires params obtainable from portal only.
			}

			if model.TokenUserDisplayName != "" {
				powerBIOutputProps.TokenUserDisplayName = utils.String(model.TokenUserDisplayName)
			}

			if model.TokenUserPrincipalName != "" {
				powerBIOutputProps.TokenUserPrincipalName = utils.String(model.TokenUserPrincipalName)
			}

			props := outputs.Output{
				Name: utils.String(model.Name),
				Properties: &outputs.OutputProperties{
					Datasource: &outputs.PowerBIOutputDataSource{
						Properties: powerBIOutputProps,
					},
				},
			}

			var opts outputs.CreateOrReplaceOperationOptions
			if _, err = client.CreateOrReplace(ctx, id, props, opts); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r OutputPowerBIResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.OutputsClient
			id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state OutputPowerBIResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			needUpdateDataSourceProps := false
			dataSourceProps := outputs.PowerBIOutputDataSourceProperties{}
			d := metadata.ResourceData

			if d.HasChange("dataset") {
				needUpdateDataSourceProps = true
				dataSourceProps.Dataset = &state.DataSet
			}

			if d.HasChange("table") {
				needUpdateDataSourceProps = true
				dataSourceProps.Table = &state.Table
			}

			if d.HasChange("group_name") {
				needUpdateDataSourceProps = true
				dataSourceProps.GroupName = &state.GroupName
			}

			if d.HasChange("group_id") {
				needUpdateDataSourceProps = true
				dataSourceProps.GroupId = &state.GroupID
			}

			if d.HasChange("token_user_principal_name") {
				needUpdateDataSourceProps = true
				dataSourceProps.TokenUserPrincipalName = &state.TokenUserPrincipalName
			}

			if d.HasChange("token_user_display_name") {
				needUpdateDataSourceProps = true
				dataSourceProps.TokenUserDisplayName = &state.TokenUserDisplayName
			}

			if !needUpdateDataSourceProps {
				return nil
			}

			updateDataSource := outputs.PowerBIOutputDataSource{
				Properties: &dataSourceProps,
			}

			props := outputs.Output{
				Properties: &outputs.OutputProperties{
					Datasource: updateDataSource,
				},
			}

			var opts outputs.UpdateOperationOptions
			if _, err = client.Update(ctx, *id, props, opts); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r OutputPowerBIResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.OutputsClient
			id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					output, ok := props.Datasource.(outputs.PowerBIOutputDataSource)
					if !ok {
						return fmt.Errorf("converting %s to a PowerBI Output", *id)
					}

					streamingJobId := streamingjobs.NewStreamingJobID(id.SubscriptionId, id.ResourceGroupName, id.StreamingJobName)

					state := OutputPowerBIResourceModel{
						Name:               id.OutputName,
						StreamAnalyticsJob: streamingJobId.ID(),
					}

					dataset := ""
					if v := output.Properties.Dataset; v != nil {
						dataset = *v
					}
					state.DataSet = dataset

					table := ""
					if v := output.Properties.Table; v != nil {
						table = *v
					}
					state.Table = table

					groupId := ""
					if v := output.Properties.GroupId; v != nil {
						groupId = *v
					}
					state.GroupID = groupId

					groupName := ""
					if v := output.Properties.GroupName; v != nil {
						groupName = *v
					}
					state.GroupName = groupName

					state.TokenUserDisplayName = metadata.ResourceData.Get("token_user_display_name").(string)
					state.TokenUserPrincipalName = metadata.ResourceData.Get("token_user_principal_name").(string)

					return metadata.Encode(&state)
				}
			}
			return nil
		},
	}
}

func (r OutputPowerBIResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.OutputsClient
			id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if resp, err := client.Delete(ctx, *id); err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}
			return nil
		},
	}
}

func (r OutputPowerBIResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return outputs.ValidateOutputID
}

func (r OutputPowerBIResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.StreamAnalytics.OutputsClient
		resp, err := client.Get(ctx, *id)
		if err != nil || resp.Model == nil || resp.Model.Properties == nil {
			return fmt.Errorf("reading %s: %+v", *id, err)
		}

		props := resp.Model.Properties
		if _, ok := props.Datasource.(outputs.PowerBIOutputDataSource); !ok {
			return fmt.Errorf("specified output is not of type")
		}
		return nil
	}
}

func (r OutputPowerBIResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsOutputPowerBiV0ToV1{},
		},
	}
}
