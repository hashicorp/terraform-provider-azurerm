// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"strings"
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

var _ sdk.Resource = AutonomousDatabaseCloneResource{}

type AutonomousDatabaseCloneResource struct{}

type AutonomousDatabaseCloneResourceModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	// Required for Clone
	Source       string `tfschema:"source"`
	SourceId     string `tfschema:"source_id"`
	CloneType    string `tfschema:"clone_type"`
	DataBaseType string `tfschema:"data_base_type"`

	// Required (inherited from base)
	AdminPassword                string  `tfschema:"admin_password"`
	BackupRetentionPeriodInDays  int64   `tfschema:"backup_retention_period_in_days"`
	CharacterSet                 string  `tfschema:"character_set"`
	ComputeCount                 float64 `tfschema:"compute_count"`
	ComputeModel                 string  `tfschema:"compute_model"`
	DataStorageSizeInTbs         int64   `tfschema:"data_storage_size_in_tbs"`
	DbVersion                    string  `tfschema:"db_version"`
	DbWorkload                   string  `tfschema:"db_workload"`
	DisplayName                  string  `tfschema:"display_name"`
	LicenseModel                 string  `tfschema:"license_model"`
	AutoScalingEnabled           bool    `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled bool    `tfschema:"auto_scaling_for_storage_enabled"`
	MtlsConnectionRequired       bool    `tfschema:"mtls_connection_required"`
	NationalCharacterSet         string  `tfschema:"national_character_set"`
	SubnetId                     string  `tfschema:"subnet_id"`
	VnetId                       string  `tfschema:"virtual_network_id"`

	// Optional for Clone
	CustomerContacts               []string `tfschema:"customer_contacts"`
	RefreshableModel               string   `tfschema:"refreshable_model"`
	TimeUntilReconnectCloneEnabled string   `tfschema:"time_until_reconnect_clone_enabled"`

	// optional for clone from backup timestamp
	Timestamp                         string `tfschema:"timestamp"`
	UseLatestAvailableBackupTimeStamp bool   `tfschema:"use_latest_available_backup_time_stamp"`
}

func (AutonomousDatabaseCloneResource) Arguments() map[string]*pluginsdk.Schema {
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
				string(autonomousdatabases.SourceTypeBackupFromId),
				string(autonomousdatabases.SourceTypeDatabase),
				string(autonomousdatabases.SourceTypeCloneToRefreshable),
			}, false),
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				// Source is create-only and not returned by Azure API
				return old != "" && new == ""
			},
		},
		"source_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: autonomousdatabases.ValidateAutonomousDatabaseID,
		},

		"clone_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(autonomousdatabases.CloneTypeFull),
				string(autonomousdatabases.CloneTypeMetadata),
			}, false),
		},
		"data_base_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(autonomousdatabases.DataBaseTypeClone),
				string(autonomousdatabases.DataBaseTypeCloneFromBackupTimestamp),
			}, false),
		},

		//optional

		"refreshable_model": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(autonomousdatabases.RefreshableModelTypeAutomatic),
				string(autonomousdatabases.RefreshableModelTypeManual),
			}, false),
		},

		"time_until_reconnect_clone_enabled": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		//optional for clone from backup time stamp

		"use_latest_available_backup_time_stamp": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"timestamp": {
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

		"data_storage_size_in_tbs": {
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
			ValidateFunc: validate.ValidateCloneWorkloadType,
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

		"tags": commonschema.Tags(),
	}
}

func (AutonomousDatabaseCloneResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (AutonomousDatabaseCloneResource) ModelObject() interface{} {
	return &AutonomousDatabaseCloneResourceModel{}
}

func (AutonomousDatabaseCloneResource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_clone"
}

func (r AutonomousDatabaseCloneResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AutonomousDatabaseCloneResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
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

			// Set properties based on database type
			if model.DataBaseType == string(autonomousdatabases.DataBaseTypeCloneFromBackupTimestamp) {
				param.Properties = &autonomousdatabases.AutonomousDatabaseFromBackupTimestampProperties{
					// Clone-specific properties
					CloneType:    autonomousdatabases.CloneType(model.CloneType),
					SourceId:     model.SourceId,
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
					DataStorageSizeInTbs:           pointer.To(model.DataStorageSizeInTbs),
					DbWorkload:                     pointer.To(autonomousdatabases.WorkloadType(model.DbWorkload)),
					DbVersion:                      pointer.To(model.DbVersion),
					DisplayName:                    pointer.To(model.DisplayName),
					IsAutoScalingEnabled:           pointer.To(model.AutoScalingEnabled),
					IsAutoScalingForStorageEnabled: pointer.To(model.AutoScalingForStorageEnabled),
					IsMtlsConnectionRequired:       pointer.To(model.MtlsConnectionRequired),
					LicenseModel:                   pointer.To(autonomousdatabases.LicenseModel(model.LicenseModel)),
					NcharacterSet:                  pointer.To(model.NationalCharacterSet),
					SubnetId:                       pointer.To(model.SubnetId),
					VnetId:                         pointer.To(model.VnetId),
				}
				cloneBackup := param.Properties.(*autonomousdatabases.AutonomousDatabaseFromBackupTimestampProperties)
				if model.Timestamp != "" {
					cloneBackup.Timestamp = pointer.To(model.Timestamp)
				}
			} else {
				// Regular clone
				param.Properties = &autonomousdatabases.AutonomousDatabaseCloneProperties{
					// Clone-specific properties
					CloneType:    autonomousdatabases.CloneType(model.CloneType),
					SourceId:     model.SourceId,
					Source:       pointer.To(autonomousdatabases.SourceType(model.Source)),
					DataBaseType: autonomousdatabases.DataBaseTypeClone,

					// Optional clone properties
					TimeUntilReconnectCloneEnabled: pointer.To(model.TimeUntilReconnectCloneEnabled),

					// Base properties
					AdminPassword:                  pointer.To(model.AdminPassword),
					BackupRetentionPeriodInDays:    pointer.To(model.BackupRetentionPeriodInDays),
					CharacterSet:                   pointer.To(model.CharacterSet),
					ComputeCount:                   pointer.To(model.ComputeCount),
					ComputeModel:                   pointer.To(autonomousdatabases.ComputeModel(model.ComputeModel)),
					CustomerContacts:               expandCloneCustomerContactsPtr(model.CustomerContacts),
					DataStorageSizeInTbs:           pointer.To(model.DataStorageSizeInTbs),
					DbWorkload:                     pointer.To(autonomousdatabases.WorkloadType(model.DbWorkload)),
					DbVersion:                      pointer.To(model.DbVersion),
					DisplayName:                    pointer.To(model.DisplayName),
					IsAutoScalingEnabled:           pointer.To(model.AutoScalingEnabled),
					IsAutoScalingForStorageEnabled: pointer.To(model.AutoScalingForStorageEnabled),
					IsMtlsConnectionRequired:       pointer.To(model.MtlsConnectionRequired),
					LicenseModel:                   pointer.To(autonomousdatabases.LicenseModel(model.LicenseModel)),
					NcharacterSet:                  pointer.To(model.NationalCharacterSet),
					SubnetId:                       pointer.To(model.SubnetId),
					VnetId:                         pointer.To(model.VnetId),
				}

				// Set optional fields if provided for regular clone
				cloneProps := param.Properties.(*autonomousdatabases.AutonomousDatabaseCloneProperties)
				if model.RefreshableModel != "" {
					cloneProps.RefreshableModel = pointer.To(autonomousdatabases.RefreshableModelType(model.RefreshableModel))
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AutonomousDatabaseCloneResource) Read() sdk.ResourceFunc {
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
			state.Name = id.AutonomousDatabaseName
			state.ResourceGroupName = id.ResourceGroupName

			if val, ok := metadata.ResourceData.GetOk("source"); ok {
				state.Source = val.(string)
			}
			if v, ok := metadata.ResourceData.GetOk("use_latest_available_backup_time_stamp"); ok {
				state.UseLatestAvailableBackupTimeStamp = v.(bool)
			}
			if v, ok := metadata.ResourceData.GetOk("timestamp"); ok {
				state.Timestamp = v.(string)
			}

			if v, ok := metadata.ResourceData.GetOk("refreshable_model"); ok {
				state.RefreshableModel = v.(string)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				// Handle both clone property types
				if cloneProps, ok := model.Properties.(autonomousdatabases.AutonomousDatabaseCloneProperties); ok {
					state.CloneType = string(cloneProps.CloneType)
					state.SourceId = cloneProps.SourceId
					state.DataBaseType = string(cloneProps.DataBaseType)
					state.TimeUntilReconnectCloneEnabled = pointer.From(cloneProps.TimeUntilReconnectCloneEnabled)

					// Base properties
					state.AdminPassword = metadata.ResourceData.Get("admin_password").(string)
					state.BackupRetentionPeriodInDays = pointer.From(cloneProps.BackupRetentionPeriodInDays)
					state.CharacterSet = pointer.From(cloneProps.CharacterSet)
					state.ComputeCount = pointer.From(cloneProps.ComputeCount)
					if cloneProps.ComputeModel != nil {
						state.ComputeModel = string(*cloneProps.ComputeModel)
					}
					state.DataStorageSizeInTbs = pointer.From(cloneProps.DataStorageSizeInTbs)
					state.DbVersion = pointer.From(cloneProps.DbVersion)
					if cloneProps.DbWorkload != nil {
						state.DbWorkload = string(*cloneProps.DbWorkload)
					}
					state.DisplayName = pointer.From(cloneProps.DisplayName)
					state.AutoScalingEnabled = pointer.From(cloneProps.IsAutoScalingEnabled)
					state.AutoScalingForStorageEnabled = pointer.From(cloneProps.IsAutoScalingForStorageEnabled)
					state.MtlsConnectionRequired = pointer.From(cloneProps.IsMtlsConnectionRequired)
					if cloneProps.LicenseModel != nil {
						state.LicenseModel = string(*cloneProps.LicenseModel)
					}
					state.NationalCharacterSet = pointer.From(cloneProps.NcharacterSet)
					state.SubnetId = pointer.From(cloneProps.SubnetId)
					state.VnetId = pointer.From(cloneProps.VnetId)

					if cloneProps.CustomerContacts != nil {
						state.CustomerContacts = flattenCloneCustomerContacts(*cloneProps.CustomerContacts)
					}
				} else if backupProps, ok := model.Properties.(autonomousdatabases.AutonomousDatabaseFromBackupTimestampProperties); ok {
					state.CloneType = string(backupProps.CloneType)
					state.SourceId = backupProps.SourceId
					state.DataBaseType = string(backupProps.DataBaseType)

					// Base properties
					state.AdminPassword = metadata.ResourceData.Get("admin_password").(string)
					state.BackupRetentionPeriodInDays = pointer.From(backupProps.BackupRetentionPeriodInDays)
					state.CharacterSet = pointer.From(backupProps.CharacterSet)
					state.ComputeCount = pointer.From(backupProps.ComputeCount)
					if backupProps.ComputeModel != nil {
						state.ComputeModel = string(*backupProps.ComputeModel)
					}
					state.DataStorageSizeInTbs = pointer.From(backupProps.DataStorageSizeInTbs)
					state.DbVersion = pointer.From(backupProps.DbVersion)
					if backupProps.DbWorkload != nil {
						state.DbWorkload = string(*backupProps.DbWorkload)
					}
					state.DisplayName = pointer.From(backupProps.DisplayName)
					state.AutoScalingEnabled = pointer.From(backupProps.IsAutoScalingEnabled)
					state.AutoScalingForStorageEnabled = pointer.From(backupProps.IsAutoScalingForStorageEnabled)
					state.MtlsConnectionRequired = pointer.From(backupProps.IsMtlsConnectionRequired)
					if backupProps.LicenseModel != nil {
						state.LicenseModel = string(*backupProps.LicenseModel)
					}
					state.NationalCharacterSet = pointer.From(backupProps.NcharacterSet)
					state.SubnetId = pointer.From(backupProps.SubnetId)
					state.VnetId = pointer.From(backupProps.VnetId)

					if backupProps.CustomerContacts != nil {
						state.CustomerContacts = flattenCloneCustomerContacts(*backupProps.CustomerContacts)
					}
				} else {
					return fmt.Errorf("%s was not of expected clone type", id)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AutonomousDatabaseCloneResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
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

func (r AutonomousDatabaseCloneResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// nothing to update
			return nil
		},
	}
}

func (r AutonomousDatabaseCloneResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabases.ValidateAutonomousDatabaseID
}

func expandCloneCustomerContacts(input []string) []autonomousdatabases.CustomerContact {
	if len(input) == 0 {
		return nil
	}

	contacts := make([]autonomousdatabases.CustomerContact, 0, len(input))
	for _, email := range input {
		if strings.TrimSpace(email) != "" {
			contacts = append(contacts, autonomousdatabases.CustomerContact{
				Email: email,
			})
		}
	}
	return contacts
}

func expandCloneCustomerContactsPtr(input []string) *[]autonomousdatabases.CustomerContact {
	if len(input) == 0 {
		return nil
	}

	contacts := expandCloneCustomerContacts(input)
	return &contacts
}

func flattenCloneCustomerContacts(input []autonomousdatabases.CustomerContact) []string {
	if len(input) == 0 {
		return nil
	}

	emails := make([]string, 0, len(input))
	for _, contact := range input {
		if contact.Email != "" {
			emails = append(emails, contact.Email)
		}
	}
	return emails
}

func (AutonomousDatabaseCloneResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Second,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			if metadata.ResourceData == nil {
				return nil
			}

			sourceId := metadata.ResourceData.Get("source_id").(string)
			dbWorkload := metadata.ResourceData.Get("db_workload").(string)

			if sourceId == "" || dbWorkload == "" {
				return nil
			}

			if metadata.ResourceData.Id() != "" {
				return nil
			}

			sourceWorkload, err := getSourceWorkload(ctx, sourceId, metadata)
			if err != nil {
				return nil
			}

			targets, exists := workloadMatrix[sourceWorkload]
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

var workloadMatrix = map[string][]string{
	"DW":   {"OLTP", "DW"},
	"OLTP": {"DW", "OLTP"},
	"AJD":  {"OLTP", "DW", "APEX"},
	"APEX": {"AJD", "OLTP", "DW"},
}

func getSourceWorkload(ctx context.Context, sourceId string, metadata sdk.ResourceMetaData) (string, error) {

	id, err := autonomousdatabases.ParseAutonomousDatabaseID(sourceId)
	if err != nil {
		return "", fmt.Errorf("invalid source_id format: %v", err)
	}

	if metadata.Client == nil || metadata.Client.Oracle == nil || metadata.Client.Oracle.OracleClient == nil {
		return "", fmt.Errorf("oracle client not available")
	}

	client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return "", fmt.Errorf("failed to get source database: %v", err)
	}

	if resp.Model == nil || resp.Model.Properties == nil {
		return "", fmt.Errorf("source database has no properties")
	}

	switch props := resp.Model.Properties.(type) {
	case autonomousdatabases.AutonomousDatabaseProperties:
		if props.DbWorkload != nil {
			return string(*props.DbWorkload), nil
		}
	case autonomousdatabases.AutonomousDatabaseCloneProperties:
		if props.DbWorkload != nil {
			return string(*props.DbWorkload), nil
		}
	case autonomousdatabases.AutonomousDatabaseFromBackupTimestampProperties:
		if props.DbWorkload != nil {
			return string(*props.DbWorkload), nil
		}
	}

	return "", fmt.Errorf("workload type not found in source database properties")
}
