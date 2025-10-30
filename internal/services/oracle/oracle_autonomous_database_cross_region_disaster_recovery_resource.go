// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = AutonomousDatabaseCrossRegionDisasterRecoveryResource{}

type AutonomousDatabaseCrossRegionDisasterRecoveryResource struct{}

type AutonomousDatabaseCrossRegionDisasterRecoveryResourceModel struct {
	Location                     string            `tfschema:"location"`
	Name                         string            `tfschema:"name"`
	ResourceGroupName            string            `tfschema:"resource_group_name"`
	Tags                         map[string]string `tfschema:"tags"`
	RemoteDisasterRecoveryType   string            `tfschema:"remote_disaster_recovery_type"`
	SourceAutonomousDatabaseId   string            `tfschema:"source_autonomous_database_id"`
	AutoScalingEnabled           bool              `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled bool              `tfschema:"auto_scaling_for_storage_enabled"`
	BackupRetentionPeriodInDays  int64             `tfschema:"backup_retention_period_in_days"`
	CharacterSet                 string            `tfschema:"character_set"`
	ComputeCount                 float64           `tfschema:"compute_count"`
	ComputeModel                 string            `tfschema:"compute_model"`
	DataStorageSizeInTb          int64             `tfschema:"data_storage_size_in_tb"`
	DatabaseVersion              string            `tfschema:"database_version"`
	DatabaseWorkload             string            `tfschema:"database_workload"`
	DisplayName                  string            `tfschema:"display_name"`
	LicenseModel                 string            `tfschema:"license_model"`
	MtlsConnectionRequired       bool              `tfschema:"mtls_connection_required"`
	NationalCharacterSet         string            `tfschema:"national_character_set"`
	SubnetId                     string            `tfschema:"subnet_id"`
	VnetId                       string            `tfschema:"virtual_network_id"`

	// Optional
	CustomerContacts                 []string `tfschema:"customer_contacts"`
	ReplicateAutomaticBackupsEnabled bool     `tfschema:"replicate_automatic_backups_enabled"`
}

func (AutonomousDatabaseCrossRegionDisasterRecoveryResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
			ForceNew:     true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},

		// Cross Region Disaster Recovery
		// Required

		"source_autonomous_database_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: autonomousdatabases.ValidateAutonomousDatabaseID,
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		// Optional
		"replicate_automatic_backups_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"tags": commonschema.TagsForceNew(),
	}
}

func (AutonomousDatabaseCrossRegionDisasterRecoveryResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"remote_disaster_recovery_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"auto_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"auto_scaling_for_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"backup_retention_period_in_days": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
		"character_set": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"compute_count": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},
		"compute_model": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"data_storage_size_in_tb": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
		"database_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"database_workload": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"license_model": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"mtls_connection_required": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"national_character_set": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"customer_contacts": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (AutonomousDatabaseCrossRegionDisasterRecoveryResource) ModelObject() interface{} {
	return &AutonomousDatabaseCrossRegionDisasterRecoveryResourceModel{}
}

func (AutonomousDatabaseCrossRegionDisasterRecoveryResource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_cross_region_disaster_recovery"
}

func (r AutonomousDatabaseCrossRegionDisasterRecoveryResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AutonomousDatabaseCrossRegionDisasterRecoveryResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			sourceId, err := autonomousdatabases.ParseAutonomousDatabaseID(model.SourceAutonomousDatabaseId)
			if err != nil {
				return err
			}
			sourceDb, err := client.Get(ctx, *sourceId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", sourceId, err)
			}
			if sourceDb.Model == nil {
				return fmt.Errorf("retrieving %s: `Model` was nil", sourceId)
			}
			sourceLocation := sourceDb.Model.Location
			if location.Normalize(model.Location) == location.Normalize(sourceLocation) {
				return fmt.Errorf("disaster Recovery database must reside in a different region from the source database (source is '%s', target is '%s')", sourceLocation, model.Location)
			}

			if sourceDb.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `Properties` was nil", sourceId)
			}
			sourceProps := sourceDb.Model.Properties.AutonomousDatabaseBaseProperties()

			param := autonomousdatabases.AutonomousDatabase{
				Name:     pointer.To(model.Name),
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
				Properties: &autonomousdatabases.AutonomousDatabaseCrossRegionDisasterRecoveryProperties{
					Source:                         autonomousdatabases.SourceCrossRegionDisasterRecovery,
					SourceId:                       model.SourceAutonomousDatabaseId,
					SourceOcid:                     sourceProps.Ocid,
					SourceLocation:                 pointer.To(sourceLocation),
					RemoteDisasterRecoveryType:     autonomousdatabases.DisasterRecoveryTypeAdg,
					IsReplicateAutomaticBackups:    pointer.To(model.ReplicateAutomaticBackupsEnabled),
					AdminPassword:                  sourceProps.AdminPassword,
					BackupRetentionPeriodInDays:    sourceProps.BackupRetentionPeriodInDays,
					CharacterSet:                   sourceProps.CharacterSet,
					ComputeCount:                   sourceProps.ComputeCount,
					ComputeModel:                   sourceProps.ComputeModel,
					CustomerContacts:               sourceProps.CustomerContacts,
					DataBaseType:                   autonomousdatabases.DataBaseTypeCrossRegionDisasterRecovery,
					DataStorageSizeInTbs:           sourceProps.DataStorageSizeInTbs,
					DbWorkload:                     sourceProps.DbWorkload,
					DbVersion:                      sourceProps.DbVersion,
					DisplayName:                    pointer.To(model.DisplayName),
					IsAutoScalingEnabled:           sourceProps.IsAutoScalingEnabled,
					IsAutoScalingForStorageEnabled: sourceProps.IsAutoScalingForStorageEnabled,
					IsMtlsConnectionRequired:       sourceProps.IsMtlsConnectionRequired,
					LicenseModel:                   sourceProps.LicenseModel,
					NcharacterSet:                  sourceProps.NcharacterSet,
					SubnetId:                       pointer.To(model.SubnetId),
					VnetId:                         pointer.To(model.VnetId),
				},
			}

			id := autonomousdatabases.NewAutonomousDatabaseID(subscriptionId,
				model.ResourceGroupName,
				model.Name)
			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (AutonomousDatabaseCrossRegionDisasterRecoveryResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := AutonomousDatabaseCrossRegionDisasterRecoveryResourceModel{
				Name:              id.AutonomousDatabaseName,
				ResourceGroupName: id.ResourceGroupName,
			}
			if model := result.Model; model != nil {
				props, ok := model.Properties.(autonomousdatabases.AutonomousDatabaseCrossRegionDisasterRecoveryProperties)
				if !ok {
					return fmt.Errorf("%s was not of type `CrossRegionDisasterRecovery`", id)
				}

				state.ReplicateAutomaticBackupsEnabled = pointer.From(props.IsReplicateAutomaticBackups)
				state.RemoteDisasterRecoveryType = string(props.RemoteDisasterRecoveryType)
				state.SourceAutonomousDatabaseId = props.SourceId
				state.AutoScalingEnabled = pointer.From(props.IsAutoScalingEnabled)
				state.BackupRetentionPeriodInDays = pointer.From(props.BackupRetentionPeriodInDays)
				state.AutoScalingForStorageEnabled = pointer.From(props.IsAutoScalingForStorageEnabled)
				state.CharacterSet = pointer.From(props.CharacterSet)
				state.ComputeCount = pointer.From(props.ComputeCount)
				state.ComputeModel = pointer.FromEnum(props.ComputeModel)
				state.CustomerContacts = flattenAdbsCustomerContacts(props.CustomerContacts)
				state.DataStorageSizeInTb = pointer.From(props.DataStorageSizeInTbs)
				state.DatabaseWorkload = pointer.FromEnum(props.DbWorkload)
				state.DatabaseVersion = pointer.From(props.DbVersion)
				state.DisplayName = pointer.From(props.DisplayName)
				state.LicenseModel = pointer.FromEnum(props.LicenseModel)
				state.Location = model.Location
				state.NationalCharacterSet = pointer.From(props.NcharacterSet)
				state.SubnetId = pointer.From(props.SubnetId)
				state.Tags = pointer.From(model.Tags)
				state.VnetId = pointer.From(props.VnetId)
			}
			return metadata.Encode(&state)
		},
	}
}

func (AutonomousDatabaseCrossRegionDisasterRecoveryResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases

			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (AutonomousDatabaseCrossRegionDisasterRecoveryResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabases.ValidateAutonomousDatabaseID
}
