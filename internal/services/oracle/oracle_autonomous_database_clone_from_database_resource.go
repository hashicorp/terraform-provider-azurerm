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

var _ sdk.ResourceWithCustomizeDiff = AutonomousDatabaseCloneFromDatabaseResource{}

type AutonomousDatabaseCloneFromDatabaseResource struct{}

type AutonomousDatabaseCloneResourceModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	// Required for Clone

	SourceAutonomousDatabaseId string `tfschema:"source_autonomous_database_id"`
	CloneType                  string `tfschema:"clone_type"`

	// Required (inherited from base)

	AdminPassword                string   `tfschema:"admin_password"`
	BackupRetentionPeriodInDays  int64    `tfschema:"backup_retention_period_in_days"`
	CharacterSet                 string   `tfschema:"character_set"`
	ComputeCount                 float64  `tfschema:"compute_count"`
	ComputeModel                 string   `tfschema:"compute_model"`
	DataStorageSizeInTb          int64    `tfschema:"data_storage_size_in_tb"`
	DatabaseVersion              string   `tfschema:"database_version"`
	DatabaseWorkload             string   `tfschema:"database_workload"`
	DisplayName                  string   `tfschema:"display_name"`
	LicenseModel                 string   `tfschema:"license_model"`
	AutoScalingEnabled           bool     `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled bool     `tfschema:"auto_scaling_for_storage_enabled"`
	MtlsConnectionRequired       bool     `tfschema:"mtls_connection_required"`
	NationalCharacterSet         string   `tfschema:"national_character_set"`
	SubnetId                     string   `tfschema:"subnet_id"`
	VnetId                       string   `tfschema:"virtual_network_id"`
	AllowedIps                   []string `tfschema:"allowed_ips"`

	// Optional for Clone

	CustomerContacts   []string `tfschema:"customer_contacts"`
	RefreshableModel   string   `tfschema:"refreshable_model"`
	TimeUntilReconnect string   `tfschema:"time_until_reconnect"`
}

func (AutonomousDatabaseCloneFromDatabaseResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
			ForceNew:     true,
		},
		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		// Clone-specific required fields

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

		"refreshable_model": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(autonomousdatabases.PossibleValuesForRefreshableModelType(), false),
		},

		"time_until_reconnect": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		// Required (inherited from base)

		"admin_password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ForceNew:     true,
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
			ForceNew:     true,
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
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 384),
		},

		"database_version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"database_workload": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(autonomousdatabases.PossibleValuesForWorkloadType(), false),
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 255),
		},

		"auto_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
			ForceNew: true,
		},

		"auto_scaling_for_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
			ForceNew: true,
		},

		"mtls_connection_required": {
			Type:     pluginsdk.TypeBool,
			Required: true,
			ForceNew: true,
		},

		"license_model": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(autonomousdatabases.PossibleValuesForLicenseModel(), false),
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
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.Any(
					validation.IsIPv4Address,
					validation.IsCIDR,
				),
			},
		},

		"tags": commonschema.TagsForceNew(),
	}
}

func (AutonomousDatabaseCloneFromDatabaseResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (AutonomousDatabaseCloneFromDatabaseResource) ModelObject() interface{} {
	return &AutonomousDatabaseCloneResourceModel{}
}

func (AutonomousDatabaseCloneFromDatabaseResource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_clone_from_database"
}

func (r AutonomousDatabaseCloneFromDatabaseResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AutonomousDatabaseCloneResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := autonomousdatabases.NewAutonomousDatabaseID(
				subscriptionId,
				model.ResourceGroupName,
				model.Name)

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

			param.Properties = &autonomousdatabases.AutonomousDatabaseCloneProperties{
				CloneType:                      autonomousdatabases.CloneType(model.CloneType),
				SourceId:                       model.SourceAutonomousDatabaseId,
				Source:                         pointer.To(autonomousdatabases.SourceTypeDatabase),
				DataBaseType:                   autonomousdatabases.DataBaseTypeClone,
				TimeUntilReconnectCloneEnabled: pointer.To(model.TimeUntilReconnect),

				// Base properties
				AdminPassword:                  pointer.To(model.AdminPassword),
				BackupRetentionPeriodInDays:    pointer.To(model.BackupRetentionPeriodInDays),
				CharacterSet:                   pointer.To(model.CharacterSet),
				ComputeCount:                   pointer.To(model.ComputeCount),
				ComputeModel:                   pointer.To(autonomousdatabases.ComputeModel(model.ComputeModel)),
				CustomerContacts:               pointer.To(expandCloneCustomerContacts(model.CustomerContacts)),
				DataStorageSizeInTbs:           pointer.To(model.DataStorageSizeInTb),
				DbWorkload:                     pointer.To(autonomousdatabases.WorkloadType(model.DatabaseWorkload)),
				DbVersion:                      pointer.To(model.DatabaseVersion),
				DisplayName:                    pointer.To(model.DisplayName),
				IsAutoScalingEnabled:           pointer.To(model.AutoScalingEnabled),
				IsAutoScalingForStorageEnabled: pointer.To(model.AutoScalingForStorageEnabled),
				IsMtlsConnectionRequired:       pointer.To(model.MtlsConnectionRequired),
				LicenseModel:                   pointer.To(autonomousdatabases.LicenseModel(model.LicenseModel)),
				NcharacterSet:                  pointer.To(model.NationalCharacterSet),
				WhitelistedIPs:                 pointer.To(model.AllowedIps),
			}

			properties := param.Properties.(*autonomousdatabases.AutonomousDatabaseCloneProperties)
			if model.RefreshableModel != "" {
				properties.RefreshableModel = pointer.To(autonomousdatabases.RefreshableModelType(model.RefreshableModel))
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

func (r AutonomousDatabaseCloneFromDatabaseResource) Read() sdk.ResourceFunc {
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

			var state AutonomousDatabaseCloneResourceModel

			if v, ok := metadata.ResourceData.GetOk("refreshable_model"); ok {
				state.RefreshableModel = v.(string)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
				state.Name = id.AutonomousDatabaseName
				state.ResourceGroupName = id.ResourceGroupName

				props, ok := model.Properties.(autonomousdatabases.AutonomousDatabaseCloneProperties)
				if !ok {
					return fmt.Errorf("%s was not of type `Clone`", id)
				}
				state.CloneType = string(props.CloneType)
				state.SourceAutonomousDatabaseId = props.SourceId
				state.TimeUntilReconnect = pointer.From(props.TimeUntilReconnectCloneEnabled)

				// Base properties
				state.AdminPassword = metadata.ResourceData.Get("admin_password").(string)
				state.BackupRetentionPeriodInDays = pointer.From(props.BackupRetentionPeriodInDays)
				state.CharacterSet = pointer.From(props.CharacterSet)
				state.ComputeCount = pointer.From(props.ComputeCount)
				state.ComputeModel = pointer.FromEnum(props.ComputeModel)
				state.DataStorageSizeInTb = pointer.From(props.DataStorageSizeInTbs)
				state.DatabaseVersion = pointer.From(props.DbVersion)
				state.DatabaseVersion = pointer.From(props.DbVersion)
				state.DatabaseWorkload = string(pointer.From(props.DbWorkload))
				state.DisplayName = pointer.From(props.DisplayName)
				state.AutoScalingEnabled = pointer.From(props.IsAutoScalingEnabled)
				state.AutoScalingForStorageEnabled = pointer.From(props.IsAutoScalingForStorageEnabled)
				state.MtlsConnectionRequired = pointer.From(props.IsMtlsConnectionRequired)
				state.LicenseModel = string(pointer.From(props.LicenseModel))
				state.NationalCharacterSet = pointer.From(props.NcharacterSet)
				state.SubnetId = pointer.From(props.SubnetId)
				state.VnetId = pointer.From(props.VnetId)
				state.AllowedIps = pointer.From(props.WhitelistedIPs)
				state.CustomerContacts = flattenAdbsCustomerContacts(props.CustomerContacts)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AutonomousDatabaseCloneFromDatabaseResource) Delete() sdk.ResourceFunc {
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

func (r AutonomousDatabaseCloneFromDatabaseResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabases.ValidateAutonomousDatabaseID
}

func (AutonomousDatabaseCloneFromDatabaseResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Second,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			sourceId := metadata.ResourceDiff.Get("source_autonomous_database_id").(string)
			dbWorkload := metadata.ResourceDiff.Get("database_workload").(string)

			if sourceId == "" || dbWorkload == "" {
				return nil
			}

			if metadata.ResourceDiff.Id() != "" {
				return nil
			}

			sourceWorkload, err := getSourceWorkloadforClone(ctx, sourceId, metadata)
			if err != nil {
				return err
			}

			targets, exists := workloadMatrixForClone[sourceWorkload]
			if !exists {
				return fmt.Errorf("unsupported source workload: %s", sourceWorkload)
			}

			for _, target := range targets {
				if dbWorkload == target {
					return nil
				}
			}

			return fmt.Errorf("invalid workload: %s->%s not allowed", sourceWorkload, dbWorkload)
		},
	}
}

var workloadMatrixForClone = map[string][]string{
	"DW":   {"OLTP", "DW"},
	"OLTP": {"DW", "OLTP"},
	"AJD":  {"OLTP", "DW", "APEX"},
	"APEX": {"AJD", "OLTP", "DW"},
}

func getSourceWorkloadforClone(ctx context.Context, sourceId string, metadata sdk.ResourceMetaData) (string, error) {
	id, err := autonomousdatabases.ParseAutonomousDatabaseID(sourceId)
	if err != nil {
		return "", err
	}

	client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return "", fmt.Errorf("retrieving %s", err)
	}

	if resp.Model == nil {
		return "", fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if resp.Model.Properties == nil {
		return "", fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
	}

	props := resp.Model.Properties.AutonomousDatabaseBaseProperties()
	if props.DbWorkload == nil {
		return "", fmt.Errorf("unable to determine workload type for %s", id)
	}

	return string(*props.DbWorkload), nil
}
