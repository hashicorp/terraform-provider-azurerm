// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/sensitivitylabels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -properties "column_name,schema_name,table_name" -compare-values "subscription_id:database_id,resource_group_name:database_id,server_name:database_id,database_name:database_id" -known-values "sensitivity_label_source:current"

type MssqlDatabaseSensitivityLabelResource struct{}

var (
	_ sdk.ResourceWithIdentity       = MssqlDatabaseSensitivityLabelResource{}
	_ sdk.ResourceWithUpdate         = MssqlDatabaseSensitivityLabelResource{}
	_ sdk.ResourceWithCustomImporter = MssqlDatabaseSensitivityLabelResource{}
)

type MssqlDatabaseSensitivityLabelModel struct {
	ColumnName        string  `tfschema:"column_name"`
	DatabaseId        string  `tfschema:"database_id"`
	SchemaName        string  `tfschema:"schema_name"`
	TableName         string  `tfschema:"table_name"`
	InformationType   *string `tfschema:"information_type"`
	InformationTypeId *string `tfschema:"information_type_id"`
	LabelId           *string `tfschema:"label_id"`
	LabelName         *string `tfschema:"label_name"`
	Rank              *string `tfschema:"rank"`
}

func (MssqlDatabaseSensitivityLabelResource) ResourceType() string {
	return "azurerm_mssql_database_sensitivity_label"
}

func (MssqlDatabaseSensitivityLabelResource) ModelObject() any {
	return &MssqlDatabaseSensitivityLabelModel{}
}

func (MssqlDatabaseSensitivityLabelResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sensitivitylabels.ValidateSensitivityLabelSourceID
}

func (MssqlDatabaseSensitivityLabelResource) Identity() resourceids.ResourceId {
	return &sensitivitylabels.SensitivityLabelSourceId{}
}

func (MssqlDatabaseSensitivityLabelResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"column_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 128),
		},

		"database_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSqlDatabaseID,
		},

		"schema_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 128),
		},

		"table_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 128),
		},

		"information_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 128),
		},

		"information_type_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"label_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"label_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 128),
		},

		"rank": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(sensitivitylabels.SensitivityLabelRankNone),
			ValidateFunc: validation.StringInSlice(sensitivitylabels.PossibleValuesForSensitivityLabelRank(), false),
		},
	}
}

func (MssqlDatabaseSensitivityLabelResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r MssqlDatabaseSensitivityLabelResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.SensitivityLabelsClient

			var config MssqlDatabaseSensitivityLabelModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			databaseId, err := commonids.ParseSqlDatabaseID(config.DatabaseId)
			if err != nil {
				return err
			}
			// `SensitivityLabelSourceRecommended` is recommended by the service instead of a "real" resource can be managed by Terraform,
			// Thus hard-code the source to `SensitivityLabelSourceCurrent` for the resource ID.
			id := sensitivitylabels.NewSensitivityLabelSourceID(databaseId.SubscriptionId, databaseId.ResourceGroupName, databaseId.ServerName, databaseId.DatabaseName, config.SchemaName, config.TableName, config.ColumnName, sensitivitylabels.SensitivityLabelSourceCurrent)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			columnId := sensitivitylabels.NewColumnID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.DatabaseName, id.SchemaName, id.TableName, id.ColumnName)
			payload := sensitivitylabels.SensitivityLabel{
				Properties: &sensitivitylabels.SensitivityLabelProperties{
					ClientClassificationSource: pointer.To(sensitivitylabels.ClientClassificationSourceNative),
					InformationType:            config.InformationType,
					InformationTypeId:          config.InformationTypeId,
					LabelId:                    config.LabelId,
					LabelName:                  config.LabelName,
					Rank:                       pointer.ToEnum[sensitivitylabels.SensitivityLabelRank](pointer.From(config.Rank)),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, columnId, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r MssqlDatabaseSensitivityLabelResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.SensitivityLabelsClient
			id, err := sensitivitylabels.ParseSensitivityLabelSourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (MssqlDatabaseSensitivityLabelResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.SensitivityLabelsClient
			id, err := sensitivitylabels.ParseSensitivityLabelSourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config MssqlDatabaseSensitivityLabelModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model` or `properties` was nil", id)
			}

			properties := existing.Model.Properties
			if metadata.ResourceData.HasChange("information_type") {
				properties.InformationType = config.InformationType
			}
			if metadata.ResourceData.HasChange("information_type_id") {
				properties.InformationTypeId = config.InformationTypeId
			}
			if metadata.ResourceData.HasChange("label_id") {
				properties.LabelId = config.LabelId
			}
			if metadata.ResourceData.HasChange("label_name") {
				properties.LabelName = config.LabelName
			}
			if metadata.ResourceData.HasChange("rank") {
				properties.Rank = pointer.ToEnum[sensitivitylabels.SensitivityLabelRank](pointer.From(config.Rank))
			}

			properties.ClientClassificationSource = pointer.To(sensitivitylabels.ClientClassificationSourceNative)
			properties.ColumnName = nil
			properties.IsDisabled = nil
			properties.SchemaName = nil
			properties.TableName = nil

			columnId := sensitivitylabels.NewColumnID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.DatabaseName, id.SchemaName, id.TableName, id.ColumnName)
			if _, err := client.CreateOrUpdate(ctx, columnId, sensitivitylabels.SensitivityLabel{Properties: properties}); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (MssqlDatabaseSensitivityLabelResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.SensitivityLabelsClient
			id, err := sensitivitylabels.ParseSensitivityLabelSourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			columnId := sensitivitylabels.NewColumnID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.DatabaseName, id.SchemaName, id.TableName, id.ColumnName)
			if resp, err := client.Delete(ctx, columnId); err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (MssqlDatabaseSensitivityLabelResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := sensitivitylabels.ParseSensitivityLabelSourceID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		if id.SensitivityLabelSource != sensitivitylabels.SensitivityLabelSourceCurrent {
			return fmt.Errorf("only the `current` sensitivity label can be imported, got %q", id.SensitivityLabelSource)
		}
		return nil
	}
}

func (MssqlDatabaseSensitivityLabelResource) flatten(metadata sdk.ResourceMetaData, id *sensitivitylabels.SensitivityLabelSourceId, model *sensitivitylabels.SensitivityLabel) error {
	state := MssqlDatabaseSensitivityLabelModel{
		ColumnName: id.ColumnName,
		DatabaseId: commonids.NewSqlDatabaseID(id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.DatabaseName).ID(),
		SchemaName: id.SchemaName,
		TableName:  id.TableName,
		Rank:       pointer.To(string(sensitivitylabels.SensitivityLabelRankNone)),
	}

	if properties := model.Properties; properties != nil {
		state.InformationType = properties.InformationType
		state.InformationTypeId = properties.InformationTypeId
		state.LabelId = properties.LabelId
		state.LabelName = properties.LabelName
		if properties.Rank != nil {
			state.Rank = pointer.To(pointer.FromEnum(properties.Rank))
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}
