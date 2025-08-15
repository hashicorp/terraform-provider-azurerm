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
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = AutonomousDatabaseCloneFromBackupResource{}

type AutonomousDatabaseCloneFromBackupResource struct{}

type AutonomousDatabaseCloneFromBackupResourceModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	// Required for Clone

	Source                     string `tfschema:"source"`
	SourceAutonomousDatabaseId string `tfschema:"source_autonomous_database_id"`
	CloneType                  string `tfschema:"clone_type"`

	// Required (inherited from base)

	AdminPassword                string   `tfschema:"admin_password"`
	AllowedIps                   []string `tfschema:"allowed_ips"`
	BackupRetentionPeriodInDays  int64    `tfschema:"backup_retention_period_in_days"`
	CharacterSet                 string   `tfschema:"character_set"`
	ComputeCount                 float64  `tfschema:"compute_count"`
	ComputeModel                 string   `tfschema:"compute_model"`
	DataStorageSizeInTb          int64    `tfschema:"data_storage_size_in_tb"`
	DbVersion                    string   `tfschema:"db_version"`
	DbWorkload                   string   `tfschema:"db_workload"`
	DisplayName                  string   `tfschema:"display_name"`
	LicenseModel                 string   `tfschema:"license_model"`
	AutoScalingEnabled           bool     `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled bool     `tfschema:"auto_scaling_for_storage_enabled"`
	MtlsConnectionRequired       bool     `tfschema:"mtls_connection_required"`
	NationalCharacterSet         string   `tfschema:"national_character_set"`
	SubnetId                     string   `tfschema:"subnet_id"`
	VnetId                       string   `tfschema:"virtual_network_id"`

	// Optional

	BackupTimestamp                   string   `tfschema:"backup_timestamp"`
	CustomerContacts                  []string `tfschema:"customer_contacts"`
	UseLatestAvailableBackupTimeStamp bool     `tfschema:"use_latest_available_backup_time_stamp"`
}

func (AutonomousDatabaseCloneFromBackupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
			ForceNew:     true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// Clone-specific required fields

		"source": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(autonomousdatabases.SourceTypeBackupFromTimestamp),
			}, false),
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				// Source is create-only and not returned by Azure API
				return old != "" && new == ""
			},
		},
		"source_autonomous_database_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: autonomousdatabases.ValidateAutonomousDatabaseID,
		},

		"clone_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(autonomousdatabases.PossibleValuesForCloneType(), false),
		},
		// optional

		"use_latest_available_backup_time_stamp": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"backup_timestamp": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		// Required (inherited from base)

		"admin_password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validate.AutonomousDatabasePassword,
		},

		"backup_retention_period_in_days": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 60),
		},

		"character_set": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"compute_count": {
			Type:         pluginsdk.TypeFloat,
			Required:     true,
			ValidateFunc: validation.FloatBetween(2.0, 512.0),
		},

		"compute_model": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AdbsComputeModel,
		},

		"data_storage_size_in_tb": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 384),
		},

		"db_version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"db_workload": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(autonomousdatabases.PossibleValuesForWorkloadType(), false),
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
		},

		"auto_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"auto_scaling_for_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"mtls_connection_required": {
			Type:     pluginsdk.TypeBool,
			Required: true,
			ForceNew: true,
		},

		"license_model": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(autonomousdatabases.LicenseModelLicenseIncluded),
				string(autonomousdatabases.LicenseModelBringYourOwnLicense),
			}, false),
		},

		"national_character_set": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// Optional clone-specific fields
		"customer_contacts": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.CustomerContactEmail,
			},
		},
		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},
		"allowed_ips": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MaxItems: 1024,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsIPv4Address,
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (AutonomousDatabaseCloneFromBackupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (AutonomousDatabaseCloneFromBackupResource) ModelObject() interface{} {
	return &AutonomousDatabaseCloneFromBackupResourceModel{}
}

func (AutonomousDatabaseCloneFromBackupResource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_clone_from_backup"
}

func (r AutonomousDatabaseCloneFromBackupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AutonomousDatabaseCloneFromBackupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := autonomousdatabases.NewAutonomousDatabaseID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := autonomousdatabases.AutonomousDatabase{
				Name:     pointer.To(model.Name),
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
			}

			param.Properties = &autonomousdatabases.AutonomousDatabaseFromBackupTimestampProperties{
				CloneType:    autonomousdatabases.CloneType(model.CloneType),
				SourceId:     model.SourceAutonomousDatabaseId,
				Source:       autonomousdatabases.Source(autonomousdatabases.SourceTypeBackupFromTimestamp),
				DataBaseType: autonomousdatabases.DataBaseTypeCloneFromBackupTimestamp,

				UseLatestAvailableBackupTimeStamp: pointer.To(model.UseLatestAvailableBackupTimeStamp),

				// Base properties
				AdminPassword:                  pointer.To(model.AdminPassword),
				BackupRetentionPeriodInDays:    pointer.To(model.BackupRetentionPeriodInDays),
				CharacterSet:                   pointer.To(model.CharacterSet),
				ComputeCount:                   pointer.To(model.ComputeCount),
				ComputeModel:                   pointer.To(autonomousdatabases.ComputeModel(model.ComputeModel)),
				CustomerContacts:               expandCloneCustomerContactsPtr(model.CustomerContacts),
				DataStorageSizeInTbs:           pointer.To(model.DataStorageSizeInTb),
				DbWorkload:                     pointer.To(autonomousdatabases.WorkloadType(model.DbWorkload)),
				DbVersion:                      pointer.To(model.DbVersion),
				DisplayName:                    pointer.To(model.DisplayName),
				IsAutoScalingEnabled:           pointer.To(model.AutoScalingEnabled),
				IsAutoScalingForStorageEnabled: pointer.To(model.AutoScalingForStorageEnabled),
				IsMtlsConnectionRequired:       pointer.To(model.MtlsConnectionRequired),
				LicenseModel:                   pointer.To(autonomousdatabases.LicenseModel(model.LicenseModel)),
				NcharacterSet:                  pointer.To(model.NationalCharacterSet),
				WhitelistedIPs:                 pointer.To(model.AllowedIps),
			}
			properties := param.Properties.(*autonomousdatabases.AutonomousDatabaseFromBackupTimestampProperties)
			if model.BackupTimestamp != "" {
				properties.Timestamp = pointer.To(model.BackupTimestamp)
			}
			if model.SubnetId != "" {
				properties.SubnetId = pointer.To(model.SubnetId)
			}

			if model.VnetId != "" {
				properties.VnetId = pointer.To(model.VnetId)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AutonomousDatabaseCloneFromBackupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var state AutonomousDatabaseCloneFromBackupResourceModel

			if val, ok := metadata.ResourceData.GetOk("source"); ok {
				state.Source = val.(string)
			}
			if v, ok := metadata.ResourceData.GetOk("use_latest_available_backup_time_stamp"); ok {
				state.UseLatestAvailableBackupTimeStamp = v.(bool)
			}
			if v, ok := metadata.ResourceData.GetOk("timestamp"); ok {
				state.BackupTimestamp = v.(string)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.Name = id.AutonomousDatabaseName
				state.ResourceGroupName = id.ResourceGroupName

				props, ok := model.Properties.(autonomousdatabases.AutonomousDatabaseFromBackupTimestampProperties)
				if !ok {
					return fmt.Errorf("%s was not of type `Regular`", id)
				}
				state.CloneType = string(props.CloneType)
				state.SourceAutonomousDatabaseId = props.SourceId

				state.AdminPassword = metadata.ResourceData.Get("admin_password").(string)
				state.BackupRetentionPeriodInDays = pointer.From(props.BackupRetentionPeriodInDays)
				state.CharacterSet = pointer.From(props.CharacterSet)
				state.ComputeCount = pointer.From(props.ComputeCount)
				state.ComputeModel = pointer.FromEnum(props.ComputeModel)
				state.CustomerContacts = flattenCloneCustomerContacts(pointer.From(props.CustomerContacts))
				state.DataStorageSizeInTb = pointer.From(props.DataStorageSizeInTbs)
				state.DbVersion = pointer.From(props.DbVersion)
				state.DbWorkload = string(pointer.From(props.DbWorkload))
				state.DisplayName = pointer.From(props.DisplayName)
				state.AutoScalingEnabled = pointer.From(props.IsAutoScalingEnabled)
				state.AutoScalingForStorageEnabled = pointer.From(props.IsAutoScalingForStorageEnabled)
				state.MtlsConnectionRequired = pointer.From(props.IsMtlsConnectionRequired)
				state.LicenseModel = string(pointer.From(props.LicenseModel))
				state.NationalCharacterSet = pointer.From(props.NcharacterSet)
				state.SubnetId = pointer.From(props.SubnetId)
				state.VnetId = pointer.From(props.VnetId)
				state.AllowedIps = pointer.From(props.WhitelistedIPs)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AutonomousDatabaseCloneFromBackupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases

			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r AutonomousDatabaseCloneFromBackupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return nil
		},
	}
}

func (r AutonomousDatabaseCloneFromBackupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabases.ValidateAutonomousDatabaseID
}
