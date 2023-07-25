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
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type DataConnectorAwsS3Resource struct{}

var _ sdk.ResourceWithUpdate = DataConnectorAwsS3Resource{}
var _ sdk.ResourceWithCustomImporter = DataConnectorAwsS3Resource{}

type DataConnectorAwsS3Model struct {
	Name                    string   `tfschema:"name"`
	LogAnalyticsWorkspaceId string   `tfschema:"log_analytics_workspace_id"`
	AwsRoleArm              string   `tfschema:"aws_role_arn"`
	DestinationTable        string   `tfschema:"destination_table"`
	SqsUrls                 []string `tfschema:"sqs_urls"`
}

func (r DataConnectorAwsS3Resource) Arguments() map[string]*pluginsdk.Schema {
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

		"aws_role_arn": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.IsARN,
		},

		"destination_table": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sqs_urls": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func (r DataConnectorAwsS3Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataConnectorAwsS3Resource) ResourceType() string {
	return "azurerm_sentinel_data_connector_aws_s3"
}

func (r DataConnectorAwsS3Resource) ModelObject() interface{} {
	return &DataConnectorAwsS3Model{}
}

func (r DataConnectorAwsS3Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DataConnectorID
}

func (r DataConnectorAwsS3Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		_, err := importSentinelDataConnector(securityinsight.DataConnectorKindAmazonWebServicesS3)(ctx, metadata.ResourceData, metadata.Client)
		return err
	}
}

func (r DataConnectorAwsS3Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.DataConnectorsClient

			var plan DataConnectorAwsS3Model
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

			params := securityinsight.AwsS3DataConnector{
				Name: &plan.Name,
				AwsS3DataConnectorProperties: &securityinsight.AwsS3DataConnectorProperties{
					DestinationTable: utils.String(plan.DestinationTable),
					SqsUrls:          &plan.SqsUrls,
					RoleArn:          utils.String(plan.AwsRoleArm),
					DataTypes: &securityinsight.AwsS3DataConnectorDataTypes{
						Logs: &securityinsight.AwsS3DataConnectorDataTypesLogs{
							State: securityinsight.DataTypeStateEnabled,
						},
					},
				},
				Kind: securityinsight.KindBasicDataConnectorKindAmazonWebServicesS3,
			}
			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataConnectorAwsS3Resource) Read() sdk.ResourceFunc {
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

			dc, ok := existing.Value.(securityinsight.AwsS3DataConnector)
			if !ok {
				return fmt.Errorf("%s was not an AWS S3 Data Connector", id)
			}

			model := DataConnectorAwsS3Model{
				Name:                    id.Name,
				LogAnalyticsWorkspaceId: workspaceId.ID(),
			}

			if prop := dc.AwsS3DataConnectorProperties; prop != nil {
				if prop.RoleArn != nil {
					model.AwsRoleArm = *prop.RoleArn
				}
				if prop.DestinationTable != nil {
					model.DestinationTable = *prop.DestinationTable
				}
				if prop.SqsUrls != nil {
					model.SqsUrls = *prop.SqsUrls
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (DataConnectorAwsS3Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.DataConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan DataConnectorAwsS3Model
			if err := metadata.Decode(&plan); err != nil {
				return err
			}

			client := metadata.Client.Sentinel.DataConnectorsClient

			resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			params, ok := resp.Value.(securityinsight.AwsS3DataConnector)
			if !ok {
				return fmt.Errorf("%s was not an AWS S3 Data Connector", id)
			}

			if props := params.AwsS3DataConnectorProperties; props != nil {
				if metadata.ResourceData.HasChange("aws_role_arn") {
					props.RoleArn = &plan.AwsRoleArm
				}
				if metadata.ResourceData.HasChange("destination_table") {
					props.DestinationTable = &plan.DestinationTable
				}
				if metadata.ResourceData.HasChange("sqs_urls") {
					props.SqsUrls = &plan.SqsUrls
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, params); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r DataConnectorAwsS3Resource) Delete() sdk.ResourceFunc {
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
