// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

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

var _ sdk.Resource = AutonomousDatabaseRegularResource{}

type AutonomousDatabaseRegularResource struct{}

type AutonomousDatabaseRegularResourceModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	// Required
	AdminPassword                string                          `tfschema:"admin_password"`
	BackupRetentionPeriodInDays  int64                           `tfschema:"backup_retention_period_in_days"`
	CharacterSet                 string                          `tfschema:"character_set"`
	ComputeCount                 float64                         `tfschema:"compute_count"`
	ComputeModel                 string                          `tfschema:"compute_model"`
	DataStorageSizeInTbs         int64                           `tfschema:"data_storage_size_in_tbs"`
	DbVersion                    string                          `tfschema:"db_version"`
	DbWorkload                   string                          `tfschema:"db_workload"`
	DisplayName                  string                          `tfschema:"display_name"`
	LicenseModel                 string                          `tfschema:"license_model"`
	LongTermBackUpSchedule       []LongTermBackUpScheduleDetails `tfschema:"long_term_backup_schedule"`
	AutoScalingEnabled           bool                            `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled bool                            `tfschema:"auto_scaling_for_storage_enabled"`
	MtlsConnectionRequired       bool                            `tfschema:"mtls_connection_required"`
	NationalCharacterSet         string                          `tfschema:"national_character_set"`
	SubnetId                     string                          `tfschema:"subnet_id"`
	VnetId                       string                          `tfschema:"virtual_network_id"`
	AllowedIps                   []string                        `tfschema:"allowed_ips"`
	Ocid                         string                          `tfschema:"ocid"`

	// Optional
	CustomerContacts []string `tfschema:"customer_contacts"`
}

func (AutonomousDatabaseRegularResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AutonomousDatabaseName,
			ForceNew:     true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// Required
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(autonomousdatabases.WorkloadTypeDW),
				string(autonomousdatabases.WorkloadTypeOLTP),
			}, false),
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

		"license_model": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(autonomousdatabases.LicenseModelLicenseIncluded),
				string(autonomousdatabases.LicenseModelBringYourOwnLicense),
			}, false),
		},

		"long_term_backup_schedule": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"repeat_cadence": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(autonomousdatabases.PossibleValuesForRepeatCadenceType(), false),
					},
					"time_of_backup": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
					"retention_period_in_days": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(90, 2558),
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},

		"national_character_set": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// Optional
		"customer_contacts": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.CustomerContactEmail,
			},
		},

		"mtls_connection_required": {
			Type:     pluginsdk.TypeBool,
			Required: true,
			ForceNew: true,
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
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsIPv4Address,
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (AutonomousDatabaseRegularResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (AutonomousDatabaseRegularResource) ModelObject() interface{} {
	return &AutonomousDatabaseRegularResource{}
}

func (AutonomousDatabaseRegularResource) ResourceType() string {
	return "azurerm_oracle_autonomous_database"
}

func (r AutonomousDatabaseRegularResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AutonomousDatabaseRegularResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := autonomousdatabases.NewAutonomousDatabaseID(subscriptionId,
				model.ResourceGroupName,
				model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}
			properties := &autonomousdatabases.AutonomousDatabaseProperties{
				AdminPassword:                  pointer.To(model.AdminPassword),
				BackupRetentionPeriodInDays:    pointer.To(model.BackupRetentionPeriodInDays),
				CharacterSet:                   pointer.To(model.CharacterSet),
				ComputeCount:                   pointer.To(model.ComputeCount),
				ComputeModel:                   pointer.To(autonomousdatabases.ComputeModel(model.ComputeModel)),
				DataBaseType:                   "Regular",
				DataStorageSizeInTbs:           pointer.To(model.DataStorageSizeInTbs),
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

			if len(model.CustomerContacts) > 0 {
				properties.CustomerContacts = pointer.To(expandAdbsCustomerContacts(model.CustomerContacts))
			}

			if model.SubnetId != "" {
				properties.SubnetId = pointer.To(model.SubnetId)
			}

			if model.VnetId != "" {
				properties.VnetId = pointer.To(model.VnetId)
			}

			param := autonomousdatabases.AutonomousDatabase{
				Name:       pointer.To(model.Name),
				Location:   location.Normalize(model.Location),
				Tags:       pointer.To(model.Tags),
				Properties: properties,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if len(model.LongTermBackUpSchedule) > 0 {
				backupUpdate := autonomousdatabases.AutonomousDatabaseUpdate{
					Properties: &autonomousdatabases.AutonomousDatabaseUpdateProperties{
						LongTermBackupSchedule: expandLongTermBackupSchedule(model.LongTermBackUpSchedule),
					},
				}
				if err := client.UpdateThenPoll(ctx, id, backupUpdate); err != nil {
					return fmt.Errorf("configuring backup schedule for %s: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AutonomousDatabaseRegularResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			var model AutonomousDatabaseRegularResourceModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			_, err = client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			// Check what needs to be updated
			needsGeneralUpdate := r.hasGeneralUpdates(metadata)
			needsBackupScheduleUpdate := metadata.ResourceData.HasChange("long_term_backup_schedule")

			// Step 1: Handle general updates (everything except backup schedule)
			if needsGeneralUpdate {
				generalUpdate := autonomousdatabases.AutonomousDatabaseUpdate{
					Properties: &autonomousdatabases.AutonomousDatabaseUpdateProperties{},
				}

				if metadata.ResourceData.HasChange("tags") {
					generalUpdate.Tags = pointer.To(model.Tags)
				}
				if metadata.ResourceData.HasChange("backup_retention_period_in_days") {
					generalUpdate.Properties.BackupRetentionPeriodInDays = pointer.To(model.BackupRetentionPeriodInDays)
				}
				if metadata.ResourceData.HasChange("data_storage_size_in_tbs") {
					generalUpdate.Properties.DataStorageSizeInTbs = pointer.To(model.DataStorageSizeInTbs)
				}
				if metadata.ResourceData.HasChange("compute_count") {
					generalUpdate.Properties.ComputeCount = pointer.To(model.ComputeCount)
				}
				if metadata.ResourceData.HasChange("auto_scaling_enabled") {
					generalUpdate.Properties.IsAutoScalingEnabled = pointer.To(model.AutoScalingEnabled)
				}
				if metadata.ResourceData.HasChange("auto_scaling_for_storage_enabled") {
					generalUpdate.Properties.IsAutoScalingForStorageEnabled = pointer.To(model.AutoScalingForStorageEnabled)
				}
				if metadata.ResourceData.HasChange("allowed_ips") {
					generalUpdate.Properties.WhitelistedIPs = pointer.To(model.AllowedIps)
				}

				if err := client.UpdateThenPoll(ctx, *id, generalUpdate); err != nil {
					return fmt.Errorf("updating general properties for %s: %+v", *id, err)
				}
			}

			// Step 2: Handle backup schedule update separately
			if needsBackupScheduleUpdate {
				backupUpdate := autonomousdatabases.AutonomousDatabaseUpdate{
					Properties: &autonomousdatabases.AutonomousDatabaseUpdateProperties{
						LongTermBackupSchedule: expandLongTermBackupSchedule(model.LongTermBackUpSchedule),
					},
				}

				if err := client.UpdateThenPoll(ctx, *id, backupUpdate); err != nil {
					return fmt.Errorf("updating backup schedule for %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (AutonomousDatabaseRegularResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases
			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := AutonomousDatabaseRegularResourceModel{
				Name:              id.AutonomousDatabaseName,
				ResourceGroupName: id.ResourceGroupName,
			}
			if model := result.Model; model != nil {
				props, ok := model.Properties.(autonomousdatabases.AutonomousDatabaseProperties)
				if !ok {
					return fmt.Errorf("%s was not of type `Regular`", id)
				}
				state.AdminPassword = metadata.ResourceData.Get("admin_password").(string)
				state.AutoScalingEnabled = pointer.From(props.IsAutoScalingEnabled)
				state.BackupRetentionPeriodInDays = pointer.From(props.BackupRetentionPeriodInDays)
				state.AutoScalingForStorageEnabled = pointer.From(props.IsAutoScalingForStorageEnabled)
				state.CharacterSet = pointer.From(props.CharacterSet)
				state.ComputeCount = pointer.From(props.ComputeCount)
				state.ComputeModel = pointer.FromEnum(props.ComputeModel)
				state.CustomerContacts = flattenAdbsCustomerContacts(props.CustomerContacts)
				state.DataStorageSizeInTbs = pointer.From(props.DataStorageSizeInTbs)
				state.DbWorkload = string(pointer.From(props.DbWorkload))
				state.DbVersion = pointer.From(props.DbVersion)
				state.DisplayName = pointer.From(props.DisplayName)
				state.LicenseModel = string(pointer.From(props.LicenseModel))
				state.Location = result.Model.Location
				state.MtlsConnectionRequired = pointer.From(props.IsMtlsConnectionRequired)
				state.Name = pointer.ToString(result.Model.Name)
				state.NationalCharacterSet = pointer.From(props.NcharacterSet)
				state.SubnetId = pointer.From(props.SubnetId)
				state.Tags = pointer.From(result.Model.Tags)
				state.VnetId = pointer.From(props.VnetId)
				state.LongTermBackUpSchedule = FlattenLongTermBackUpScheduleDetails(props.LongTermBackupSchedule)
				state.AllowedIps = pointer.From(props.WhitelistedIPs)
				state.Ocid = pointer.From(props.Ocid)
			}
			return metadata.Encode(&state)
		},
	}
}

func (AutonomousDatabaseRegularResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases

			id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (AutonomousDatabaseRegularResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabases.ValidateAutonomousDatabaseID
}

func expandAdbsCustomerContacts(customerContactsList []string) []autonomousdatabases.CustomerContact {
	customerContacts := make([]autonomousdatabases.CustomerContact, 0, len(customerContactsList))
	for _, customerContact := range customerContactsList {
		customerContacts = append(customerContacts, autonomousdatabases.CustomerContact{
			Email: customerContact,
		})
	}
	return customerContacts
}

func flattenAdbsCustomerContacts(customerContactsList *[]autonomousdatabases.CustomerContact) []string {
	var customerContacts []string
	if customerContactsList != nil {
		for _, customerContact := range *customerContactsList {
			customerContacts = append(customerContacts, customerContact.Email)
		}
	}
	return customerContacts
}

func expandLongTermBackupSchedule(input []LongTermBackUpScheduleDetails) *autonomousdatabases.LongTermBackUpScheduleDetails {
	if len(input) == 0 {
		return nil
	}
	schedule := input[0]
	return &autonomousdatabases.LongTermBackUpScheduleDetails{
		RepeatCadence:         pointer.To(autonomousdatabases.RepeatCadenceType(schedule.RepeatCadence)),
		TimeOfBackup:          pointer.To(schedule.TimeOfBackup),
		RetentionPeriodInDays: pointer.To(schedule.RetentionPeriodInDays),
		IsDisabled:            pointer.To(!schedule.Enabled),
	}
}

func (r AutonomousDatabaseRegularResource) hasGeneralUpdates(metadata sdk.ResourceMetaData) bool {
	return metadata.ResourceData.HasChange("tags") ||
		metadata.ResourceData.HasChange("backup_retention_period_in_days") ||
		metadata.ResourceData.HasChange("data_storage_size_in_tbs") ||
		metadata.ResourceData.HasChange("compute_count") ||
		metadata.ResourceData.HasChange("auto_scaling_enabled") ||
		metadata.ResourceData.HasChange("auto_scaling_for_storage_enabled") ||
		metadata.ResourceData.HasChange("allowed_ips")
}
